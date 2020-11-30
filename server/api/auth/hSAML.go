package auth

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
	logger "github.com/sirupsen/logrus"

	saml2 "github.com/russellhaering/gosaml2"
	"github.com/russellhaering/gosaml2/types"
	dsig "github.com/russellhaering/goxmldsig"

	"github.com/go-chi/chi"

	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"

	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/providers/uidp"
)

// SAMLLoginHandler handles SAML login request. IF validation is failed, return 403 response.
// If validation succeeds, respond with TRASA session response (csrf and session tokens)
func SAMLLoginHandler(w http.ResponseWriter, r *http.Request) {

	authlog := logs.NewLog(r, "dashboard")

	err := r.ParseForm()
	if err != nil {
		logger.Error(err)
		setFailedResponse(w, consts.SAML_INVALID_POSTFORM)
		err = logs.Store.LogLogin(&authlog, consts.SAML_INVALID_POSTFORM, false)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	orgID := chi.URLParam(r, "orgid")
	vendorName := chi.URLParam(r, "vendorname")

	logger.Debugf("orgID: %s , vendorName: %s", orgID, vendorName)

	storedIdp, err := uidp.Store.GetByName(orgID, vendorName)
	if err != nil {
		logger.Error(err)
		setFailedResponse(w, consts.UIDP_NOT_REGISTERED)
		err = logs.Store.LogLogin(&authlog, consts.UIDP_NOT_REGISTERED, false)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	metadata := &types.EntityDescriptor{}
	err = xml.Unmarshal([]byte(storedIdp.IdpMeta), metadata)
	if err != nil {
		logger.Error(err)
		setFailedResponse(w, consts.SAML_METADATA_ERROR)
		err = logs.Store.LogLogin(&authlog, consts.SAML_METADATA_ERROR, false)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	certStore := dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}

	for _, kd := range metadata.IDPSSODescriptor.KeyDescriptors {
		for _, xcert := range kd.KeyInfo.X509Data.X509Certificates {
			if xcert.Data == "" {
				logger.Error(err)
				setFailedResponse(w, consts.SAML_PARSE_CERT)
				err = logs.Store.LogLogin(&authlog, consts.SAML_PARSE_CERT, false)
				if err != nil {
					logrus.Error(err)
				}
				return
			}
			certData, err := base64.StdEncoding.DecodeString(xcert.Data)
			if err != nil {
				logger.Error(err)
				setFailedResponse(w, consts.SAML_PARSE_CERT)
				err = logs.Store.LogLogin(&authlog, consts.SAML_PARSE_CERT, false)
				if err != nil {
					logrus.Error(err)
				}
				return
			}

			idpCert, err := x509.ParseCertificate(certData)
			if err != nil {
				logger.Error(err)
				setFailedResponse(w, consts.SAML_PARSE_CERT)
				err = logs.Store.LogLogin(&authlog, consts.SAML_PARSE_CERT, false)
				if err != nil {
					logrus.Error(err)
				}
				return
			}

			certStore.Roots = append(certStore.Roots, idpCert)
		}
	}

	// We sign the AuthnRequest with a random key because Okta doesn't seem
	// to verify these.
	randomKeyStore := dsig.RandomKeyStoreForTest()

	sp := &saml2.SAMLServiceProvider{
		IdentityProviderSSOURL:      metadata.IDPSSODescriptor.SingleSignOnServices[0].Location,
		IdentityProviderIssuer:      metadata.EntityID,
		ServiceProviderIssuer:       storedIdp.RedirectURL,
		AssertionConsumerServiceURL: storedIdp.RedirectURL,
		SignAuthnRequests:           true,
		AudienceURI:                 storedIdp.AudienceURI,
		IDPCertificateStore:         &certStore,
		SPKeyStore:                  randomKeyStore,
	}

	assertionInfo, err := sp.RetrieveAssertionInfo(r.FormValue("SAMLResponse"))
	if err != nil {
		setFailedResponse(w, consts.SAML_INVALID_ASSERTION_INFO)
		logger.Error(err)
		err = logs.Store.LogLogin(&authlog, consts.SAML_INVALID_ASSERTION_INFO, false)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	if assertionInfo.WarningInfo.InvalidTime {
		setFailedResponse(w, consts.SAML_INVALID_TIME)
		err = logs.Store.LogLogin(&authlog, consts.SAML_INVALID_TIME, false)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	if assertionInfo.WarningInfo.NotInAudience {
		setFailedResponse(w, consts.SAML_AUDIENCE_MISMATCH)
		err = logs.Store.LogLogin(&authlog, consts.SAML_AUDIENCE_MISMATCH, false)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	// if we are here, it means saml validation is success.
	// we should now extract user specefic info's, build one time session token and bind user detail to this token.

	// get user info from database
	userWoPass, err := Store.GetLoginDetails(assertionInfo.NameID, orgID)
	if err != nil {
		logrus.Error(err)
		err = logs.Store.LogLogin(&authlog, consts.REASON_USER_NOT_FOUND, false)
		if err != nil {
			logrus.Error(err)
		}
		setFailedResponse(w, "USER_NOT_FOUND")
		return
	}

	authlog.UpdateUser(userWoPass)

	userWoPass.Email = assertionInfo.NameID

	userWoPass.Groups = assertionInfo.Values.GetAll(("groups"))

	userWoPass.OrgID = orgID

	user := models.CopyUserWithoutPass(*userWoPass)
	//Update authlog value with user fields
	authlog.UpdateUser(&models.UserWithPass{User: user})

	var uc models.UserContext
	uc.User = &user

	org, err := orgs.Store.Get(user.OrgID)
	if err != nil {
		setFailedResponse(w, "FAILED_TO_GET_ORG_DETAILS")
		return
	}

	uc.Org = org
	uc.DeviceID = ""
	uc.BrowserID = ""

	// we do not need to handle intent response here since all request will be for dashboard login?
	sessionToken, csrfToken, err := SetSession(uc)
	if err != nil {
		setFailedResponse(w, "could not set user session")
		err = logs.Store.LogLogin(&authlog, "failed to create session", false)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	var response userAuthSessionResp
	response.User = user
	response.CSRFToken = csrfToken

	err = logs.Store.LogLogin(&authlog, "", true)
	if err != nil {
		logrus.Error(err)
	}

	// we set session token in HTTPonly cookie and expect csrf token in http header.
	xSESSION := http.Cookie{
		Name:     "X-SESSION",
		Value:    sessionToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
	}

	// we set csrf token in cookie. Unlike PasswordLogin flow, we cant send this token in http response body.
	// Dashboard should check for this cookie in response interceptor, store it in local storage and delete cookie from browser.
	xCSRF := http.Cookie{
		Name:     "X-CSRF",
		Value:    csrfToken,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Expires:  time.Now().Add(20 * time.Second),
		Path:     "/",
	}

	http.SetCookie(w, &xSESSION)
	http.SetCookie(w, &xCSRF)

	// check user role and redirect to authorized page (overview for admin, my page for normal user)
	redirectURL := fmt.Sprintf("https://%s/overview", global.GetConfig().Trasa.ListenAddr)
	if user.UserRole == utils.NormalizeString("selfuser") {
		redirectURL = fmt.Sprintf("https://%s/my", global.GetConfig().Trasa.ListenAddr)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
	return

}

var idpVendors = map[string]string{
	"Okta":      "saml",
	"Jumpcloud": "saml",
}

// setFailedResponse returns failed saml login page
func setFailedResponse(w http.ResponseWriter, reason string) {

	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(403)
	t := template.New("samlFailedRespTemplate")

	t, err := t.Parse(samlFailedRespTemplate)
	if err != nil {
		logrus.Error(err)
	}
	type tmplVal struct {
		Reason string
	}

	err = t.Execute(w, tmplVal{
		Reason: reason,
	})
	if err != nil {
		logrus.Error(err)
	}

	return

}

var samlFailedRespTemplate string = `<!DOCTYPE html>
<html>
<head>
<style>
.item1 { grid-area: header; }
.item2 { grid-area: menu; }
.item3 { grid-area: main; font-size: 25px;}
.item4 { grid-area: right; }
.item5 { grid-area: footer;  font-size: 15px;}
.reason {font-size: 15px; color: 'navy' }

.grid-container {
  display: grid;
  grid-template-areas:
    'header header header header header header'
    'main main main main main main'
    'footer footer footer footer footer footer';
  grid-gap: 10px;
  /* background-color: #2196F3; */
  padding: 10px;
}

.grid-container > div {
  /* background-color: rgba(255, 255, 255, 0.8); */
  text-align: center;
  padding: 20px 0;
  /* font-size: 30px; */
}
</style>
</head>
<body>
<br /><br /><br /><br />

<div class="grid-container">
  <div class="item1">
      <img src="data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+CjwhLS0gQ3JlYXRlZCB3aXRoIElua3NjYXBlIChodHRwOi8vd3d3Lmlua3NjYXBlLm9yZy8pIC0tPgoKPHN2ZwogICB4bWxuczpkYz0iaHR0cDovL3B1cmwub3JnL2RjL2VsZW1lbnRzLzEuMS8iCiAgIHhtbG5zOmNjPSJodHRwOi8vY3JlYXRpdmVjb21tb25zLm9yZy9ucyMiCiAgIHhtbG5zOnJkZj0iaHR0cDovL3d3dy53My5vcmcvMTk5OS8wMi8yMi1yZGYtc3ludGF4LW5zIyIKICAgeG1sbnM6c3ZnPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIKICAgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIgogICB4bWxuczpzb2RpcG9kaT0iaHR0cDovL3NvZGlwb2RpLnNvdXJjZWZvcmdlLm5ldC9EVEQvc29kaXBvZGktMC5kdGQiCiAgIHhtbG5zOmlua3NjYXBlPSJodHRwOi8vd3d3Lmlua3NjYXBlLm9yZy9uYW1lc3BhY2VzL2lua3NjYXBlIgogICB3aWR0aD0iNy43MTcxMjM1aW4iCiAgIGhlaWdodD0iMS4wMjQzOTMyaW4iCiAgIHZpZXdCb3g9IjAgMCAxOTYuMDE0OTQgMjYuMDE5NTg3IgogICB2ZXJzaW9uPSIxLjEiCiAgIGlkPSJzdmc1MzUzIgogICBpbmtzY2FwZTp2ZXJzaW9uPSIwLjkyLjMgKGQyNDRiOTUsIDIwMTgtMDgtMDIpIgogICBzb2RpcG9kaTpkb2NuYW1lPSJ0cmFzYS5zdmciPgogIDxkZWZzCiAgICAgaWQ9ImRlZnM1MzQ3IiAvPgogIDxzb2RpcG9kaTpuYW1lZHZpZXcKICAgICBpZD0iYmFzZSIKICAgICBwYWdlY29sb3I9IiNmZmZmZmYiCiAgICAgYm9yZGVyY29sb3I9IiM2NjY2NjYiCiAgICAgYm9yZGVyb3BhY2l0eT0iMS4wIgogICAgIGlua3NjYXBlOnBhZ2VvcGFjaXR5PSIwLjAiCiAgICAgaW5rc2NhcGU6cGFnZXNoYWRvdz0iMiIKICAgICBpbmtzY2FwZTp6b29tPSIxLjQiCiAgICAgaW5rc2NhcGU6Y3g9IjM3MC4wMjk5MSIKICAgICBpbmtzY2FwZTpjeT0iNTYuOTcxMDY1IgogICAgIGlua3NjYXBlOmRvY3VtZW50LXVuaXRzPSJtbSIKICAgICBpbmtzY2FwZTpjdXJyZW50LWxheWVyPSJsYXllcjEiCiAgICAgc2hvd2dyaWQ9ImZhbHNlIgogICAgIGlua3NjYXBlOnNob3dwYWdlc2hhZG93PSJmYWxzZSIKICAgICB1bml0cz0iaW4iCiAgICAgc2hvd2d1aWRlcz0idHJ1ZSIKICAgICBpbmtzY2FwZTpndWlkZS1iYm94PSJ0cnVlIgogICAgIGlua3NjYXBlOnNuYXAtb3RoZXJzPSJmYWxzZSIKICAgICBpbmtzY2FwZTpvYmplY3QtcGF0aHM9ImZhbHNlIgogICAgIGlua3NjYXBlOnNuYXAtbm9kZXM9ImZhbHNlIgogICAgIGlua3NjYXBlOnNuYXAtZ2xvYmFsPSJmYWxzZSIKICAgICBmaXQtbWFyZ2luLXRvcD0iMCIKICAgICBmaXQtbWFyZ2luLWxlZnQ9IjAiCiAgICAgZml0LW1hcmdpbi1yaWdodD0iMCIKICAgICBmaXQtbWFyZ2luLWJvdHRvbT0iMCIKICAgICBpbmtzY2FwZTp3aW5kb3ctd2lkdGg9IjE5MjAiCiAgICAgaW5rc2NhcGU6d2luZG93LWhlaWdodD0iMTA1MiIKICAgICBpbmtzY2FwZTp3aW5kb3cteD0iMCIKICAgICBpbmtzY2FwZTp3aW5kb3cteT0iMCIKICAgICBpbmtzY2FwZTp3aW5kb3ctbWF4aW1pemVkPSIxIj4KICAgIDxzb2RpcG9kaTpndWlkZQogICAgICAgcG9zaXRpb249Ijk2Ljc1MzYxLDI1Ljk4Mjg0NCIKICAgICAgIG9yaWVudGF0aW9uPSIwLDEiCiAgICAgICBpZD0iZ3VpZGU1MzU1IgogICAgICAgaW5rc2NhcGU6bG9ja2VkPSJmYWxzZSIgLz4KICAgIDxzb2RpcG9kaTpndWlkZQogICAgICAgcG9zaXRpb249Ijk4LjM1NzIzLDEuMTI2Nzc3OCIKICAgICAgIG9yaWVudGF0aW9uPSIwLDEiCiAgICAgICBpZD0iZ3VpZGU1MzU3IgogICAgICAgaW5rc2NhcGU6bG9ja2VkPSJmYWxzZSIgLz4KICAgIDxzb2RpcG9kaTpndWlkZQogICAgICAgcG9zaXRpb249IjAuMDAyMDM4NjQ4NCwxMy40MjExNzQiCiAgICAgICBvcmllbnRhdGlvbj0iMSwwIgogICAgICAgaWQ9Imd1aWRlNTM1OSIKICAgICAgIGlua3NjYXBlOmxvY2tlZD0iZmFsc2UiIC8+CiAgICA8c29kaXBvZGk6Z3VpZGUKICAgICAgIHBvc2l0aW9uPSIzNS4yODE2MTksMy43OTk0NzEiCiAgICAgICBvcmllbnRhdGlvbj0iMSwwIgogICAgICAgaWQ9Imd1aWRlNTM2MSIKICAgICAgIGlua3NjYXBlOmxvY2tlZD0iZmFsc2UiIC8+CiAgICA8c29kaXBvZGk6Z3VpZGUKICAgICAgIHBvc2l0aW9uPSI3NS42MzkzMjIsMC4zMjQ5NjY3NSIKICAgICAgIG9yaWVudGF0aW9uPSIxLDAiCiAgICAgICBpZD0iZ3VpZGU1MzYzIgogICAgICAgaW5rc2NhcGU6bG9ja2VkPSJmYWxzZSIgLz4KICAgIDxzb2RpcG9kaTpndWlkZQogICAgICAgcG9zaXRpb249IjQwLjM1OTc0LC0xNC4xMDc1ODYiCiAgICAgICBvcmllbnRhdGlvbj0iMSwwIgogICAgICAgaWQ9Imd1aWRlNTM2NSIKICAgICAgIGlua3NjYXBlOmxvY2tlZD0iZmFsc2UiIC8+CiAgICA8c29kaXBvZGk6Z3VpZGUKICAgICAgIHBvc2l0aW9uPSI4MC43MTc0NDMsMTMuMTUzOTA1IgogICAgICAgb3JpZW50YXRpb249IjEsMCIKICAgICAgIGlkPSJndWlkZTUzODEiCiAgICAgICBpbmtzY2FwZTpsb2NrZWQ9ImZhbHNlIiAvPgogICAgPHNvZGlwb2RpOmd1aWRlCiAgICAgICBwb3NpdGlvbj0iMTE1LjcyOTc0LC0xLjAxMTM4MTMiCiAgICAgICBvcmllbnRhdGlvbj0iMSwwIgogICAgICAgaWQ9Imd1aWRlNTM4MyIKICAgICAgIGlua3NjYXBlOmxvY2tlZD0iZmFsc2UiIC8+CiAgICA8c29kaXBvZGk6Z3VpZGUKICAgICAgIHBvc2l0aW9uPSIxMjAuODA3ODYsMTkuMDMzODM1IgogICAgICAgb3JpZW50YXRpb249IjEsMCIKICAgICAgIGlkPSJndWlkZTUzOTUiCiAgICAgICBpbmtzY2FwZTpsb2NrZWQ9ImZhbHNlIiAvPgogICAgPHNvZGlwb2RpOmd1aWRlCiAgICAgICBwb3NpdGlvbj0iMTU1LjU1MjkxLDkuNDEyMTMwMiIKICAgICAgIG9yaWVudGF0aW9uPSIxLDAiCiAgICAgICBpZD0iZ3VpZGU1Mzk3IgogICAgICAgaW5rc2NhcGU6bG9ja2VkPSJmYWxzZSIgLz4KICAgIDxzb2RpcG9kaTpndWlkZQogICAgICAgcG9zaXRpb249IjE2MC42MzEwNCwxMS4wMTU3NDgiCiAgICAgICBvcmllbnRhdGlvbj0iMSwwIgogICAgICAgaWQ9Imd1aWRlNTQwOSIKICAgICAgIGlua3NjYXBlOmxvY2tlZD0iZmFsc2UiIC8+CiAgICA8c29kaXBvZGk6Z3VpZGUKICAgICAgIHBvc2l0aW9uPSIxOTUuNjQzMzQsLTEuMjc4NjUwNiIKICAgICAgIG9yaWVudGF0aW9uPSIxLDAiCiAgICAgICBpZD0iZ3VpZGU1NDExIgogICAgICAgaW5rc2NhcGU6bG9ja2VkPSJmYWxzZSIgLz4KICAgIDxzb2RpcG9kaTpndWlkZQogICAgICAgcG9zaXRpb249IjEwNy4yNzc4MiwyMC42NzQ2NSIKICAgICAgIG9yaWVudGF0aW9uPSIwLDEiCiAgICAgICBpZD0iZ3VpZGU1NDIxIgogICAgICAgaW5rc2NhcGU6bG9ja2VkPSJmYWxzZSIgLz4KICAgIDxzb2RpcG9kaTpndWlkZQogICAgICAgcG9zaXRpb249IjEyMC42OTU5NywxNi4zMjc5MjQiCiAgICAgICBvcmllbnRhdGlvbj0iMCwxIgogICAgICAgaWQ9Imd1aWRlNTQzMSIKICAgICAgIGlua3NjYXBlOmxvY2tlZD0iZmFsc2UiIC8+CiAgICA8c29kaXBvZGk6Z3VpZGUKICAgICAgIHBvc2l0aW9uPSIxMzEuMjc5MzIsMTEuNzkyMjA5IgogICAgICAgb3JpZW50YXRpb249IjAsMSIKICAgICAgIGlkPSJndWlkZTU0MzMiCiAgICAgICBpbmtzY2FwZTpsb2NrZWQ9ImZhbHNlIiAvPgogICAgPHNvZGlwb2RpOmd1aWRlCiAgICAgICBwb3NpdGlvbj0iNTkuNzE0MzM3LDE2LjM5NDcxOCIKICAgICAgIG9yaWVudGF0aW9uPSIwLDEiCiAgICAgICBpZD0iZ3VpZGU1Njc1IgogICAgICAgaW5rc2NhcGU6bG9ja2VkPSJmYWxzZSIgLz4KICA8L3NvZGlwb2RpOm5hbWVkdmlldz4KICA8bWV0YWRhdGEKICAgICBpZD0ibWV0YWRhdGE1MzUwIj4KICAgIDxyZGY6UkRGPgogICAgICA8Y2M6V29yawogICAgICAgICByZGY6YWJvdXQ9IiI+CiAgICAgICAgPGRjOmZvcm1hdD5pbWFnZS9zdmcreG1sPC9kYzpmb3JtYXQ+CiAgICAgICAgPGRjOnR5cGUKICAgICAgICAgICByZGY6cmVzb3VyY2U9Imh0dHA6Ly9wdXJsLm9yZy9kYy9kY21pdHlwZS9TdGlsbEltYWdlIiAvPgogICAgICAgIDxkYzp0aXRsZT48L2RjOnRpdGxlPgogICAgICA8L2NjOldvcms+CiAgICA8L3JkZjpSREY+CiAgPC9tZXRhZGF0YT4KICA8ZwogICAgIGlua3NjYXBlOmxhYmVsPSJMYXllciAxIgogICAgIGlua3NjYXBlOmdyb3VwbW9kZT0ibGF5ZXIiCiAgICAgaWQ9ImxheWVyMSIKICAgICB0cmFuc2Zvcm09InRyYW5zbGF0ZSgtMTQuNDMwNTE2LC0xOTcuNTM5MDIpIj4KICAgIDxyZWN0CiAgICAgICBzdHlsZT0iZmlsbDojMGIxNzI4O2ZpbGwtb3BhY2l0eToxO3N0cm9rZS13aWR0aDowLjI5MzI1Nzc0IgogICAgICAgaWQ9InJlY3Q1MzY3IgogICAgICAgd2lkdGg9IjM1LjI4Njk5MSIKICAgICAgIGhlaWdodD0iNS4yOTIxMjMzIgogICAgICAgeD0iMTQuNDMwNTE2IgogICAgICAgeT0iMTk3LjU4MzE5IiAvPgogICAgPHJlY3QKICAgICAgIHN0eWxlPSJmaWxsOiMwYjE3Mjg7ZmlsbC1vcGFjaXR5OjE7c3Ryb2tlLXdpZHRoOjAuMzE2NzU3MTQiCiAgICAgICBpZD0icmVjdDUzNjkiCiAgICAgICB3aWR0aD0iNS44Nzk5Mjk1IgogICAgICAgaGVpZ2h0PSIyMC4zMzc0OTgiCiAgICAgICB4PSIyOS42NjY5MTYiCiAgICAgICB5PSIyMDIuNjUzOSIgLz4KICAgIDxyZWN0CiAgICAgICBzdHlsZT0iZmlsbDojMGIxNzI4O2ZpbGwtb3BhY2l0eToxO3N0cm9rZS13aWR0aDowLjI2NDU4MzMyIgogICAgICAgaWQ9InJlY3Q1Mzc1IgogICAgICAgd2lkdGg9IjMuNDc0NTAzOCIKICAgICAgIGhlaWdodD0iNy4yMTYyNzcxIgogICAgICAgeD0iODYuNTk1MzI5IgogICAgICAgeT0iMjAyLjkyMTE3IiAvPgogICAgPHBhdGgKICAgICAgIGlkPSJyZWN0NTM3OSIKICAgICAgIHRyYW5zZm9ybT0ibWF0cml4KDAuMjY0NTgzMzMsMCwwLDAuMjY0NTgzMzMsMTQuNDMwNTE2LDE5Ny41MzkwMikiCiAgICAgICBzdHlsZT0iZmlsbDojMGIxNzI4O2ZpbGwtb3BhY2l0eToxO3N0cm9rZS13aWR0aDoxLjc4MDA2Mzg3IgogICAgICAgZD0ibSAyMjkuMzI4MTIsMzYuMzU1NDY5IGMgLTAuMjAwNzYsMC4wMjkxOSAtMC40MTY4OCwwLjA0MzczIC0wLjY1MDM5LDAuMDQxMDIgLTAuMDE5MiwwLjAwMTYgLTAuMDMyNSwwLjAwOTEgLTAuMDUyNywwLjAwOTggLTAuMDUwMSwwLjAwMTYgLTAuMTAwMjcsMC4wMDMyIC0wLjE1MDM5LDAuMDAzOSAtMC4wNzQ0LDAuMDAxNSAtMC4xNDgyOCwwLjAwMzIgLTAuMjIyNjYsMC4wMDM5IC0wLjA3MzcsMC4wMDE4IC0wLjE0NywwLjAwMzggLTAuMjIwNywwLjAwNzggLTAuMDcyMywwLjAwMjIgLTAuMTQ0NDYsMC4wMDQgLTAuMjE2OCwwLjAwNzggLTAuMDc3NCwwLjAwMjYgLTAuMTU0OTQsMC4wMDM5IC0wLjIzMjQyLDAuMDAzOSAtMC4wNzQsLTMuOGUtNSAtMC4xNDg2NiwtMC4wMDIgLTAuMjIyNjUsLTAuMDAyIGggLTAuMjMwNDcgLTAuMjI4NTIgYyAtMC4wNzcyLC0zLjhlLTUgLTAuMTU1MjEsMy43ZS01IC0wLjIzMjQyLDAgaCAtMC4yMTg3NSBjIC0wLjA3MDIsLTMuOGUtNSAtMC4xNDA3MiwwIC0wLjIxMDk0LDAgaCAtMC4yMDExNyAtMC4yMDMxMyBjIC0wLjA2NzQsLTMuOGUtNSAtMC4xMzU2OSw3LjVlLTUgLTAuMjAzMTIsMCAtMC4wNjc1LC0zLjhlLTUgLTAuMTMzNzEsMy43ZS01IC0wLjIwMTE3LDAgaCAtMC4yMDUwOCAtMC4yMDg5OSBjIC0wLjA3MTYsLTEuMTRlLTQgLTAuMTQzMTgsMS4xM2UtNCAtMC4yMTQ4NCwwIC0wLjA2OTYsLTEuMTRlLTQgLTAuMTM5MzYsMS41MWUtNCAtMC4yMDg5OCwwIC0wLjA2NzksLTEuNTJlLTQgLTAuMTM3MiwtMC4wMDE4IC0wLjIwNTA4LC0wLjAwMiAtMC4wNjg3LC0xLjUyZS00IC0wLjEzNjM3LDEuODdlLTQgLTAuMjA1MDgsMCAtMC4wNzAzLC0xLjg5ZS00IC0wLjE0MDY4LDIuMjVlLTQgLTAuMjEwOTQsMCAtMC4wNjgsLTIuMjdlLTQgLTAuMTM3MDQsLTAuMDAxNyAtMC4yMDUwOCwtMC4wMDIgLTAuMDY1MSwtMy4wMmUtNCAtMC4xMzAxOSwzLjQxZS00IC0wLjE5NTMxLDAgLTAuMDYzMSwtMy43OGUtNCAtMC4xMjQ0MiwtMC4wMDE1IC0wLjE4NzUsLTAuMDAyIC0wLjA2NDEsLTQuMTVlLTQgLTAuMTI5MjksLTAuMDAxNCAtMC4xOTMzNiwtMC4wMDIgLTAuMDY0MiwtNS4yOWUtNCAtMC4xMjcyNiw2LjA1ZS00IC0wLjE5MTQsMCAtMC4wNjM3LC02LjQyZS00IC0wLjEyNzY5LC0wLjAwMTIgLTAuMTkxNDEsLTAuMDAyIC0wLjA2MzYsLTcuNTZlLTQgLTAuMTI3NzYsLTAuMDAzIC0wLjE5MTQxLC0wLjAwMzkgLTAuMDYxOSwtOS40NWUtNCAtMC4xMjM2LC04LjJlLTQgLTAuMTg1NTQsLTAuMDAyIC0wLjA2MTksLTAuMDAxMSAtMC4xMjM2OCwtMC4wMDI1IC0wLjE4NTU1LC0wLjAwMzkgLTAuMDYwNCwtMC4wMDE0IC0wLjEyMTI0LC0wLjAwMjkgLTAuMTgxNjQsLTAuMDAzOSAtMC4wNTg0LC0wLjAwMTcgLTAuMTE3NDMsLTAuMDAzNSAtMC4xNzU3OCwtMC4wMDM5IC0wLjA1NTcsLTAuMDAyIC0wLjExMDMxLC0wLjAwNCAtMC4xNjYwMiwtMC4wMDc4IC0wLjA1NDUsLTAuMDAyNSAtMC4xMDk2LC0wLjAwNCAtMC4xNjQwNiwtMC4wMDc4IC0wLjA0MywtMC4wMDI0IC0wLjA4NTksLTAuMDA2NSAtMC4xMjg5MSwtMC4wMDk4IC0wLjA0MDksMC4wMDI0IC0wLjA4MTcsMC4wMDU5IC0wLjEyMzA0LDAuMDA1OSAtMC4wNTE3LDAuMDAxMiAtMC4xMDQ1NSwwLjAwMjYgLTAuMTU2MjUsMC4wMDM5IC0wLjA1MDksMC4wMDE1IC0wLjEwMTQ0LDAuMDAzMSAtMC4xNTIzNSwwLjAwMzkgLTAuMDQ5NiwwLjAwMTcgLTAuMDk4OCwwLjAwMzggLTAuMTQ4NDMsMC4wMDc4IC0wLjA0OTQsMC4wMDIyIC0wLjA5OTEsMC4wMDQgLTAuMTQ4NDQsMC4wMDc4IC0wLjA1MjksMC4wMDIgLTAuMTA1MzMsMC4wMDE1IC0wLjE1ODIsMC4wMDIgLTAuMDQ4NSwtMy4wMmUtNCAtMC4wOTgxLDUuMTNlLTQgLTAuMTQ2NDksLTAuMDAyIC0wLjA0NTEsLTAuMDAxOSAtMC4wODk3LC0wLjAwNCAtMC4xMzQ3NiwtMC4wMDc4IC0wLjA1MjEsLTAuMDAzOCAtMC4xMDI4NiwtMC4wMTE5NyAtMC4xNTQzLC0wLjAxOTUzIC0wLjAwMSwtMi41NWUtNCAtMC4wMDMsLTAuMDAxNyAtMC4wMDQsLTAuMDAyIC0wLjE2ODc3LDAuMDEwNjEgLTAuMzM4ODYsMC4wMTc4MyAtMC41MDc4MSwwLjAyMTQ4IC0wLjE4MDcsMC4wMDQgLTAuMzYwMjcsMC4wMDc5IC0wLjU0MTAyLDAuMDAzOSAtMC4wODQxLC0wLjAwMyAtMC4xNjc4OCwtMC4wMTI0OCAtMC4yNTE5NSwtMC4wMTc1OCBsIDIuNjE1MjQsMTguNTkzNzUgNjYuNjQwNjIsNDMuMzYxMzI4IC0zLjY3OTY5LC0yNi4xNjc5NjkgeiIgLz4KICAgIDxyZWN0CiAgICAgICBzdHlsZT0iZmlsbDojMGIxNzI4O2ZpbGwtb3BhY2l0eToxO3N0cm9rZS13aWR0aDowLjI2ODYwMzM1IgogICAgICAgaWQ9InJlY3Q1Mzg1IgogICAgICAgd2lkdGg9IjUuODg4MTE1NCIKICAgICAgIGhlaWdodD0iMjQuNzU2MyIKICAgICAgIHg9IjEzOC4xOTU1NCIKICAgICAgIHk9IjE3MS4wNTk5NyIKICAgICAgIHRyYW5zZm9ybT0ibWF0cml4KDAuOTc1OTE2NCwwLjIxODE0NDg2LC0wLjIwNTk2MzkyLDAuOTc4NTU5NTksMCwwKSIgLz4KICAgIDxyZWN0CiAgICAgICBzdHlsZT0iZmlsbDojMGIxNzI4O2ZpbGwtb3BhY2l0eToxO3N0cm9rZS13aWR0aDowLjI2NDU4MzMyIgogICAgICAgaWQ9InJlY3Q1Mzg3IgogICAgICAgd2lkdGg9IjUuMDc4MTIwNyIKICAgICAgIGhlaWdodD0iMjUuMTIzMzM1IgogICAgICAgeD0iNzIuNTkwNzkiCiAgICAgICB5PSIyMjAuOTAzNyIKICAgICAgIHRyYW5zZm9ybT0icm90YXRlKC0xMy4wODAzMDQpIiAvPgogICAgPHJlY3QKICAgICAgIHN0eWxlPSJmaWxsOiMwYjE3Mjg7ZmlsbC1vcGFjaXR5OjE7c3Ryb2tlLXdpZHRoOjAuMjY0NTgzMzIiCiAgICAgICBpZD0icmVjdDUzOTMiCiAgICAgICB3aWR0aD0iMTQuOTY3MDkzIgogICAgICAgaGVpZ2h0PSI0LjU0MzU4MiIKICAgICAgIHg9IjEwNS41NTY2MyIKICAgICAgIHk9IjIwNy4yNTA3NSIgLz4KICAgIDxyZWN0CiAgICAgICBzdHlsZT0iZmlsbDojMGIxNzI4O2ZpbGwtb3BhY2l0eToxO3N0cm9rZS13aWR0aDowLjI2NDU4MzMyIgogICAgICAgaWQ9InJlY3Q1NDAxIgogICAgICAgd2lkdGg9IjQuMDA5MDQyNyIKICAgICAgIGhlaWdodD0iNi45NDkwMDc1IgogICAgICAgeD0iMTM1LjI0Mzc2IgogICAgICAgeT0iMjAyLjExOTM1IiAvPgogICAgPHJlY3QKICAgICAgIHN0eWxlPSJmaWxsOiMwYjE3Mjg7ZmlsbC1vcGFjaXR5OjE7c3Ryb2tlLXdpZHRoOjAuMjY3NjgzNzEiCiAgICAgICBpZD0icmVjdDU0MDMiCiAgICAgICB3aWR0aD0iMzQuODM5NTMxIgogICAgICAgaGVpZ2h0PSI0LjYzODA3NTgiCiAgICAgICB4PSIxMzUuMjIyMTUiCiAgICAgICB5PSIyMDcuMjU5NTUiIC8+CiAgICA8cmVjdAogICAgICAgdHJhbnNmb3JtPSJtYXRyaXgoMC45NzYwMDIzNCwwLjIxNzc2MDAxLC0wLjIwNjMyOTg0LDAuOTc4NDgyNSwwLDApIgogICAgICAgeT0iMTUzLjkzMTI2IgogICAgICAgeD0iMjE1LjY4MzQiCiAgICAgICBoZWlnaHQ9IjI0LjcxMjM5NSIKICAgICAgIHdpZHRoPSI1Ljg4NzU5NzEiCiAgICAgICBpZD0icmVjdDU0MTMiCiAgICAgICBzdHlsZT0iZmlsbDojMGIxNzI4O2ZpbGwtb3BhY2l0eToxO3N0cm9rZS13aWR0aDowLjI2ODM1MzI1IiAvPgogICAgPHJlY3QKICAgICAgIHRyYW5zZm9ybT0icm90YXRlKC0xMy4wODAzMDQpIgogICAgICAgeT0iMjM4LjgwOCIKICAgICAgIHg9IjE0OS42NDk5MiIKICAgICAgIGhlaWdodD0iMjUuMTIzMzM1IgogICAgICAgd2lkdGg9IjUuMDc4MTIwNyIKICAgICAgIGlkPSJyZWN0NTQxNSIKICAgICAgIHN0eWxlPSJmaWxsOiMwYjE3Mjg7ZmlsbC1vcGFjaXR5OjE7c3Ryb2tlLXdpZHRoOjAuMjY0NTgzMzIiIC8+CiAgICA8cmVjdAogICAgICAgeT0iMjA3LjIyNzExIgogICAgICAgeD0iMTg0LjUyNjY5IgogICAgICAgaGVpZ2h0PSI0LjU0MzU4MiIKICAgICAgIHdpZHRoPSIxNC45NjcwOTMiCiAgICAgICBpZD0icmVjdDU0MTkiCiAgICAgICBzdHlsZT0iZmlsbDojMGIxNzI4O2ZpbGwtb3BhY2l0eToxO3N0cm9rZS13aWR0aDowLjI2NDU4MzMyIiAvPgogICAgPHJlY3QKICAgICAgIHk9IjE5Ny41ODMxOSIKICAgICAgIHg9IjU0Ljc1NTg0OCIKICAgICAgIGhlaWdodD0iNS4zMTU3NDY4IgogICAgICAgd2lkdGg9IjM1LjMzNDIzNiIKICAgICAgIGlkPSJyZWN0NTQyMyIKICAgICAgIHN0eWxlPSJmaWxsOiMwYjE3Mjg7ZmlsbC1vcGFjaXR5OjE7c3Ryb2tlLXdpZHRoOjAuMjk0MTA4MjQiIC8+CiAgICA8cmVjdAogICAgICAgc3R5bGU9ImZpbGw6IzBiMTcyODtmaWxsLW9wYWNpdHk6MTtzdHJva2Utd2lkdGg6MC4yNTA3MDk0NCIKICAgICAgIGlkPSJyZWN0NTQyNSIKICAgICAgIHdpZHRoPSIyNS43OTAzMzciCiAgICAgICBoZWlnaHQ9IjUuMjkyMTIzMyIKICAgICAgIHg9Ijk5LjcxMTM5NSIKICAgICAgIHk9IjE5Ny41ODMxOSIgLz4KICAgIDxyZWN0CiAgICAgICBzdHlsZT0iZmlsbDojMGIxNzI4O2ZpbGwtb3BhY2l0eToxO3N0cm9rZS13aWR0aDowLjI5MDk5MTI1IgogICAgICAgaWQ9InJlY3Q1NDI3IgogICAgICAgd2lkdGg9IjM0Ljc0MzY0OSIKICAgICAgIGhlaWdodD0iNS4yOTIxMjMzIgogICAgICAgeD0iMTM1LjI0MTE1IgogICAgICAgeT0iMTk3LjU4MzE5IiAvPgogICAgPHJlY3QKICAgICAgIHk9IjE5Ny42MDY4MSIKICAgICAgIHg9IjE3OC43MTc4NSIKICAgICAgIGhlaWdodD0iNS4yNjg0OTk5IgogICAgICAgd2lkdGg9IjI1LjkwODQ1NSIKICAgICAgIGlkPSJyZWN0NTQyOSIKICAgICAgIHN0eWxlPSJmaWxsOiMwYjE3Mjg7ZmlsbC1vcGFjaXR5OjE7c3Ryb2tlLXdpZHRoOjAuMjUwNzIxNDIiIC8+CiAgICA8cGF0aAogICAgICAgaWQ9InJlY3Q1NDM1IgogICAgICAgdHJhbnNmb3JtPSJtYXRyaXgoMC4yNjQ1ODMzMywwLDAsMC4yNjQ1ODMzMywxNC40MzA1MTYsMTk3LjUzOTAyKSIKICAgICAgIHN0eWxlPSJmaWxsOiMwYjE3Mjg7ZmlsbC1vcGFjaXR5OjE7c3Ryb2tlLXdpZHRoOjEuMDEyMzA1MzgiCiAgICAgICBkPSJNIDE1NS44MjAzMSwzNi4zODA4NTkgViA1NC4xNDg0MzggSCAyODUuODkwNjIgViAzNi4zODA4NTkgSCAyMjkuMDYyNSBjIC0wLjEyNjUsMC4wMTAzMSAtMC4yNTU1MSwwLjAxODU0IC0wLjM5MjU4LDAuMDE3NTggLTAuMDE2NCwwLjAwMTIgLTAuMDI3NywwLjAwNzMgLTAuMDQ0OSwwLjAwNzggLTAuMDUwMSwwLjAwMTYgLTAuMTAwMjcsMC4wMDMyIC0wLjE1MDM5LDAuMDAzOSAtMC4wNzQ0LDAuMDAxNSAtMC4xNDgyOCwwLjAwMzIgLTAuMjIyNjYsMC4wMDM5IC0wLjA3MzcsMC4wMDE4IC0wLjE0NywwLjAwMzggLTAuMjIwNywwLjAwNzggLTAuMDcyMywwLjAwMjIgLTAuMTQ0NDYsMC4wMDQgLTAuMjE2OCwwLjAwNzggLTAuMDc3NCwwLjAwMjYgLTAuMTU0OTQsMC4wMDM5IC0wLjIzMjQyLDAuMDAzOSAtMC4wNzQsLTMuOGUtNSAtMC4xNDg2NiwtMC4wMDIgLTAuMjIyNjUsLTAuMDAyIGggLTAuMjMwNDcgLTAuMjI4NTIgYyAtMC4wNzcyLC0zLjhlLTUgLTAuMTU1MjEsMy43ZS01IC0wLjIzMjQyLDAgaCAtMC4yMTg3NSBjIC0wLjA3MDIsLTMuOGUtNSAtMC4xNDA3MiwwIC0wLjIxMDk0LDAgaCAtMC4yMDExNyAtMC4yMDMxMyBjIC0wLjA2NzQsLTMuOGUtNSAtMC4xMzU2OSw3LjVlLTUgLTAuMjAzMTIsMCAtMC4wNjc1LC0zLjhlLTUgLTAuMTMzNzEsMy43ZS01IC0wLjIwMTE3LDAgaCAtMC4yMDUwOCAtMC4yMDg5OSBjIC0wLjA3MTYsLTEuMTRlLTQgLTAuMTQzMTgsMS4xM2UtNCAtMC4yMTQ4NCwwIC0wLjA2OTYsLTEuMTRlLTQgLTAuMTM5MzYsMS41MWUtNCAtMC4yMDg5OCwwIC0wLjA2NzksLTEuNTJlLTQgLTAuMTM3MiwtMC4wMDE4IC0wLjIwNTA4LC0wLjAwMiAtMC4wNjg3LC0xLjUyZS00IC0wLjEzNjM3LDEuODdlLTQgLTAuMjA1MDgsMCAtMC4wNzAzLC0xLjg5ZS00IC0wLjE0MDY4LDIuMjVlLTQgLTAuMjEwOTQsMCAtMC4wNjgsLTIuMjdlLTQgLTAuMTM3MDQsLTAuMDAxNyAtMC4yMDUwOCwtMC4wMDIgLTAuMDY1MSwtMy4wMmUtNCAtMC4xMzAxOSwzLjQxZS00IC0wLjE5NTMxLDAgLTAuMDYzMSwtMy43OGUtNCAtMC4xMjQ0MiwtMC4wMDE1IC0wLjE4NzUsLTAuMDAyIC0wLjA2NDEsLTQuMTVlLTQgLTAuMTI5MjksLTAuMDAxNCAtMC4xOTMzNiwtMC4wMDIgLTAuMDY0MiwtNS4yOWUtNCAtMC4xMjcyNiw2LjA1ZS00IC0wLjE5MTQsMCAtMC4wNjM3LC02LjQyZS00IC0wLjEyNzY5LC0wLjAwMTIgLTAuMTkxNDEsLTAuMDAyIC0wLjA2MzYsLTcuNTZlLTQgLTAuMTI3NzYsLTAuMDAzIC0wLjE5MTQxLC0wLjAwMzkgLTAuMDYxOSwtOS40NWUtNCAtMC4xMjM2LC04LjJlLTQgLTAuMTg1NTQsLTAuMDAyIC0wLjA2MTksLTAuMDAxMSAtMC4xMjM2OCwtMC4wMDI1IC0wLjE4NTU1LC0wLjAwMzkgLTAuMDYwNCwtMC4wMDE0IC0wLjEyMTI0LC0wLjAwMjkgLTAuMTgxNjQsLTAuMDAzOSAtMC4wNTg0LC0wLjAwMTcgLTAuMTE3NDMsLTAuMDAzNSAtMC4xNzU3OCwtMC4wMDM5IC0wLjA1NTcsLTAuMDAyIC0wLjExMDMxLC0wLjAwNCAtMC4xNjYwMiwtMC4wMDc4IC0wLjA1NDUsLTAuMDAyNSAtMC4xMDk2LC0wLjAwNCAtMC4xNjQwNiwtMC4wMDc4IC0wLjAzNTgsLTAuMDAyIC0wLjA3MTYsLTAuMDA1IC0wLjEwNzQyLC0wLjAwNzggaCAtMC4wNTQ3IGMgLTAuMDI5OCwwLjAwMTUgLTAuMDU5NiwwLjAwMzkgLTAuMDg5OCwwLjAwMzkgLTAuMDUxNywwLjAwMTIgLTAuMTA0NTUsMC4wMDI2IC0wLjE1NjI1LDAuMDAzOSAtMC4wNTA5LDAuMDAxNSAtMC4xMDE0NCwwLjAwMzEgLTAuMTUyMzUsMC4wMDM5IC0wLjA0OTYsMC4wMDE3IC0wLjA5ODgsMC4wMDM4IC0wLjE0ODQzLDAuMDA3OCAtMC4wNDk0LDAuMDAyMiAtMC4wOTkxLDAuMDA0IC0wLjE0ODQ0LDAuMDA3OCAtMC4wNTI5LDAuMDAyIC0wLjEwNTMzLDAuMDAxNSAtMC4xNTgyLDAuMDAyIC0wLjA0ODUsLTMuMDJlLTQgLTAuMDk4MSw1LjEzZS00IC0wLjE0NjQ5LC0wLjAwMiAtMC4wNDUxLC0wLjAwMTkgLTAuMDg5NywtMC4wMDQgLTAuMTM0NzYsLTAuMDA3OCAtMC4wNTIxLC0wLjAwMzggLTAuMTAyODYsLTAuMDExOTcgLTAuMTU0MywtMC4wMTk1MyBoIC0wLjA4MiBjIC0wLjE0MzEsMC4wMDc2IC0wLjI4NjczLDAuMDE1NzUgLTAuNDI5NjksMC4wMTk1MyAtMC4xODA3LDAuMDAzOCAtMC4zNjAyNiwwLjAwNzcgLTAuNTQxMDIsMC4wMDM5IC0wLjEwMTEsLTAuMDAzOCAtMC4yMDE3NCwtMC4wMTU4OCAtMC4zMDI3MywtMC4wMjM0NCAtMC4xNjI5NywwLjAxODkgLTAuMzMxMDQsMC4wMzAyNCAtMC41MTk1MywwIHoiIC8+CiAgICA8cmVjdAogICAgICAgeT0iMjAyLjYzMDI4IgogICAgICAgeD0iNTQuNzY0NTM4IgogICAgICAgaGVpZ2h0PSIxOS43NzA1MzUiCiAgICAgICB3aWR0aD0iNS44Nzk5Mjk1IgogICAgICAgaWQ9InJlY3Q1NDM3IgogICAgICAgc3R5bGU9ImZpbGw6IzBiMTcyODtmaWxsLW9wYWNpdHk6MTtzdHJva2Utd2lkdGg6MC4zMTIzMTA3IiAvPgogICAgPHJlY3QKICAgICAgIHN0eWxlPSJmaWxsOiMwYjE3Mjg7ZmlsbC1vcGFjaXR5OjE7c3Ryb2tlLXdpZHRoOjAuMjE2MTE4MTkiCiAgICAgICBpZD0icmVjdDU0MzkiCiAgICAgICB3aWR0aD0iNS44Nzk5Mjk1IgogICAgICAgaGVpZ2h0PSI5LjQ2NzMyNTIiCiAgICAgICB4PSIxNjQuMjI2NDQiCiAgICAgICB5PSIyMTMuMDA0MzYiIC8+CiAgICA8cmVjdAogICAgICAgeT0iMjE3LjMwODgyIgogICAgICAgeD0iMTM1LjEyMzAzIgogICAgICAgaGVpZ2h0PSI1LjEyNjc1ODYiCiAgICAgICB3aWR0aD0iMzQuNzkwODk3IgogICAgICAgaWQ9InJlY3Q1NDQxIgogICAgICAgc3R5bGU9ImZpbGw6IzBiMTcyODtmaWxsLW9wYWNpdHk6MTtzdHJva2Utd2lkdGg6MC4yODY2MDM0OCIgLz4KICA8L2c+Cjwvc3ZnPgo=" alt="TRASA LOGO" width="400" />
  </div>
  <!-- <div class="item2">Menu</div> -->
  <div class="item3">
      We could not verify your login request. 
      <div class="reason"> reason= {{.Reason}}</div>
    </div>  
  <!-- <div class="item4">Right</div> -->
  <div class="item5">
    
    Here's what you can do :<br />
    <ol>
    <li>Retry login again</li>    
		<li>If this problem persists, contact your administrator</li>  
    </ol>
   
  </div>

</div>

</body>
</html>
`

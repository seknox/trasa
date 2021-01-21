package uidp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/providers/vault/tsxvault"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// GetAllIdps retrieves all idps configured for organization
func GetAllIdps(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	storedIdp, err := Store.GetAllIdps(uc.User.OrgID)
	if err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to fetch IDP", "GetExternalIdpsForLogin-storedIdp", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "IDPs", "StoredIdps", storedIdp)
}

// GetAllIdpsWoa.
// TODO This should be rate limited when #63 is implemented
func GetAllIdpsWoa(w http.ResponseWriter, r *http.Request) {

	storedIdp, err := Store.GetAllIdpsWoa()
	if err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to fetch IDP", "GetExternalIdpsForLogin-storedIdp", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "IDPs", "StoredIdps", storedIdp)
}

// CreateIdp created Identity Provider
func CreateIdp(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	var idp models.IdentityProvider
	if err := json.NewDecoder(r.Body).Decode(&idp); err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "Could not create idp", nil)
		return
	}

	idp = PreConfiguredIdps(idp, uc)

	err := Store.CreateIDP(&idp)
	if err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to create IDP", "Could  not create IDP", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "IDP created", fmt.Sprintf(`IDP "%s" created`, idp.IdpName), idp)
}

// UpdateIDP
func UpdateIdp(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	var idp models.IdentityProvider
	if err := json.NewDecoder(r.Body).Decode(&idp); err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "Could not update IDP", nil)
		return
	}

	idp.CreatedBy = uc.User.ID
	idp.OrgID = uc.User.OrgID
	idp.LastUpdated = time.Now().Unix()
	idp.IsEnabled = true

	//	var err error
	if idp.IdpType == "saml" {
		idp.SCIMEndpoint = fmt.Sprintf("https://%s/scim/v2", global.GetConfig().Trasa.ListenAddr)
		idp.RedirectURL = fmt.Sprintf("https://%s/auth/external/saml/%s/%s", global.GetConfig().Trasa.ListenAddr, uc.User.OrgID, idp.IdpName)

		err := Store.UpdateSAMLIDP(&idp)
		if err != nil {
			logger.Error(err)
			utils.TrasaResponse(w, 200, "failed", "failed to update IDP", "Could not update IDP", nil)
			return
		}
	} else if idp.IdpType == "ldap" {
		// first check if ldap configuration is valid.
		// we can try binding(authenticating) ldap server with client id and secret.
		// If it fails, return invalid detail response.
		if idp.IdpName == "ad" {
			fullUserPath := fmt.Sprintf("CN=%s,%s", idp.ClientID, idp.IdpMeta)
			err := BindLdap(fullUserPath, idp.ClientSecret, idp.Endpoint)
			if err != nil {
				logger.Error(err)
				utils.TrasaResponse(w, 200, "failed", err.Error(), "Could not update IDP", nil)
				return
			}
		} else {
			fullUserPath := fmt.Sprintf("uid=%s, %s", idp.ClientID, idp.IdpMeta)
			err := BindLdap(fullUserPath, idp.ClientSecret, idp.Endpoint)
			if err != nil {
				logger.Error(err)
				utils.TrasaResponse(w, 200, "failed", err.Error(), "Could not update IDP", nil)
				return
			}
		}
		err := Store.UpdateLDAPIDP(&idp)
		if err != nil {
			logger.Error(err)
			utils.TrasaResponse(w, 200, "failed", "failed to update IDP", "Could not update IDP", nil)
			return
		}

		var key models.KeysHolder

		key.OrgID = uc.User.OrgID
		key.KeyID = idp.IdpID
		key.KeyTag = "*******"
		key.AddedBy = uc.User.ID
		key.AddedAt = time.Now().UnixNano()
		//fmt.Println("keys: ", string(dbstore.TsxvKey.Key[:]))
		ct, err := tsxvault.Store.AesEncrypt([]byte(idp.ClientSecret))
		if err != nil {
			logger.Error(err)
			utils.TrasaResponse(w, 200, "failed", "failed to store token. Vault is sealed", "Could not update IDP", nil)
			return
		}
		key.KeyVal = ct
		key.KeyName = idp.IdpName

		err = tsxvault.Store.StoreKeyOrTokens(key)
		if err != nil {
			logger.Error(err)
			utils.TrasaResponse(w, 200, "failed", "Failed to update IDP. Error in StoreKeyOrTokens.", "Could not update IDP", nil)
			return
		}

	}

	utils.TrasaResponse(w, 200, "success", "IDP updated", "IDP updated", idp)
}

type idpVendor struct {
	Name            string `json:"name"`
	IdpType         string `json:"idpType"`
	IntegrationType string `json:"string"`
	SCIMEndpoint    string `json:"scimEndpoint"`
	APIKey          string `json:"apiKey"`
}

var suppoertedIdps = []string{"okta", "jumpcloud", "gsuite"}

// PreConfiguredIdps is hardcoded pre-configurations for supported indetity providers
func PreConfiguredIdps(idp models.IdentityProvider, uc models.UserContext) models.IdentityProvider {
	var v models.IdentityProvider

	switch true {
	case idp.IdpType == "okta":
		v.IdpName = "okta"
		v.IdpType = "saml"
		v.IntegrationType = "scim"
		v.ApiKey = ""
		v.SCIMEndpoint = fmt.Sprintf("https://%s/scim/v2", global.GetConfig().Trasa.ListenAddr)
		v.ClientID = ""
		v.ClientSecret = ""
		v.AudienceURI = utils.GetRandomString(7)
		v.RedirectURL = fmt.Sprintf("https://%s/auth/external/saml/%s/%s", global.GetConfig().Trasa.ListenAddr, uc.User.OrgID, v.IdpName)
		v.OrgID = uc.User.OrgID
		v.CreatedBy = uc.User.ID
		v.LastUpdated = time.Now().Unix()
		v.IdpID = utils.GetUUID()
		v.IsEnabled = true
		return v
	case idp.IdpType == "onelogin":
		v.IdpName = "onelogin"
		v.IdpType = "saml"
		v.IntegrationType = "scim"
		v.ApiKey = ""
		v.SCIMEndpoint = fmt.Sprintf("https://%s/scim/v2", global.GetConfig().Trasa.ListenAddr)
		v.ClientID = ""
		v.ClientSecret = ""
		v.AudienceURI = utils.GetRandomString(7)
		v.RedirectURL = fmt.Sprintf("https://%s/auth/external/saml/%s/%s", global.GetConfig().Trasa.ListenAddr, uc.User.OrgID, v.IdpName)
		v.OrgID = uc.User.OrgID
		v.CreatedBy = uc.User.ID
		v.LastUpdated = time.Now().Unix()
		v.IdpID = utils.GetUUID()
		v.IsEnabled = true
		return v
	case idp.IdpType == "jumpcloud":
		v.IdpName = "jumpcloud"
		v.IdpType = "saml"
		v.IntegrationType = "scim"
		v.ApiKey = ""
		v.SCIMEndpoint = ""
		v.ClientID = ""
		v.ClientSecret = ""
		v.AudienceURI = utils.GetRandomString(7)
		v.RedirectURL = fmt.Sprintf("https://%s/auth/external/saml/%s/%s", global.GetConfig().Trasa.ListenAddr, uc.User.OrgID, v.IdpName)
		v.OrgID = uc.User.OrgID
		v.CreatedBy = uc.User.ID
		v.LastUpdated = time.Now().Unix()
		v.IdpID = utils.GetUUID()
		v.IsEnabled = true
		return v
	case idp.IdpType == "freeipa":
		v.IdpName = "freeipa"
		v.IdpType = "ldap"
		v.IntegrationType = "manual-import"
		v.ApiKey = ""
		v.SCIMEndpoint = ""
		v.ClientID = ""
		v.ClientSecret = ""
		v.AudienceURI = ""
		v.RedirectURL = ""
		v.OrgID = uc.User.OrgID
		v.CreatedBy = uc.User.ID
		v.LastUpdated = time.Now().Unix()
		v.IdpID = utils.GetUUID()
		v.IsEnabled = true
		return v
	case idp.IdpType == "ad":
		v.IdpName = "ad"
		v.IdpType = "ldap"
		v.IntegrationType = "manual-import"
		v.ApiKey = ""
		v.SCIMEndpoint = ""
		v.ClientID = ""
		v.ClientSecret = ""
		v.AudienceURI = ""
		v.RedirectURL = ""
		v.OrgID = uc.User.OrgID
		v.CreatedBy = uc.User.ID
		v.LastUpdated = time.Now().Unix()
		v.IdpID = utils.GetUUID()
		v.IsEnabled = true
		return v
	case idp.IdpType == "ldap":
		v.IdpName = idp.IdpName
		v.IdpType = "ldap"
		v.IntegrationType = "manual-import"
		v.ApiKey = ""
		v.SCIMEndpoint = ""
		v.ClientID = ""
		v.ClientSecret = ""
		v.AudienceURI = ""
		v.RedirectURL = ""
		v.OrgID = uc.User.OrgID
		v.CreatedBy = uc.User.ID
		v.LastUpdated = time.Now().Unix()
		v.IdpID = utils.GetUUID()
		v.IsEnabled = true
		return v
	default:
		return v

	}

	//return v

}

// GenerateSCIMAuthToken creates auth token for scim connector.
// token is basically passwaord with format "password:orgID" which is hashed and stored in database.
// password is returned to user only once. first 4 characters of password
// is stored as key tag which will be returned in subsequent request.
func GenerateSCIMAuthToken(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	idpID := chi.URLParam(r, "idpID")

	// create password
	pass := utils.GetRandomString(14)

	orgpass := fmt.Sprintf("%s:%s", uc.User.OrgID, pass)

	hashedpass, err := bcrypt.GenerateFromPassword([]byte(orgpass), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed generate hashed password", "GenerateSCIMAuthToken", nil)
		return
	}

	var req models.KeysHolder

	req.OrgID = uc.User.OrgID
	req.KeyID = idpID
	req.KeyTag = fmt.Sprintf("%s-xxxx-xxxx...", pass[0:4])
	req.AddedBy = uc.User.ID
	req.AddedAt = time.Now().Unix()

	// we do not need to encrypt the key here since we are storing bcrypt hash
	req.KeyVal = hashedpass
	req.KeyName = consts.SCIMKEY

	err = tsxvault.Store.StoreKeyOrTokens(req)
	if err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to store token.", "GenerateSCIMAuthToken", nil)
		return
	}

	encodedPass := utils.EncodeBase64([]byte(orgpass))
	utils.TrasaResponse(w, 200, "success", "IDPs", "GenerateSCIMAuthToken", encodedPass)
}

//ActivateOrDisableIdp enable or disable IDP
func ActivateOrDisableIdp(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	type idpState struct {
		IdpID  string `json:"idpID"`
		Active bool   `json:"active"`
	}

	var req idpState
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "ActivateOrDisableIdp", nil)
		return
	}

	err := Store.activateOrDisableIdp(uc.User.OrgID, req.IdpID, time.Now().Unix(), req.Active)
	if err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to disable IDP", "ActivateOrDisableIdp-ActivateOrDisableIdp", nil)
		return
	}

	var resp string
	if req.Active == true {
		resp = "IDP activated"
	} else {
		resp = "IDP disabled"
	}
	utils.TrasaResponse(w, 200, "success", resp, "StoredIdps", nil)
}

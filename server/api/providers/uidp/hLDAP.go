package uidp

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"

	"github.com/seknox/trasa/server/api/providers/vault"
	"github.com/seknox/trasa/server/api/providers/vault/tsxvault"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/models"

	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/ldap.v2"
)

// ImportLdapUsers search and imports ldap users
func ImportLdapUsers(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	var req models.IdentityProvider
	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "failed to import users", nil)
		return
	}

	// get idp detail from idpID
	idpDetail, err := Store.GetByID(uc.User.OrgID, req.IdpID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Invalid IdP ID", "failed to import users", nil)
		return
	}

	// get client secret from keyholders.
	ct, err := vault.Store.GetKeyOrTokenWithKeyvalAndID(uc.User.OrgID, consts.KEY_LDAP, req.IdpID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "IFailed to fetch secret value for ldap system user", "failed to import users", nil)
		return
	}

	password, err := tsxvault.Store.AesDecrypt(ct.KeyVal)

	ldapBindUsername := fmt.Sprintf("uid=%s,%s", idpDetail.ClientID, idpDetail.IdpMeta)
	if idpDetail.IdpName == "ad" {
		ldapBindUsername = fmt.Sprintf("CN=%s,%s", idpDetail.ClientID, idpDetail.IdpMeta)
	}

	// search users in ldap
	ldapUsers, err := bindSearchImportLdapUsers(uc, ldapBindUsername, string(password), idpDetail.Endpoint, req.AudienceURI, idpDetail.IdpName)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", err.Error(), "failed to import users", nil)
		return
	}

	// we can now start importing users. username and emails are unqiue in TRASA. failed insert users should be colleceted and informed back to the administrators.
	var failedusers = make([]string, 0)
	for k, v := range ldapUsers {
		err := users.Store.Create(&v)
		if err != nil {
			failedusers = append(failedusers, fmt.Sprintf("%d) user with email: %s and username: %s.\n", k+1, v.Email, v.UserName))
		}
	}

	var resp importeUsers
	resp.TotalUsers = len(ldapUsers)
	resp.SuccessCount = resp.TotalUsers - len(failedusers)
	resp.FailedCount = len(failedusers)
	resp.FailedUsers = failedusers

	// return response to user along with failed users (if any)
	utils.TrasaResponse(w, 200, "success", "users created", "Users imported from IDP", resp)
}

type importeUsers struct {
	TotalUsers   int      `json:"totalUsers"`
	SuccessCount int      `json:"successCount"`
	FailedCount  int      `json:"failedCount"`
	FailedUsers  []string `json:"failedUsers"`
}

func BindLdap(uname, pass, domain string) error {

	//logrus.Trace("LDAP Login details: ", uname, pass, domain)
	tc := &tls.Config{
		ServerName:         domain,
		InsecureSkipVerify: true,
	}

	l, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", domain, 636), tc)
	if err != nil {
		return err
	}
	defer l.Close()

	err = l.Bind(uname, pass)
	if err != nil {
		return err
	}

	return nil
}

func searchLdap(l *ldap.Conn, query string, queryFields []string) (*ldap.SearchResult, error) {
	searchRequest := ldap.NewSearchRequest(
		query,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",
		queryFields,
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return sr, err
	}

	return sr, nil
}

func bindSearchImportLdapUsers(uc models.UserContext, uname, pass, domain, searchquery, idpName string) ([]models.UserWithPass, error) {

	var lus = make([]models.UserWithPass, 0)

	tc := &tls.Config{
		ServerName:         domain,
		InsecureSkipVerify: true,
	}

	l, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", domain, 636), tc)
	if err != nil {
		return nil, err
	}
	defer l.Close()

	err = l.Bind(uname, pass)
	if err != nil {
		return nil, err
	}

	// Search for the given base
	sr, err := searchLdap(l, searchquery, []string{"member"})
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range sr.Entries {

		mem := v.GetAttributeValues("member")
		fmt.Println("member: ", mem)

		for _, k := range mem {

			sr, err := searchLdap(l, k, []string{"objectGUID", "sAMAccountName", "userPrincipalName", "name"})
			if err != nil {
				logrus.Error(err)
			}
			for _, u := range sr.Entries {
				guid := u.GetRawAttributeValue("objectGUID")
				username := u.GetAttributeValue("sAMAccountName")
				email := u.GetAttributeValue("userPrincipalName")
				name := u.GetAttributeValue("name")
				fmt.Println("user: ", fmt.Sprintf("%x", guid), username, email, name)

				for _, v := range sr.Entries {
					email := v.GetAttributeValue("userPrincipalName")
					if email != "" {
						guid := u.GetRawAttributeValue("objectGUID")
						username := u.GetAttributeValue("sAMAccountName")
						name := u.GetAttributeValue("name")
						fullName := strings.Split(name, " ")
						firstName := ""
						lastName := ""
						if len(fullName) > 1 {
							firstName = fullName[0]
							lastName = fullName[1]
						} else {
							firstName = name
						}

						var lu models.UserWithPass
						lu.OrgID = uc.User.OrgID
						lu.UserName = utils.NormalizeString(username)
						lu.FirstName = utils.NormalizeString(firstName)
						lu.LastName = utils.NormalizeString(lastName)
						lu.Email = utils.NormalizeString(email)
						lu.UserRole = "selfUser"
						lu.ID = utils.GetUUID()
						lu.ExternalID = utils.NormalizeString(fmt.Sprintf("%x", guid))
						lu.IdpName = idpName
						lu.Password = ""
						lu.Status = true
						lu.CreatedAt = time.Now().Unix()
						lu.UpdatedAt = lu.CreatedAt

						lus = append(lus, lu)
					}

				}

			}
		}

	}

	//timezone, _ := time.LoadLocation(uc.Org.Timezone)

	return lus, nil
}

// GetUsersAll returns json array of user list.
func GetAllUsersForIdp(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	idpName := chi.URLParam(r, "idpname")

	val, err := users.Store.GetAllByIdp(userContext.Org.ID, idpName)

	if err != nil {
		logrus.Debug(err)
		utils.TrasaResponse(w, 200, "failed", "users not fetched", "GetUsersAll", nil)
		return
	}
	utils.TrasaResponse(w, 200, "success", "users fetched", "GetAllUsersForIdp", val)
}

type transferReq struct {
	IdpName  string     `json:"idpName"`
	Userlist [][]string `json:"userList"`
}

// GetUsersAll returns json array of user list.
func TransferUserToGivenIdp(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	var req transferReq
	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "TransferUserToGivenIdp")
		return
	}

	for _, v := range req.Userlist {
		err := users.Store.TransferUser(uc.Org.ID, v[2], req.IdpName)
		if err != nil {
			logrus.Debug(err)
		}
	}

	utils.TrasaResponse(w, 200, "success", "user(s) updated", "GetAllUsersForIdp")
}

package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/seknox/trasa/server/api/crypt/vault"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type creds struct {
	Username   string `json:"username"`
	Credential string `json:"credential"`
	ServiceID  string `json:"serviceID"`
	Type       string `json:"type"`
}

// StoreServiceCredentials takes username password from client (trasa-dashboard for now) and stores it in vault.
// It will also store the event in trasadb. This will come handy for in-app audit logs.
// storing it separate will also decouples our core dependency in vault
func StoreServiceCredentials(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	var req creds

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "failed to save password")
		return
	}

	if req.Type == "key" {
		_, err := ssh.ParsePrivateKey([]byte(req.Credential))
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "Invalid SSH key", "failed to save password")
			return
		}
	}

	var s models.ServiceSecretVault
	s.KeyID = utils.GetRandomString(7)
	s.ServiceID = req.ServiceID
	s.SecretType = req.Type
	s.OrgID = userContext.Org.ID

	s.Secret = []byte(req.Credential)
	s.SecretID = req.Username
	s.AddedAt = time.Now().Unix()
	s.LastUpdated = time.Now().Unix()

	err := vault.Store.StoreSecret(s)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not save password", "Could not save password")
		return
	}

	err = Store.AddManagedAccounts(req.ServiceID, userContext.Org.ID, req.Username)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not save password", "Could not save password")
		return
	}

	utils.TrasaResponse(w, 200, "success", "cred stored", fmt.Sprintf(`Password saved for "%s" user `, req.Username), req.Username)

	// we also store user names that has been enrolled in secret store in cockroachdb to reference
	// managed accounts in that app.

}

func ViewCreds(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Got GetPassword")
	userContext := r.Context().Value("user").(models.UserContext)
	var req creds

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "Could not view password", nil, nil)
		return
	}

	passCred, err1 := vault.Store.GetSecret(userContext.User.OrgID, req.ServiceID, req.Type, req.Username)
	if err1 != nil {
		logrus.Error(err1)
		utils.TrasaResponse(w, 200, "failed", "Could not view password", "Could not view password", nil, nil)
		return
	}

	req.Credential = passCred

	service, err := Store.GetFromID(req.ServiceID)
	if err != nil {
		logrus.Error(err, "invalid service ID in view creds")
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Invalid service", "Could not view password")
		return
	}

	utils.TrasaResponse(w, 200, "success", "creds fetched", fmt.Sprintf(`Viewed password for "%s" user in "%s" app`, req.Username, service.Name), req)
}

// DeleteCreds deletes stored creds from both database and vault.
func DeleteCreds(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	//	fmt.Println("Got deletepass")
	var req creds

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "false", "invalid request", "Could not delete password")
		return
	}

	err := vault.Store.TsxvDeleteSecret(userContext.User.OrgID, req.ServiceID, "password", req.Username)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "DeleteCreds", "Could not delete password")
		return
	}

	err = vault.Store.TsxvDeleteSecret(userContext.User.OrgID, req.ServiceID, "key", req.Username)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "DeleteCreds", "Could not delete password")
		return
	}

	// we also need to delete username from managed accounts from service table.

	err = Store.RemoveManagedAccounts(req.ServiceID, userContext.Org.ID, req.Username)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "DeleteCreds", "Could not delete password")
		return
	}

	utils.TrasaResponse(w, 200, "success", "creds deleted", fmt.Sprintf(`Password deleted for user "%s"`, req.Username))

}

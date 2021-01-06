package vault

import (
	"encoding/json"
	"fmt"
	"time"

	hcvault "github.com/hashicorp/vault/api"
	"github.com/seknox/trasa/server/models"
	"github.com/sirupsen/logrus"
)

// initclient initialize vault client
func(s cryptStore) initclient() (*hcvault.Client, error) {

	//  configure address
	conf := &hcvault.Config{
		Address: s.TsxvKey.CredProv.ProviderAddr,
	}

	// init new client
	vc, err := hcvault.NewClient(conf)
	if err != nil {
		return vc, err
	}

	// set access token
	vc.SetToken(s.TsxvKey.CredProv.ProviderAccessToken)

	return vc, nil
}

// StoreCred stores credential to vault
func (s cryptStore) HCVStoreCred(cred models.ServiceSecretVault) error {


	// init client
	client, err := Store.initclient()
	if err != nil {
		return err
	}

	conn := client.Logical()

	AddedAt := time.Now().UTC().String()
	LastUpdated := time.Now().UTC().String()

	var data = make(map[string]interface{})

	// vault secret kv path for trasa
	vaultSecretPath := fmt.Sprintf("/trasa/data/%s:%s", cred.OrgID, cred.ServiceID)

	var resp vaultResponse

	// First we check is secret for host is already stored
	secret, err := conn.Read(vaultSecretPath)
	if err != nil {
		return err
	}

	// if secret is nil it means vault does not have any secret yet stored for this host and type.
	if secret == nil {
		logrus.Debug("not set yet")

		var user []secrets

		credOne := secrets{
			SecretID:    cred.SecretID,
			Secret:    string(cred.Secret),
			AddedAt:     AddedAt,
			LastUpdated: LastUpdated,
		}

		user = append(user, credOne)

		finalData := hostData{
			Secrets: user,
		}

		data["data"] = finalData

		// write to vault
		_, err = conn.Write(vaultSecretPath, data)
		if err != nil {
			return err
		}
		return nil
	}

	// if we are here, it means vault has already piece of secret stored for this host.
	// we now have to check if request user already exist in vault. if not, we add new user with
	// current added and updated time.
	val, err := json.Marshal(secret.Data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(val, &resp)
	if err != nil {
		return err
	}


	for i := 0; i < len(resp.Data.Secrets); i++ {
		v := &resp.Data.Secrets[i]
		if v.SecretID == cred.SecretID {
			logrus.Debug("found same user: ", v.SecretID)
			v.Secret = string(cred.Secret)
			v.LastUpdated = time.Now().UTC().String()

			finalData := hostData{
				Secrets: resp.Data.Secrets,
			}

			data["data"] = finalData
			// write to vault
			_, err = conn.Write(vaultSecretPath, data)
			if err != nil {
				return err
			}
			return nil
		}
	}

	logrus.Debug("found new account: ", cred.SecretID)

	credOne := secrets{
		SecretID:    cred.SecretID,
		Secret:    string(cred.Secret),
		AddedAt:     AddedAt,
		LastUpdated: LastUpdated,
	}

	resp.Data.Secrets = append(resp.Data.Secrets, credOne)
	//	user = append(user, credOne)

	finalData := hostData{
		Secrets: resp.Data.Secrets,
	}

	data["data"] = finalData

	// write to vault
	_, err = conn.Write(vaultSecretPath, data)
	if err != nil {
		return err
	}
	return nil
}



// ReadCred reads credential stored in vault
func (s cryptStore) HCVReadCred(orgID, serviceID, secretID string) (string, error) {

	// init client
	client, err := Store.initclient()
	if err != nil {
		return "", err
	}


	conn := client.Logical()

	var resp vaultResponse

	vaultSecretPath := fmt.Sprintf("/trasa/data/%s:%s", orgID, serviceID)

	// read from vault
	secret, err := conn.Read(vaultSecretPath)
	if err != nil {
		return "", err
	}

	if secret == nil {
		return "", fmt.Errorf("%d", 1)
	}



	val, err := json.Marshal(secret.Data)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(val, &resp)
	if err != nil {
		return "", err
	}

	var cred string
	for _, v := range resp.Data.Secrets {
		if v.SecretID == secretID {
			cred = v.Secret
		}
	}

	return cred, nil

}

// RemoveCred removed credential from vault
func (s cryptStore) HCVRemoveCred(orgID, serviceID, secretID string) error {
	client, err := Store.initclient()
	if err != nil {
		return err
	}

	conn := client.Logical()

	var data = make(map[string]interface{})

	vaultSecretPath := fmt.Sprintf("/trasa/data/%s:%s", orgID, serviceID)

	var resp vaultResponse

	// First we check is secret for host is already stored
	secret, err := conn.Read(vaultSecretPath)
	if err != nil {
		return err
	}

	// if secret is nil it means user tried to delete creds from structure that does not exist yet.
	// we will return error
	if secret == nil {
		err := fmt.Errorf("%d", 1)
		return err
	}

	// if we are here, it means vaut has secret structure stored for provided host.
	// we marshal the secret response to slice structure and remove element(user) from this structure.
	val, err := json.Marshal(secret.Data)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(val, &resp)
	if err != nil {
		return err
	}


	for i := 0; i < len(resp.Data.Secrets); i++ {
		v := &resp.Data.Secrets[i]
		if v.SecretID == secretID {
			resp.Data.Secrets = sliceUserCred(i, resp.Data.Secrets)

			finalData := hostData{
				Secrets: resp.Data.Secrets,
			}

			data["data"] = finalData

			_, err = conn.Write(vaultSecretPath, data)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}


// RemoveCred removed credential from vault
func (s cryptStore) HCVDeleteForService(orgID, serviceID string) error {
	client, err := Store.initclient()
	if err != nil {
		return err
	}

	conn := client.Logical()


	vaultSecretPath := fmt.Sprintf("/trasa/data/%s:%s", orgID, serviceID)



	// Delete from vault. Note this will not destroy!!
	_, err = conn.Delete(vaultSecretPath)
	if err != nil {
		return err
	}


	return nil
}

// sliceUserCred deletes user secret element from vauls secret structure.
// code taken from https://yourbasic.org/golang/delete-element-slice/
func sliceUserCred(index int, array []secrets) []secrets {
	var u secrets
	array[index] = array[len(array)-1]
	array[len(array)-1] = u
	array = array[:len(array)-1]

	return array
}


type vaultResponse struct {
	Data struct {
		Secrets []secrets
	}
	MetaData struct {
		CreatedTime  string
		DeletionTime string
		Destroyed    bool
		Version      int
	}
}

type vaultRequestData struct {
	Secrets interface{}
}

type secrets struct {
	SecretID    string
	Secret    string
	AddedAt     string
	LastUpdated string
}

type hostData struct {
	Secrets interface{}
}

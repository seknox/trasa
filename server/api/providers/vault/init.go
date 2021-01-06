package vault

import (
	hcvault "github.com/hashicorp/vault/api"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	Store = cryptStore{State: state}
}

//Store is the package state variable which contains database connections
var Store adapter

type adapter interface {
	StoreKey(k models.KeysHolder) error
	GetKeyOrTokenWithTag(orgID string, keyName string) (*models.KeysHolder, error)
	GetKeyOrTokenWithKeyval(orgID, keyName string) (*models.KeysHolder, error)
	GetKeyOrTokenWithKeyvalAndID(orgID, keyName, keyID string) (*models.KeysHolder, error)

	// generic credential crud ops
	StoreCred(key models.ServiceSecretVault) error 
	ReadCred(orgID, serviceID, secretType, secretID string) (string, error)
	RemoveCred(orgID, serviceID, secretType, secretID string) error
	DeleteCreds(orgID, serviceID string) error

	// hashicorp vault functions for credential crud ops
	initclient() (*hcvault.Client, error) 
	HCVStoreCred(cred models.ServiceSecretVault) error 
	HCVReadCred(orgID, serviceID, secretID string) (string, error)
	HCVRemoveCred(orgID, serviceID, secretID string) error
	HCVDeleteForService(orgID, serviceID string) error
}

type cryptStore struct {
	*global.State
}

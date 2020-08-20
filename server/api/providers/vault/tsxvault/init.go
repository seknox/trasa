package tsxvault

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	Store = vaultStore{State: state}
}

//InitStoreMock will init mock state of this package
func InitStoreMock() *vaultMock {
	m := new(vaultMock)
	Store = m
	return m
}

//Store is the package state variable which contains database connections
var Store adapter

type vaultStore struct {
	*global.State
}

//TODO @sshahcodes rename the functions
//
// Expose only generic functions like StoreSecret and hide functions like TsxvStoreSecret
// Implementation of TRASA vault/ separate vault should be hidden to the package caller
type adapter interface {
	StoreSecret(key models.ServiceSecretVault) error
	GetSecret(orgID, serviceID, secretType, secretid string) (string, error)
	TsxvStoreSecret(secret models.ServiceSecretVault) error
	TsxvGetSecretDetail(orgID, serviceID, appType, secretid string) (*models.ServiceSecretVault, error)
	TsxvGetSecret(orgID, serviceID, appType, secretid string) ([]byte, error)
	TsxvDeleteSecret(orgID, serviceID, secretType, secretid string) error
	TsxvDeleteAllSecret(orgID string) error
	TsxvStoreEncKeyHash(secret models.EncryptionKeyLog) error
	TsxvGetEncKeyHash(orgID, keyHash string) (models.EncryptionKeyLog, error)
	TsxvUpdateEncKeyHashLog(orgID, keyHash string, time int64, status bool) error
	TsxvdeactivateAllKeys(orgID string, time int64) error
	TsxvTestEncrypter(key []byte) error
	TsxvTestDecrypter(key []byte) error
	TsxVaultTester() error
	AesEncrypt(message []byte) ([]byte, error)
	AesDecrypt(message []byte) ([]byte, error)
	GetKeyOrTokenWithTag(orgID string, keyName string) (*models.KeysHolder, error)
	StoreKeyOrTokens(k models.KeysHolder) error
	GenAndStoreKey(orgID string) (*[32]byte, error)
	GetTsxVaultKey() (*[32]byte, bool)
	SetTsxVaultKey(key *[32]byte, status bool)
}

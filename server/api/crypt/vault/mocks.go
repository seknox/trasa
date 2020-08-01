package vault

import (
	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

type VaultMock struct {
	mock.Mock
}

func (v VaultMock) StoreSecret(key models.ServiceSecretVault) error {
	panic("implement me")
}

func (v VaultMock) GetSecret(orgID, serviceID, secretType, secretid string) (string, error) {
	args := v.Called(orgID, serviceID, secretType, secretid)
	return args.String(0), args.Error(1)
}

func (v VaultMock) TsxvStoreSecret(secret models.ServiceSecretVault) error {
	panic("implement me")
}

func (v VaultMock) TsxvGetSecretDetail(orgID, serviceID, appType, secretid string) (*models.ServiceSecretVault, error) {
	panic("implement me")
}

func (v VaultMock) TsxvGetSecret(orgID, serviceID, appType, secretid string) ([]byte, error) {
	panic("implement me")
}

func (v VaultMock) TsxvDeleteSecret(orgID, serviceID, secretType, secretid string) error {
	panic("implement me")
}

func (v VaultMock) TsxvDeleteAllSecret(orgID string) error {
	panic("implement me")
}

func (v VaultMock) TsxvStoreEncKeyHash(secret models.EncryptionKeyLog) error {
	panic("implement me")
}

func (v VaultMock) TsxvGetEncKeyHash(orgID, keyHash string) (models.EncryptionKeyLog, error) {
	panic("implement me")
}

func (v VaultMock) TsxvUpdateEncKeyHashLog(orgID, keyHash string, time int64, status bool) error {
	panic("implement me")
}

func (v VaultMock) TsxvdeactivateAllKeys(orgID string, time int64) error {
	panic("implement me")
}

func (v VaultMock) TsxvTestEncrypter(key []byte) error {
	panic("implement me")
}

func (v VaultMock) TsxvTestDecrypter(key []byte) error {
	panic("implement me")
}

func (v VaultMock) TsxVaultTester() error {
	panic("implement me")
}

func (v VaultMock) AesEncrypt(message []byte) ([]byte, error) {
	panic("implement me")
}

func (v VaultMock) AesDecrypt(message []byte) ([]byte, error) {
	panic("implement me")
}

func (v VaultMock) GetKeyOrTokenWithTag(orgID string, keyName string) (*models.KeysHolder, error) {
	panic("implement me")
}

func (v VaultMock) StoreKeyOrTokens(k models.KeysHolder) error {
	panic("implement me")
}

func (v VaultMock) GenAndStoreKey(orgID string) (*[32]byte, error) {
	panic("implement me")
}

func (v VaultMock) GetTsxVaultKey() (*[32]byte, bool) {
	panic("implement me")
}

func (v VaultMock) SetTsxVaultKey(key *[32]byte, status bool) {
	panic("implement me")
}

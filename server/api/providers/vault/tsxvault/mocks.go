package tsxvault

import (
	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

type vaultMock struct {
	mock.Mock
}

func (v *vaultMock) StoreSecret(key models.ServiceSecretVault) error {
	panic("implement me")
}

func (v *vaultMock) GetSecret(orgID, serviceID, secretType, secretid string) (string, error) {
	args := v.Called(orgID, serviceID, secretType, secretid)
	return args.String(0), args.Error(1)
}

func (v *vaultMock) TsxvStoreSecret(secret models.ServiceSecretVault) error {
	panic("implement me")
}

func (v *vaultMock) TsxvGetSecretDetail(orgID, serviceID, appType, secretid string) (*models.ServiceSecretVault, error) {
	panic("implement me")
}

func (v *vaultMock) TsxvGetSecret(orgID, serviceID, appType, secretid string) ([]byte, error) {
	panic("implement me")
}

func (v *vaultMock) TsxvDeleteSecret(orgID, serviceID, secretType, secretid string) error {
	panic("implement me")
}

func (v *vaultMock) TsxvDeleteAllSecret(orgID string) error {
	panic("implement me")
}

func (v *vaultMock) TsxvStoreEncKeyHash(secret models.EncryptionKeyLog) error {
	panic("implement me")
}

func (v *vaultMock) TsxvGetEncKeyHash(orgID, keyHash string) (models.EncryptionKeyLog, error) {
	panic("implement me")
}

func (v *vaultMock) TsxvUpdateEncKeyHashLog(orgID, keyHash string, time int64, status bool) error {
	panic("implement me")
}

func (v *vaultMock) TsxvdeactivateAllKeys(orgID string, time int64) error {
	panic("implement me")
}

func (v *vaultMock) TsxvTestEncrypter(key []byte) error {
	panic("implement me")
}

func (v *vaultMock) TsxvTestDecrypter(key []byte) error {
	panic("implement me")
}

func (v *vaultMock) TsxVaultTester() error {
	panic("implement me")
}

func (v *vaultMock) AesEncrypt(message []byte) ([]byte, error) {
	panic("implement me")
}

func (v *vaultMock) AesDecrypt(message []byte) ([]byte, error) {
	panic("implement me")
}

func (v *vaultMock) GetKeyOrTokenWithTag(orgID string, keyName string) (*models.KeysHolder, error) {
	panic("implement me")
}

func (v *vaultMock) StoreKeyOrTokens(k models.KeysHolder) error {
	panic("implement me")
}

func (v *vaultMock) GenAndStoreKey(orgID string) (*[32]byte, error) {
	panic("implement me")
}

func (v *vaultMock) GetTsxVaultKey() (*[32]byte, bool) {
	panic("implement me")
}

func (v *vaultMock) SetTsxVaultKey(key *[32]byte, status bool,  credprov models.CredProvProps) {
	panic("implement me")
}

func (v *vaultMock) UpdateTsxVaultKeyCredProvConfig(credprov models.CredProvProps)  {
	panic("implement me")
}


// SetTsxCPxyKey assigns retreived cloud prxy api key in global state
func (v *vaultMock) SetTsxCPxyKey(key string) {

	panic("implement me")

}

// GetTsxCPxyKey retreives cloud prxy api key from global state
func (v *vaultMock) GetTsxCPxyKey() string {

	panic("implement me")

}

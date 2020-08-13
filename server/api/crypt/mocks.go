package crypt

import (
	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

type cryptMock struct {
	mock.Mock
}

func (c *cryptMock) StoreKey(k models.KeysHolder) error {
	panic("implement me")
}

func (c *cryptMock) GetKeyOrTokenWithTag(orgID string, keyName string) (*models.KeysHolder, error) {
	panic("implement me")
}

func (c *cryptMock) GetKeyOrTokenWithKeyval(orgID, keyName string) (*models.KeysHolder, error) {
	panic("implement me")
}

func (c *cryptMock) GetKeyOrTokenWithKeyvalAndID(orgID, keyName, keyID string) (*models.KeysHolder, error) {
	panic("implement me")
}

func (c *cryptMock) StoreCert(ch models.CertHolder) error {
	panic("implement me")
}

func (c *cryptMock) DelCA(userID, orgID string) error {
	panic("implement me")
}

func (c *cryptMock) GetCertDetail(orgID, entityID, certType string) (*models.CertHolder, error) {
	panic("implement me")
}

func (c *cryptMock) GetAllCAs(orgID string) ([]models.CertHolder, error) {
	panic("implement me")
}

func (c *cryptMock) GetCAkey(orgID string) ([]byte, error) {
	panic("implement me")
}

func (c *cryptMock) GetCertHolder(certType, entityID, orgID string) (models.CertHolder, error) {
	args := c.Called(certType, entityID, orgID)
	return args.Get(0).(models.CertHolder), args.Error(1)
}

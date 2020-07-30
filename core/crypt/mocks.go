package crypt

import (
	"github.com/seknox/trasa/models"
	"github.com/stretchr/testify/mock"
)

type CryptMock struct {
	mock.Mock
}

func (c CryptMock) StoreKey(k models.KeysHolder) error {
	panic("implement me")
}

func (c CryptMock) GetKeyOrTokenWithTag(orgID string, keyName string) (*models.KeysHolder, error) {
	panic("implement me")
}

func (c CryptMock) GetKeyOrTokenWithKeyval(orgID, keyName string) (*models.KeysHolder, error) {
	panic("implement me")
}

func (c CryptMock) GetKeyOrTokenWithKeyvalAndID(orgID, keyName, keyID string) (*models.KeysHolder, error) {
	panic("implement me")
}

func (c CryptMock) StoreCert(ch models.CertHolder) error {
	panic("implement me")
}

func (c CryptMock) DelCA(userID, orgID string) error {
	panic("implement me")
}

func (c CryptMock) GetCertDetail(orgID, entityID, certType string) (*models.CertHolder, error) {
	panic("implement me")
}

func (c CryptMock) GetAllCAs(orgID string) ([]models.CertHolder, error) {
	panic("implement me")
}

func (c CryptMock) GetCAkey(orgID string) ([]byte, error) {
	panic("implement me")
}

func (c CryptMock) GetCertHolder(certType, entityID, orgID string) (models.CertHolder, error) {
	args := c.Called(certType, entityID, orgID)
	return args.Get(0).(models.CertHolder), args.Error(1)
}

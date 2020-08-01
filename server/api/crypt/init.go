package crypt

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

func InitStore(state *global.State) {
	Store = CryptStore{State: state}
}

func InitStoreMock() *CryptMock {
	m := new(CryptMock)
	Store = m
	return m
}

var Store CryptAdapter

type CryptAdapter interface {
	StoreKey(k models.KeysHolder) error
	GetKeyOrTokenWithTag(orgID string, keyName string) (*models.KeysHolder, error)
	GetKeyOrTokenWithKeyval(orgID, keyName string) (*models.KeysHolder, error)
	GetKeyOrTokenWithKeyvalAndID(orgID, keyName, keyID string) (*models.KeysHolder, error)

	StoreCert(ch models.CertHolder) error
	DelCA(userID, orgID string) error
	GetCertDetail(orgID, entityID, certType string) (*models.CertHolder, error)
	GetAllCAs(orgID string) ([]models.CertHolder, error)
	GetCAkey(orgID string) ([]byte, error)
	GetCertHolder(certType, entityID, orgID string) (models.CertHolder, error)
}

type CryptStore struct {
	*global.State
}

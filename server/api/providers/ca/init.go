package ca

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	Store = cryptStore{State: state}
}

//InitStoreMock will init mock state of this package
func InitStoreMock() *cryptMock {
	m := new(cryptMock)
	Store = m
	return m
}

//Store is the package state variable which contains database connections
var Store adapter

type adapter interface {
	StoreCert(ch models.CertHolder) error
	DelCA(userID, orgID string) error
	GetCertDetail(orgID, entityID, certType string) (*models.CertHolder, error)
	GetAllCAs(orgID string) ([]models.CertHolder, error)
	GetCAkey(orgID string) ([]byte, error)
	GetCertHolder(certType, entityID, orgID string) (models.CertHolder, error)
}

type cryptStore struct {
	*global.State
}

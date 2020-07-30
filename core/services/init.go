package services

import (
	"github.com/seknox/trasa/global"
	"github.com/seknox/trasa/models"
)

func InitStore(con *global.State) {
	// initialize local state
	Store = ServiceStore{
		State: con,
	}

}

//func InitServiceStoreMock() *servicemocks.ServiceStoreMock {
//	// initialize local state
//	mockApp := new(servicemocks.ServiceStoreMock)
//	Store = mockApp
//	return mockApp
//}

var Store Adapter

type ServiceStore struct {
	*global.State
}

type Adapter interface {

	//CRUD
	GetFromID(serviceID string) (*models.Service, error)
	//GetFromName(appName, orgID string) (*models.Service, error)
	GetFromHostname(hostname, serviceType, remoteAppName, orgID string) (*models.Service, error)
	GetAllByType(serviceType, orgID string) (services []models.Service, err error)

	Create(service *models.Service) error
	Update(service *models.Service) error
	Delete(serviceID, orgID string) (string, error)

	//Appusers

	//Settinngs
	// GetHttpProxy(serviceID, orgID string) (models.HttpProxyOps, error)
	updateHttpProxy(serviceID, orgID string, time int64, proxyConfig models.ReverseProxy) error
	AddManagedAccounts(serviceID, orgID string, username string) error
	RemoveManagedAccounts(serviceID, orgIDstring, username string) error

	//stats
	GetTotalServiceUsers(serviceID string, orgID string) (int64, error)

	GetServiceNameFromID(serviceID string, orgID string) (string, error)

	UpdateSSLCerts(caCert, caKey, clientCert, clientKey, serviceID, orgID string) error
	UpdateHostCert(hostCert, serviceID, orgID string) error
}

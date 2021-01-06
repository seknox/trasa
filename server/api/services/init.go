package services

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

//InitStore initialises package state
func InitStore(con *global.State) {
	// initialize local state
	Store = serviceStore{
		State: con,
	}

}

//Store is the package state variable which contains database connections
var Store adapter

type serviceStore struct {
	*global.State
}

type adapter interface {

	//CRUD
	GetFromID(serviceID string) (*models.Service, error)
	GetFromServiceName(serviceName, orgID string) (service *models.Service, err error)
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

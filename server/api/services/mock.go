package services

import (
	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

func InitServiceStoreMock() *serviceMock {
	// initialize local state
	mockApp := new(serviceMock)
	Store = mockApp
	return mockApp
}

type serviceMock struct {
	mock.Mock
}

func (s serviceMock) GetFromID(serviceID string) (*models.Service, error) {
	args := s.Called(serviceID)
	return args.Get(0).(*models.Service), args.Error(1)
	panic("implement me")
}

func (s serviceMock) GetFromHostname(hostname, serviceType, remoteAppName, orgID string) (*models.Service, error) {
	panic("implement me")
}

func (s serviceMock) GetAllByType(serviceType, orgID string) (services []models.Service, err error) {
	panic("implement me")
}

func (s serviceMock) Create(service *models.Service) error {
	args := s.Called(service)
	return args.Error(0)
}

func (s serviceMock) Update(service *models.Service) error {
	panic("implement me")
}

func (s serviceMock) Delete(serviceID, orgID string) (string, error) {
	panic("implement me")
}

func (s serviceMock) updateHttpProxy(serviceID, orgID string, time int64, proxyConfig models.ReverseProxy) error {
	panic("implement me")
}

func (s serviceMock) AddManagedAccounts(serviceID, orgID string, username string) error {
	panic("implement me")
}

func (s serviceMock) RemoveManagedAccounts(serviceID, orgIDstring, username string) error {
	panic("implement me")
}

func (s serviceMock) GetTotalServiceUsers(serviceID string, orgID string) (int64, error) {
	panic("implement me")
}

func (s serviceMock) GetServiceNameFromID(serviceID string, orgID string) (string, error) {
	panic("implement me")
}

func (s serviceMock) UpdateSSLCerts(caCert, caKey, clientCert, clientKey, serviceID, orgID string) error {
	panic("implement me")
}

func (s serviceMock) UpdateHostCert(hostCert, serviceID, orgID string) error {
	panic("implement me")
}

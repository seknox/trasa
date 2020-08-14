package orgs

import (
	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

type orgMock struct {
	mock.Mock
}

func (o *orgMock) Create(org *models.Org) error {
	panic("implement me")
}

func (o *orgMock) Get(orgID string) (models.Org, error) {
	args := o.Called(orgID)
	return args.Get(0).(models.Org), args.Error(1)
}

func (o *orgMock) CheckOrgExists() (string, error) {
	args := o.Called()
	return args.String(0), args.Error(1)
}

func (o *orgMock) CreateOrg(org *models.Org) error {

	args := o.Called(org)
	return args.Error(0)
}
func (o *orgMock) GetIDP(orgID, idpName string) (models.IdentityProvider, error) {
	args := o.Called(orgID, idpName)
	return args.Get(0).(models.IdentityProvider), args.Error(0)
}

func (o *orgMock) GetLicense(orgID string) (license *models.License, err error) {
	panic("implement me")
}

func (o *orgMock) RemoveAllManagedAccounts(orgID string) error {
	panic("implement me")
}

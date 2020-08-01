package orgs

import (
	"github.com/seknox/trasa/models"
	"github.com/stretchr/testify/mock"
)

type OrgMock struct {
	mock.Mock
}

func (o OrgMock) Create(org *models.Org) error {
	panic("implement me")
}

func (o OrgMock) Get(orgID string) (models.Org, error) {
	args := o.Called(orgID)
	return args.Get(0).(models.Org), args.Error(1)
}

func (o OrgMock) CheckOrgExists() (string, error) {
	args := o.Called()
	return args.String(0), args.Error(1)
}

func (o OrgMock) CreateOrg(org *models.Org) error {

	args := o.Called(org)
	return args.Error(0)
}
func (o OrgMock) GetIDP(orgID, idpName string) (models.IdentityProvider, error) {
	args := o.Called(orgID, idpName)
	return args.Get(0).(models.IdentityProvider), args.Error(0)
}

func (o OrgMock) GetLicense(orgID string) (license *models.License, err error) {
	panic("implement me")
}

func (o OrgMock) RemoveAllManagedAccounts(orgID string) error {
	panic("implement me")
}

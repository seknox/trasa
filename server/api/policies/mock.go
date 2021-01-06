package policies

import (
	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

//InitStoreMock will init mock state of this package
func InitStoreMock() *policyMock {
	m := new(policyMock)
	Store = m
	return m
}

type policyMock struct {
	mock.Mock
}

func (p policyMock) GetPolicy(policyID, orgID string) (models.Policy, error) {
	args := p.Called(policyID, orgID)
	return args.Get(0).(models.Policy), args.Error(1)
}

func (p policyMock) GetAllPolicies(orgID string) ([]models.Policy, error) {
	panic("implement me")
}

func (p policyMock) CreatePolicy(policy models.Policy) error {
	panic("implement me")
}

func (p policyMock) UpdatePolicy(policy models.Policy) error {
	panic("implement me")
}

func (p policyMock) DeletePolicy(policyID, orgID string) error {
	panic("implement me")
}

func (p policyMock) GetAccessPolicy(userID, serviceID, privilege, orgID string) (policy *models.Policy, adhoc bool, err error) {
	panic("implement me")
}

func (p policyMock) GetUserGroupAccessPolicyFromGroupNames(groups []string, serviceID, privilege, orgID string) (policy *models.Policy, adhoc bool, err error) {
	panic("implement me")
}

func (p policyMock) GetServiceUserGroupAccessPolicyFromGroupNames(groups []string, serviceID, privilege, orgID string) (policy *models.Policy, adhoc bool, err error) {
	panic("implement me")
}

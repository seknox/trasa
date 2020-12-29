package accessmap

import (
	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

//InitStoreMock will init mock state of this package
func InitStoreMock() *accessmapMock {
	m := new(accessmapMock)
	Store = m
	return m
}

type accessmapMock struct {
	mock.Mock
}

func (a accessmapMock) GetAllDynamicAccessRules(orgID string) ([]models.DynamicAccess, error) {
	panic("implement me")
}

func (a accessmapMock) CreateDynamicAccessRule(setting models.DynamicAccess) error {
	panic("implement me")
}

func (a accessmapMock) DeleteDynamicAccessRule(id, orgID string) error {
	panic("implement me")
}

func (a accessmapMock) GetDynamicAccessPolicy(groups []string, userID, orgID string) (*models.Policy, error) {
	panic("implement me")
}

func (a accessmapMock) GetServiceUserMaps(serviceID, orgID string) (appusers []models.AccessMapDetail, err error) {
	panic("implement me")
}

func (a accessmapMock) CreateServiceUserMap(appUser *models.ServiceUserMap) error {
	args := a.Called(appUser)
	return args.Error(0)
}

func (a accessmapMock) UpdateServiceUserMap(mapID, orgID, privilege string) error {
	panic("implement me")
}

func (a accessmapMock) DeleteServiceUserMap(mapID, orgID string) (string, error) {
	panic("implement me")
}

func (a accessmapMock) CheckIfPrivilegeExist(privilege, orgID, serviceID string) bool {
	panic("implement me")
}

func (a accessmapMock) CreateServiceGroupUserGroupMap(data *models.ServiceGroupUserGroupMap) error {
	panic("implement me")
}

func (a accessmapMock) UpdateServiceGroupUserGroupMap(mapID, orgID, privilege string) error {
	panic("implement me")
}

func (a accessmapMock) DeleteServiceGroupUserGroupMap(mapID, orgID string) (string, string, error) {
	panic("implement me")
}

func (a accessmapMock) GetAssignedUserGroupsWithPolicies(serviceGroupID, orgID string) ([]UserGroupOfServiceGroup, error) {
	panic("implement me")
}

func (a accessmapMock) GetUserGroupsToAddInServiceGroup(orgID string) ([]models.Group, error) {
	panic("implement me")
}

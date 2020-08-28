package groups

import (
	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

//InitStoreMock will init mock state of this package
func InitStoreMock() *groupMock {
	m := new(groupMock)
	Store = m
	return m
}

type groupMock struct {
	mock.Mock
}

func (g groupMock) Get(groupID, orgID string) (models.Group, error) {
	panic("implement me")
}

func (g groupMock) GetAllServiceGroups(orgID string) ([]models.Group, error) {
	panic("implement me")
}

func (g groupMock) GetAllUserGroups(orgID string) ([]models.Group, error) {
	panic("implement me")
}

func (g groupMock) Create(group *models.Group) error {
	panic("implement me")
}

func (g groupMock) Update(group *models.Group) error {
	panic("implement me")
}

func (g groupMock) Delete(groupID, orgID string) (name string, err error) {
	panic("implement me")
}

func (g groupMock) CheckIfUserInGroup(userID, orgID string, groupIDs []string) (bool, error) {
	args := g.Called(userID, orgID, groupIDs)
	return args.Bool(0), args.Error(1)
}

func (g groupMock) GetUsersInGroup(groupID, org string) ([]models.User, error) {
	panic("implement me")
}

func (g groupMock) GetUsersNotInGroup(groupID, orgID string) ([]models.User, error) {
	panic("implement me")
}

func (g groupMock) GetServicesNotInGroup(groupID, orgID string) ([]models.Service, error) {
	panic("implement me")
}

func (g groupMock) GetServicesInGroup(groupID, org string) ([]models.Service, error) {
	panic("implement me")
}

func (g groupMock) AddServicesToGroup(group models.Group, serviceIDs []string) (err error) {
	panic("implement me")
}

func (g groupMock) RemoveServicesFromGroup(groupID, orgID string, serviceIDs []string) error {
	panic("implement me")
}

func (g groupMock) AddUsersToGroup(group models.Group, userIDs []string) (err error) {
	panic("implement me")
}

func (g groupMock) RemoveUsersFromGroup(groupID, orgID string, userIDs []string) (err error) {
	panic("implement me")
}

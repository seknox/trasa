package groups

import (
	"github.com/seknox/trasa/global"
	"github.com/seknox/trasa/models"
)

func InitStore(state *global.State) {
	Store = GroupStore{state}
}

var Store Adapter

type GroupStore struct {
	*global.State
}

type Adapter interface {
	//CRUD

	Get(groupID, orgID string) (models.Group, error)
	GetAllServiceGroups(orgID string) ([]models.Group, error)
	GetAllUserGroups(orgID string) ([]models.Group, error)

	Create(group *models.Group) error
	Update(group *models.Group) error
	Delete(groupID, orgID string) (name string, err error)

	//maps

	CheckIfUserInGroup(userID, orgID string, groupIDs []string) (bool, error)

	GetUsersInGroup(groupID, org string) ([]models.User, error)
	GetUsersNotInGroup(groupID, orgID string) ([]models.User, error)

	GetServicesNotInGroup(groupID, orgID string) ([]models.Service, error)
	GetServicesInGroup(groupID, org string) ([]models.Service, error)

	AddServicesToGroup(group models.Group, serviceIDs []string) (err error)
	RemoveServicesFromGroup(groupID, orgID string, serviceIDs []string) error

	AddUsersToGroup(group models.Group, userIDs []string) (err error)
	RemoveUsersFromGroup(groupID, orgID string, userIDs []string) (err error)
}

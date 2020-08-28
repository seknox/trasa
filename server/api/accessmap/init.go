package accessmap

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	Store = accessMapStore{state}
}

//Store is the package state variable which contains database connections
var Store accessMapAdapter

type accessMapStore struct {
	*global.State
}

type accessMapAdapter interface {
	GetServiceUserMaps(serviceID, orgID string) (appusers []models.AccessMapDetail, err error)
	CreateServiceUserMap(appUser *models.ServiceUserMap) error
	UpdateServiceUserMap(mapID, orgID, privilege string) error
	DeleteServiceUserMap(mapID, orgID string) (string, error)

	CheckIfPrivilegeExist(privilege, orgID, serviceID string) bool
	//DeleteAppUserbasedOnUserID(userID, orgID string) error

	CreateServiceGroupUserGroupMap(data *models.ServiceGroupUserGroupMap) error
	UpdateServiceGroupUserGroupMap(mapID, orgID, privilege string) error
	DeleteServiceGroupUserGroupMap(mapID, orgID string) (string, string, error)

	GetAssignedUserGroupsWithPolicies(serviceGroupID, orgID string) ([]UserGroupOfServiceGroup, error)
	GetUserGroupsToAddInServiceGroup(orgID string) ([]models.Group, error)
}

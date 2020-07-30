package accessmap

import (
	"github.com/seknox/trasa/global"
	"github.com/seknox/trasa/models"
)

func InitStore(state *global.State) {
	Store = AccessMapStore{state}
}

var Store AccessMapAdapter

type AccessMapStore struct {
	*global.State
}

type AccessMapAdapter interface {
	GetServiceUserMaps(serviceID, orgID string) (appusers []models.AccessMapDetail, err error)
	CreateServiceUserMap(appUser *models.ServiceUserMap) error
	UpdateServiceUserMap(mapID, orgID, privilege string) error
	DeleteServiceUserMap(mapID, orgID string) (string, error)

	CheckIfPrivilegeExist(privilege, orgID, serviceID string) bool
	//DeleteAppUserbasedOnUserID(userID, orgID string) error

	CreateServiceGroupUserGroupMap(data *models.ServiceGroupUserGroupMap) error
	UpdateServiceGroupUserGroupMap(mapID, orgID, privilege string) error
	DeleteServiceGroupUserGroupMap(mapID, orgID string) (string, string, error)

	GetAssignedUserGroupsWithPolicies(serviceGroupID, orgID string) ([]userGroupOfServiceGroup, error)
	GetUserGroupsToAddInServiceGroup(orgID string) ([]models.Group, error)
}

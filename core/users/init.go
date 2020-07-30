package users

import (
	"github.com/seknox/trasa/core/users/userstoremocks"
	"github.com/seknox/trasa/global"
	"github.com/seknox/trasa/models"
)

func InitStore(con *global.State) {
	// initialize local state
	Store = UserStore{
		State: con,
	}

}

func InitStoreMock() *userstoremocks.UserStoreMock {
	// initialize local state
	mockUser := new(userstoremocks.UserStoreMock)
	Store = mockUser
	return mockUser
}

var Store UserAdapter

type UserStore struct {
	*global.State
}

type UserAdapter interface {
	//CRUD
	GetFromID(userID, orgID string) (*models.User, error)
	GetAll(orgID string) ([]models.User, error)
	GetAdminEmails(orgID string) ([]string, error)
	//GetTotalForOrg(orgID string) (int64, error)
	Create(user *models.UserWithPass) error
	Delete(userID, orgID string) (string, string, error)
	Update(user models.User) error

	UpdatePublicKey(userID string, publicKey string) error

	UpdatePassword(userID, password string) error
	GetEnforcedPolicy(userID, orgID, enforceType string) (policy models.PolicyEnforcer, err error)
	DeleteActivePolicy(userID, orgID, enforceType string) error
	GetPasswordState(userID, orgID string) (models.PasswordState, error)
	UpdatePasswordState(userID, orgID, oldpassword string, time int64) error
	EnforcePolicy(policy models.PolicyEnforcer) error

	//Groups

	GetGroups(userID, orgID string) ([]models.Group, error)

	//Appusers
	GetAccessMapDetails(userID, orgID string) ([]models.AccessMapDetail, error)

	GetAssignedServices(userID, orgID string) (services []models.Service, err error)

	DeleteAllUserAccessMaps(userID, orgID string) error

	//Devices

	GetAllDevices(userID, orgID string) ([]models.UserDevice, error)
	DeregisterUserDevices(userID, orgID string) error
	//GetTOTPDevices retrieves user devices with totp secrets
	GetTOTPDevices(userID, orgID string) ([]models.UserDevice, error)
}

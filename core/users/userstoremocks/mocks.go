package userstoremocks

import (
	"github.com/seknox/trasa/models"
	"github.com/test-go/testify/mock"
)

type UserStoreMock struct {
	mock.Mock
}

func (us UserStoreMock) GetFromID(userID, orgID string) (*models.User, error) {
	args := us.Called(userID, orgID)

	return args.Get(0).(*models.User), args.Error(1)
}

func (us UserStoreMock) GetAdminEmails(orgID string) ([]string, error) {
	args := us.Called(orgID)

	return args.Get(0).([]string), args.Error(1)
}

func (us UserStoreMock) GetAll(orgID string) ([]models.User, error) {
	panic("implement me")
}

func (us UserStoreMock) Create(user *models.UserWithPass) error {
	panic("implement me")
}

func (us UserStoreMock) Delete(userID, orgID string) (string, string, error) {
	panic("implement me")
}

func (us UserStoreMock) Update(user models.User) error {
	panic("implement me")
}

func (us UserStoreMock) UpdatePublicKey(userID string, publicKey string) error {
	panic("implement me")
}

func (us UserStoreMock) UpdatePassword(userID, password string) error {
	panic("implement me")
}

func (us UserStoreMock) GetEnforcedPolicy(userID, orgID, enforceType string) (policy models.PolicyEnforcer, err error) {
	return models.PolicyEnforcer{}, nil
	panic("implement me")
}

func (us UserStoreMock) DeleteActivePolicy(userID, orgID, enforceType string) error {
	panic("implement me")
}

func (us UserStoreMock) GetPasswordState(userID, orgID string) (models.PasswordState, error) {
	args := us.Called(userID, orgID)

	return args.Get(0).(models.PasswordState), args.Error(0)
}

func (us UserStoreMock) UpdatePasswordState(userID, orgID, oldpassword string, time int64) error {
	panic("implement me")
}

func (us UserStoreMock) EnforcePolicy(policy models.PolicyEnforcer) error {
	return us.Called(policy).Error(0)
}

func (us UserStoreMock) GetGroups(userID, orgID string) ([]models.Group, error) {
	panic("implement me")
}

func (us UserStoreMock) GetAccessMapDetails(userID, orgID string) ([]models.AccessMapDetail, error) {
	panic("implement me")
}

func (us UserStoreMock) GetAssignedServices(userID, orgID string) (services []models.Service, err error) {
	panic("implement me")
}

func (us UserStoreMock) DeleteAllUserAccessMaps(userID, orgID string) error {
	panic("implement me")
}

func (us UserStoreMock) GetAllDevices(userID, orgID string) ([]models.UserDevice, error) {
	panic("implement me")
}

func (us UserStoreMock) DeregisterUserDevices(userID, orgID string) error {
	panic("implement me")
}

func (us UserStoreMock) GetTOTPDevices(userID, orgID string) ([]models.UserDevice, error) {
	return []models.UserDevice{{
		UserID:        userID,
		OrgID:         orgID,
		DeviceID:      "21321jk3h12i3",
		MachineID:     "bcs8djcsicojsdm",
		DeviceType:    "mobile",
		FcmToken:      "",
		TotpSec:       "VVSOVFTGXE",
		DeviceHygiene: models.DeviceHygiene{},
	}}, nil
}

var dataMap = map[string]models.User{
	"asdasd-6c09305f-4080-47f3-b66a-fe4dd4efb827": models.User{
		ID:         "asdasd",
		OrgID:      "6c09305f-4080-47f3-b66a-fe4dd4efb827",
		UserName:   "bha",
		FirstName:  "Bhargab",
		MiddleName: "",
		LastName:   "Acharya",
		Email:      "bhrg3se@gmail.com",
		Groups:     nil,
		UserRole:   "",
		Status:     false,
		IdpName:    "",
		ExternalID: "",
		CreatedAt:  0,
		UpdatedAt:  0,
	},

	"6c09305f-4080-47f3-b66a-fe4dd4efb824-6c09305f-4080-47f3-b66a-fe4dd4efb827": models.User{
		ID:         "6c09305f-4080-47f3-b66a-fe4dd4efb824",
		OrgID:      "6c09305f-4080-47f3-b66a-fe4dd4efb827",
		UserName:   "tree1",
		FirstName:  "sakshyam",
		MiddleName: "",
		LastName:   "shah",
		Email:      "sakshyam@seknox.com",
		Groups:     nil,
		UserRole:   "orgAdmin",
		Status:     true,
		IdpName:    "trasa",
		ExternalID: "",
		CreatedAt:  1534973524,
		UpdatedAt:  1586413415,
	},
}

//func (us UserStoreMock) GetByID(userID, orgID string) (*models.User, error) {
//
//	//args:=us.Called(userID,orgID)
//	//
//	//return args.Get(0).(*models.User),args.Error(0)
//	//user, ok := dataMap[userID+"-"+orgID]
//	//if !ok {
//	//	return nil, sql.ErrNoRows
//	//}
//	//return &user, nil
//}

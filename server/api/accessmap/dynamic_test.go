package accessmap

import (
	"github.com/seknox/trasa/server/api/groups"
	"github.com/seknox/trasa/server/api/notif"
	"github.com/seknox/trasa/server/api/policies"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/test-go/testify/mock"
	"reflect"
	"testing"
)

func TestCreateDynamicService(t *testing.T) {
	systemstore := system.InitStoreMock()
	groupstore := groups.InitStoreMock()
	servicestore := services.InitServiceStoreMock()
	accessstore := InitStoreMock()
	userstore := users.InitStoreMock()
	notifstore := notif.InitStoreMock()

	notifstore.On("SendEmail", "orgWithDynamicAccessEnabled", consts.EMAIL_DYNAMIC_ACCESS, mock.Anything).Return(nil)
	userstore.On("GetAdminEmails", mock.Anything).Return([]string{"admin@domain"}, nil)
	accessstore.On("CreateServiceUserMap", mock.Anything).Return(nil)
	servicestore.On("Create", mock.Anything).Return(nil)

	groupstore.On("CheckIfUserInGroup", "authorisedUserID", "orgWithDynamicAccessEnabled", []string{"someGroupID"}).
		Return(true, nil)
	groupstore.On("CheckIfUserInGroup", "unauthorisedUserID", "orgWithDynamicAccessEnabled", []string{"someGroupID"}).
		Return(false, nil)

	systemstore.On("GetGlobalSetting", "orgWithDynamicAccessEnabled", consts.GLOBAL_DYNAMIC_ACCESS).Return(models.GlobalSettings{
		SettingID:    "123213123",
		OrgID:        "orgWithDynamicAccessEnabled",
		Status:       true,
		SettingType:  consts.GLOBAL_DYNAMIC_ACCESS,
		SettingValue: `{"userGroups": ["someGroupID"], "policyID": "f022d753-5f5f-4035-b3d4-59db0079d634"}`,
	}, nil)

	systemstore.On("GetGlobalSetting", "orgWithDynamicAccessDisabled", consts.GLOBAL_DYNAMIC_ACCESS).Return(models.GlobalSettings{
		SettingID:    "123213123",
		OrgID:        "orgWithDynamicAccessEnabled",
		Status:       false,
		SettingType:  consts.GLOBAL_DYNAMIC_ACCESS,
		SettingValue: `{"userGroups": ["someGroupID"], "policyID": "f022d753-5f5f-4035-b3d4-59db0079d634"}`,
	}, nil)

	type args struct {
		hostname    string
		serviceType string
		userID      string
		userEmail   string
		privilege   string
		orgID       string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Service
		wantErr bool
	}{
		{
			"should fail if dynamic access is disabled",
			args{
				hostname:    "someHost",
				serviceType: "ssh",
				userID:      "authorisedUserID",
				privilege:   "somePrivilege",
				orgID:       "orgWithDynamicAccessDisabled",
			},
			nil,
			true,
		},
		{
			"should fail if user is not allowed to dynamic access",
			args{
				hostname:    "someHost",
				serviceType: "ssh",
				userID:      "unauthorisedUserID",
				privilege:   "somePrivilege",
				orgID:       "orgWithDynamicAccessEnabled",
			},
			nil,
			true,
		},

		{
			"should create service if service is not created",
			args{
				hostname:    "someHost",
				serviceType: "ssh",
				userID:      "authorisedUserID",
				privilege:   "somePrivilege",
				orgID:       "orgWithDynamicAccessEnabled",
			},
			&models.Service{
				OrgID:    "orgWithDynamicAccessEnabled",
				Name:     "someHost",
				Hostname: "someHost",
				Type:     "ssh",
			},
			false,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateDynamicService(tt.args.hostname, tt.args.serviceType, tt.args.userID, tt.args.userEmail, tt.args.privilege, tt.args.orgID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDynamicService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want != nil {
				t.Errorf("CreateDynamicService() got = %v, want %v", got, tt.want)
			}
			if got != nil && tt.want == nil {
				t.Errorf("CreateDynamicService() got = %v, want %v", got, tt.want)
			}

			if got == nil {
				return
			}
			if tt.want == nil {
				return
			}

			if got.Name != tt.want.Name {
				t.Errorf("CreateDynamicService().Name got = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Hostname != tt.want.Hostname {
				t.Errorf("CreateDynamicService().Hostname got = %v, want %v", got.Hostname, tt.want.Hostname)
			}
			if got.OrgID != tt.want.OrgID {
				t.Errorf("CreateDynamicService().OrgID got = %v, want %v", got.OrgID, tt.want.OrgID)
			}
			if got.Type != tt.want.Type {
				t.Errorf("CreateDynamicService().Type got = %v, want %v", got.Type, tt.want.Type)
			}

		})
	}
}

func TestCreateDynamicAccessMap(t *testing.T) {
	systemstore := system.InitStoreMock()
	groupstore := groups.InitStoreMock()
	servicestore := services.InitServiceStoreMock()
	accessstore := InitStoreMock()
	userstore := users.InitStoreMock()
	notifstore := notif.InitStoreMock()
	policystore := policies.InitStoreMock()

	fullP := models.Policy{
		PolicyID:         "f022d753-5f5f-4035-b3d4-59db0079d634",
		OrgID:            "orgWithDynamicAccessEnabled",
		PolicyName:       "full",
		DayAndTime:       []models.DayAndTimePolicy{{Days: []string{"Sunday"}}},
		TfaRequired:      false,
		RecordSession:    false,
		FileTransfer:     false,
		IPSource:         "",
		AllowedCountries: "",
		DevicePolicy:     models.DevicePolicy{},
		RiskThreshold:    0,
		CreatedAt:        0,
		UpdatedAt:        0,
		Expiry:           "",
		IsExpired:        false,
		UsedBy:           0,
	}

	policystore.On("GetPolicy", "f022d753-5f5f-4035-b3d4-59db0079d634", "orgWithDynamicAccessEnabled").Return(fullP, nil)
	notifstore.On("SendEmail", "orgWithDynamicAccessEnabled", consts.EMAIL_DYNAMIC_ACCESS, mock.Anything).Return(nil)
	userstore.On("GetAdminEmails", mock.Anything).Return([]string{"admin@domain"}, nil)
	accessstore.On("CreateServiceUserMap", mock.Anything).Return(nil)
	servicestore.On("GetFromID", "someServiceID").Return(&models.Service{}, nil)

	groupstore.On("CheckIfUserInGroup", "authorisedUserID", "orgWithDynamicAccessEnabled", []string{"someGroupID"}).
		Return(true, nil)
	groupstore.On("CheckIfUserInGroup", "unauthorisedUserID", "orgWithDynamicAccessEnabled", []string{"someGroupID"}).
		Return(false, nil)

	systemstore.On("GetGlobalSetting", "orgWithDynamicAccessEnabled", consts.GLOBAL_DYNAMIC_ACCESS).Return(models.GlobalSettings{
		SettingID:    "123213123",
		OrgID:        "orgWithDynamicAccessEnabled",
		Status:       true,
		SettingType:  consts.GLOBAL_DYNAMIC_ACCESS,
		SettingValue: `{"userGroups": ["someGroupID"], "policyID": "f022d753-5f5f-4035-b3d4-59db0079d634"}`,
	}, nil)

	systemstore.On("GetGlobalSetting", "orgWithDynamicAccessDisabled", consts.GLOBAL_DYNAMIC_ACCESS).Return(models.GlobalSettings{
		SettingID:    "123213123",
		OrgID:        "orgWithDynamicAccessEnabled",
		Status:       false,
		SettingType:  consts.GLOBAL_DYNAMIC_ACCESS,
		SettingValue: `{"userGroups": ["someGroupID"], "policyID": "f022d753-5f5f-4035-b3d4-59db0079d634"}`,
	}, nil)

	type args struct {
		serviceID string
		userID    string
		userEmail string
		privilege string
		orgID     string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Policy
		wantErr bool
	}{
		{
			"should fail if  dynamic access is disabled",
			args{
				serviceID: "someServiceID",
				userID:    "authorisedUserID",
				userEmail: "",
				privilege: "somePrivilege",
				orgID:     "orgWithDynamicAccessDisabled",
			},
			nil,
			true,
		}, {
			"should fail if user is not allowed to dynamic access",
			args{
				serviceID: "someServiceID",
				userID:    "unauthorisedUserID",
				userEmail: "admin@domain",
				privilege: "somePrivilege",
				orgID:     "orgWithDynamicAccessDisabled",
			},
			nil,
			true,
		},

		{
			"should pass",
			args{
				serviceID: "someServiceID",
				userID:    "authorisedUserID",
				userEmail: "admin@dmain",
				privilege: "somePrivilege",
				orgID:     "orgWithDynamicAccessEnabled",
			},
			&fullP,
			false,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateDynamicAccessMap(tt.args.serviceID, tt.args.userID, tt.args.userEmail, tt.args.privilege, tt.args.orgID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDynamicAccessMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateDynamicAccessMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

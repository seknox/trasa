package rdpproxy

import (
	"reflect"
	"testing"

	"github.com/seknox/trasa/core/crypt"
	"github.com/seknox/trasa/server/api/providers/vault"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
)

func Test_handlePass(t *testing.T) {
	_ = crypt.InitStoreMock()
	systemstore := system.InitStoreMock()
	vaultstore := vault.InitStoreMock()

	type args struct {
		params *models.ConnectionParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.UpstreamCreds
		wantErr bool
	}{
		{
			"when password is passed in params, password from vault should be ignored",
			args{&models.ConnectionParams{
				ServiceID:   "123",
				OrgID:       "abc",
				Privilege:   "admin",
				Password:    "abc123",
				UserID:      "345f6tg7yhui",
				ServiceType: "",
			}},
			&models.UpstreamCreds{
				Password:          "abc123",
				HostCert:          "",
				HostCaCert:        "",
				UserCaCert:        "",
				ClientCert:        "",
				ClientKey:         "",
				SkipHostVerify:    false,
				MinimumChar:       8,
				ZxcvbnScore:       3,
				EnforceStrongPass: false,
			},
			false,
		},
		{
			"when password in params is blank, password from vault should be used",
			args{&models.ConnectionParams{
				ServiceID:   "123",
				OrgID:       "abc",
				Privilege:   "admin",
				Password:    "",
				UserID:      "345f6tg7yhui",
				ServiceType: "",
			}},
			&models.UpstreamCreds{
				Password:          "passwordfromvault",
				HostCert:          "",
				HostCaCert:        "",
				UserCaCert:        "",
				ClientCert:        "",
				ClientKey:         "",
				SkipHostVerify:    false,
				MinimumChar:       8,
				ZxcvbnScore:       3,
				EnforceStrongPass: false,
			},
			false,
		},
		{
			"when enforce password is true, password policy should be checked",
			args{&models.ConnectionParams{
				ServiceID:   "123",
				OrgID:       "xyz",
				Privilege:   "admin",
				Password:    "weakpas",
				UserID:      "345f6tg7yhui",
				ServiceType: "",
			}},
			nil,
			true,
		},
		{
			"when enforce password is true, password policy should be checked",
			args{&models.ConnectionParams{
				ServiceID:   "123",
				OrgID:       "xyz",
				Privilege:   "admin",
				Password:    "sTronGp@@@sword123",
				UserID:      "345f6tg7yhui",
				ServiceType: "",
			}},
			&models.UpstreamCreds{
				Password:          "sTronGp@@@sword123",
				HostCert:          "",
				HostCaCert:        "",
				UserCaCert:        "",
				ClientCert:        "",
				ClientKey:         "",
				SkipHostVerify:    false,
				MinimumChar:       8,
				ZxcvbnScore:       3,
				EnforceStrongPass: true,
			},
			false,
		},

		// TODO: Add test cases.
	}

	systemstore.On("GetGlobalSetting", "abc", consts.GLOBAL_PASSWORD_CONFIG).Return(models.GlobalSettings{
		SettingID:    "123213123",
		OrgID:        "abc",
		Status:       false,
		SettingType:  consts.GLOBAL_PASSWORD_CONFIG,
		SettingValue: `{"enforceStrongPass": false, "expiry": "never", "minimumChars": 8, "zxcvbnScore": 3}`,
		UpdatedBy:    "",
		UpdatedOn:    0,
	}, nil)
	systemstore.On("GetGlobalSetting", "xyz", consts.GLOBAL_PASSWORD_CONFIG).Return(models.GlobalSettings{
		SettingID:    "123213123",
		OrgID:        "xyz",
		Status:       false,
		SettingType:  consts.GLOBAL_PASSWORD_CONFIG,
		SettingValue: `{"enforceStrongPass": true, "expiry": "never", "minimumChars": 8, "zxcvbnScore": 3}`,
		UpdatedBy:    "",
		UpdatedOn:    0,
	}, nil)

	vaultstore.On("GetSecret", "xyz", "123", "key", "admin").
		Return("", nil)
	vaultstore.On("GetSecret", "abc", "123", "key", "admin").
		Return("", nil)

	vaultstore.On("GetSecret", "xyz", "123", "password", "admin").
		Return("passwordfromvault", nil)
	vaultstore.On("GetSecret", "abc", "123", "password", "admin").
		Return("passwordfromvault", nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handlePass(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("handlePass() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handlePass() got = %v, want %v", got, tt.want)
			}
		})
	}
}

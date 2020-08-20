package services

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/seknox/trasa/server/api/providers/ca"
	"github.com/seknox/trasa/server/api/providers/vault/tsxvault"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
)

func TestGetUpstreamCreds(t *testing.T) {
	cryptstore := ca.InitStoreMock()
	systemstore := system.InitStoreMock()
	vaultstore := tsxvault.InitStoreMock()

	type args struct {
		user        string
		serviceID   string
		serviceType string
		orgID       string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.UpstreamCreds
		wantErr bool
	}{
		{
			name: "get creds with password and cser key",
			args: args{"root", "123asdaxbas5werty", "ssh", "abc"},
			want: &models.UpstreamCreds{
				Password:          "testpass",
				HostCert:          "",
				HostCaCert:        "",
				UserCaCert:        "",
				ClientCert:        "",
				ClientKey:         "testkey",
				SkipHostVerify:    false,
				MinimumChar:       0,
				ZxcvbnScore:       0,
				EnforceStrongPass: false,
			},
			wantErr: false,
		}, {
			name: "get credds with host CA and user CA",
			args: args{"root", "12345werty", "ssh", "abc"},
			want: &models.UpstreamCreds{
				Password:          "testpass",
				HostCert:          "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcOknjd1n3dWnqtDWHbZyZUhK5Sr13f7demahHZ+HSLoWmz2TVfq1djhoME8wk581mOZnAbhhDZtL2JLIfYpoe8gZE9oqVqgCdSg3h9B4GqDN0Dwje7GVgVnv8Zbad+Tswcgn8j8NRoU7h6gfPW+kAvDuk0eB+liaurIx7+J5vejAHiidQ7ZmA+TomaP0ZIlI1QXPyu2oQjPkSuk8M3EFnaI4NRF/QSUe4LBgf113hNf3mBcRJUFdT1b3aLtpp7UlONiS2bhnuV53mTg8aigeQ2ixUe+MYc8hsdXPfRpjjdPa9wAyIq4lt1m96d0aodqhOZgPt5QnRbExqIGvgPSIB bhrg3se@bananana",
				HostCaCert:        "host ca",
				UserCaCert:        "user ca",
				ClientCert:        "",
				ClientKey:         "testkey",
				SkipHostVerify:    false,
				MinimumChar:       0,
				ZxcvbnScore:       0,
				EnforceStrongPass: false,
			},
			wantErr: false,
		},
		{
			name: "get rdp creds",
			args: args{"root", "12345werty", "rdp", "abc"},
			want: &models.UpstreamCreds{
				Password:          "testpass",
				HostCert:          "",
				HostCaCert:        "",
				UserCaCert:        "",
				ClientCert:        "",
				ClientKey:         "testkey",
				SkipHostVerify:    false,
				MinimumChar:       0,
				ZxcvbnScore:       0,
				EnforceStrongPass: false,
			},
			wantErr: false,
		},

		// TODO: Add more test cases.
	}

	for _, tt := range tests {

		if tt.args.serviceType == "ssh" {
			fmt.Println(tt.args.serviceType)
			cryptstore.
				On("GetCertHolder", consts.CERT_TYPE_SSH_HOST_KEY, tt.args.serviceID, tt.args.orgID).
				Return(models.CertHolder{
					CertID:   "12321312",
					OrgID:    tt.args.orgID,
					EntityID: tt.args.serviceID,
					Cert:     []byte(tt.want.HostCert),
					CertType: consts.CERT_TYPE_SSH_HOST_KEY,
				}, nil).Times(1)
			cryptstore.
				On("GetCertHolder", consts.CERT_TYPE_SSH_CA, "user", tt.args.orgID).
				Return(models.CertHolder{
					CertID:   "12321312",
					OrgID:    tt.args.orgID,
					EntityID: "user",
					Cert:     []byte(tt.want.UserCaCert),
					CertType: consts.CERT_TYPE_SSH_CA,
				}, nil).Times(2)
			cryptstore.
				On("GetCertHolder", consts.CERT_TYPE_SSH_CA, "host", tt.args.orgID).
				Return(models.CertHolder{
					CertID:   "1232skamdaskmd1312",
					OrgID:    tt.args.orgID,
					EntityID: "user",
					Cert:     []byte(tt.want.HostCaCert),
					CertType: consts.CERT_TYPE_SSH_CA,
				}, nil).Times(1)

		}

		systemstore.On("GetGlobalSetting", "abc", consts.GLOBAL_PASSWORD_CONFIG).Return(models.GlobalSettings{
			SettingID:    "123213123",
			OrgID:        "abc",
			Status:       false,
			SettingType:  consts.GLOBAL_PASSWORD_CONFIG,
			SettingValue: "{}",
			UpdatedBy:    "",
			UpdatedOn:    0,
		}, nil).Times(1)

		vaultstore.On("GetSecret", tt.args.orgID, tt.args.serviceID, "key", tt.args.user).
			Return(tt.want.ClientKey, nil).Times(1)

		vaultstore.On("GetSecret", tt.args.orgID, tt.args.serviceID, "password", tt.args.user).
			Return(tt.want.Password, nil).Times(1)

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUpstreamCreds(tt.args.user, tt.args.serviceID, tt.args.serviceType, tt.args.orgID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUpstreamCreds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUpstreamCreds() got = %v, want %v", got, tt.want)
			}
		})
	}

	vaultstore.AssertExpectations(t)
	cryptstore.AssertExpectations(t)
}

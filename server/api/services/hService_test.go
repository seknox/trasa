package services

import (
	"testing"

	"github.com/seknox/trasa/server/models"
)

func Test_fillInitialFields(t *testing.T) {
	type args struct {
		req *models.Service
	}
	type testc struct {
		name string
		args args
	}

	tests := []testc{
		{
			name: "Create Service with test fields",
			args: args{
				req: &models.Service{
					ID:                    "1c09305f-4080-47f3-b66b-fe4dd4efb820",
					OrgID:                 "6c09305f-4080-47f3-b66a-fe4dd4efb827",
					Name:                  "testServiceName",
					SecretKey:             "4c09305f-4080-77f3-b66a-fe4dd4efb821",
					Passthru:              false,
					Hostname:              "1.1.1.1",
					Type:                  "rdp",
					ManagedAccounts:       "",
					NativeLog:             false,
					ProxyConfig:           models.ReverseProxy{},
					ExternalProviderName:  "",
					ExternalID:            "",
					ExternalSecurityGroup: "",
					DistroName:            "",
					DistroVersion:         "",
					IPDetails:             models.IPDetails{},
					CreatedAt:             0,
					UpdatedAt:             0,
					DeletedAt:             0,
				},
			},
		},
		{
			name: "",
			args: args{
				req: &models.Service{
					ID:                    "",
					OrgID:                 "",
					Name:                  "",
					SecretKey:             "",
					Passthru:              false,
					Hostname:              "",
					Type:                  "",
					ManagedAccounts:       "",
					RemoteAppName:         "",
					Adhoc:                 false,
					NativeLog:             false,
					RdpProtocol:           "",
					ProxyConfig:           models.ReverseProxy{},
					PublicKey:             "",
					ExternalProviderName:  "",
					ExternalID:            "",
					ExternalSecurityGroup: "",
					DistroName:            "",
					DistroVersion:         "",
					IPDetails:             models.IPDetails{},
					CreatedAt:             0,
					UpdatedAt:             0,
					DeletedAt:             0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := fillInitialFields(tt.args.req)
			if got.ID == "" {
				t.Errorf("fillInitialFields().ID = %v, want %v", got.ID, "<not nil>")
			}
			if got.Name != tt.args.req.Name {
				t.Errorf("fillInitialFields().Name = %v, want %v", got.Name, tt.args.req.Name)
			}
			if got.SecretKey == "" {
				t.Errorf("fillInitialFields().SecretKey = %v, want %v", got.SecretKey, "<not nil>")
			}
			if got.CreatedAt == 0 {
				t.Errorf("fillInitialFields().CreatedAt = %v, want %v", got.CreatedAt, "<not nil>")
			}
			if got.UpdatedAt == 0 {
				t.Errorf("fillInitialFields().UpdatedAt = %v, want %v", got.UpdatedAt, "<not nil>")
			}

		})
	}
}

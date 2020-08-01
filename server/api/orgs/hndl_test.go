package orgs

import (
	"testing"

	"github.com/seknox/trasa/server/models"
)

func Test_getoOrgUserFromRequest(t *testing.T) {
	type args struct {
		req models.InitSignup
	}
	tests := []struct {
		name     string
		args     args
		wantOrg  models.Org
		wantUser models.UserWithPass
	}{
		{
			name: "",
			args: args{models.InitSignup{
				OrgName:     "Some Org",
				UserName:    "root",
				FirstName:   "Bhargab",
				MiddleName:  "",
				LastName:    "Acharya",
				Email:       "me@bhargab.com.np",
				PhoneNumber: "987654321",
				Country:     "Nepal",
				Timezone:    "Asia/Kathmandu",
			}},
			wantOrg: models.Org{
				OrgName:        "Some Org",
				Domain:         "bhargab.com.np",
				PrimaryContact: "me@bhargab.com.np",
				Timezone:       "Asia/Kathmandu",
				PhoneNumber:    "987654321",
			},
			wantUser: models.UserWithPass{
				User: models.User{
					UserName:   "root",
					FirstName:  "Bhargab",
					MiddleName: "",
					LastName:   "Acharya",
					Email:      "me@bhargab.com.np",
					UserRole:   "orgAdmin",
					Status:     true,
					IdpName:    "trasa",
				},
			},
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOrg, gotUser := getoOrgUserFromRequest(tt.args.req)

			if gotOrg.Timezone != tt.wantOrg.Timezone {
				t.Errorf("getoOrgUserFromRequest() gotOrg.Timezone = %v, want %v", gotOrg.Timezone, tt.wantOrg.Timezone)
			}

			if gotOrg.Domain != tt.wantOrg.Domain {
				t.Errorf("getoOrgUserFromRequest() gotOrg.Domain = %v, want %v", gotOrg.Domain, tt.wantOrg.Domain)
			}

			if gotOrg.PhoneNumber != tt.wantOrg.PhoneNumber {
				t.Errorf("getoOrgUserFromRequest() gotOrg.PhoneNumber = %v, want %v", gotOrg.PhoneNumber, tt.wantOrg.PhoneNumber)
			}

			if gotOrg.PrimaryContact != tt.wantOrg.PrimaryContact {
				t.Errorf("getoOrgUserFromRequest() gotOrg.PrimaryContact = %v, want %v", gotOrg.PrimaryContact, tt.wantOrg.PrimaryContact)
			}

			if gotUser.Email != tt.wantUser.Email {
				t.Errorf("getoOrgUserFromRequest() gotUser.Email = %v, want %v", gotUser.Email, tt.wantUser.Email)
			}
			if gotUser.FirstName != tt.wantUser.FirstName {
				t.Errorf("getoOrgUserFromRequest() gotUser.FirstName = %v, want %v", gotUser.FirstName, tt.wantUser.FirstName)
			}
			if gotUser.LastName != tt.wantUser.LastName {
				t.Errorf("getoOrgUserFromRequest() gotUser.LastName = %v, want %v", gotUser.LastName, tt.wantUser.LastName)
			}
			if gotUser.IdpName != tt.wantUser.IdpName {
				t.Errorf("getoOrgUserFromRequest() gotUser.IdpName = %v, want %v", gotUser.IdpName, tt.wantUser.IdpName)
			}
			if gotUser.CreatedAt == 0 {
				t.Errorf("getoOrgUserFromRequest() gotUser.CreatedAt = %v, want %v", gotUser.CreatedAt, tt.wantUser.CreatedAt)
			}
			if gotUser.UserName != tt.wantUser.UserName {
				t.Errorf("getoOrgUserFromRequest() gotUser.UserName = %v, want %v", gotUser.UserName, tt.wantUser.UserName)
			}

		})
	}
}

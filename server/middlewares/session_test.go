package middlewares

import (
	"reflect"
	"testing"

	"github.com/seknox/trasa/server/api/auth"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/models"
)

func Test_getUserContext(t *testing.T) {
	userstore := users.InitStoreMock()
	orgstore := orgs.InitStoreMock()
	_ = redis.InitStoreMock()

	mockUser := models.User{
		ID:         "someUserID",
		OrgID:      "someOrgID",
		UserName:   "testUname",
		FirstName:  "Bha",
		MiddleName: "",
		LastName:   "Ach",
		Email:      "user@example.com",
		UserRole:   "orgAdmin",
		Status:     true,
		IdpName:    "trasa",
	}
	mockOrg := models.Org{
		ID:             "someOrgID",
		OrgName:        "testOrg",
		Domain:         "example.com",
		PrimaryContact: "user@example.com",
		Timezone:       "Asia/Kathmandu",
		PhoneNumber:    "12345678",
	}

	userstore.
		On("GetFromID", "someUserID", "someOrgID").
		Return(&mockUser, nil).Times(2)

	orgstore.
		On("Get", "someOrgID").
		Return(mockOrg, nil).Times(2)

	var nulUC models.UserContext

	type args struct {
		orgID  string
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    models.UserContext
		wantErr bool
	}{
		{
			name:    "blank userID orgID",
			args:    args{"", ""},
			want:    nulUC,
			wantErr: true,
		},
		{
			name:    "with valid userID orgID",
			args:    args{"someOrgID", "someUserID"},
			want:    models.UserContext{&mockUser, mockOrg, "", ""},
			wantErr: false,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setSession, setCsrf, err := auth.SetSession(tt.args.userID, tt.args.orgID, "", "")
			if err != nil {
				t.Errorf("failed setting session: %v", err)
			}

			got, err := getUserContext(setSession, setCsrf)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUserContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUserContext() got = %v, want %v", got, tt.want)
			}
		})
	}
}

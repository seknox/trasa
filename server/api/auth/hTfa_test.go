package auth

import (
	"reflect"
	"testing"

	"github.com/seknox/trasa/server/api/devices"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"golang.org/x/crypto/bcrypt"
)

func Test_handleIntentResponse(t *testing.T) {

	_ = users.InitStoreMock()
	redis.InitStoreMock()
	authstore := InitStoreMock()
	//systemstore := system.InitStoreMock()
	orgstore := orgs.InitStoreMock()
	//_ = logs.InitStoreMock()
	//_ = misc.InitMock()

	pass, _ := bcrypt.GenerateFromPassword([]byte("testpass@123"), bcrypt.DefaultCost)
	authstore.
		On("GetLoginDetails", "rot", "").
		Return(&models.UserWithPass{
			User: models.User{
				ID:         "123",
				OrgID:      "abc",
				UserName:   "rot",
				FirstName:  "B",
				MiddleName: "",
				LastName:   "Acharya",
				Email:      "user@example.com",
				UserRole:   "orgAdmin",
				Status:     true,
				IdpName:    "trasa",
			},
			OrgName:  "testOrg",
			Password: string(pass),
		}, nil)

	orgstore.On("Get", "abc").Return(models.Org{
		ID:       "abc",
		OrgName:  "someOrg",
		Timezone: "Asia/Kathmandu",
	}, nil)

	type args struct {
		req  TfaRequest
		user *models.User
	}
	tests := []struct {
		name         string
		args         args
		wantStatus   string
		wantReason   consts.FailedReason
		wantIntent   string
		wantSession  bool
		wantRespType reflect.Type
	}{
		{
			name: "successful dash login ",
			args: args{TfaRequest{
				Token:     "1234",
				TfaMethod: "totp",
				Totp:      "4432",
				Intent:    consts.AUTH_REQ_DASH_LOGIN,
			}, &models.User{
				ID:       "123456789",
				OrgID:    "abc",
				Email:    "user@example.com",
				UserName: "rot",
				UserRole: "orgAdmin",
				Status:   true,
				IdpName:  "trasa",
			}},
			wantStatus:   "success",
			wantReason:   "",
			wantSession:  true,
			wantIntent:   "DashboardLogin",
			wantRespType: reflect.TypeOf(userAuthSessionResp{}),
		},

		{
			name: "successful enrol device ",
			args: args{TfaRequest{
				Token:     "1234",
				TfaMethod: "totp",
				Totp:      "4432",
				Intent:    consts.AUTH_REQ_ENROL_DEVICE,
			}, &models.User{
				ID:       "123456789",
				OrgID:    "abc",
				Email:    "user@example.com",
				UserName: "rot",
				UserRole: "orgAdmin",
				Status:   true,
				IdpName:  "trasa",
			}},
			wantStatus:   "success",
			wantReason:   "",
			wantSession:  false,
			wantIntent:   consts.AUTH_RESP_ENROL_DEVICE,
			wantRespType: reflect.TypeOf(devices.EnrolDeviceStruct{}),
		},

		{
			name: "successful change pass",
			args: args{TfaRequest{
				Token:     "1234",
				TfaMethod: "totp",
				Totp:      "4432",
				Intent:    consts.AUTH_REQ_CHANGE_PASS,
			}, &models.User{
				ID:       "123456789",
				OrgID:    "abc",
				Email:    "user@example.com",
				UserName: "rot",
				UserRole: "orgAdmin",
				Status:   true,
				IdpName:  "trasa",
			}},
			wantStatus:   "success",
			wantReason:   "",
			wantSession:  false,
			wantIntent:   consts.AUTH_RESP_CHANGE_PASS,
			wantRespType: reflect.TypeOf(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, gotReason, gotIntent, gotSession, gotResp := handleIntentResponse(tt.args.req, tt.args.user, "", "")
			if gotStatus != tt.wantStatus {
				t.Errorf("handleIntentResponse() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
			if gotReason != tt.wantReason {
				t.Errorf("handleIntentResponse() gotReason = %v, want %v", gotReason, tt.wantReason)
			}
			if gotIntent != tt.wantIntent {
				t.Errorf("handleIntentResponse() gotIntent = %v, want %v", gotIntent, tt.wantIntent)
			}
			if tt.wantSession && gotSession == "" {
				t.Error("handleIntentResponse() got blank Session while status is success ")
			}

			//if !reflect.DeepEqual(gotResp, tt.wantResp) {
			//		t.Errorf("handleIntentResponse() gotResp = %v, want %v", gotResp, tt.wantResp)
			//}

			if reflect.TypeOf(gotResp) != tt.wantRespType {
				t.Errorf("handleIntentResponse() got resp type = %v, wanted of type %v", reflect.TypeOf(gotResp), tt.wantRespType)
			}

		})
	}
}

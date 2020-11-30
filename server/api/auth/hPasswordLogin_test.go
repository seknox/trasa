package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/misc"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginHandler(t *testing.T) {
	_ = users.InitStoreMock()
	redis.InitStoreMock()
	authstore := InitStoreMock()
	systemstore := system.InitStoreMock()
	orgstore := orgs.InitStoreMock()
	_ = logs.InitStoreMock()
	_ = misc.InitMock()

	systemstore.On("GetGlobalSetting", "abc", consts.GLOBAL_PASSWORD_CONFIG).Return(models.GlobalSettings{
		SettingID:    "123213123",
		OrgID:        "abc",
		Status:       false,
		SettingType:  consts.GLOBAL_PASSWORD_CONFIG,
		SettingValue: "{}",
		UpdatedBy:    "",
		UpdatedOn:    0,
	}, nil)
	systemstore.On("GetGlobalSetting", "abc", consts.GLOBAL_DEVICE_HYGIENE_CHECK).Return(models.GlobalSettings{
		SettingID:    "123213123",
		OrgID:        "abc",
		Status:       false,
		SettingType:  consts.GLOBAL_DEVICE_HYGIENE_CHECK,
		SettingValue: "{}",
		UpdatedBy:    "",
		UpdatedOn:    0,
	}, nil)

	orgstore.On("Get", "abc").Return(models.Org{
		ID:       "abc",
		OrgName:  "someOrg",
		Timezone: "Asia/Kathmandu",
	}, nil)

	pass, _ := bcrypt.GenerateFromPassword([]byte("testpass@123"), bcrypt.DefaultCost)
	authstore.
		On("GetLoginDetails", "user@example.com", "").
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

	//userstore.On()

	handler := http.HandlerFunc(LoginHandler)

	tests := []struct {
		arg          LoginRequest
		wantedStatus string
		wantedIntent string
		wantedData   []interface{}
	}{
		{
			arg: LoginRequest{
				Email:    "user@example.com",
				Password: "testpass@123",
			},
			wantedStatus: "success",
			wantedIntent: consts.AUTH_RESP_TFA_REQUIRED,
			wantedData:   nil,
		},

		{
			arg: LoginRequest{
				Email:    "user@example.com",
				Password: "wrongpass",
			},
			wantedStatus: "failed",
			wantedIntent: "Dashboard Login",
			wantedData:   nil,
		},
	}

	for _, test := range tests {

		req, err := http.NewRequest("GET", "/api/v1/login", bytes.NewBuffer(utils.MarshallStructByte(test.arg)))
		if err != nil {
			t.Fatal(err)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)
		//fmt.Println(string(rr.Body.Bytes()))

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		resp, err := utils.ParseTrasaResponse(rr.Body.Bytes())
		if err != nil {
			t.Fatalf("handler returned non nil error : got %v want %v",
				err, nil)
		}

		if resp.Status != test.wantedStatus {
			t.Errorf("handler returned wrong status : got %v want %v",
				resp.Status, test.wantedStatus)
		}

		if resp.Intent != test.wantedIntent {
			t.Errorf("handler returned wrong intent : got %v want %v",
				resp.Intent, test.wantedIntent)
		}

	}

	//// Check the response body is what we expect.
	//expected := `{"alive": true}`
	//if rr.Body.String() != expected {
	//	t.Errorf("handler returned unexpected body: got %v want %v",
	//		rr.Body.String(), expected)
	//}
}

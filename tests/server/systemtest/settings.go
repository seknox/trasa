package systemtest

import (
	"encoding/json"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func UpdateSettings(t *testing.T) {

	dynamicAccess := models.GlobalDynamicAccessSettings{
		Status:     false,
		PolicyID:   "somePolicyID",
		UserGroups: []string{`somegroupID`, `somegroupID2`},
	}

	pPolicy := models.PasswordPolicy{
		Expiry:            "2021-03-05",
		MinimumChars:      9,
		EnforceStrongPass: false,
		ZxcvbnScore:       4,
	}

	devHygReq := struct {
		EnableDeviceHygieneCheck bool `json:"enableDeviceHygieneCheck"`
	}{
		false,
	}

	emailSett := models.EmailIntegrationConfig{
		IntegrationType: "smtp",
		AuthKey:         "trasa.io",
		AuthPass:        "somepass",
		ServerAddress:   "127.0.0.1",
		ServerPort:      "4567",
		SenderAddress:   "noreply@trasa.io",
	}

	type args struct {
		reqData interface{}
		handler http.HandlerFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "update dynamic access",
			args:    args{dynamicAccess, system.UpdateDynamicAccessSetting},
			wantErr: false,
		},
		{
			name:    "update password policy",
			args:    args{pPolicy, system.UpdatePasswordPolicy},
			wantErr: false,
		},
		{
			name:    "update device hygiene Check",
			args:    args{devHygReq, system.UpdateDeviceHygieneSetting},
			wantErr: false,
		},
		{
			name:    "update email setting",
			args:    args{emailSett, system.UpdateEmailSetting},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := testutils.GetReqWithBody(t, tt.args.reqData)

			testutils.AddTestUserContext(tt.args.handler).ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			var resp struct {
				models.TrasaResponseStruct
			}

			err := json.Unmarshal(rr.Body.Bytes(), &resp)
			if err != nil {
				t.Fatal(err)
			}

			if tt.wantErr != (resp.Status != "success") {
				t.Fatalf(`  wantErr=%t gotStatus=%s  gotReason=%s`, tt.wantErr, resp.Status, resp.Reason)
			}

		})
	}

	gotResp := getSystemSettings(t)

	if gotResp.DynamicAccess.Status != dynamicAccess.Status {
		t.Errorf(`getSystemSettings.DynamicAccess got=%v want=%v`, gotResp.DynamicAccess.Status, dynamicAccess.Status)
	}

	//TODO add more expectations

}

func getSystemSettings(t *testing.T) system.GlobalSettingsResp {
	req := testutils.GetReqWithBody(t, nil)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(system.GlobalSettings))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []system.GlobalSettingsResp `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

	if len(resp.Data) == 0 {
		t.Fatalf(`response data is blank, resp: %s`, string(rr.Body.Bytes()))
	}
	data := resp.Data[0]

	return data
}

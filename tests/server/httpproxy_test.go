package server_test

import (
	"encoding/json"
	"github.com/seknox/trasa/server/api/auth/serviceauth"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHTTPAccessProxy(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name        string
		args        args
		wantSuccess bool
	}{
		{
			"should fail when hostname is incorrect",
			args{testutils.GetReqWithBody(t, serviceauth.NewSession{
				HostName:  "gitlab01.trasa.io",
				TfaMethod: "totp",
				TotpCode:  testutils.GetTotpCode(testutils.MocktotpSEC),
				ExtToken:  "cb6dd3f6-54c2-4cb0-b294-e22c2aa708e4",
			})},
			false,
		},

		{
			"should fail when ext token is incorrect",
			args{testutils.GetReqWithBody(t, serviceauth.NewSession{
				HostName:  "gitlab01.trasa.io",
				TfaMethod: "totp",
				TotpCode:  testutils.GetTotpCode(testutils.MocktotpSEC),
				ExtToken:  "db6dd3f6-54c2-4cb0-b294-e22c2aa708e4",
			})},
			false,
		},

		{
			"should fail when totp is incorrect",
			args{testutils.GetReqWithBody(t, serviceauth.NewSession{
				HostName:  "gitlab01.trasa.io",
				TfaMethod: "totp",
				TotpCode:  "123456",
				ExtToken:  "cb6dd3f6-54c2-4cb0-b294-e22c2aa708e4",
			})},
			false,
		},

		{
			"should fail if service is not authorised",
			args{testutils.GetReqWithBody(t, serviceauth.NewSession{
				HostName:  "test00.trasa.io",
				TfaMethod: "totp",
				TotpCode:  testutils.GetTotpCode(testutils.MocktotpSEC),
				ExtToken:  "cb6dd3f6-54c2-4cb0-b294-e22c2aa708e4",
			})},
			false,
		},

		{
			"should pass",
			args{testutils.GetReqWithBody(t, serviceauth.NewSession{
				HostName:  "gitlab01.trasa.io",
				TfaMethod: "totp",
				TotpCode:  testutils.GetTotpCode(testutils.MocktotpSEC),
				ExtToken:  "cb6dd3f6-54c2-4cb0-b294-e22c2aa708e4",
			})},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(serviceauth.AuthHTTPAccessProxy)

			handler.ServeHTTP(rr, tt.args.req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			var resp models.TrasaResponseStruct
			err := json.Unmarshal(rr.Body.Bytes(), &resp)
			if err != nil {
				t.Fatal(err)
			}

			if tt.wantSuccess && resp.Status != "success" {
				t.Errorf("AuthHTTPAccessProxy() wanted success, got:%s reason %s", resp.Status, resp.Reason)
				return
			}

		})
	}

}

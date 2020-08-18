package integration_tests

import (
	"bytes"
	"encoding/json"
	"github.com/seknox/trasa/server/api/auth/serviceauth"
	"github.com/seknox/trasa/server/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAgentAuth(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name        string
		args        args
		wantSuccess bool
	}{
		{
			"should fail when serviceID/service Key is incorrect",
			args{getreqWithBody(t, serviceauth.ServiceAgentLogin{
				ServiceID:  "2fef188a-cc12-438b-8564-2803a072f650",
				ServiceKey: "sasd76asd67asd67asgd7asnskadasd",
				User:       "admin",
				TfaMethod:  "",
				TotpCode:   "",
				UserIP:     "",
				TrasaID:    "",
			})},
			false,
		},

		{
			"should fail when trasaID is incorrect",
			args{getreqWithBody(t, serviceauth.ServiceAgentLogin{
				ServiceID:  "2fef188a-cc13-438b-8564-2803a072f650",
				ServiceKey: "d9ef5359f13f6f6f6c89b4a9be9958ed13",
				User:       "admin",
				TfaMethod:  "",
				TotpCode:   "",
				UserIP:     "",
				TrasaID:    "incorrect",
			})},
			false,
		},

		{
			"should fail when privilege is incorrect",
			args{getreqWithBody(t, serviceauth.ServiceAgentLogin{
				ServiceID:  "2fef188a-cc13-438b-8564-2803a072f650",
				ServiceKey: "d9ef5359f13f6f6f6c89b4a9be9958ed13",
				User:       "admin",
				TfaMethod:  "",
				TotpCode:   "",
				UserIP:     "",
				TrasaID:    "root",
			})},
			false,
		},
		{
			"should fail when totp code is incorrect",
			args{getreqWithBody(t, serviceauth.ServiceAgentLogin{
				ServiceID:  "2fef188a-cc13-438b-8564-2803a072f650",
				ServiceKey: "d9ef5359f13f6f6f6c89b4a9be9958ed13",
				User:       "bhrg3se",
				TotpCode:   "1323214",
				UserIP:     "",
				TrasaID:    "root",
			})},
			false,
		},

		{
			"should fail when ip is invalid",
			args{getreqWithBody(t, serviceauth.ServiceAgentLogin{
				ServiceID:  "2fef188a-cc13-438b-8564-2803a072f650",
				ServiceKey: "d9ef5359f13f6f6f6c89b4a9be9958ed13",
				User:       "bhrg3se",
				TotpCode:   getTotpCode(totpSEC),
				UserIP:     "",
				TrasaID:    "root",
			})},
			false,
		},
		{
			"should fail when adhoc is enabled",
			args{getreqWithBody(t, serviceauth.ServiceAgentLogin{
				ServiceID:  "08d97469-4a2f-46d3-86bc-3005b4c99c6c",
				ServiceKey: "2a094b1fba624b26eaa02f5e2b9f5755ea",
				User:       "sakshyam",
				TotpCode:   getTotpCode(totpSEC),
				UserIP:     "",
				TrasaID:    "root",
			})},
			false,
		},
		{
			"should pass",
			args{getreqWithBody(t, serviceauth.ServiceAgentLogin{
				ServiceID:  "2fef188a-cc13-438b-8564-2803a072f650",
				ServiceKey: "d9ef5359f13f6f6f6c89b4a9be9958ed13",
				User:       "bhrg3se",
				TotpCode:   getTotpCode(totpSEC),
				UserIP:     "127.0.0.1",
				TrasaID:    "root",
			})},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(serviceauth.AgentLogin)

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
				t.Errorf("AgentLogin() wanted success, got:%s reason %s", resp.Status, resp.Reason)
				return
			}

		})
	}

}

func getreqWithBody(t *testing.T, body interface{}) *http.Request {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	return req
}

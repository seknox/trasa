package system

import (
	"encoding/json"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func UpdateSecurityRules(t *testing.T) {
	//UpdateSecurityRule
	reqData := system.UpdateSecurityRulesReq{
		Status: "enabled",
		RuleID: "someRuleID",
	}
	req := testutils.GetReqWithBody(t, reqData)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(system.UpdateSecurityRule))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		//Data []system.SysStatus `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

}

func GetSecurityRules(t *testing.T) {

	req := testutils.GetReqWithBody(t, nil)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(system.SecurityRules))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]models.SecurityRule `json:"data"`
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

	found := false
	for _, rule := range resp.Data[0] {
		//t.Log(rule.RuleID,rule.ConstName)
		if rule.RuleID == "someRuleID" {
			if !rule.Status {
				t.Error(`sec rule expected to be enabled`)
			}
			found = true
		}
	}

	if !found {
		t.Error("someRuleID not found")
	}
}

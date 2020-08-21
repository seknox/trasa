package crudtest

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/policies"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func CreatePolicy(t *testing.T) models.Policy {
	creatPolReq := models.Policy{
		PolicyName: "some-policy",
		DayAndTime: []models.DayAndTimePolicy{{
			Days:     []string{`Sunday`, "Monday", "Tuesday", "Saturday"},
			FromTime: "01:00",
			ToTime:   "22:00",
		}},
		TfaRequired:   false,
		RecordSession: false,
		FileTransfer:  false,
		IPSource:      "0.0.0.0/0",
		DevicePolicy:  models.DevicePolicy{},
		Expiry:        "2050-01-01",
		IsExpired:     false,
	}

	req := testutils.GetReqWithBody(t, creatPolReq)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(policies.CreatePolicy))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []models.Policy `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

	if len(resp.Data) == 0 {
		t.Fatal(resp.Reason)
	}

	return resp.Data[0]

}

func UpdatePolicy(t *testing.T, policyID string) {
	creatPolReq := models.Policy{
		PolicyID:   policyID,
		PolicyName: "some-full-policy",
		DayAndTime: []models.DayAndTimePolicy{{
			Days:     []string{`Sunday`, "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			FromTime: "01:00",
			ToTime:   "22:00",
		}},
		TfaRequired:   true,
		RecordSession: true,
		FileTransfer:  true,
		IPSource:      "0.0.0.0/0",
		DevicePolicy:  models.DevicePolicy{},
		Expiry:        "2050-01-01",
		IsExpired:     false,
	}
	req := testutils.GetReqWithBody(t, creatPolReq)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(policies.UpdatePolicy))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

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

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

}

func GetPolicy(t *testing.T, expected models.Policy) {

	req := testutils.GetReqWithBody(t, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("policyID", expected.PolicyID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	t.Log(expected.PolicyID)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(policies.GetPolicy))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []models.Policy `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}
	if len(resp.Data) == 0 {
		t.Fatal(resp.Reason)
	}

	if !reflect.DeepEqual(resp.Data[0], expected) {
		t.Errorf("newly created policy  \n got = %v,\n want = %v", resp.Data[0], expected)
	}

}

func GetPolicies(t *testing.T, expected models.Policy) {

	req := testutils.GetReqWithBody(t, nil)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(policies.GetPolicies))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]models.Policy `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}
	if len(resp.Data) == 0 {
		t.Fatal(resp.Reason)
	}

	exists := false
	for _, s := range resp.Data[0] {
		if s.PolicyID == expected.PolicyID {
			exists = true
			if !reflect.DeepEqual(s, expected) {
				t.Errorf("InAppNotif  \n got = %v,\n want = %v", s, expected)
			}
		}
	}

	if !exists {
		t.Errorf(`newly created service not got in services.GetAllServices`)
	}

}

func DeletePolicy(t *testing.T, policyID string) {

	deleteServReq := struct {
		PolicyID []string `json:"policyID"`
	}{PolicyID: []string{policyID}}

	req := testutils.GetReqWithBody(t, deleteServReq)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(policies.DeletePolicies))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

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

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

}

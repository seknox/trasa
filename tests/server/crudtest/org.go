package crudtest

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func UpdateOrg(t *testing.T) models.Org {
	reqdata := models.Org{
		OrgName:        "differentOrgName",
		Domain:         "different.domain",
		PrimaryContact: "primary.contact.com",
		Timezone:       "UTC",
		PhoneNumber:    "9999999999",
	}
	req := testutils.GetReqWithBody(t, reqdata)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(orgs.Update))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []models.Org `json:"data"`
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

func GetOrg(t *testing.T, expected models.Org) {

	req := testutils.GetReqWithBody(t, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", expected.ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(orgs.Get))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []models.Org `json:"data"`
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

	resp.Data[0].CreatedAt = 0
	if !reflect.DeepEqual(resp.Data[0], expected) {
		t.Errorf("GetUserDetails  \n got = %v,\n want = %v", resp.Data[0], expected)
	}

}

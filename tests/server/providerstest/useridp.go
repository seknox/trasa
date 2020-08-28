package providerstest

import (
	"encoding/json"
	"github.com/seknox/trasa/server/api/providers/uidp"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func CreateIdp(t *testing.T) {
	//id, org_id, name,type, meta, is_enabled, redirect_url, audience_uri,client_id, endpoint, created_by , integration_type,scim_endpoint, last_updated )
	idp := models.IdentityProvider{
		IdpName:      "freeipa",
		IdpType:      "freeipa",
		IsEnabled:    true,
		ClientID:     "someClientID",
		ClientSecret: "clientSec",
		SCIMEndpoint: "someEndpoint",
		ApiKey:       "someAPIKey",
	}
	req := testutils.GetReqWithBody(t, idp)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(uidp.CreateIdp))

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

func UpdateIdp(t *testing.T) {
	//id, org_id, name,type, meta, is_enabled, redirect_url, audience_uri,client_id, endpoint, created_by , integration_type,scim_endpoint, last_updated )
	idp := models.IdentityProvider{
		IdpName:      "ldap",
		IdpType:      "ldap",
		IsEnabled:    true,
		ClientID:     "someClientID",
		ClientSecret: "someClientSec",
		SCIMEndpoint: "someEndpoint",
		ApiKey:       "someAPIKey",
	}
	req := testutils.GetReqWithBody(t, idp)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(uidp.UpdateIdp))

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

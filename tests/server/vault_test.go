package server_test

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVault(t *testing.T) {
	initVault(t)
	getStatus(t)
	testStoreVault(t)
	testGetKey(t)
	storeSecret(t)
	getSecret(t)
	getUpstreamCredsTest(t)
	deleteSecret(t)

}

func initVault(t *testing.T) {

	req := testutils.GetReqWithBody(t, system.VaultInit{
		SecretShares:    5,
		SecretThreshold: 3,
	})

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(system.TsxvaultInit))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []system.VaultInitResp `json:"data"`
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

	if len(data.UnsealKeys) != 5 {
		t.Fatalf(`len DecryptKeys got:%d want:%d`, len(data.UnsealKeys), 5)
	}

}

func testStoreVault(t *testing.T) {
	reqdata := models.KeysHolderReq{
		KeyTag:  "someKeyTag",
		KeyName: "someKeyName",
		KeyVal:  "someKeyValue",
	}

	req := testutils.GetReqWithBody(t, reqdata)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(system.StoreKey))

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

func testGetKey(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(system.Getkey))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("vendorID", "someKeyName")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []models.KeysHolderReq `json:"data"`
	}

	err = json.Unmarshal(rr.Body.Bytes(), &resp)
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

	//if(data.OrgID!="153f7582-5ae2-46ba-8c1c-79ef73fe296e"){
	//	t.Errorf(`data.OrgID got:%s want:%s`,data.OrgID,"153f7582-5ae2-46ba-8c1c-79ef73fe296e")
	//
	//}

	if data.AddedBy != "13c45cfb-72ca-4177-b968-03604cab6a27" {
		t.Errorf(`data.AddedBy got:%s want:%s`, data.AddedBy, "13c45cfb-72ca-4177-b968-03604cab6a27")
	}

	if data.KeyTag != "some-xxxx-xxxx..." {
		t.Errorf(`data.KeyTag got:%s want:%s`, data.KeyTag, "some-xxxx-xxxx...")
	}

	if data.KeyName != "someKeyName" {
		t.Errorf(`data.KeyName got:%s want:%s`, data.KeyName, "someKeyName")
	}

	if data.KeyVal != "some-xxxx-xxxx..." {
		t.Errorf(`data.KeyVal got:%s want:%s`, data.KeyVal, "some-xxxx-xxxx...")
	}

}

func getStatus(t *testing.T) {

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(system.Status))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []system.VaultStatus `json:"data"`
	}

	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}
	if len(resp.Data) == 0 {
		t.Fatal(resp.Reason)
	}

	data := resp.Data[0]

	if !data.InitStatus.Status {
		t.Fatalf(`vault not initialised`)
	}

}

func storeSecret(t *testing.T) {
	reqdata := services.ServiceCreds{
		Username:   "root",
		Credential: "somePass",
		ServiceID:  "4ea851b8-6299-4c61-8137-58771aaa8899",
		Type:       "password",
	}

	req := testutils.GetReqWithBody(t, reqdata)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(services.StoreServiceCredentials))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []string `json:"data"`
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

	if resp.Data[0] != "root" {
		t.Errorf(`username got=%s want %s`, resp.Data[0], "root")
	}

}

func getSecret(t *testing.T) {

	reqdata := services.ServiceCreds{
		Username:  "root",
		ServiceID: "4ea851b8-6299-4c61-8137-58771aaa8899",
		Type:      "password",
	}

	req := testutils.GetReqWithBody(t, reqdata)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(services.ViewCreds))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []services.ServiceCreds `json:"data"`
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

	data := resp.Data[0]

	if data.ServiceID != "4ea851b8-6299-4c61-8137-58771aaa8899" {
		t.Errorf(`data.ServiceID got:%s want:%s`, data.ServiceID, "4ea851b8-6299-4c61-8137-58771aaa8899")
	}
	if data.Credential != "somePass" {
		t.Errorf(`data.Credential got:%s want:%s`, data.Credential, "somePass")
	}
	if data.Username != "root" {
		t.Errorf(`data.Username got:%s want:%s`, data.Username, "root")
	}
	if data.Type != "password" {
		t.Errorf(`data.Type got:%s want:%s`, data.Type, "password")
	}

}

func deleteSecret(t *testing.T) {

	reqdata := services.ServiceCreds{
		Username:  "root",
		ServiceID: "4ea851b8-6299-4c61-8137-58771aaa8899",
		Type:      "password",
	}

	req := testutils.GetReqWithBody(t, reqdata)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(services.DeleteCreds))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []services.ServiceCreds `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

}

func getUpstreamCredsTest(t *testing.T) {

	creds, err := services.GetUpstreamCreds("root", "4ea851b8-6299-4c61-8137-58771aaa8899", "ssh", "153f7582-5ae2-46ba-8c1c-79ef73fe296e")
	if err != nil {
		t.Fatal(err)
	}

	if creds.Password != "somePass" {
		t.Errorf(`creds.Password got:%s want:%s`, creds.Password, "somePass")
	}
}

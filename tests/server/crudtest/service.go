package crudtest

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func CreateService(t *testing.T) {
	creatServReq := models.Service{
		Name:     "some Service",
		Hostname: "some.host",
		Type:     "ssh",
	}
	req := testutils.GetReqWithBody(t, creatServReq)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(services.CreateService))

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

func UpdateService(t *testing.T, serviceID string) {
	creatServReq := models.Service{
		ID:       serviceID,
		Name:     "someService",
		Hostname: "some.host",
		Type:     "rdp",
	}
	req := testutils.GetReqWithBody(t, creatServReq)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(services.UpdateService))

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

func UpdateHTTPConfig(t *testing.T, serviceID string) {
	creatServReq := services.ReverseProxyReq{
		ServiceID: serviceID,
		Name:      "someService",
		Proxy: models.ReverseProxy{
			RouteRule:           "someRouteRule",
			PassHostheader:      false,
			UpstreamServer:      "http://grilab.com",
			StrictTLSValidation: false,
		},
	}
	req := testutils.GetReqWithBody(t, creatServReq)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(services.UpdateHTTPProxy))

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

func UpdateHostCerts(t *testing.T, serviceID, certVal string) {
	creatServReq := services.UpdateHostCertsReq{
		CertVal:   certVal,
		ServiceID: serviceID,
	}
	req := testutils.GetReqWithBody(t, creatServReq)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(services.UpdateHostCerts))

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

func UpdateSSLCerts(t *testing.T, serviceID, certVal string) {
	creatServReq := services.UpdateSSLCertsReq{
		SslKey:  "5ard65safdtyasudbjkansd",
		SslCert: "asdfasdg7tasgdayusdasd",
		CaCert:  "6asfd6asfd65asfd7atsfd6asdguybasjbdjas",
	}
	req := testutils.GetReqWithBody(t, creatServReq)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("serviceID", serviceID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(services.UpdateSSLCerts))

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

func DownloadHostCerts(t *testing.T, serviceID string) {

	req := testutils.GetReqWithBody(t, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("serviceID", serviceID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(services.DownloadHostCerts))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if len(rr.Body.Bytes()) == 0 {
		t.Error("DownloadHostCerts returned empty")
	}

}

func GetAllServices(t *testing.T) (serviceID string) {

	req := testutils.GetReqWithBody(t, nil)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(services.GetAllServices))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []services.AllServicesByType `json:"data"`
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
	for _, s := range resp.Data[0].SSH {
		if s.Name == "some service" && s.Hostname == "some.host" {
			exists = true
			serviceID = s.ID
		}
	}

	if !exists {
		t.Errorf(`newly created service not got in services.GetAllServices`)
	}

	return serviceID

}

func GetService(t *testing.T, serviceID string) {

	req := testutils.GetReqWithBody(t, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("serviceID", serviceID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(services.GetServiceDetail))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []models.Service `json:"data"`
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

	if resp.Data[0].ID != serviceID {
		t.Errorf(`services.GetServiceDetail got=%v want=%v`, resp.Data[0].ID, serviceID)
	}

}

func DeleteService(t *testing.T, sID string) {

	deleteServReq := models.Service{
		ID: sID,
	}

	req := testutils.GetReqWithBody(t, deleteServReq)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(services.DeleteService))

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

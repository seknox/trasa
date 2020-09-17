package providerstest

import (
	"context"
	"encoding/json"
	"github.com/cloudflare/cfssl/csr"
	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/providers/ca"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func CreateHTTPCA(t *testing.T) {
	reqData := new(csr.CertificateRequest)
	reqData.CN = "example.com"

	req := testutils.GetReqWithBody(t, reqData)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(ca.InitCA))

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

func CreateSSHCA(t *testing.T) {
	req := testutils.GetReqWithBody(t, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "user")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(ca.InitSSHCA))

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

func GetAllCAs(t *testing.T) {
	req := testutils.GetReqWithBody(t, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "user")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(ca.GetAllCAs))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]ca.CertHolderResponse `json:"data"`
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

	sshCAExists := false
	httpcaExists := false

	for _, c := range data {
		if c.OrgID != testutils.MockOrgID {
			t.Errorf("GetAllCAs[].OrgID: got %v want %v",
				c.OrgID, testutils.MockOrgID)
		}

		if c.CertType == consts.CERT_TYPE_SSH_CA {
			sshCAExists = true
		}

		if c.CertType == consts.CERT_TYPE_HTTP_CA {
			httpcaExists = true
		}

	}

	if !sshCAExists {
		t.Error("GetAllCAs[]: ssh  CA not found")
	}

	if !httpcaExists {
		t.Error("GetAllCAs[]: http  CA not found")
	}

}

func GetHttpCADetail(t *testing.T) {
	req := testutils.GetReqWithBody(t, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "user")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(ca.GetHttpCADetail))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []ca.CertHolderResponse `json:"data"`
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

	if data.OrgID != testutils.MockOrgID {
		t.Errorf("GetAllCAs[].OrgID: got %v want %v",
			data.OrgID, testutils.MockOrgID)
	}

	if data.CertType != consts.CERT_TYPE_HTTP_CA {
		t.Errorf("GetAllCAs[].CertType: got %v want %v",
			data.CertType, consts.CERT_TYPE_HTTP_CA)
	}

}

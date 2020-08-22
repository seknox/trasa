package server_test

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/my"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthLogs(t *testing.T) {
	//entitytype

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	///stats/ips/{entitytype}/{entityid}/{timeFilter}/{statusFilter}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("entitytype", "org")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(logs.GetLoginEvents))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]logs.AuthLog `json:"data"`
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

	if len(data) != 100 {
		t.Fatalf(`incorrect number of logs: got=%d want=%d`, len(data), 100)
	}

	for _, al := range data {
		if al.OrgID != "" && al.OrgID != "153f7582-5ae2-46ba-8c1c-79ef73fe296e" {
			t.Errorf(`incorrect orgID: got=%s want=%s`, al.OrgID, "153f7582-5ae2-46ba-8c1c-79ef73fe296e")
		}
	}

}

func TestMyAuthLogs(t *testing.T) {
	//entitytype

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(my.GetMyEvents))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]logs.AuthLog `json:"data"`
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

	if len(data) != 100 {
		t.Fatalf(`incorrect number of logs: got=%d want=%d`, len(data), 100)
	}

	for _, al := range data {
		if al.OrgID != "" && al.OrgID != "153f7582-5ae2-46ba-8c1c-79ef73fe296e" {
			t.Errorf(`incorrect orgID: got=%s want=%s`, al.OrgID, "153f7582-5ae2-46ba-8c1c-79ef73fe296e")
		}

		if al.UserID != "" && al.UserID != testutils.MockUserID {
			t.Errorf(`incorrect userID: got=%s want=%s`, al.UserID, testutils.MockUserID)

		}
	}

}

func TestAuthLogsByPage(t *testing.T) {
	type args struct {
		entityType string
		entityID   string
		dateFrom   string
		dateTo     string
		size       string
		page       string
	}
	tests := []struct {
		name     string
		args     args
		wantSize int
	}{
		{
			name:     "get logs with any params",
			args:     args{"org", "org", "2006-01-02", "2222-02-02", "100", "2"},
			wantSize: 100,
		},

		{
			name:     "get logs with any params",
			args:     args{"org", "org", "2006-01-02", "2222-02-02", "100", "1"},
			wantSize: 100,
		},

		{
			name:     "get logs with any params",
			args:     args{"org", "org", "2006-01-02", "2222-02-02", "200", "1"},
			wantSize: 200,
		},

		//{
		//	name:     "get logs with any params",
		//	args:     args{"org", "org", "2006-01-02", "2222-02-02", "500", "1"},
		//	wantSize: 236,
		//},

		//TODO add entity specific tests

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testAuthLog(t, tt.args.entityType, tt.args.entityID, tt.args.dateFrom, tt.args.dateTo, tt.args.page, tt.args.size)

			if len(got) != tt.wantSize {
				t.Errorf("GetLoginEventsByPage() = %v, want %v", len(got), tt.wantSize)
			}
			data := got
			for _, al := range data {
				if al.OrgID != "" && al.OrgID != "153f7582-5ae2-46ba-8c1c-79ef73fe296e" {
					t.Errorf(`incorrect orgID: got=%s want=%s`, al.OrgID, "153f7582-5ae2-46ba-8c1c-79ef73fe296e")
				}
			}

		})
	}

}

func testAuthLog(t *testing.T, entityType, entityID, from, to, page, size string) []logs.AuthLog {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	///stats/ips/{entitytype}/{entityid}/{timeFilter}/{statusFilter}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("entitytype", entityType)
	rctx.URLParams.Add("entityid", entityID)
	rctx.URLParams.Add("dateFrom", from)
	rctx.URLParams.Add("dateTo", to)
	rctx.URLParams.Add("page", page)
	rctx.URLParams.Add("size", size)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(logs.GetLoginEventsByPage))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]logs.AuthLog `json:"data"`
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

	return data

}

func TestMyAuthLogsByPage(t *testing.T) {
	type args struct {
		dateFrom string
		dateTo   string
		size     string
		page     string
	}
	tests := []struct {
		name     string
		args     args
		wantSize int
	}{
		{
			name:     "get logs with any params",
			args:     args{"2006-01-02", "2222-02-02", "100", "2"},
			wantSize: 100,
		},

		{
			name:     "get logs with any params",
			args:     args{"2006-01-02", "2222-02-02", "100", "1"},
			wantSize: 100,
		},

		{
			name:     "get logs with any params",
			args:     args{"2006-01-02", "2222-02-02", "200", "1"},
			wantSize: 200,
		},

		//{
		//	name:     "get logs with any params",
		//	args:     args{"2006-01-02", "2222-02-02", "500", "1"},
		//	wantSize: 219,
		//},

		//TODO add entity specific tests

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testMyAuthLog(t, tt.args.dateFrom, tt.args.dateTo, tt.args.page, tt.args.size)

			if len(got) != tt.wantSize {
				t.Errorf("GetMyEventsByPage() = %v, want %v", len(got), tt.wantSize)
			}
			data := got
			for _, al := range data {
				if al.OrgID != "" && al.OrgID != "153f7582-5ae2-46ba-8c1c-79ef73fe296e" {
					t.Errorf(`incorrect orgID: got=%s want=%s`, al.OrgID, "153f7582-5ae2-46ba-8c1c-79ef73fe296e")
				}
				if al.UserID != "" && al.UserID != testutils.MockUserID {
					t.Errorf(`incorrect userID: got=%s want=%s`, al.UserID, testutils.MockUserID)

				}
			}

		})
	}

}

func testMyAuthLog(t *testing.T, from, to, page, size string) []logs.AuthLog {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	///stats/ips/{entitytype}/{entityid}/{timeFilter}/{statusFilter}
	rctx := chi.NewRouteContext()

	rctx.URLParams.Add("dateFrom", from)
	rctx.URLParams.Add("dateTo", to)
	rctx.URLParams.Add("page", page)
	rctx.URLParams.Add("size", size)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(my.GetMyEventsByPage))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]logs.AuthLog `json:"data"`
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

	return data

}

func TestInAppTrailLogs(t *testing.T) {
	type args struct {
		dateFrom string
		dateTo   string
		size     string
		page     string
	}
	tests := []struct {
		name     string
		args     args
		wantSize int
	}{
		{
			name:     "get logs with size limit",
			args:     args{"2006-01-02", "2222-02-02", "10", "2"},
			wantSize: 10,
		},

		{
			name:     "get logs with lize limit",
			args:     args{"2006-01-02", "2222-02-02", "10", "1"},
			wantSize: 10,
		},

		{
			name:     "get logs with over size",
			args:     args{"2006-01-02", "2222-02-02", "200", "1"},
			wantSize: 68,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testInAppTrail(t, tt.args.dateFrom, tt.args.dateTo, tt.args.page, tt.args.size)

			if len(got) != tt.wantSize {
				t.Errorf("GetLoginEventsByPage() = %v, want %v", len(got), tt.wantSize)
			}
			data := got
			for _, al := range data {
				if al.OrgID != "" && al.OrgID != "153f7582-5ae2-46ba-8c1c-79ef73fe296e" {
					t.Errorf(`incorrect orgID: got=%s want=%s`, al.OrgID, "153f7582-5ae2-46ba-8c1c-79ef73fe296e")
				}
			}

		})
	}

}

func testInAppTrail(t *testing.T, from, to, page, size string) []models.InAppTrail {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	///stats/ips/{entitytype}/{entityid}/{timeFilter}/{statusFilter}
	rctx := chi.NewRouteContext()

	rctx.URLParams.Add("dateFrom", from)
	rctx.URLParams.Add("dateTo", to)
	rctx.URLParams.Add("page", page)
	rctx.URLParams.Add("size", size)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(logs.GetAllInAppTrails))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]models.InAppTrail `json:"data"`
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

	return data

}

package server_test

import (
	"context"
	"encoding/json"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/stats"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAggregatedUsers(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(stats.GetAggregatedUsers))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []stats.TotalUsers `json:"data"`
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

	if data.Users != 3 {
		t.Errorf(`incorrect number of users, expected:%d got %d`, 3, data.Users)
	}
	if data.Admins != 1 {
		t.Errorf(`incorrect number of admins, expected:%d got %d`, 1, data.Admins)
	}
	if data.DisabledUsers != 0 {
		t.Errorf(`incorrect number of disabled users, expected:%d got %d`, 0, data.DisabledUsers)
	}

}

// Our handlers satisfy http.Handler, so we can call their ServeHTTP method

func TestGetAggregatedDevices(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(stats.GetAggregatedDevices))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []stats.AllUserDevices `json:"data"`
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

	if data.TotalBrowsers < 3 {
		t.Errorf(`incorrect number of browsers, expected:%d got %d`, 3, data.TotalBrowsers)
	}
	if len(data.BrowserByType) < 1 {
		t.Errorf(`incorrect number of browser types, expected:%d got %d`, 1, len(data.BrowserByType))
	}

	if data.TotalMobiles != 2 {
		t.Errorf(`incorrect number of mobiles, expected:%d got %d`, 2, data.TotalMobiles)
	}
	if data.TotalWorkstations < 2 {
		t.Errorf(`incorrect number of workstations, expected:%d got %d`, 2, data.TotalWorkstations)
	}

}

func TestGetAggregatedFailedReasons(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(stats.GetAggregatedFailedReasons))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]stats.FailedReasonsByType `json:"data"`
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

	if len(data) != 7 {
		t.Errorf(`incorrect number of failed reasons, expected:%d got %d`, 7, len(data))
	}

	//TODO add more checks

}

func TestGetAggregatedServices(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(stats.GetAggregatedServices))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []stats.AllServices `json:"data"`
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

	if data.DynamicService != false {
		t.Errorf(`incorrect dynamic service status, expected:%t got %t`, false, data.DynamicService)
	}
	//if data.SessionRecordingDisabledCount != 1 {
	//	t.Errorf(`incorrect SessionRecordingDisabledCount, expected:%d got %d`, 1, data.SessionRecordingDisabledCount)
	//}
	if data.TotalGroups != 2 {
		t.Errorf(`incorrect TotalGroups, expected:%d got %d`, 2, data.TotalGroups)
	}
	if data.TotalServices < 3 {
		t.Errorf(`incorrect TotalServices, expected:%d got %d`, 3, data.TotalServices)
	}

}

func TestGetAggregatedLoginHours(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(stats.GetAggregatedLoginHours))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]stats.LoginsByHour `json:"data"`
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

	t.Log(string(rr.Body.Bytes()))
	if len(data) != 24 {
		t.Errorf(`incorrect login hours length, expected:%d got %d`, 24, len(data))
	}

}

func TestGetTotalManagedUsers(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/api/v1/stats/totalmanagedusers/service/2fef188a-cc13-438b-8564-2803a072f650", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(stats.GetTotalManagedUsers))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("entityid", "2fef188a-cc13-438b-8564-2803a072f650")
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
		Data []int `json:"data"`
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

	if data != 1 {
		t.Errorf(`incorrect managed users, expected:%d got %d`, 1, data)
	}

}

func TestGetServicePermStats(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/api/v1/stats/appperms/2fef188a-cc13-438b-8564-2803a072f650", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(stats.GetServicePermStats))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("serviceID", "2fef188a-cc13-438b-8564-2803a072f650")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []stats.PermStats `json:"data"`
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

	if data.Users != 1 {
		t.Errorf(`incorrect assigned users, expected:%d got %d`, 1, data.Users)
	}
	if data.Groups != 0 {
		t.Errorf(`incorrect assigned groups, expected:%d got %d`, 0, data.Groups)
	}
	if data.Policies != 2 {
		t.Errorf(`incorrect assigned policies, expected:%d got %d`, 2, data.Policies)
	}
	if data.Privileges != 2 {
		t.Errorf(`incorrect assigned priviliges, expected:%d got %d`, 2, data.Privileges)
	}
	if data.Secrets != 1 {
		t.Errorf(`incorrect stored secrets, expected:%d got %d`, 1, data.Secrets)
	}

}

func TestGetIPAggs(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/api/v1/stats/appperms/2fef188a-cc13-438b-8564-2803a072f650", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(stats.GetIPAggs))

	///stats/ips/{entitytype}/{entityid}/{timeFilter}/{statusFilter}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("entityid", "2fef188a-cc13-438b-8564-2803a072f650")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []stats.AggIps `json:"data"`
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

	_ = data
	//TODO add test expectations

	//if data.Value != 233 {
	//	t.Errorf(`incorrect assigned users, expected:%d got %d`, 233, data.Value)
	//}
	//if data.Groups != 0 {
	//	t.Errorf(`incorrect assigned groups, expected:%d got %d`, 0, data.Groups)
	//}
	//if data.Policies != 2 {
	//	t.Errorf(`incorrect assigned policies, expected:%d got %d`, 2, data.Policies)
	//}
	//if data.Privileges != 2 {
	//	t.Errorf(`incorrect assigned priviliges, expected:%d got %d`, 2, data.Privileges)
	//}
	//if data.Secrets != 2 {
	//	t.Errorf(`incorrect stored secrets, expected:%d got %d`, 2, data.Secrets)
	//}

}

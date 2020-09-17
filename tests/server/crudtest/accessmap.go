package crudtest

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/accessmap"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func CreateUserAccessMap(t *testing.T, serviceID, userID, policyID, privilege string) {

	reqData := accessmap.AssignUserToApp{
		ServiceID: serviceID,
		Privilege: privilege,
		Users:     []string{userID},
		PolicyID:  []string{policyID},
	}

	req := testutils.GetReqWithBody(t, reqData)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(accessmap.CreateServiceUserMap))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		//Data [][]string `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

}

func CreateUserGroupAccessMap(t *testing.T, serviceID, userGroupID, policyID, privilege string) {

	reqData := accessmap.ServiceGroupUserGroupMapRequest{
		ServiceGroupID: serviceID,
		MapType:        "service",
		UserGroupID:    []string{userGroupID},
		Privilege:      privilege,
		PolicyID:       []string{policyID},
	}

	req := testutils.GetReqWithBody(t, reqData)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(accessmap.CreateServiceGroupUserGroupMap))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		//Data [][]string `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

}

func CreateUserGroupServiceGroupAccessMap(t *testing.T, serviceGroupID, userGroupID, policyID, privilege string) {

	reqData := accessmap.ServiceGroupUserGroupMapRequest{
		ServiceGroupID: serviceGroupID,
		MapType:        "servicegroup",
		UserGroupID:    []string{userGroupID},
		Privilege:      privilege,
		PolicyID:       []string{policyID},
	}

	req := testutils.GetReqWithBody(t, reqData)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(accessmap.CreateServiceGroupUserGroupMap))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		//Data [][]string `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

}

func GetUserAccessMap(t *testing.T, serviceID, userID, policyID, privilege string) (mapID string) {

	req := testutils.GetReqWithBody(t, nil)

	//serviceID
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("serviceID", serviceID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(accessmap.GetUserAccessMaps))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]models.AccessMapDetail `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

	if len(resp.Status) == 0 {
		t.Fatal(resp.Reason)
	}

	data := resp.Data[0]
	found := false

	for _, am := range data {
		if serviceID == am.ServiceID && userID == am.UserID && policyID == am.Policy.PolicyID && privilege == am.Privilege {
			found = true
			mapID = am.MapID
		}
	}

	if !found {
		t.Error("user access map not found")
	}

	return mapID

}

func GetUserGroupsAssignedToService(t *testing.T, serviceID, userGroupID, policyID, privilege string) (mapID string) {

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
		Data [][]accessmap.UserGroupOfServiceGroup `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

	if len(resp.Status) == 0 {
		t.Fatal(resp.Reason)
	}

	data := resp.Data[0]
	found := false

	for _, am := range data {
		if userGroupID == am.UsergroupID && policyID == am.PolicyID && privilege == am.Privilege {
			found = true
			mapID = am.MapID
		}
	}

	if !found {
		t.Error("user group service group access map not found")
	}

	return mapID

}

func GetUserGroupsAssignedToServiceGroups(t *testing.T, serviceID, userGroupID, policyID, privilege string) (mapID string) {

	req := testutils.GetReqWithBody(t, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupID", serviceID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(accessmap.GetUserGroupsAssignedToServiceGroups))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]accessmap.UserGroupOfServiceGroup `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

	if len(resp.Status) == 0 {
		t.Fatal(resp.Reason)
	}

	data := resp.Data[0]
	found := false

	for _, am := range data {
		if userGroupID == am.UsergroupID && policyID == am.PolicyID && privilege == am.Privilege {
			found = true
			mapID = am.MapID
		}
	}

	if !found {
		t.Error("user group service group access map not found")
	}

	return mapID

}

func GetUserGroupServiceGroupAccessMaps(t *testing.T, serviceGroupID, userGroupID, policyID, privilege string) (mapID string) {

	req := testutils.GetReqWithBody(t, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("serviceGroupID", serviceGroupID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(accessmap.GetUserGroupServiceGroupAccessMaps))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]accessmap.UserGroupOfServiceGroup `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

	if len(resp.Status) == 0 {
		t.Fatal(resp.Reason)
	}

	data := resp.Data[0]
	found := false

	for _, am := range data {
		if userGroupID == am.UsergroupID && policyID == am.PolicyID && privilege == am.Privilege {
			found = true
			mapID = am.MapID
		}
	}

	if !found {
		t.Error("user groups access map not found")
	}

	return mapID
}

func UpdateUserAccessMap(t *testing.T, mapID string) {

	reqData := accessmap.UpdatePrivilege{
		MapID:     mapID,
		Privilege: "newPrivilege",
	}

	req := testutils.GetReqWithBody(t, reqData)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(accessmap.UpdateServiceUserMap))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		//Data [][]string `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

}

func UpdateUserGroupAccessMap(t *testing.T, mapID string) {

	reqData := accessmap.UpdatePrivilege{
		MapID:     mapID,
		Privilege: "newPrivilege",
	}

	req := testutils.GetReqWithBody(t, reqData)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(accessmap.UpdateServiceGroupUserGroup))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		//Data [][]string `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

}

func DeleteUserAccessMap(t *testing.T, mapID string) {

	reqData := accessmap.DeleteServiceUserMapReq{
		MapIDs: []string{mapID},
	}

	req := testutils.GetReqWithBody(t, reqData)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(accessmap.DeleteServiceUserMap))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		//Data [][]string `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

}

func DeleteUserGroupServiceGroupAccessMap(t *testing.T, mapID string) {

	reqData := accessmap.RmGroupMap{
		MapID: []string{mapID},
	}

	req := testutils.GetReqWithBody(t, reqData)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(accessmap.DeleteServiceGroupUserGroupMap))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		//Data [][]string `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

}

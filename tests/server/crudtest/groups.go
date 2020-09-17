package crudtest

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/groups"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func CreateGroup(t *testing.T, gType string) string {
	creatServReq := models.Group{
		GroupType: gType,
		GroupName: "some-group" + gType,
	}
	req := testutils.GetReqWithBody(t, creatServReq)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(groups.CreateGroup))

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
	return resp.Data[0]

}

func UpdateGroup(t *testing.T, groupID string) {
	creatServReq := models.Group{
		GroupID:   groupID,
		GroupName: "some-user-group",
	}
	req := testutils.GetReqWithBody(t, creatServReq)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(groups.UpdateGroup))

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

func UpdateServiceGroup(t *testing.T, serviceID, groupID, updateType string) {

	reqData := groups.UpdateServiceGroupReq{
		GroupID:    groupID,
		UpdateType: updateType,
		ServiceIDs: []string{serviceID},
	}

	req := testutils.GetReqWithBody(t, reqData)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(groups.UpdateServiceGroup))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]string `json:"data"`
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

}

func UpdateUserGroup(t *testing.T, userID, groupID, updateType string) {

	reqData := groups.UpdateUsersGroupReq{
		GroupID:    groupID,
		UpdateType: updateType,
		UserIDs:    []string{userID},
	}

	req := testutils.GetReqWithBody(t, reqData)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(groups.UpdateUsersGroup))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]string `json:"data"`
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

}

func GetAllGroups(t *testing.T, gType, expectedID string) {

	req := testutils.GetReqWithBody(t, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupType", gType)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(groups.GetAllGroups))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]models.Group `json:"data"`
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

	found := false
	for _, gr := range resp.Data[0] {
		//	t.Log(gr.GroupID, ":", expectedID)
		if gr.GroupID == expectedID {
			found = true
		}
	}

	if !found {
		t.Errorf("newly created %s group not found", gType)
	}

}

func GetUserGroupDetail(t *testing.T, expectedUserID, expectedGroupID string) {

	req := testutils.GetReqWithBody(t, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupid", expectedGroupID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(groups.GetUserGroup))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []groups.GroupUsers `json:"data"`
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

	if resp.Data[0].GroupMeta.GroupID != expectedGroupID {
		t.Errorf("getUserGroup().groupMeta  got=%v want=%v", resp.Data[0].GroupMeta.GroupID, expectedUserID)

	}

	found := false
	for _, us := range resp.Data[0].AddedUsers {
		//t.Log(rule.RuleID,rule.ConstName)
		if us.ID == expectedUserID {
			found = true
		}
	}

	if !found {
		t.Error("user not found in user group")
	}

}

func GetServiceGroupDetail(t *testing.T, expectedServiceID, expectedGroupID string) {

	req := testutils.GetReqWithBody(t, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupID", expectedGroupID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(groups.GetServiceGroup))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []groups.GroupApps `json:"data"`
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

	if resp.Data[0].GroupMeta.GroupID != expectedGroupID {
		t.Errorf("getServiceGroup().groupMeta  got=%v want=%v", resp.Data[0].GroupMeta.GroupID, expectedGroupID)

	}

	found := false
	for _, us := range resp.Data[0].AddedServices {
		//t.Log(us.ID , expectedServiceID)
		//t.Log(us.ID == expectedServiceID)
		if us.ID == expectedServiceID {
			found = true
		}
	}

	if !found {
		t.Error("service not found in service group")
	}

}

func DeleteGroup(t *testing.T, groupID string) {

	req := testutils.GetReqWithBody(t, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupID", groupID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(groups.DeleteGroup))

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

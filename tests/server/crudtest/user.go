package crudtest

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func CreateUser(t *testing.T) models.User {
	reqData := users.CreateUserReq{User: models.UserWithPass{User: models.User{
		UserName:   "someusername",
		FirstName:  "somefirstname",
		MiddleName: "",
		LastName:   "somelastname",
		Email:      "somename@somedomain.com",
		UserRole:   "orgAdmin",
		Status:     true,
		IdpName:    "trasa",
	}}}

	req := testutils.GetReqWithBody(t, reqData)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(users.CreateUser))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []users.CreateUserResp `json:"data"`
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
	if resp.Data[0].ConfirmationLink == "" {
		t.Errorf(`user create confirmation link blank`)
	}

	return resp.Data[0].User

}

func UpdateUser(t *testing.T, userID string) models.User {
	creatServReq := users.CreateUserReq{User: models.UserWithPass{User: models.User{
		ID:        userID,
		UserName:  "someuser",
		FirstName: "different",
		LastName:  "second",
		Email:     "someuser@diffrentDomain.com",
		UserRole:  "orgAdmin",
		Status:    true,
		IdpName:   "trasa",
	}}}
	req := testutils.GetReqWithBody(t, creatServReq)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(users.UpdateUser))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []users.CreateUserReq `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

	t.Log(rr.Body.String())
	if len(resp.Data) == 0 {
		t.Fatal(resp.Reason)
	}
	return models.CopyUserWithoutPass(resp.Data[0].User)

}

func GetUser(t *testing.T, expected models.User) {

	req := testutils.GetReqWithBody(t, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", expected.ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(users.GetUserDetails))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []users.UserDetails `json:"data"`
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

	resp.Data[0].User.ExternalID = expected.ExternalID
	if !reflect.DeepEqual(resp.Data[0].User, expected) {
		t.Errorf("GetUserDetails  \n got = %v,\n want = %v", resp.Data[0].User, expected)
	}

}

func GetAllUsers(t *testing.T, expected models.User) {

	req := testutils.GetReqWithBody(t, nil)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(users.GetAllUsers))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]models.User `json:"data"`
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
		if s.ID == expected.ID {
			exists = true
			s.ExternalID = expected.ExternalID
			if !reflect.DeepEqual(s, expected) {
				t.Errorf("GetAllUsers  \n got = %v,\n want = %v", s, expected)
			}
		}
	}

	if !exists {
		t.Errorf(`newly created user not got in users.GetAllUsers`)
	}

}

func DeleteUser(t *testing.T, userID string) {

	req := testutils.GetReqWithBody(t, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", userID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(users.DeleteUser))

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

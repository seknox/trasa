package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/seknox/ssh"
	"github.com/seknox/trasa/server/accessproxy/sshproxy"
	"github.com/seknox/trasa/server/api/my"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestGetMyServicesDetail(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(my.GetMyServicesDetail))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []my.MyServiceDetail `json:"data"`
	}

	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

	data := resp.Data[0]

	if data.User.ID != "13c45cfb-72ca-4177-b968-03604cab6a27" {
		t.Errorf(`GetMyServicesDetail returned incorrect userID. want: %s got: %s`, "13c45cfb-72ca-4177-b968-03604cab6a27", data.User.ID)
	}

	myServices := data.MyServices

	if len(myServices) != 3 {
		t.Errorf(`GetMyServicesDetail returned incorrect number of services. want: %d got: %d`, 3, len(data.MyServices))
	}

	//TODO add more expectations

}

func TestMyAccountDetails(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(my.MyAccountDetails))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []my.SingleUserDetailV2 `json:"data"`
	}

	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

	data := resp.Data[0]

	if data.User.ID != "13c45cfb-72ca-4177-b968-03604cab6a27" {
		t.Errorf(`MyAccountDetails returned incorrect userID. want: %s got: %s`, "13c45cfb-72ca-4177-b968-03604cab6a27", data.User.ID)
	}

	myServices := data.AssignedServices

	if len(myServices) != 3 {
		t.Errorf(`MyAccountDetails returned incorrect number of services. want: %d got: %d`, 3, len(data.AssignedServices))
	}

	if len(data.UserDevices) < 3 {
		t.Errorf(`MyAccountDetails returned incorrect number of devices. want: %d got: %d`, 3, len(data.UserDevices))
	}

	if len(data.UserGroups) != 1 {
		t.Errorf(`MyAccountDetails returned incorrect number of groups. want: %d got: %d`, 1, len(data.UserGroups))
	}

	//TODO add more expectations

}

func TestAuthMeta(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("appID", "2fef188a-cc13-438b-8564-2803a072f650")
	rctx.URLParams.Add("username", "admin")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(my.GetAuthMeta))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []my.AuthMetaResp `json:"data"`
	}

	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

	data := resp.Data[0]

	if data.IsDeviceHygeneRequired {
		t.Errorf(`GetAuthMeta.IsDeviceHygeneRequired  want: %v got: %v`, false, data.IsDeviceHygeneRequired)
	}

	if data.IsPasswordRequired {
		t.Errorf(`GetAuthMeta.IsPasswordRequired  want: %v got: %v`, false, true)
	}

}

func TestDownloadKey(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(my.GenerateKeyPair))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	testutils.MockPrivateKey2 = string(rr.Body.Bytes())

	k, err := ssh.ParsePrivateKey(rr.Body.Bytes())
	if err != nil {
		t.Errorf(`invalid user key`)
	}

	user, err := sshproxy.SSHStore.GetUserFromPublicKey(k.PublicKey(), testutils.MockOrgID)
	if err != nil {
		t.Errorf(`incorrect user key`)
	}

	if user.ID != testutils.MockUserID {
		t.Errorf(`incorrect user ID, want=%v got=%v`, testutils.MockUserID, user.ID)
	}

}

func TestMyFiles(t *testing.T) {

	fileUpload(t)
	fileDownloadList(t, 1)
	token := getFileDownloadToken(t)
	fileDownload(t, token)
	fileDelete(t)
	fileDownloadList(t, 0)

}

func fileUpload(t *testing.T) {
	file, err := os.Create("somefile.test")
	if err != nil {
		t.Error(err)
	}
	file.WriteString("test data, test data")
	file.Seek(0, 0)

	filename := file.Name()
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		writer.Close()
		t.Error(err)
	}

	io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(my.FileUploadHandler))

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

	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

}

func fileDownloadList(t *testing.T, wantLen int) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(my.GetDownloadableFileList))

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

	fileList := resp.Data[0]

	if len(fileList) != wantLen {
		t.Fatalf("incorrect file list length, got=%d want=%d", len(fileList), wantLen)
	}

}

func getFileDownloadToken(t *testing.T) string {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(my.GetFileDownloadToken))

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

	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

	if rr.Header().Get("sskey") == "" {
		t.Fatal("download token is empty")
	}

	return rr.Header().Get("sskey")

}

func fileDownload(t *testing.T, token string) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("sskey", token)
	rctx.URLParams.Add("fileName", "somefile.test")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(my.FileDownloadHandler))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body := rr.Body.String()

	if body != "test data, test data" {
		t.Errorf("handler returned wrong body: got %v want %v",
			body, "test data, test data")
	}

}

func fileDelete(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("fileName", "somefile.test")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(my.FileDeleteHandler))

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

	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

}

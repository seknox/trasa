package accessproxytest

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/hostrouter"
	webproxy "github.com/seknox/trasa/server/accessproxy/http"
	"github.com/seknox/trasa/server/api/auth/serviceauth"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func prepareProxy(t *testing.T, done chan bool) {
	webproxy.PrepareProxyConfig()

	prouter := chi.NewRouter()
	prouter.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		pxy := webproxy.Proxy()

		pxy.ServeHTTP(w, r)
	})

	hr := hostrouter.New()

	r := chi.NewRouter()
	hr.Map("*", prouter)
	r.Mount("/", hr)

	s := http.Server{
		Addr:    ":3339",
		Handler: r,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			t.Error(err)
		}
	}()

	<-done
	s.Close()

}

func prepareUpstream(t *testing.T, done chan bool) {

	r := chi.NewRouter()
	r.Get("/test", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("someKey", "someValue")
		writer.WriteHeader(200)
		writer.Write([]byte("Hello"))
	})

	s := http.Server{
		Addr:    "localhost:3338",
		Handler: r,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			t.Log(err)
		}
	}()

	<-done
	s.Close()

}

func TestHTTPProxy(t *testing.T) {
	done := make(chan bool)
	go prepareProxy(t, done)
	go prepareUpstream(t, done)
	defer func() { done <- true }()
	time.Sleep(time.Second * 3)
	sess := auth(t)

	req, err := http.NewRequest("GET", `http://localhost:3339/test`, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Host", "localhost:3339")
	req.Header.Set("TRASA-X-CSRF", sess.CsrfToken)
	req.Header.Set("TRASA-X-SESSION", sess.SessionID)

	cl := http.Client{}

	resp, err := cl.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(respBody) != "Hello" {
		t.Errorf("respBody got=%v want=%v", string(respBody), "Hello")
	}

	if hVal := resp.Header.Get("someKey"); hVal != "someValue" {
		t.Errorf("headers.someKey got=%v want=%v", hVal, "someValue")
	}

}

func auth(t *testing.T) serviceauth.Session {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serviceauth.AuthHTTPAccessProxy)

	handler.ServeHTTP(rr, testutils.GetReqWithBody(t, serviceauth.NewSession{
		HostName:  "localhost:3339",
		TfaMethod: "totp",
		TotpCode:  testutils.GetTotpCode(testutils.MocktotpSEC),
		ExtToken:  "cb6dd3f6-54c2-4cb0-b294-e22c2aa708e4",
	}))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []serviceauth.Session `json:"data"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatalf("AuthHTTPAccessProxy() wanted success, got:%s reason %s", resp.Status, resp.Reason)

	}

	if len(resp.Data) == 0 {
		t.Fatalf("AuthHTTPAccessProxy() wanted success, got:%s reason %s", resp.Status, resp.Reason)

	}

	return resp.Data[0]

}

package integration_tests

import (
	"encoding/json"
	"github.com/seknox/trasa/server/api/my"
	"github.com/seknox/trasa/server/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMyServicesDetail(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddTestUserContext(my.GetMyServicesDetail))

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

}

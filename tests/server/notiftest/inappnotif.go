package notiftest

import (
	"encoding/json"
	"github.com/seknox/trasa/server/api/notif"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func GetPendingNotif(t *testing.T, expextedNotif models.InAppNotification) {

	req := testutils.GetReqWithBody(t, nil)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(notif.GetPendingNotif))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data [][]models.InAppNotification `json:"data"`
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

	found := false
	for _, gotNotif := range data {
		if gotNotif.NotificationID == expextedNotif.NotificationID {
			if !reflect.DeepEqual(gotNotif, expextedNotif) {
				t.Errorf("InAppNotif  \n got = %v,\n want = %v", gotNotif, expextedNotif)
			}
			found = true
		}
	}

	if !found {
		t.Error("Cound not get newly added notif")

	}

}

func AddNotif(t *testing.T) models.InAppNotification {
	testNotif := models.InAppNotification{
		NotificationID:    utils.GetUUID(),
		UserID:            testutils.MockUserID,
		EmitterID:         testutils.MockUserID,
		OrgID:             testutils.MockOrgID,
		NotificationLabel: "Some Label",
		NotificationText:  "Some event has occurred",
		CreatedOn:         time.Now().Unix(),
		IsResolved:        false,
	}

	err := notif.Store.StoreNotif(testNotif)
	if err != nil {
		t.Fatalf("could not store notif: %v", err)
	}

	return testNotif
}

func ResolvNotif(t *testing.T, notifID string) {
	reqData := struct {
		NotifID string `json:"notifID"`
	}{
		NotifID: notifID,
	}

	req := testutils.GetReqWithBody(t, reqData)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(notif.ResolveNotif))

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

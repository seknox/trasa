package systemtest

import (
	"encoding/json"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func SystemStatus(t *testing.T) {
	req := testutils.GetReqWithBody(t, nil)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(system.SystemStatus))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []system.SysStatus `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
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

	if data.DiskStat == nil {
		t.Fatalf(`DiskStat is nil, resp: %s`, string(rr.Body.Bytes()))
	}
	if data.HostStat == nil {
		t.Fatalf(`DiskStat is nil, resp: %s`, string(rr.Body.Bytes()))
	}
	if data.MemStat == nil {
		t.Fatalf(`DiskStat is nil, resp: %s`, string(rr.Body.Bytes()))
	}

}

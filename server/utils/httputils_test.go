package utils

import (
	"net/http/httptest"
	"testing"
)

func TestTrasaResponse(t *testing.T) {
	//	t.Log("testing..")
	w := httptest.NewRecorder()
	TrasaResponse(w, 200, "success", "succeed", "test TrasaResponse", nil)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	expectedBody := `{"status":"success","reason":"succeed","intent":"test TrasaResponse","data":[null]}`
	if w.Body.String() != expectedBody {
		t.Fatalf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expectedBody)
	}

}

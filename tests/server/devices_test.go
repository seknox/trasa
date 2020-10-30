package server_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/auth"
	"github.com/seknox/trasa/server/api/crypt"
	"github.com/seknox/trasa/server/api/my"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/seknox/trasa/tests/server/testutils"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"
)

//TODO

func TestRegisterWorkstation(t *testing.T) {
	key := trasadaKex(t, "enrolDevice")
	deviceRegisterReq(t, key)
}

func deviceRegisterReq(t *testing.T, key []byte) {
	//RegisterUserDevice

	dh := auth.DeviceDetail{
		DeviceBrowser:     models.DeviceBrowser{},
		BrowserExtensions: nil,
		DeviceHygiene: models.DeviceHygiene{
			DeviceInfo: models.DeviceInfo{
				DeviceName:    "some Device",
				DeviceVersion: "1.1.1",
				MachineID:     "askdnkjs87asyd7yaushduiasnd asd",
				Brand:         "some brand",
				Manufacturer:  "some Manufact",
				DeviceModel:   "some model",
			},
			DeviceOS: models.DeviceOS{
				OSName:              "some Os",
				OSVersion:           "2.3.4",
				KernelType:          "linux",
				KernelVersion:       "1.1.1",
				ReadableVersion:     "",
				LatestSecurityPatch: "",
				AutoUpdate:          false,
				PendingUpdates:      nil,
				JailBroken:          false,
				DebugModeEnabled:    false,
				IsEmulator:          false,
			},
			LoginSecurity: models.LoginSecurity{
				AutologinEnabled:         false,
				LoginMethod:              "",
				PasswordLastUpdated:      "",
				TfaConfigured:            false,
				IdleDeviceScreenLockTime: "",
				IdleDeviceScreenLock:     false,
				RemoteLoginEnabled:       false,
			},
			NetworkInfo: models.NetworkInfo{
				Hostname:         "somehost",
				DomainControlled: false,
				DomainName:       "some domain",
				InterfaceName:    "some int",
				IPAddress:        "192.168.0.1",
				MacAddress:       "",
				WirelessNetwork:  false,
				OpenWifiConn:     false,
				NetworkName:      "",
				NetworkSecurity:  "",
			},
			EndpointSecurity: models.EndpointSecurity{
				EpsConfigured:           false,
				EpsVendorName:           "",
				EpsVersion:              "",
				EpsMeta:                 "",
				FirewallEnabled:         false,
				FirewallPolicy:          "",
				DeviceEncryptionEnabled: false,
				DeviceEncryptionMeta:    "",
			},
			LastCheckedTime: 0,
		},
	}

	dhbytes, _ := json.Marshal(dh)
	encdh, err := utils.AESEncrypt(key, dhbytes)
	if err != nil {
		t.Fatal(err)
	}

	enrolReq := auth.RegisterDeviceReq{
		TfaMethod:     "totp",
		TotpCode:      testutils.GetTotpCode(testutils.MocktotpSEC),
		TrasaID:       testutils.MockTrasaID,
		OrgID:         testutils.MockOrgID,
		DeviceName:    "Some device name",
		DeviceHygiene: hex.EncodeToString(encdh),
	}

	req := testutils.GetReqWithBody(t, enrolReq)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(auth.RegisterUserDevice))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data string `json:"data"`
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
}

func trasadaKex(t *testing.T, intent string) []byte {

	// gen keypair for trasaExtNative.
	priv, pub, err := utils.ECDHGenKeyPair()
	if err != nil {
		logrus.Error(err)
	}
	enrolReq := crypt.KexRequest{
		Intent:    intent,
		IntentID:  testutils.MockTrasaID,
		DeviceID:  "",
		PublicKey: hex.EncodeToString(pub[:]),
	}

	req := testutils.GetReqWithBody(t, enrolReq)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(crypt.Kex))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data string `json:"data"`
	}

	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("failed to read kex response: %v", err)
	}

	// update secret key
	pubBytes, err := hex.DecodeString(resp.Data)
	if err != nil {
		t.Fatal(err)
	}

	var serverPublicKey [32]byte

	copy(serverPublicKey[:], pubBytes)

	sec := utils.ECDHComputeSecret(priv, &serverPublicKey)

	t.Log("our secret key: ", hex.EncodeToString(sec))

	return sec
}

func TestDeleteMyDevice(t *testing.T) {
	t.Run("deleting own device should pass", func(t *testing.T) {

		req, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatal(err)
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("deviceID", "db8dc6b0-84c1-43b2-929f-fa0479a25441")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(testutils.AddTestSelfUserContext(my.RemoveMyDevice))

		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var resp models.TrasaResponseStruct

		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		if err != nil {
			t.Fatal(err)
		}

		if resp.Status != "success" {
			t.Fatal(resp.Reason)
		}

	})

	t.Run("deleting other user's device should fail", func(t *testing.T) {

		req, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatal(err)
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("deviceID", "aac1e00f-bcd0-4fb1-87b2-8f3cbb746b14")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(testutils.AddTestSelfUserContext(my.RemoveMyDevice))

		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusForbidden {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusForbidden)
		}

		var resp models.TrasaResponseStruct

		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		if err != nil {
			t.Fatal(err)
		}

		if resp.Status != "failed" {
			t.Errorf(`expected status %s got %s`, "failed", resp.Reason)
		}

	})

}

package devicestest

import (
	"bytes"
	"encoding/json"
	"github.com/seknox/trasa/server/api/auth"
	"github.com/seknox/trasa/server/api/devices"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/tests/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func EnrollTFADeviceTest(t *testing.T) {

	t.Log("enroll device test with trasa authenticator")

	totpSec, deviceID := sendCreds(t)
	sendDeviceDetailsToCloudProxy(t, deviceID)
	confirmTOTPCode(t, totpSec, deviceID)
	token := login(t)
	sendTFA(t, token, totpSec)
	time.Sleep(time.Second)
	deleteMyDevice(t, deviceID)

	t.Log("enroll device test with 3rd party authenticator")
	totpSec, deviceID = sendCreds(t)
	confirmTOTPCode(t, totpSec, deviceID)
	token = login(t)
	sendTFA(t, token, totpSec)
	deleteMyDevice(t, deviceID)
}

func sendCreds(t *testing.T) (totpSec, deviceID string) {
	req := testutils.GetReqWithBody(t, auth.LoginRequest{
		Email:    testutils.MockTrasaID2,
		Password: testutils.MocktrasaPass2,
	})

	req.RemoteAddr = "127.0.0.1:4567"

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.LoginHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []devices.EnrolDeviceStruct `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("could not unmarshall: %v,%v", err, rr.Body.String())
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

	if len(resp.Data) != 1 {
		t.Fatal(resp.Reason)
	}

	data := resp.Data[0]

	if data.Account != testutils.MockTrasaID2 {
		t.Errorf("handler returned wrong account name: got %v want %v",
			data.Account, testutils.MockTrasaID2)
	}

	if data.CloudProxyURL != "https://sg.cpxy.trasa.io" {
		t.Errorf("handler returned wrong account name: got %v want %v",
			data.CloudProxyURL, testutils.MockTrasaID2)
	}

	if data.TotpSSC == "" {
		t.Errorf("handler returned empty totp secret got %v ",
			data.TotpSSC)
	}
	if data.DeviceID == "" {
		t.Errorf("handler returned empty DeviceID got %v ",
			data.DeviceID)
	}

	if data.OrgName != "Trasa" {
		t.Errorf("handler returned wrong org name: got %v want %v",
			data.OrgName, "Trasa")
	}

	return data.TotpSSC, data.DeviceID

}

func sendDeviceDetailsToCloudProxy(t *testing.T, deviceID string) {

	dh := models.DeviceHygiene{
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
	}
	dhBytes, err := json.Marshal(dh)

	dataBytes, err := json.Marshal(devices.DeviceEnrollResp{
		DeviceID:     deviceID,
		FCMToken:     "kasndjbasiduaysgduasbdjhasbdhjasudvasutdvjhc xnzmcxct765c67as57dgasudjabsjdhb23;l4k23lj43iasugdbhjs==//\\\\sadjksa",
		PublicKey:    "-----BEGIN RSA PUBLIC KEY-----\nMIICCgKCAgEA3OaU59JGmotiY7TQGKSSukdEuvSpPIclI9/lG/z8WiJs9axm1xRo\nobFHyKITEr1Kdfhj8YezNxgCK12OKptPM+U3kn8MlHpZVw7H7VIf74kp8/C/0HZ2\n6fGE0bYYo3ZqIyzyFHUFvDlja17KYKuTUOuVXtmIq0fEb9cLIdqXANnFwbL5brji\niyjdLhSo9v/FadYJP2FDmtFD/3ymeujRAKc33RTKmd6YEZbeSWXST7ddO0SZIOax\nN6FWlgx13vs0IpL9pbsTm+jEbJeuYZOTAfI1kXh+tViv/p4lStnSt1X2eujD1EvG\nStzdw5Q6zy7VCBgDloBgGKyLjwWx1FToByNHFG+qFgbmjm6kCAsZ6R1dFiH2dRxS\n5CYFAYzgau2tgQHFBmkWfA1U3IvDeeA6zDKiRsRh+YAb3rvBD86zW70Fce1uyqwH\ntVqUUMee2csF1RolgRs9zgaCbjcRvOe7JYHF6HYgvcD3pWzXeDPgjTqzGoF850KX\nujIWg3l8ad1tIBypx3M+3Y0SV8w5vBUhy5Tf2poeDSf+PsiNuj/jWfGzk7IJtQFR\nsgQGV4G61G3o5k/A0np77ix93tqY6INnxW8l7jgBQj+zIo0xEDOACryNo3PWaNXD\nciPS7csjtng2kwZT53Rm8sPi83VQRDRcU6G0JP7TYB4VacBCoxEK76sCAwEAAQ==\n-----END RSA PUBLIC KEY-----",
		DeviceFinger: string(dhBytes),
	})

	r, err := http.Post("https://sg.cpxy.trasa.io/api/v1/passmydevicedetail", "application/json", bytes.NewBuffer(dataBytes))

	if err != nil {
		t.Fatal(err)
	}

	// Check the status code is what we expect.
	if status := r.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

}

func confirmTOTPCode(t *testing.T, totpSec, deviceID string) {
	totpCode := testutils.GetTotpCode(totpSec)

	req := testutils.GetReqWithBody(t, auth.ConfirmTOTPPreq{
		TOTPCode: totpCode,
		DeviceID: deviceID,
	})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.ConfirmTOTPAndSave)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []devices.EnrolDeviceStruct `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Fatal(resp.Reason)
	}

}

func login(t *testing.T) (token string) {
	req := testutils.GetReqWithBody(t, auth.LoginRequest{
		Email:    testutils.MockTrasaID2,
		Password: testutils.MocktrasaPass2,
		Intent:   consts.AUTH_REQ_DASH_LOGIN,
	})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.LoginHandler)

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
		t.Errorf("handler returned wrong status : got %v want %v",
			resp.Status, "success")
	}

	if resp.Intent != consts.AUTH_RESP_TFA_REQUIRED {
		t.Errorf("handler returned wrong intent : got %v want %v",
			resp.Intent, consts.AUTH_RESP_TFA_REQUIRED)
	}
	if len(resp.Data) != 1 {
		t.Error("login token empty")
	}

	return resp.Data[0]
}

func sendTFA(t *testing.T, token, totpSec string) {
	req := testutils.GetReqWithBody(t, auth.TfaRequest{
		Token:     token,
		TfaMethod: "totp",
		Totp:      testutils.GetTotpCode(totpSec),
		Intent:    consts.AUTH_REQ_DASH_LOGIN,
	})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.TfaHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp struct {
		models.TrasaResponseStruct
		Data []struct {
			User      models.User `json:"user"`
			CSRFToken string      `json:"CSRFToken"`
		} `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "success" {
		t.Errorf("handler returned wrong status with reason %s : got %v want %v",
			resp.Reason, resp.Status, "success")
	}

	if len(resp.Data) != 1 {
		t.Error("login token empty")
	}
	cookies := rr.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "X-SESSION" && cookie.Value != "" {
			return
		}
	}
	t.Error("session cookie not found")

}

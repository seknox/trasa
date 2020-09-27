package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/seknox/fireser/wrkstn/hygiene"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

func enrolOrSyncDevice(extData IncomingMessage, intent string) (string, error) {

	return sendHygiene(extData, intent)
}

func sendHygiene(extData IncomingMessage, intent string) (string, error) {
	// perform key exchange
	err := kex(extData.Data.AuthData.TrasaID, "", "KEX_ENROL_DEVICE", extData.Data.Hostname)
	if err != nil {
		logrus.Error("Failed to Kex: ", err)
	}

	// secret key should already be retreived. lets encrypt data.
	dh := getHygiene()
	var deviceEncData deviceDetail
	deviceEncData.BrowserExtensions = extData.Data.BrowserExtensions
	deviceEncData.DeviceBrowser = extData.Data.DeviceBrowser
	deviceEncData.DeviceHygiene = dh
	dhbytes, err := getDHJSONBytes(deviceEncData)
	if err != nil {
		logrus.Error(err)
		return "", fmt.Errorf("unable to get hyegiene")
	}

	key, ok := keyHolder["KEX_ENROL_DEVICE"]
	if !ok {
		logrus.Errorf("secret key not found")
		return "", fmt.Errorf("unable to retreive kex data")
	}
	cipherText, err := utils.AESEncrypt(key, dhbytes)
	if err != nil {
		logrus.Error(err)
		return "", fmt.Errorf("unable to encrypt: %v", err)
	}

	var enrolReq registerDeviceReq
	enrolReq.DeviceHygiene = hex.EncodeToString(cipherText)
	enrolReq.DeviceName = fmt.Sprintf("%s (%s)", dh.NetworkInfo.Hostname, dh.DeviceOS.OSName)
	enrolReq.TfaMethod = extData.Data.AuthData.TfaMethod
	enrolReq.TotpCode = extData.Data.AuthData.TotpCode
	enrolReq.TrasaID = extData.Data.AuthData.TrasaID

	respData, err := sendRegisterDeviceReq(enrolReq, extData.Data.Hostname)
	if err != nil {
		logrus.Error(err)
	}

	delete(keyHolder, "KEX_ENROL_DEVICE")

	return respData, nil
}

func sendRegisterDeviceReq(req registerDeviceReq, hostName string) (string, error) {
	client := getHTTPClient(true)
	url := fmt.Sprintf("%s/auth/device/register", hostName)

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	hresp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}
	defer hresp.Body.Close()

	var resp trasaResponse
	err = json.NewDecoder(hresp.Body).Decode(&resp)
	if err != nil {
		return "", fmt.Errorf("failed to get device detail: %v", err)
	}

	if resp.Status == "failed" {
		return "", fmt.Errorf("failed to register")
	}

	return resp.Data, nil
}

func getEncryptedHygiene(extData IncomingMessage) (string, error) {
	logrus.Debug("received getEncryptedHygiene")

	// server public key
	pubBytes, err := hex.DecodeString(extData.Data.PubKey)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	logrus.Debug("before ECDHGenKeyPair")
	// gen our keypair
	priv, pub, err := utils.ECDHGenKeyPair()
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	logrus.Debug("before copy serverPublicBytes")
	var serverPublicBytes [32]byte
	copy(serverPublicBytes[:], pubBytes)

	logrus.Debug("before ECDHComputeSecret")

	// gen secret key
	sec := utils.ECDHComputeSecret(priv, &serverPublicBytes)
	logrus.Debug("our secret: ", hex.EncodeToString(sec))
	// prepare dh
	dh := getHygiene()
	var deviceEncData deviceDetail
	deviceEncData.BrowserExtensions = extData.Data.BrowserExtensions
	deviceEncData.DeviceBrowser = extData.Data.DeviceBrowser
	deviceEncData.DeviceHygiene = dh
	dhbytes, err := getDHJSONBytes(deviceEncData)
	if err != nil {
		logrus.Error(err)
		return "", fmt.Errorf("unable to get hyegiene")
	}

	logrus.Debug("before AESEncrypt ")

	// encrypt dh
	cipherText, err := utils.AESEncrypt(sec, dhbytes)
	if err != nil {
		logrus.Error(err)
		return "", fmt.Errorf("unable to encrypt: %v", err)
	}

	type getHygieneResp struct {
		ClientPubKey string `json:"clientPubKey"`
		EncryptedDH  string `json:"encryptedDH"`
	}

	var resp getHygieneResp
	resp.EncryptedDH = hex.EncodeToString(cipherText)
	resp.ClientPubKey = hex.EncodeToString(pub[:])

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logrus.Error(err)
		return "", fmt.Errorf("unable to marshal: %v", err)
	}

	logrus.Debug("before returning ")

	return string(respBytes), nil
}

type extResponse struct {
	ExtID      string   `json:"extID"`
	RootDomain string   `json:"rootDomain"`
	SsoDomain  string   `json:"ssoDomain"`
	WSPath     string   `json:"wsPath"`
	Hosts      []string `json:"hosts"`
	TrasaDACom bool     `json:"trasaDACom"`
}

type registerDeviceReq struct {
	TfaMethod string `json:"tfaMethod"`
	TotpCode  string `json:"totpCode"`
	TrasaID   string `json:"trasaID"`
	OrgID     string `json:"orgID"`
	// Device name to be used before we decrypt deviceHygiene
	DeviceName    string `json:"deviceName"`
	DeviceHygiene string `json:"deviceHygiene"`
}

type browserData struct {
	Hostname          string                     `json:"hostName"`
	PubKey            string                     `json:"pubKey"`
	DeviceBrowser     models.DeviceBrowser       `json:"deviceBrowser"`
	BrowserExtensions []models.BrowserExtensions `json:"browserExtensions"`
	AuthData          tfaData                    `json:"authData"`
}

type tfaData struct {
	TfaMethod string `json:"tfaMethod"`
	TotpCode  string `json:"totpCode"`
	TrasaID   string `json:"trasaID"`
}

type deviceDetail struct {
	DeviceBrowser     models.DeviceBrowser       `json:"deviceBrowser"`
	BrowserExtensions []models.BrowserExtensions `json:"browserExtensions"`
	DeviceHygiene     hygiene.DeviceHygiene      `json:"deviceHygiene"`
}

func getHygiene() hygiene.DeviceHygiene {
	return hygiene.GetDeviceHygiene(runtime.GOOS)
}

func getDHJSONBytes(dh deviceDetail) ([]byte, error) {
	mar, err := json.Marshal(dh)
	if err != nil {
		return mar, err
	}

	return mar, nil
}

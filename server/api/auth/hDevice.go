package auth

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/providers/ca"
	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/api/users"
	"golang.org/x/crypto/ssh"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/seknox/trasa/server/api/auth/tfa"
	"github.com/seknox/trasa/server/api/devices"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

type RegisterDeviceReq struct {
	TfaMethod string `json:"tfaMethod"`
	TotpCode  string `json:"totpCode"`
	TrasaID   string `json:"trasaID"`
	OrgID     string `json:"orgID"`
	// Device name to be used before we decrypt deviceHygiene
	DeviceName    string `json:"deviceName"`
	DeviceHygiene string `json:"deviceHygiene"`
}

type DeviceDetail struct {
	DeviceBrowser     models.DeviceBrowser       `json:"deviceBrowser"`
	BrowserExtensions []models.BrowserExtensions `json:"browserExtensions"`
	DeviceHygiene     models.DeviceHygiene       `json:"deviceHygiene"`
}

// RegisterUserDevice registers new user device and stores device hygiene, device browser and browser extensions details.
func RegisterUserDevice(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("RegisterUserDevice request received")

	var req RegisterDeviceReq
	err := utils.ParseAndValidateRequest(r, &req)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "success", "invalid request", "RegisterUserDevice", nil)
		return
	}

	authlog := logs.NewLog(r, "regDevice")

	// get user info from database
	userDetailFromDB, err := Store.GetLoginDetails(req.TrasaID, req.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "User not found", "Ext login", nil)
		return
	}
	authlog.UpdateUser(userDetailFromDB)

	switch req.TfaMethod {
	case "totp":
		check, _, err := tfa.VerifyTotpCode(req.TotpCode, userDetailFromDB.ID, userDetailFromDB.OrgID)
		if err != nil || !check {
			logs.Store.LogLogin(&authlog, consts.REASON_INVALID_TOTP, false)
			utils.TrasaResponse(w, 200, "failed", "totp failed", "Device registration", nil)
			return
		}
	default:
		status, msg := tfa.SendU2F(userDetailFromDB.ID, userDetailFromDB.OrgID, fmt.Sprintf("Device registration: %s", req.DeviceName), utils.GetIp(r))
		if !status {
			logs.Store.LogLogin(&authlog, consts.REASON_U2F_FAILED, false)
			utils.TrasaResponse(w, 200, "failed", msg, "Device registration", nil)
			return
		}
	}

	// IF we reach here, all required authentication has been passed.
	// We will decrypt device hygiene and then store details in database.
	// Note: Do not forget to delete secret key before returning from this handler.

	//  retrieve secret key for this request.
	secretKeyFromKex, ok := global.ECDHKexDerivedKey[req.TrasaID]
	if !ok {
		logrus.Tracef("key not found in Kex store for %s ", req.TrasaID)
		utils.TrasaResponse(w, 200, "failed", "key not found in Kex store", "Device registration", nil)
		return
	}

	dhBytes, err := hex.DecodeString(req.DeviceHygiene)
	if err != nil {
		logrus.Debug("cannot decode hex string", err)
		logs.Store.LogLogin(&authlog, consts.REASON_SPOOFED_LOGIN, false)
		utils.TrasaResponse(w, 200, "failed", "failed to decrypt data", "Device registration", nil)
		return
	}

	// decrypt the device details.
	decryptedBytes, err := utils.AESDecrypt(secretKeyFromKex.Secretkey[:], dhBytes)
	if err != nil {
		logrus.Debug(err)
		logs.Store.LogLogin(&authlog, consts.REASON_SPOOFED_LOGIN, false)
		utils.TrasaResponse(w, 200, "failed", "failed to decrypt data", "Device registration", nil)
		return
	}

	//delete secret key from store
	delete(global.ECDHKexDerivedKey, req.TrasaID)

	// unmarshall decryptedBytes to deviceDetail struct
	var dh DeviceDetail
	err = json.Unmarshal(decryptedBytes, &dh)
	if err != nil {
		//TODO return error here??
		logrus.Debug("cannot unmarshall decrypted Device Hygiene", err)
	}

	// begin store process

	//register device if its not already registered
	deviceID, err := devices.Store.GetDeviceIDFromExtID(dh.DeviceHygiene.DeviceInfo.MachineID)
	if err != nil || deviceID == "" {
		logrus.Debug(err)
		deviceID = utils.GetUUID()
		// Create new workstation and store dh.
		var rd models.UserDevice
		rd.DeviceID = deviceID
		rd.UserID = userDetailFromDB.ID
		rd.OrgID = userDetailFromDB.OrgID
		rd.MachineID = dh.DeviceHygiene.DeviceInfo.MachineID
		rd.DeviceType = "workstation"
		rd.FcmToken = ""
		rd.PublicKey = ""
		rd.DeviceFinger = "{}"

		//dh.DeviceHygiene.NetworkInfo.IPAddress = utils.GetIp(r)
		rd.DeviceHygiene = dh.DeviceHygiene

		rd.AddedAt = time.Now().Unix()

		err = devices.Store.Register(rd)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "Could not register device", "ExtLogin", nil)
			return
		}

	}

	// Store Browser reference device_id of rd.DeviceID
	var brsr models.DeviceBrowser
	brsr.ID = utils.GetUUID()
	brsr.OrgID = userDetailFromDB.OrgID
	brsr.DeviceID = deviceID
	brsr.Version = dh.DeviceBrowser.Version
	brsr.Name = dh.DeviceBrowser.Name
	brsr.Build = dh.DeviceBrowser.Build
	brsr.UserAgent = dh.DeviceBrowser.UserAgent
	//  TODO hardcoded false value here. Do we really need to check isBot?
	brsr.IsBot = false
	brsr.Extensions = dh.BrowserExtensions

	logrus.Debug("IDD: ", deviceID, brsr.DeviceID)

	err = devices.Store.RegisterBrowser(brsr)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not register device", "ExtLogin", nil)
		return
	}

	// store browser extensions
	// for _, v := range dh.BrowserExtensions {
	// 	check := devices.IsKnownExts(v.ExtensionID)
	// 	if check != true {
	// 		err := devices.Store.BrowserStoreExtensionDetails(v, userDetailFromDB.OrgID, userDetailFromDB.ID, brsr.ID)
	// 		if err != nil {
	// 			// if we get error here, it means extensiondetails could not be store, alert admins here. TODO
	// 			logrus.Trace(err)
	// 		}
	// 	}

	// }

	logs.Store.LogLogin(&authlog, "", true)
	logrus.Trace("RegisterUserDevice- Sending Response")
	// Respond with success
	regDeviceRes(w, userDetailFromDB.OrgID, brsr.ID)
}

type extResponse struct {
	ExtID      string   `json:"extID"`
	RootDomain string   `json:"rootDomain"`
	SsoDomain  string   `json:"ssoDomain"`
	WSPath     string   `json:"wsPath"`
	Hosts      []string `json:"hosts"`
	TrasaDACom bool     `json:"trasaDACom"`
}

func regDeviceRes(w http.ResponseWriter, orgID, deviceID string) {

	var resp extResponse
	resp.ExtID = deviceID
	resp.RootDomain = global.GetConfig().Trasa.Rootdomain
	resp.SsoDomain = global.GetConfig().Trasa.Ssodomain
	u, err := url.Parse(global.GetConfig().Trasa.Dashboard)
	if err != nil {
		logrus.Error(err)
	}
	resp.WSPath = fmt.Sprintf("wss://%s", u.Host)

	allservices, err := services.Store.GetAllByType("http", orgID)
	if err != nil {
		utils.TrasaResponse(w, 200, "failed", "No http services available", "SyncExtension", nil)
		return
	}

	resp.Hosts = make([]string, 0)
	for _, v := range allservices {
		resp.Hosts = append(resp.Hosts, v.Hostname)
	}

	globalDeviceCheck, err := system.Store.GetGlobalSetting(orgID, consts.GLOBAL_DEVICE_HYGIENE_CHECK)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error fetching global settings", "GlobalSettings", nil)
		return
	}

	resp.TrasaDACom = globalDeviceCheck.Status

	r, err := json.Marshal(resp)
	if err != nil {
		logrus.Error(err)
	}
	utils.TrasaResponseWithDataString(w, 200, "success", "authorized and token garnted", "ExtLogin", string(r))
}

type syncExtReq struct {
	ExtID string `json:"extID"`
}

// SyncExtension extension does not requires sending in device hygiene.
func SyncExtension(w http.ResponseWriter, r *http.Request) {
	var req syncExtReq

	logrus.Trace("SYnc Extension REQ")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "SyncExtension", nil)
		return
	}

	// verify if extID is already registered
	orgID, err := devices.Store.CheckIfExtIsRegistered(req.ExtID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid DeviceID", "SyncExtension", nil)
		return
	}

	var resp extResponse
	resp.ExtID = req.ExtID
	resp.RootDomain = global.GetConfig().Trasa.Rootdomain
	resp.SsoDomain = global.GetConfig().Trasa.Ssodomain
	u, err := url.Parse(global.GetConfig().Trasa.Dashboard)
	if err != nil {
		logrus.Error(err)
	}
	resp.WSPath = fmt.Sprintf("wss://%s", u.Host)

	allservices, err := services.Store.GetAllByType("http", orgID)
	if err != nil {
		utils.TrasaResponse(w, 200, "failed", "No http services available", "SyncExtension", nil)
		return
	}

	resp.Hosts = make([]string, 0)
	for _, v := range allservices {
		resp.Hosts = append(resp.Hosts, v.Hostname)
	}

	globalDeviceCheck, err := system.Store.GetGlobalSetting(orgID, consts.GLOBAL_DEVICE_HYGIENE_CHECK)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error fetching global settings", "SyncExtension", nil)
		return
	}

	resp.TrasaDACom = globalDeviceCheck.Status

	utils.TrasaResponse(w, 200, "success", "ext synced", "SyncExtension", resp)
}

type UpdateHygienereq struct {
	TrasaID       string `json:"trasaID"`
	DeviceHygiene string `json:"deviceHygiene"`
	ClientKey     string `json:"clientKey"`
	Token         string `json:"token"`
}

func UpdateHygiene(w http.ResponseWriter, r *http.Request) {
	var req UpdateHygienereq
	err := utils.ParseAndValidateRequest(r, &req)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "Update hygiene", nil)
		return
	}

	authlog := logs.NewLog(r, "updateHyg")

	userDetailFromDB, err := Store.GetLoginDetails(req.TrasaID, "")
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "User not found", "Ext login", nil)
		return
	}

	dhBytes, err := hex.DecodeString(req.DeviceHygiene)
	if err != nil {
		logrus.Debug("cannot decode hex string", err)
		utils.TrasaResponse(w, 200, "failed", "failed to decrypt data", "Device hygiene update", nil)
		return
	}

	logrus.Info(req.Token, "+")

	privKey, err := redis.Store.Get(req.Token, "priv")

	if err != nil || privKey == "" {
		err := logs.Store.LogLogin(&authlog, consts.REASON_INVALID_TOKEN, false)
		if err != nil {
			logrus.Error(err)
		}

		utils.TrasaResponse(w, 200, "failed", "invalid token", "Device hygiene update", nil)
		return
	}

	privFromHexStr, err := hex.DecodeString(privKey)
	if err != nil {
		logrus.Errorf("privFromHexStr: %v ", err)
		utils.TrasaResponse(w, 200, "failed", "could not decode private key", "Device hygiene update", nil)
		return
	}

	pubFromHexStr, err := hex.DecodeString(req.ClientKey)
	if err != nil {
		logrus.Errorf("pubFromHexStr: %v ", err)
		utils.TrasaResponse(w, 200, "failed", "could not decode public key", "Device hygiene update", nil)
		return
	}

	var privBytes [32]byte
	var pubBytes [32]byte

	copy(privBytes[:], privFromHexStr)
	copy(pubBytes[:], pubFromHexStr)

	dhBytes, err = hex.DecodeString(req.DeviceHygiene)
	if err != nil {
		logrus.Errorf("dhBytes: %v ", err)
		utils.TrasaResponse(w, 200, "failed", "could not decode string", "Device hygiene update", nil)
		return
	}

	sec := utils.ECDHComputeSecret(&privBytes, &pubBytes)

	plainText, err := utils.AESDecrypt(sec, dhBytes)

	var dh DeviceDetail
	err = json.Unmarshal(plainText, &dh)
	if err != nil {
		logrus.Errorf("json.Unmarshal(plainText, &dh): %v ", err)
		utils.TrasaResponse(w, 200, "failed", "could not unmarshall", "Device hygiene update", nil)
		return
	}

	deviceID, err := devices.Store.UpdateDeviceHygiene(dh.DeviceHygiene, userDetailFromDB.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not update device hygiene", "Update hygiene", nil)
		return
	}

	//pass userID in context
	privateKeyBytes, publicKeyBytes, certBytes, err := generateTempCertificateForDeviceAgent(deviceID, userDetailFromDB.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not generate TempCertificate ForDeviceAgent", "Update hygiene", nil)
		return
	}

	err = users.Store.UpdatePublicKey(userDetailFromDB.ID, strings.TrimSpace(string(certBytes)))
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not update public key", "Update hygiene", nil)
		return
	}

	// Create a new zip archive.
	buff := &bytes.Buffer{}
	zipWriter := zip.NewWriter(buff)

	// Add some files to the archive.
	var files = []struct {
		Name string
		Body []byte
	}{
		{"id_rsa", privateKeyBytes},
		{"id_rsa.pub", publicKeyBytes},
		{"id_rsa-cert.pub", certBytes},
	}
	for _, file := range files {

		var zipFile io.Writer
		zipFile, err = zipWriter.Create(file.Name)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "could not generate zipWriter.Create", "Update hygiene", nil)
			return
		}
		_, err = zipFile.Write(file.Body)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "could not generate zipWriter.Create", "Update hygiene", nil)
			return
		}

	}

	// Make sure to check the error on Close.
	err = zipWriter.Close()
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not generate zipWriter.Create", "Update hygiene", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "", "Update hygiene", buff.Bytes())
	return

}

func generateTempCertificateForDeviceAgent(deviceID, orgID string) (privateKeyBytes, publicKeyBytes, certBytes []byte, err error) {

	bitSize := 4096
	privateKey, err := utils.GeneratePrivateKey(bitSize)
	if err != nil {
		logrus.Errorf(`could not generate private key: %v`, err)
		return
	}

	certHolder, err := ca.Store.GetCertHolder(consts.CERT_TYPE_SSH_CA, "system", orgID)
	if err != nil {
		logrus.Debugf(`could not get CA key: %v`, err)
		return
	}

	publicKeySSH, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		logrus.Errorf(`could not generate public key: %v`, err)
		return
	}

	publicKeyBytes = ssh.MarshalAuthorizedKey(publicKeySSH)

	caKey, err := ssh.ParsePrivateKey(certHolder.Key)
	if err != nil {
		logrus.Errorf(`Could not parse CA private key: %v`, err)
		return
	}

	buf := make([]byte, 8)
	_, err = rand.Read(buf)
	if err != nil {
		logrus.Errorf("failed to read random bytes: %v", err)

		return
	}
	serial := binary.LittleEndian.Uint64(buf)

	//extentions := make(map[string]string)
	extentions := map[string]string{
		"permit-X11-forwarding":   "",
		"permit-agent-forwarding": "",
		"permit-port-forwarding":  "",
		"permit-pty":              "",
		"permit-user-rc":          "",
		"trasa-hygiene":           "true",
		"trasa-device-id":         deviceID,
	}

	//principals := []string{}

	cert := ssh.Certificate{
		Key:             publicKeySSH,
		Serial:          serial,
		CertType:        ssh.UserCert,
		KeyId:           utils.GetRandomString(10),
		ValidPrincipals: nil,
		ValidAfter:      uint64(time.Now().UTC().Unix()),
		ValidBefore:     uint64(time.Now().UTC().Add(time.Minute * 5).Unix()),
		Permissions: ssh.Permissions{
			Extensions: extentions,
		},
	}

	err = cert.SignCert(rand.Reader, caKey)
	if err != nil {
		logrus.Errorf(`could not sign public key: %v`, err)
		return
	}
	//
	//err = dbstore.Connect.SavePublicKey(userID, strings.TrimSpace(string(publicKeyBytes)))
	//if err != nil {
	//	logrus.Errorf(`could not save public key: %v`, err)
	//	return
	//}

	privateKeyBytes = utils.EncodePrivateKeyToPEM(privateKey)
	certBytes = ssh.MarshalAuthorizedKey(&cert)
	if len(certBytes) == 0 {
		logrus.Errorf("failed to marshal signed certificate, empty result")
		err = errors.New("failed to marshal signed certificate, empty result")
		return
	}

	// Create a buffer to write our archive to.
	//buffer := new(bytes.Buffer)

	return privateKeyBytes, publicKeyBytes, certBytes, nil

}

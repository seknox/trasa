package auth

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/mssola/user_agent"
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

type registerDeviceReq struct {
	TfaMethod string `json:"tfaMethod"`
	TotpCode  string `json:"totpCode"`
	TrasaID   string `json:"trasaID"`
	OrgID     string `json:"orgID"`
	// Device name to be used before we decrypt deviceHygiene
	DeviceName    string `json:"deviceName"`
	DeviceHygiene string `json:"deviceHygiene"`
}

type deviceDetail struct {
	DeviceBrowser     models.DeviceBrowser       `json:"deviceBrowser"`
	BrowserExtensions []models.BrowserExtensions `json:"browserExtensions"`
	DeviceHygiene     models.DeviceHygiene       `json:"deviceHygiene"`
}

// RegisterUserDevice registers new user device and stores device hygiene, device browser and browser extensions details.
func RegisterUserDevice(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("RegisterUserDevice request received")

	var req registerDeviceReq
	err := utils.ParseAndValidateRequest(r, &req)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "success", "invalid request", "RegisterUserDevice", nil)
		return
	}

	// get user info from database
	userDetailFromDB, err := Store.GetLoginDetails(req.TrasaID, req.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "User not found", "Ext login", nil)
		return
	}

	switch req.TfaMethod {
	case "totp":
		check, _, err := tfa.VerifyTotpCode(req.TotpCode, userDetailFromDB.ID, userDetailFromDB.OrgID)
		if err != nil || !check {
			// TODO @bhrg3se auth log?
			//err := logs.Store.LogLogin(authlog, consts.REASON_INVALID_TOTP, false)
			utils.TrasaResponse(w, 200, "failed", "totp failed", "Device registration", nil)
			return
		}
	default:
		status, msg := tfa.SendU2F(userDetailFromDB.ID, userDetailFromDB.OrgID, fmt.Sprintf("Device registration: %s", req.DeviceName), utils.GetIp(r))
		if !status {
			// TODO @bhrg3se auth log?
			//err := logs.Store.LogLogin(authlog, consts.REASON_U2F_FAILED, false)
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
		logrus.Trace("key not found in Kex store")
		utils.TrasaResponse(w, 200, "failed", "key not found in Kex store", "Device registration", nil)
		return
	}

	dhBytes, err := hex.DecodeString(req.DeviceHygiene)
	if err != nil {
		logrus.Debug("cannot decode hex string", err)
		utils.TrasaResponse(w, 200, "failed", "failed to decrypt data", "Device registration", nil)
		return
	}

	// decrypt the device details.
	decryptedBytes, err := utils.AESDecrypt(secretKeyFromKex.Secretkey[:], dhBytes)
	if err != nil {
		logrus.Debug(err)
		utils.TrasaResponse(w, 200, "failed", "failed to decrypt data", "Device registration", nil)
		return
	}

	//delete secret key from store
	delete(global.ECDHKexDerivedKey, req.TrasaID)

	// unmarshall decryptedBytes to deviceDetail struct
	var dh deviceDetail
	err = json.Unmarshal(decryptedBytes, &dh)
	if err != nil {
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

		dh.DeviceHygiene.NetworkInfo.IPAddress = utils.GetIp(r)
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

type enrolExt struct {
	Email          string `json:"email"`
	TfaMethod      string `json:"tfaMethod"`
	TotpCode       string `json:"totpCode"`
	UA             string `json:"UA"`
	IP             string `json:"IP"`
	AddedAt        string `json:"addedAt"`
	OrgID          string `json:"orgID"`
	BrowserName    string `json:"browserName"`
	BrowserVersion string `json:"browserVersion"`
	OSName         string `json:"osName"`
	OSVersion      string `json:"osVersion"`
}

// EnrolBrowserExtension is used to enrol web browser (via web extension)
func EnrolBrowserExtension(w http.ResponseWriter, r *http.Request) {

	var extLogin enrolExt
	var registerDevice models.UserDevice

	if err := json.NewDecoder(r.Body).Decode(&extLogin); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request format", "ExtLogin", nil)
		return
	}

	// get user info from database
	userDetailFromDB, err := Store.GetLoginDetails(extLogin.Email, extLogin.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "User not found", "Ext login", nil)
		return
	}

	//if len(userArr) > 1 {
	//	var tempOrgs []models.Org = make([]models.Org, 0)
	//	for _, tempUser := range userArr {
	//		var tempOrg models.Org
	//		tempOrg.ID = tempUser.OrgID
	//		tempOrg.OrgName = tempUser.OrgName
	//		tempOrgs = append(tempOrgs, tempOrg)
	//	}
	//	utils.TrasaResponse(w, 200, "selectOrg", "got multiple accounts", "Ext login", nil, tempOrgs)
	//	return
	//}
	//
	//userDetailFromDB := userArr[0]

	switch extLogin.TfaMethod {
	case "totp":
		check, _, err := tfa.VerifyTotpCode(extLogin.TotpCode, userDetailFromDB.ID, userDetailFromDB.OrgID)
		if err != nil || !check {
			// TODO @bhrg3se auth log?
			//err := logs.Store.LogLogin(authlog, consts.REASON_INVALID_TOTP, false)
			utils.TrasaResponse(w, 200, "failed", "totp failed", "Browser registration", nil)
			return
		}
	default:
		status, msg := tfa.SendU2F(userDetailFromDB.ID, userDetailFromDB.OrgID, fmt.Sprintf("Browser registration: %s", extLogin.BrowserName), utils.GetIp(r))
		if !status {
			// TODO @bhrg3se auth log?
			//err := logs.Store.LogLogin(authlog, consts.REASON_U2F_FAILED, false)
			utils.TrasaResponse(w, 200, "failed", msg, "Browser registration", nil)
			return
		}
	}

	//var uaFields utils.DeviceFinger
	//TODO ua fields are empty
	ua := user_agent.New(extLogin.UA)

	// brsrName, brsrVer := ua.Browser()

	// ngn, ngnVer := ua.Engine()

	registerDevice.DeviceID = utils.GetUUID()
	registerDevice.UserID = userDetailFromDB.ID
	registerDevice.OrgID = userDetailFromDB.OrgID
	registerDevice.DeviceType = "browser"
	registerDevice.FcmToken = ""
	registerDevice.PublicKey = ""
	registerDevice.DeviceFinger = "{}"
	registerDevice.DeviceHygiene = models.DeviceHygiene{
		DeviceInfo: models.DeviceInfo{},
		DeviceOS: models.DeviceOS{
			OSName:     ua.OSInfo().Name,
			OSVersion:  ua.OSInfo().Version,
			KernelType: "",
		},
		// DeviceBrowser: models.DeviceBrowser{
		// 	BrowserName:    brsrName,
		// 	BrowserVersion: brsrVer,
		// 	EngineName:     ngn,
		// 	EngineVersion:  ngnVer,
		// 	IsBot:          ua.Bot(),
		// 	UserAgent:      extLogin.UA,
		// },
		LoginSecurity: models.LoginSecurity{},
		NetworkInfo: models.NetworkInfo{
			IPAddress: utils.GetIp(r),
		},
		EndpointSecurity: models.EndpointSecurity{},
		LastCheckedTime:  time.Now().Unix(),
	}
	//registerDevice.Brand = extLogin.BrowserName
	//registerDevice.OsName = ua.OSInfo().Name
	//registerDevice.OsVersion = ua.OSInfo().Version
	registerDevice.AddedAt = time.Now().Unix() //.In(nep).String()

	err = devices.Store.Register(registerDevice)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not register device", "ExtLogin", nil)
		return
	}

	type extResponse struct {
		Token      string   `json:"token"`
		RootDomain string   `json:"rootDomain"`
		SsoDomain  string   `json:"ssoDomain"`
		WSPath     string   `json:"wsPath"`
		Hosts      []string `json:"hosts"`
		TrasaDACom bool     `json:"trasaDACom"`
	}

	var resp extResponse
	resp.Token = registerDevice.DeviceID
	resp.RootDomain = global.GetConfig().Trasa.Rootdomain
	resp.SsoDomain = global.GetConfig().Trasa.Ssodomain
	u, err := url.Parse(global.GetConfig().Trasa.Dashboard)
	if err != nil {
		logrus.Error(err)
	}
	resp.WSPath = fmt.Sprintf("wss://%s", u.Host)

	allservices, err := services.Store.GetAllByType("http", userDetailFromDB.OrgID)
	if err != nil {
		utils.TrasaResponse(w, 200, "failed", "No http services available", "SyncExtension", nil)
		return
	}

	for _, v := range allservices {
		resp.Hosts = append(resp.Hosts, v.Hostname)
	}

	globalDeviceCheck, err := system.Store.GetGlobalSetting(userDetailFromDB.OrgID, consts.GLOBAL_DEVICE_HYGIENE_CHECK)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error fetching global settings", "GlobalSettings", nil)
		return
	}

	resp.TrasaDACom = globalDeviceCheck.Status

	utils.TrasaResponse(w, 200, "success", "authorized and token garnted", "ExtLogin", resp)

}

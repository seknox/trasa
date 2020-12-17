package auth

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/seknox/trasa/server/api/auth/tfa"
	"github.com/seknox/trasa/server/api/devices"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"github.com/tstranex/u2f"
)

type TfaRequest struct {
	// Token is unique to tfarequest which is tied to specific user profile
	Token string `json:"token"`
	// TfaMethod can be u2f, totp or htoken
	TfaMethod string `json:"tfaMethod"`
	// Totp is value if TfaMethod is totp. otherwise it is nil.
	Totp string `json:"totpCode"`
	// Intent identifies where the tfa request is in context. Value can be login, forget password, appconnect.
	Intent          string `json:"intent"`
	HTTPProxyDomain string `json:"httpProxyDomain"`
	ExtID           string `json:"extID"`
	DeviceHygiene   string `json:"deviceHygiene"`
	ClientPubKey    string `json:"clientPubKey"`
}

// TfaHandler handles two factor authentication from TRASA ui
func TfaHandler(w http.ResponseWriter, r *http.Request) {

	var req TfaRequest
	//var service dbstore.App
	//var respBody userLoginStatus

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	authlog := logs.NewLog(r, "dashboard")

	// check if request intent matches on of supported one.
	if getIntentMatch(req.Intent) == false {
		utils.TrasaResponse(w, 200, "failed", "unsupported login", "TfaHandler", nil)
		return
	}

	// query redis and verify token based on intent (same intent should be attached during creds validation. )
	orgUserStr, err := redis.Store.Get(req.Token, "orgUser")

	if err != nil || orgUserStr == "" {
		err := logs.Store.LogLogin(&authlog, consts.REASON_INVALID_TOKEN, false)
		if err != nil {
			logrus.Error(err)
		}

		utils.TrasaResponse(w, 200, "failed", "invalid token", "Dashboard Login", nil)
		return
	}

	intent, err := redis.Store.Get(req.Token, "intent")

	if err != nil || intent == "" {
		err := logs.Store.LogLogin(&authlog, consts.REASON_SPOOFED_LOGIN, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", "Spoofed login", "Dashboard Login", nil)
		return
	}

	orgUser := strings.Split(orgUserStr, ":")

	if len(orgUser) != 2 {
		logrus.Error("orgUser is invalid")
		err := logs.Store.LogLogin(&authlog, consts.REASON_UNKNOWN, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", "Something is wrong", "Dashboard Login", nil)
		return
	}

	//fmt.Println("orgUser: ", orgUser[0], orgUser[1])

	userDetails, err := users.Store.GetFromID(orgUser[1], orgUser[0])
	if err != nil {
		logrus.Error(err)

		err := logs.Store.LogLogin(&authlog, consts.REASON_UNKNOWN, false)
		if err != nil {
			logrus.Error(err)
		}

		utils.TrasaResponse(w, 200, "failed", "unable to verify user", "Dashboard Login")
		return
	}

	//Update authlog value with user fields
	authlog.UpdateUser(&models.UserWithPass{User: *userDetails})

	// decrypt and update hygiene
	globalDeviceCheck, err := system.Store.GetGlobalSetting(userDetails.OrgID, consts.GLOBAL_DEVICE_HYGIENE_CHECK)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error fetching global settings", "SyncExtension", nil)
		return
	}
	deviceID := ""
	// only update device hygiene if setting is true
	if globalDeviceCheck.Status == true {

		//if req.Intent != consts.AUTH_REQ_TFA_DH {
		//	utils.TrasaResponse(w, 200, "failed", "Device Hygiene Required", "Tfa", nil)
		//	return
		//}

		privKey, err := redis.Store.Get(req.Token, "priv")

		if err != nil || privKey == "" {
			err := logs.Store.LogLogin(&authlog, consts.REASON_INVALID_TOKEN, false)
			if err != nil {
				logrus.Error(err)
			}

			utils.TrasaResponse(w, 200, "failed", "invalid token", "Dashboard Login", nil)
			return
		}

		deviceID, err = decryptAndUpdateDH(privKey, req.ClientPubKey, req.DeviceHygiene, req.ExtID)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", err.Error(), "Dashboard Login", nil)
			return
		}
	}

	status, reason, response := handleTFAMethod(req, userDetails, &authlog)
	if status != "success" || req.TfaMethod == "u2fy" {
		utils.TrasaResponse(w, 200, status, reason, "Dashboard Login", response)
		return
	}

	var uc models.UserContext
	uc.User = userDetails

	org, err := orgs.Store.Get(userDetails.OrgID)
	if err != nil {
		utils.TrasaResponse(w, 200, status, reason, "Dashboard Login", response)
		return
	}

	uc.Org = org
	uc.DeviceID = deviceID
	uc.BrowserID = req.ExtID

	status, failedReason, intent, sessionToken, respData := handleIntentResponse(req, uc)

	err = logs.Store.LogLogin(&authlog, failedReason, status == "success")
	if err != nil {
		logrus.Error(err)
	}

	if req.Intent == consts.AUTH_REQ_DASH_LOGIN {
		// we set session token in HTTPonly cookie and expect csrf token in http header.
		xSESSION := http.Cookie{
			Name:     "X-SESSION",
			Value:    sessionToken,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Secure:   true,
			Path:     "/",
		}

		http.SetCookie(w, &xSESSION)
	}

	utils.TrasaResponse(w, 200, status, reason, intent, respData)
	return

}

type tfaSign struct {
	SignResp  u2f.SignResponse `json:"signResp"`
	Challenge string           `json:"challenge"`
	// Intent identifies where the tfa request is in context. Value can be login, forget password, appconnect.
	Intent string `json:"intent"`
}

// AUTH_REQ_DASH_LOGIN   = "AUTH_REQ_DASH_LOGIN"
// AUTH_REQ_CHANGE_PASS  = "AUTH_REQ_CHANGE_PASS"
// AUTH_REQ_ENROL_DEVICE = "AUTH_REQ_ENROL_DEVICE"
// AUTH_REQ_FORGOT_PASS  = "AUTH_REQ_FORGOT_PASS"
// AUTH_HTTP_ACCESS_PROXY = "AUTH_HTTP_ACCESS_PROXY"
// AUTH_REQ_TFA_DH = "AUTH_REQ_TFA_DH"
func getIntentMatch(intent string) bool {
	retVal := false
	availableIntents := []string{consts.AUTH_REQ_DASH_LOGIN, consts.AUTH_REQ_CHANGE_PASS, consts.AUTH_REQ_ENROL_DEVICE, consts.AUTH_REQ_FORGOT_PASS, consts.AUTH_HTTP_ACCESS_PROXY, consts.AUTH_REQ_TFA_DH}
	for _, v := range availableIntents {
		if intent == v {
			retVal = true
		}
	}

	return retVal
}

func handleTFAMethod(req TfaRequest, user *models.User, authlog *logs.AuthLog) (status, reason string, resp interface{}) {
	switch req.TfaMethod {
	// in case of u2fy, we do not need to generate login credentials here but process it in another signed response request from client
	case "u2fy":
		webSignrequest, err := tfa.SignReq(user.ID, user.OrgID)
		if err != nil {

			err := logs.Store.LogLogin(authlog, consts.REASON_UNKNOWN, false)
			if err != nil {
				logrus.Error(err)
			}
			return "failed", "unable to verify user", nil
		}

		return "success", "proceed signing request", webSignrequest

	case "totp":
		check, tfaDeviceID, err := tfa.VerifyTotpCode(req.Totp, user.ID, user.OrgID)
		authlog.TfaDeviceID = tfaDeviceID
		if err != nil || !check {
			err := logs.Store.LogLogin(authlog, consts.REASON_INVALID_TOTP, false)
			if err != nil {
				logrus.Error(err)
			}

			return "failed", "TOTP failed", nil
		}
		return "success", "", nil

	default:
		status, msg := tfa.SendU2F(user.ID, user.OrgID, "Dashboard", authlog.UserIP)
		if !status {
			err := logs.Store.LogLogin(authlog, consts.REASON_U2F_FAILED, false)
			if err != nil {
				logrus.Error(err)
			}
			return "failed", msg, nil
		}
		return "success", "", nil

	}

}

func handleIntentResponse(req TfaRequest, uc models.UserContext) (status string, reason consts.FailedReason, intent, sessionToken string, resp interface{}) {
	orgUserStr := fmt.Sprintf("%s:%s", uc.User.OrgID, uc.User.ID)
	switch req.Intent {
	// in case of u2fy, we do not need to generate login credentials here but process it in another signed response request from client
	case consts.AUTH_REQ_DASH_LOGIN:

		sessionToken, resp, err := sessionResponse(uc)
		if err != nil {
			return "failed", consts.REASON_TRASA_ERROR, "DashboardLogin", sessionToken, nil
		}
		return "success", "", "DashboardLogin", sessionToken, resp

	case consts.AUTH_REQ_ENROL_DEVICE:
		//todo this is a temporary fix
		userWithPass, err := Store.GetLoginDetails(uc.User.UserName, "")
		if err != nil {
			logrus.Error(err)
			return "failed", consts.REASON_USER_NOT_FOUND, "DashboardLogin", "", ""
		}
		resp := devices.EnrolDeviceFunc(*uc.User)
		resp.OrgName = userWithPass.OrgName
		resp.Account = userWithPass.Email
		if resp.Account == "" {
			resp.Account = userWithPass.UserName
		}
		return "success", "", consts.AUTH_RESP_ENROL_DEVICE, "", resp
	case consts.AUTH_REQ_CHANGE_PASS:
		verifyToken := utils.GetRandomString(7)
		// store token to redis
		err := redis.Store.Set(verifyToken,
			consts.TOKEN_EXPIRY_CHANGEPASS,
			"orguser", orgUserStr,
			"intent", string(consts.VERIFY_TOKEN_CHANGEPASS),
			"createdAt", time.Now().String())

		if err != nil {
			logrus.Error(err)
			return "failed", consts.REASON_TRASA_ERROR, consts.AUTH_RESP_CHANGE_PASS, "", verifyToken
		}
		return "success", "", consts.AUTH_RESP_CHANGE_PASS, "", verifyToken

	case consts.AUTH_REQ_FORGOT_PASS:
		err := forgotPassTfaResp(*uc.User)
		if err != nil {
			logrus.Error(err)
			return "failed", consts.REASON_TRASA_ERROR, consts.AUTH_RESP_FORGOT_PASS, "", nil
		}
		return "success", "", consts.AUTH_RESP_FORGOT_PASS, "", nil
	default:
		return "failed", "default", "DashboardLogin", "", nil
	}

}

func decryptAndUpdateDH(ourPriv, clientPub, clientDH, extID string) (string, error) {
	var privBytes [32]byte
	var pubBytes [32]byte

	privFromHexStr, err := hex.DecodeString(ourPriv)
	if err != nil {
		return "", fmt.Errorf("privFromHexStr: %v ", err)
	}

	pubFromHexStr, err := hex.DecodeString(clientPub)
	if err != nil {
		return "", fmt.Errorf("pubFromHexStr: %v ", err)
	}

	copy(privBytes[:], privFromHexStr)
	copy(pubBytes[:], pubFromHexStr)

	dhBytes, err := hex.DecodeString(clientDH)
	if err != nil {
		return "", fmt.Errorf("dhBytes: %v ", err)
	}

	sec := utils.ECDHComputeSecret(&privBytes, &pubBytes)

	plainText, err := utils.AESDecrypt(sec, dhBytes)

	var dh DeviceDetail
	err = json.Unmarshal(plainText, &dh)
	if err != nil {
		return "", fmt.Errorf("json.Unmarshal(plainText, &dh): %v ", err)
	}

	// get deviceID, orgID from extID
	orgID, deviceID, _, err := devices.Store.GetDeviceAndOrgIDFromExtID(extID)
	if err != nil {
		return deviceID, fmt.Errorf("GetDeviceAndOrgIDFromExtID: %v ", err)
	}

	// begin update process
	err = devices.Store.UpdateWorkstationHygiene(dh.DeviceHygiene, deviceID, orgID)
	if err != nil {
		return deviceID, fmt.Errorf("updateWorkstationHygiene: %v ", err)
	}

	// begin update process
	err = devices.Store.UpdateBrowserHygiene(dh.DeviceBrowser, extID, orgID)
	if err != nil {
		return deviceID, fmt.Errorf("updateBrowserHygiene: %v ", err)
	}

	return deviceID, nil
}

type ConfirmTOTPPreq struct {
	TOTPCode string `json:"totpCode"`
	DeviceID string `json:"deviceID"`
}

//Check newly added TOTP to complete device registration process.
//This function will also create http session
func ConfirmTOTPAndSave(w http.ResponseWriter, r *http.Request) {
	var request ConfirmTOTPPreq
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "ConfirmTOTPAndSave")
		return
	}

	authlog := logs.NewLog(r, "dashboard")

	userID_orgID_Totpsec, err := redis.Store.MGet(request.DeviceID,
		"userID",
		"orgID",
		"totpSec",
	)

	if err != nil {
		logrus.Error(err)
		err := logs.Store.LogLogin(&authlog, consts.REASON_INVALID_TOKEN, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", "invalid deviceID", "ConfirmTOTPAndSave")
		return
	}

	userID := userID_orgID_Totpsec[0]
	orgID := userID_orgID_Totpsec[1]
	totpSec := userID_orgID_Totpsec[2]

	authlog.UserID = userID
	authlog.OrgID = orgID

	prevCode, nowCode, nextCode := utils.CalculateTotp(totpSec)
	if request.TOTPCode != prevCode && request.TOTPCode != nowCode && request.TOTPCode != nextCode {
		logrus.Error("invalid TOTP code")
		err := logs.Store.LogLogin(&authlog, consts.REASON_INVALID_TOTP, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", "invalid TOTP code", "ConfirmTOTPAndSave")
		return
	}

	dev := models.UserDevice{
		UserID:     userID,
		OrgID:      orgID,
		DeviceID:   request.DeviceID,
		MachineID:  "",
		DeviceType: "mobile",
		TotpSec:    totpSec,
		Trusted:    false,
		AddedAt:    time.Now().Unix(),
	}

	fcm_publick_devHyg, err := redis.Store.MGet(request.DeviceID,
		"fcmToken",
		"publicKey",
		"deviceHygiene",
	)

	if err == nil {
		var devHyg models.DeviceHygiene
		err = json.Unmarshal([]byte(fcm_publick_devHyg[2]), &devHyg)
		if err != nil {
			logrus.Error(err)
		}

		dev.FcmToken = fcm_publick_devHyg[0]
		dev.PublicKey = fcm_publick_devHyg[1]
		dev.DeviceHygiene = devHyg

	}

	userDetails, err := users.Store.GetFromID(userID, orgID)
	if err != nil {
		logrus.Error(err)
		err := logs.Store.LogLogin(&authlog, consts.REASON_USER_NOT_FOUND, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", "could not find user", "ConfirmTOTPAndSave")
		return
	}

	authlog.UpdateUser(&models.UserWithPass{User: *userDetails})

	err = devices.Store.Register(dev)
	if err != nil {
		logrus.Error(err)
		err := logs.Store.LogLogin(&authlog, consts.REASON_UNKNOWN, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", "could not register device", "ConfirmTOTPAndSave")
		return
	}

	var uc models.UserContext
	uc.User = userDetails

	org, err := orgs.Store.Get(userDetails.OrgID)
	if err != nil {
		logrus.Error(err)
		err := logs.Store.LogLogin(&authlog, consts.REASON_ORG_NOT_FOUND, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", err.Error(), "could not register device", "ConfirmTOTPAndSave")
		return
	}

	uc.Org = org
	uc.DeviceID = ""
	uc.BrowserID = ""

	//TODO add deviceID and browserID
	sessionToken, resp, err := sessionResponse(uc)
	if err != nil {
		logrus.Error(err)
		err := logs.Store.LogLogin(&authlog, consts.REASON_FAILED_TO_GENERATE_TOKEN, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", "could not get session", "ConfirmTOTPAndSave")
		return
	}

	xSESSION := http.Cookie{
		Name:     "X-SESSION",
		Value:    sessionToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
	}

	http.SetCookie(w, &xSESSION)
	err = logs.Store.LogLogin(&authlog, "", true)
	if err != nil {
		logrus.Error(err)
	}
	utils.TrasaResponse(w, 200, "success", "", "", resp)
}

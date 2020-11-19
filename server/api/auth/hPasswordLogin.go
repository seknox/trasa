package auth

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/seknox/trasa/server/global"

	uuid "github.com/satori/go.uuid"
	"github.com/seknox/trasa/server/api/devices"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/providers/uidp"
	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/api/users/passwordpolicy"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	OrgID     string `json:"orgId"`
	UserID    string `json:"userId"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	TfaMethod string `json:"tfaMethod"`
	Totp      string `json:"totp"`
	PublicKey []byte `json:"publicKey"`
	DeviceID  string `json:"deviceID"`
	IdpName   string `json:"idpName"`
	Intent    string `json:"intent"`
}

// LoginHandler authenticates user for configured identity provider.
// successful authentication should respond with tfarequired intent. If user has not enrolled any 2fa device,
// this handler should respond with enroll device intent.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest

	if err := utils.ParseAndValidateRequest(r, &loginRequest); err != nil {
		http.Error(w, http.StatusText(400), 200)
		return
	}

	authlog := logs.NewLog(r, "dashboard")

	// get user info from database
	userDetails, err := Store.GetLoginDetails(loginRequest.Email, loginRequest.OrgID)
	if err != nil {
		logrus.Error(err)
		err = logs.Store.LogLogin(&authlog, consts.REASON_USER_NOT_FOUND, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", "User not found", "LoginHandler")
		return
	}

	authlog.UpdateUser(userDetails)

	if !userDetails.Status {
		err = logs.Store.LogLogin(&authlog, consts.REASON_USER_DISABLED, false)
		if err != nil {
			logrus.Error(err)
		}

		utils.TrasaResponse(w, 200, "failed", "User Disabled", "Dashboard Login")
		return
	}

	reason, err := CheckPassword(userDetails, loginRequest.Email, loginRequest.Password)
	if err != nil {
		logrus.Error(err)
		err = logs.Store.LogLogin(&authlog, reason, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", "invalid username or password", "Dashboard Login")
		return
	}

	// We set "changeme" as a default passsword for root account during installation. This password should be changed in first login
	// and should never be used within TRASA accounts. If password validation succeeds and password value is "changeme" we will force
	// user to change password.

	if loginRequest.Password == "changeme" {
		// we now will generate a short lived token which will be sent in following success response.
		// this will let the user setup password only when token is validated.
		verifyToken, err := redis.SetVerifyToken(userDetails.OrgID, userDetails.ID)

		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "Could not set change pass token", "set token", nil)
			return
		}

		utils.TrasaResponse(w, 200, "success", "change default password", consts.AUTH_RESP_RESET_PASS, verifyToken)
		return
	}

	// if we are here, it means username password validation succeeded.
	// 1st) lets check is the user password has expired. If true, we set enforce password change policy here.
	// @@@ pending TODO
	if userDetails.IdpName == "trasa" {
		check, err := passwordpolicy.CheckPendingPasswordRotationForUser(userDetails.ID, userDetails.OrgID)
		if err != nil {
			logrus.Error(err)
		}

		if check {
			// enforce password change policy for this user.
			err := passwordpolicy.EnforceChangePassword(userDetails.ID, userDetails.OrgID)
			if err != nil {
				logrus.Debug(err)
			}
		}
	}

	// intent tells client helps identify if tfa request is for login, or forgotpassword.
	// This also tells client(dashboard) to either proceed with tfa or enroll mobile device.

	enrolled, err := checkDeviceEnrolled(userDetails.ID, userDetails.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "enrol device", "")
		return
	}

	if !enrolled {
		resp := devices.EnrolDeviceFunc(models.CopyUserWithoutPass(*userDetails))
		resp.OrgName = userDetails.OrgName
		utils.TrasaResponse(w, http.StatusOK, "success", "enrol device", consts.AUTH_RESP_ENROL_DEVICE, resp)
		return
	}

	// Now we will generate a unique token attached to this user. and send this token in response.
	// this token will be used to validate tfa request and retrieve userID for this user.
	// this userID will be used to retrieve user detail and send in response to successful authentication.
	// @sshahcodes 18 dec, 2019. changing second parameter from hardcoded device-id to login. This will ensure that subsequent tfa requests are tied to this login authentication event.
	// token := utils.GetRandomString(15) // using public key for token

	priv, pub, err := utils.ECDHGenKeyPair()
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to generate key pair", "Kex", nil)
		return
	}

	orgUser := fmt.Sprintf("%s:%s", userDetails.OrgID, userDetails.ID)

	// If user has already enroled 2fa device, proceed with tfa process.
	err = redis.Store.Set(hex.EncodeToString(pub[:]), time.Second*400, "orgUser", orgUser, "intent", loginRequest.Intent, "priv", hex.EncodeToString(priv[:]))
	if err != nil {
		err = logs.Store.LogLogin(&authlog, reason, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, http.StatusOK, "failed", "failed setting token for polling", "", "")
		return
	}

	respIntent := consts.AUTH_RESP_TFA_REQUIRED
	globalDeviceCheck, err := system.Store.GetGlobalSetting(userDetails.OrgID, consts.GLOBAL_DEVICE_HYGIENE_CHECK)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error fetching global settings", "SyncExtension", nil)
		return
	}
	if globalDeviceCheck.Status == true {
		respIntent = consts.AUTH_RESP_TFA_DH_REQUIRED
	}

	// return AUTH_RESP_TFA_DH_REQUIRED inten with both server public key as tfa Token.
	// At this point, we shold now only return TfaRequired response.
	utils.TrasaResponse(w, http.StatusOK, "success", "second step verification required", respIntent, hex.EncodeToString(pub[:]))

}

type EnrolDeviceStruct struct {
	DeviceID      string `json:"deviceID"`
	TotpSSC       string `json:"totpSSC"`
	OrgName       string `json:"orgName"`
	CloudProxyURL string `json:"cloudProxyURL"`
}

// Enrol2FADevice primary function is to enrol user mobile device for 2FA. While previously this function was used to enrol U2F only part,
// this function now also syncs totp shared secret key for user for their particular organization.
// This is a Four step process. 1) handle user login, 2) generate device, get totpssc ID 3) send GetDeviceDetail Request to trasa cloud 4) respond with device ID and otpauth url.
func Enrol2FADevice(w http.ResponseWriter, r *http.Request) {
	// remote auth request
	var req LoginRequest
	// user struct
	var getUser models.User

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		utils.TrasaResponse(w, 200, "failed", "error parsing request", "Device Enroll", nil)
		return
	}

	authlog := logs.NewLog(r, "dashboard")

	// 1) validate uname pass
	userDetailWithPass, err := Store.GetLoginDetails(req.Email, req.OrgID)
	if err != nil {
		logrus.Error(err)
		err = logs.Store.LogLogin(&authlog, consts.REASON_USER_NOT_FOUND, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", "User not found", "Device Enroll", nil)
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
	//	utils.TrasaResponse(w, 200, "selectOrg", "got multiple users", "EnrolDevice", tempOrgs)
	//	return
	//}
	//userDetailWithPass := userArr[0]
	authlog.UpdateUser(userDetailWithPass)
	//	userDetails := utils.CopyUserWithoutPass(userDetailWithPass)
	// check password hash
	reason, err := CheckPassword(userDetailWithPass, req.Email, req.Password)

	if err != nil {

		err = logs.Store.LogLogin(&authlog, reason, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", "Invalid creds", "Device Enroll")
		return

	} else {
		// get user device from database
		//var userDevice dbstore.UserDevice

		getUser.Email = req.Email

		orguser := fmt.Sprintf("%s:%s", userDetailWithPass.OrgID, userDetailWithPass.ID)
		// (2) Generate Device ID
		deviceID, _ := uuid.NewV4()

		// (3) Send GiveMeDeviceDetail request.
		totpSec := utils.GenerateTotpSecret()

		respVal := EnrolDeviceStruct{
			DeviceID:      deviceID.String(),
			TotpSSC:       totpSec,
			OrgName:       userDetailWithPass.OrgName,
			CloudProxyURL: global.GetConfig().Trasa.CloudServer,
		}

		utils.TrasaResponse(w, 200, "success", "", "EnrolDevice", respVal)

		go devices.GiveMeDeviceDetail(orguser, deviceID.String(), totpSec)
	}

}

func CheckPassword(userDetails *models.UserWithPass, email, password string) (reason consts.FailedReason, err error) {
	// check if login request is for other IDP. FreeIpa Ldap for now.
	// TODO make idp name constant value?? Also if you change here, do not forget to change same in dashboard login intent

	if userDetails.IdpName != "trasa" {

		// get idpDetail
		idp, err := orgs.Store.GetIDP(userDetails.OrgID, userDetails.IdpName)
		if err != nil {
			return consts.REASON_IDENTITY_PROVIDER_NOT_FOUND, err
		}

		//  Authenticate user aka  Bind user

		fullUserPath := fmt.Sprintf("uid=%s,%s", email, idp.IdpMeta)
		if userDetails.IdpName == "ad" {
			fullUserPath = fmt.Sprintf("CN=%s, %s", email, idp.IdpMeta)
		}

		err = uidp.BindLdap(fullUserPath, password, idp.Endpoint)
		if err != nil {

			return consts.REASON_LDAP_AUTH_FAILED, err

		}

		return "", nil

	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(password))
	if err != nil {
		return consts.REASON_INVALID_USER_CREDS, err
	}
	return "", nil

}

// checkDeviceEnrolled checks if user has any prior mobile devices enrolled with TRASA. if device is found
// TRASA should proceed with login tfa flow, else should perform enroll device.
// @bhrg3se check if below logic is okay to check if user has already registered mobile device
func checkDeviceEnrolled(userID, orgID string) (bool, error) {
	// get user device detail. if user has not enrolled any mobile device, send with enrol device response.
	device, err := users.Store.GetTOTPDevices(userID, orgID)
	if err != nil {
		return false, err
	}

	if len(device) == 0 {
		return false, nil
	}

	return true, nil
}

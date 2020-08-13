package tfa

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/api/devices"
	"github.com/seknox/trasa/server/api/notif"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"github.com/tstranex/u2f"
)

//returns check,deviceID and error
func VerifyTotpCode(totpCode, userID, orgID string) (bool, string, error) {
	userDevices, err := users.Store.GetTOTPDevices(userID, orgID)
	if err != nil {
		logrus.Error(err)
		return false, "", err
	}

	for _, dev := range userDevices {
		totpThen, totpNow, totpAfter := utils.CalculateTotp(dev.TotpSec)
		if totpAfter == totpCode || totpNow == totpCode || totpThen == totpCode {
			return true, dev.DeviceID, nil
		}
		//logrus.Trace(totpCode," ",totpThen," ", totpNow," ", totpAfter)
	}
	return false, "", nil
}

//Deprecated
//Use HandleTfaAndGetDeviceID instead
func SendU2F(userID, orgID, appName, ip string) (bool, string) {

	// If code reaches here, its time to generate random challenge, store it in redis with 2 min timer, send remote notification to user.
	challenge := utils.GetRandomString(5)

	//fmt.Printf("challenge: %s\n", hex.EncodeToString(challenge))
	var userDevice models.UserDevice

	//userDetails := dbstore.Connect.GetUserDetail(getUser)
	userDevice.UserID = userID
	// Device ID is unknown here. how do we get one??
	//userDevice.DeviceID =

	// Retrieve fcm token. GetUserDevices returns slice of fcm tokens
	totpDevices, err := users.Store.GetTOTPDevices(userID, orgID)
	if err != nil {
		logrus.Error(err)
		return false, string(consts.REASON_DEVICE_NOT_ENROLLED)
	}
	//fmt.Printf("device fcm: %s\n", getDeviceFromDb.FcmToken)

	// here we iterate over every available fcm tokens
	timezone := "UTC"
	org, err := orgs.Store.Get(orgID)
	if err != nil || org.Timezone == "" {
		logrus.Error(err, "Timezone not defined")
	} else {
		timezone = org.Timezone
	}
	nep, err := time.LoadLocation(timezone)
	if err != nil {
		logrus.Error(err, "could not load location")
	}
	current := time.Now().In(nep)
	format := current.Format("Mon Jan _2 15:04:05 2006")

	//////////////////////////////////////////////////////////////////////
	///////////		Check if trasa is run in cloud or onpremise      /////
	//////////////////////////////////////////////////////////////////////
	platformBase := global.GetConfig().Platform.Base

	if platformBase == "private" {
		// unlike prior U2F handling here, now we just enumerate fcm tokens, store it in array and send this to cloud server to notify user.
		// we dont handle redis polling here now since status will depend on cloud server response.
		var fcmTokens []string
		for _, v := range totpDevices {
			// Send Notification to User Mobile App
			fcmTokens = append(fcmTokens, v.FcmToken)

		}
		//logrus.Trace(fcmTokens)

		tfaResp, err := sendNotificationThroughCloudProxy(fcmTokens, org.OrgName, appName, ip, format)

		if err != nil {
			logrus.Debug(err)
			return false, string(consts.REASON_U2F_FAILED)
		}

		// deviceDetail, err := devices.Store.GetFromID(tfaResp.DeviceID)
		// if err != nil {
		// 	return false, string(consts.REASON_DEVICE_NOT_ENROLLED)
		// }
		// ok, err := rsaVerify(tfaResp.Signature, tfaResp.Challenge, deviceDetail.PublicKey)
		// if err != nil || !ok {
		// 	logrus.Debug(err)
		// 	return false, string(consts.REASON_U2F_FAILED)
		// }

		//TODO Use HandleTfaAndGetDeviceID
		//_, err = devices.Store.UpdateDeviceHygiene(device.DeviceHygiene, orgID)
		//if err != nil {
		//	//TODO return false after service update is complete
		//	logrus.Debug(err)
		//}

		if tfaResp.Answer != "YES" {
			return false, string(consts.REASON_U2F_FAILED)
		}

		return true, ""
	} else {
		for _, v := range totpDevices {
			// Send Notification to User Mobile App
			err := notif.Store.SendPushNotification(v.FcmToken, org.OrgName, appName, ip, format, challenge)
			if err != nil {
				logrus.Tracef("failed notifying user: %v", err)
			}

		}

		//TODO for ugly workaround, we have escaped setting device id this time.

		//TODO @sshah check if deviceID is actually necessary
		err = redis.Store.Set(challenge, time.Second*400, "status", "false")
		//setU2f, err := dbstore.Connect.SetTokenForPolling(hex.EncodeToString(challenge), "status", "false", "device", "device-id")
		if err != nil {
			logrus.Errorf("error setting u2f: %v", err)
			return false, ""
		}

		//	fmt.Printf("successfully notified user \n")

		// Wait and verify U2f challenge
		ok, _ := redis.Store.WaitForStatusAndGet(challenge, "device")
		return ok, ""

	}
	return true, ""

}

/////////////////////////////////////////////////
//////U2F Notificatioins///////////
////////////////////////////
//TODO check and rename fields
type notifFederation struct {
	serviceID string `json:"serviceID"`
	AppName   string `json:"appSecret"`
	OrgName   string `json:"orgName"`
	IpAddress string `json:"IPAddress"`
	Time      string `json:"time"`
	//DeviceID  string   `json:"deviceId"`
	FcmTokens []string `json:"fcmTokens"`
}

func sendNotificationThroughCloudProxy(fcmTokens []string, orgName, appName, ipAddr, time string) (U2f, error) {
	var requestConfig notifFederation
	requestConfig.FcmTokens = fcmTokens

	//TODO check if this deviceID is needed

	//requestConfig.DeviceID = deviceID
	requestConfig.OrgName = orgName
	requestConfig.AppName = appName
	requestConfig.Time = time
	requestConfig.IpAddress = ipAddr

	//urlPath := "https:///onprem2fa"
	urlPath := global.GetConfig().Trasa.CloudServer + "/api/v3/onprem2fa"
	//urlPath := "http://localhost:3339" + "/api/v1/onprem2fa"

	insecureSkip := global.GetConfig().Security.InsecureSkipVerify
	result, err := utils.CallTrasaAPI(urlPath, requestConfig, insecureSkip)
	if err != nil {
		return U2f{}, err
	}

	if result.Status != "success" {
		return U2f{}, errors.New(result.Reason)
	}

	mar, err := json.Marshal(result.Data)
	if err != nil {
		logrus.Error(err)
	}
	var unMar []U2f
	err = json.Unmarshal(mar, &unMar)
	if err != nil {
		logrus.Trace("No device detail in 2fa response")
		return U2f{}, errors.Errorf("invalid 2fa response: %v", err)
	}

	//logrus.Tracef("device detail got from U2F: %v", unMar)

	if len(unMar) != 1 {
		logrus.Trace("No device detail in 2fa response")
		return U2f{}, errors.Errorf("no 2fa response: %v", err)
	}

	return unMar[0], nil
}

// WaitForStatusAndGet polls redis to get U2F challenge resopnse. when user allows U2f login, /remote/auth/u2f handler updates SetU2f
// storage with status=true value. true here means user successfully verified 2fa.
//func WaitForStatusAndGet(key string) bool {
//
//	//TODO @sshah check this logic
//	timeout := time.After(60 * time.Second)
//	tick := time.Tick(1000 * time.Millisecond)
//	//var ret string
//	//var err error
//	// Keep trying until we're timed out or got a result or got an error
//	for {
//		select {
//		// Got a timeout! fail with a timeout error
//		case <-timeout:
//			return false
//			// Got a tick, we should check on doSomething()
//		case <-tick:
//			status := redis.Store.Get(key, "status")
//			//ret, err := dbstore.Connect.GetU2f(key, status, device) //GetU2f(key, status)
//			if status == "true" {
//				return true
//			}
//		}
//
//	}
//
//}

//This function will handle all 2fa process,
// update device hygiene from u2f
// and return device ID of 2fa device
func HandleTfaAndGetDeviceID(signResponse *u2f.SignResponse, tfaMethod, totpCode, userID, clientIP, appName, timezone, orgName, orgID string) (deviceID string, reason consts.FailedReason, ok bool) {

	// First we check if tfaMethod is U2F or totp.
	if strings.ToLower(tfaMethod) == "u2fy" {
		deviceID, err := checkU2FY(signResponse, userID, orgID)
		if err != nil {
			logrus.Error(err)
			return "", consts.REASON_U2FY_FAILED, false
		}
		return deviceID, "", true
		//logrus.Tracef("VerifySignResponse error: %v", err)

	} else if strings.ToLower(tfaMethod) == "totp" {
		totpCheck, deviceID, err := VerifyTotpCode(totpCode, userID, orgID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", consts.REASON_DEVICE_NOT_ENROLLED, false
			} else {
				logrus.Error(err)
				logrus.Error(err, "")
				return "", consts.REASON_UNKNOWN, false
			}
		} else if !totpCheck {
			return deviceID, consts.REASON_INVALID_TOTP, false
		} else {
			return deviceID, "", true
		}
	}

	// If code reaches here, its time to generate random challenge, store it in redis with 2 min timer, send remote notification to user.

	// Retrieve fcm token. GetUserDevices returns slice of fcm tokens
	userDevices, err := users.Store.GetTOTPDevices(userID, orgID)
	if err != nil {
		logrus.Error(err, "could not get totp devices")
		return "", "", false
	}
	//fmt.Printf("device fcm: %s\n", getDeviceFromDb.FcmToken)

	// here we iterate over every available fcm tokens

	if timezone == "" {
		timezone = "UTC"
	}
	nep, err := time.LoadLocation(timezone)
	if err != nil {
		logrus.Error(err, "")
	}
	current := time.Now().In(nep)
	format := current.Format("Mon Jan _2 15:04:05 2006")

	//////////////////////////////////////////////////////////////////////
	///////////		Check if trasa is run in cloud or onpremise      /////
	//////////////////////////////////////////////////////////////////////
	platformBase := global.GetConfig().Platform.Base

	var tfaResp U2f

	if strings.Compare(platformBase, "private") == 0 {
		// unlike prior U2F handling here, now we just enumerate fcm tokens, store it in array and send this to cloud server to notify user.
		// we dont handle redis polling here now since status will depend on cloud server response.
		var fcmTokens []string
		for _, v := range userDevices {
			// Send Notification to User Mobile App
			fcmTokens = append(fcmTokens, v.FcmToken)

		}

		tfaResp, err = sendNotificationThroughCloudProxy(fcmTokens, orgName, appName, clientIP, format)
		if err != nil {
			logrus.Error(err)
			return deviceID, consts.REASON_U2F_FAILED, false
		}

	} else {
		challenge := utils.GetRandomString(5)

		c := make(chan U2f, 1)

		u2fChanMap[challenge] = c
		defer func() {
			close(c)
			delete(u2fChanMap, challenge)
		}()

		for _, v := range userDevices {
			// Send Notification to User Mobile App
			err := notif.Store.SendPushNotification(v.FcmToken, orgName, appName, clientIP, format, challenge)
			if err != nil {
				logrus.Tracef("failed notifying user: %v", err)
			}

		}

		timeoutChannel := time.Tick(time.Second * 60)

		select {
		case <-timeoutChannel:
			return deviceID, consts.REASON_U2F_FAILED, false
		case tfaResp = <-c:

		}

	}

	deviceDetail, err := devices.Store.GetFromID(tfaResp.DeviceID)
	if err != nil {
		logrus.Errorf("device not found: %v", err)
		return deviceID, consts.REASON_DEVICE_NOT_ENROLLED, false
	}
	ok, err = rsaVerify(tfaResp.Signature, tfaResp.Challenge, deviceDetail.PublicKey)
	if err != nil || !ok {
		logrus.Debug(err)
		return deviceID, consts.REASON_U2F_FAILED, false
	}

	var deviceHyg models.DeviceHygiene
	err = json.Unmarshal([]byte(tfaResp.DeviceInfo), &deviceHyg)
	if err != nil {
		logrus.Errorf("device detail from u2f not valid: %v", err)
		return deviceID, consts.REASON_U2F_FAILED, false
	}

	deviceID, err = devices.Store.UpdateDeviceHygiene(deviceHyg, orgID)
	if err != nil {
		//TODO return false after service update is complete
		logrus.Error(err)
	}

	if tfaResp.Answer != "YES" {
		return deviceID, consts.REASON_U2F_FAILED, false
	}
	return deviceID, "", true

}

func rsaVerify(signedChallenge string, originalChallenge string, publicKeyPEM string) (bool, error) {

	decodedSignedChallenge, decodeErr := base64.StdEncoding.DecodeString(signedChallenge)
	//challenge:="abcxyz"
	if decodeErr != nil {
		return false, decodeErr
	}

	originalChallengeBytes := []byte(originalChallenge)

	//data, pemDecodeErr := base64.StdEncoding.DecodeString(publicKeyPEM)
	//
	//if pemDecodeErr != nil {
	//	return false, pemDecodeErr
	//}

	publicKeyBytes := []byte(publicKeyPEM)

	block, _ := pem.Decode(publicKeyBytes)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return false, errors.Errorf("failed to decode PEM block containing public key")
	}
	// key, err := x509.ParsePKIXPublicKey(block.Bytes)  | ParsePKCS1PublicKey(block.Bytes)
	publicKey, pubParseErr := x509.ParsePKCS1PublicKey(block.Bytes)
	if pubParseErr != nil {
		return false, errors.Errorf("failed parsing piublic key: %v", pubParseErr)
	}

	hashed := sha512.Sum512(originalChallengeBytes)

	resErr := rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, hashed[:], decodedSignedChallenge)
	if resErr != nil {
		return false, errors.Errorf("not verified: %v", resErr)

	} else {
		return true, nil
	}

}

//Check newly enrolled device via U2F or TOTP
func CheckDeviceEnroll(deviceID, clientIP, orgName, timezone, orgID string) (bool, error) {
	deviceDetail, err := devices.Store.GetFromID(deviceID)
	if err != nil {
		return false, err
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		logrus.Error(err)
		return false, err
	}
	now := time.Now().In(loc)

	if global.GetConfig().Platform.Base == "private" {
		tfaResp, err := sendNotificationThroughCloudProxy([]string{deviceDetail.FcmToken}, orgName, "Test", clientIP, now.String())
		if err != nil || tfaResp.Answer != "YES" {
			return false, err
		}
	} else {
		challenge := utils.GetRandomString(5)
		err := notif.Store.SendPushNotification(deviceDetail.FcmToken, orgName, "Test", clientIP, now.String(), challenge)
		if err != nil {
			return false, err
		}

		//TODO for ugly workaround, we have escaped setting device id this time.
		err = redis.Store.Set(challenge, time.Second*400, "status", "false")
		if err != nil {
			return false, err
		}

		// Wait and verify U2f challenge
		checkU2FChallenge, _ := redis.Store.WaitForStatusAndGet(challenge, "device")

		if checkU2FChallenge == false {
			return false, nil
		}

		return true, nil
	}
	return false, nil
}

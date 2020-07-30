package devices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/seknox/trasa/core/notif"
	"github.com/seknox/trasa/core/redis"
	"github.com/seknox/trasa/global"
	"github.com/seknox/trasa/models"
	"github.com/seknox/trasa/utils"
	"github.com/sirupsen/logrus"
)

type giveMeDeviceDetail struct {
	DeviceID     string `json:"deviceID"`
	FcmToken     string `json:"fcmToken"`
	PublicKey    string `json:"publicKey"`
	DeviceFinger string `json:"deviceFinger"`
}

// DeviceDetailPipe
func DeviceDetailPipe(w http.ResponseWriter, r *http.Request) {
	var request giveMeDeviceDetail
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "malformed request", "DeviceDetailPipe")
	}

	//fmt.Println("got enroll request for device: ", request.DeviceID)
	// mars, _ := json.Marshal(request)
	// fmt.Println("got DeviceDetailPipe: ", string(mars))

	// storeDeviceDetailInRedis, err := dbstore.Connect.SetDeviceDetail(request.DeviceID, "status", "false", "fcmToken", request.FcmToken, "publicKey", request.PublicKey)
	// if storeDeviceDetailInRedis == false {
	// 	logger.Error(err)
	// }

	// start polling for device detail which will be updated by pipe handler. we can use same poller from Services
	check, _ := redis.Store.WaitForStatusAndGet(request.DeviceID, "fcmToken")

	// if check is true, this means we had success getting user device detail. fetch redis store again for user device and send it as a response to onprem server.
	if check {
		deviceDetailFromRedis, err := redis.Store.MGet(request.DeviceID, "status", "fcmToken", "publicKey", "deviceFinger")
		if err != nil || len(deviceDetailFromRedis) != 4 {
			logrus.Errorf("could not get device details, got %v : %v", deviceDetailFromRedis, err)
			utils.TrasaResponse(w, 200, "failed", "could not get device details", "DeviceDetailPipe")
		}

		var response giveMeDeviceDetail

		response.DeviceID = request.DeviceID
		response.FcmToken = deviceDetailFromRedis[1]
		response.PublicKey = deviceDetailFromRedis[2]
		response.DeviceFinger = deviceDetailFromRedis[3]

		jsonifiedResponse, err := json.Marshal(response)
		if err != nil {
			notif.SendErrorReport(err, "Cannot marshall to json")
		}

		//TODO use trasa response??
		w.Write(jsonifiedResponse)
		return
	}

	// fmt.Println("could not get device details: ", request.DeviceID)

	utils.TrasaResponse(w, 200, "failed", "could not get device details", "DeviceDetailPipe", nil, nil)
	return

}

// PassMyDeviceDetail receives device detial from mobile client and passes this data to deviceDetailPipe by storing values in redis.
func PassMyDeviceDetail(w http.ResponseWriter, r *http.Request) {
	var request giveMeDeviceDetail
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "malformed request", "DeviceDetailPipe", nil, nil)
		return
	}

	// mars, _ := json.Marshal(request)
	// fmt.Println("got PassMyDeviceDetail: ", string(mars))

	// fmt.Println("got request from mobile for device: ", request.DeviceID)

	err := redis.Store.Set(request.DeviceID, time.Second*400, "status", "true", "fcmToken", request.FcmToken, "publicKey", request.PublicKey, "deviceFinger", request.DeviceFinger)
	if err != nil {
		logrus.Error(err)
		return
	}

	utils.TrasaResponse(w, 200, "success", "shared your device details", "PassMyDeviceDetail", nil, nil)
	// fmt.Println("successfully passes", request.DeviceID)
	return
}

type EnrolDeviceStruct struct {
	DeviceID      string `json:"deviceID"`
	TotpSSC       string `json:"totpSSC"`
	OrgName       string `json:"orgName"`
	CloudProxyURL string `json:"cloudProxyURL"`
}

// EnrolDeviceFunc should only be called after successful authentication and authorization.
func EnrolDeviceFunc(userDetail models.User) EnrolDeviceStruct {

	deviceID, _ := uuid.NewV4()
	totpSec := utils.GenerateTotpSecret()
	orguser := fmt.Sprintf("%s:%s", userDetail.OrgID, userDetail.ID)

	respVal := EnrolDeviceStruct{
		DeviceID:      deviceID.String(),
		TotpSSC:       totpSec,
		CloudProxyURL: global.GetConfig().Trasa.CloudServer,
	}

	go GiveMeDeviceDetail(orguser, deviceID.String(), totpSec)

	return respVal
}

// GiveMeDeviceDetail calls trasa cloud server asking for users device details.
// It does by first calling cloud server by providing device ID and expects fcm token, public keys in return.
// once it receives device details, it calls redis to get user ID for that device id then stores users device detail in database.
func GiveMeDeviceDetail(orguser, deviceID, totpSec string) {
	defer func() {
		if r := recover(); r != nil {
			notif.SendErrorReport(errors.New(fmt.Sprintf(`%v:%s`, r, string(debug.Stack()))), "Panic in GiveMeDeviceDetail")
		}

	}()

	// (1) store device id and orguser into redis.
	err := redis.Store.Set(deviceID, time.Second*400, "orguser", orguser)
	if err != nil {
		logrus.Error(err)
		return
	}

	// (2) call cloud server asking for device detail.
	deviceDetail, err := callServerForDeviceDetail(deviceID)
	if err != nil {
		logrus.Error(err)
		return
	}

	// (3) fetch user detail from redis (based on deviceID).
	orguser, err = redis.Store.Get(deviceID, "orguser")
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Debug(orguser)
	orgUserArray := strings.Split(orguser, ":")

	if len(orgUserArray) != 2 {
		logrus.Errorf("length of orgUserArray is not 2: %v", orgUserArray)
		return
	}

	// (4) store user device detail in database
	var userDevice models.UserDevice

	userDevice.OrgID = orgUserArray[0]
	userDevice.UserID = orgUserArray[1]
	userDevice.DeviceID = deviceDetail.DeviceID
	userDevice.TotpSec = totpSec
	userDevice.FcmToken = deviceDetail.FcmToken
	userDevice.DeviceType = "mobile"
	userDevice.PublicKey = deviceDetail.PublicKey

	userDevice.AddedAt = time.Now().Unix() //.In(nep).String()

	//logrus.Debug(deviceDetail.DeviceFinger,"\n\n\n\n")

	var devHyg models.DeviceHygiene
	err = json.Unmarshal([]byte(deviceDetail.DeviceFinger), &devHyg)
	if err != nil {
		//TODO handle error after mobile app is also updated in ios
		// ignoring error for backward compatibility
		logrus.Trace(deviceDetail.DeviceFinger)
		logrus.Error(err)
	} else {
		userDevice.DeviceHygiene = devHyg
		userDevice.MachineID = devHyg.DeviceInfo.MachineID
	}

	//Also remove this
	if deviceDetail.DeviceFinger == "" {
		userDevice.DeviceFinger = "{}"
	} else {
		userDevice.DeviceFinger = deviceDetail.DeviceFinger //"{}"
	}

	// we register device id and fcm tokens here.
	if strings.Compare(userDevice.DeviceID, "") == 0 {
		logrus.Error("device id not found GiveMeDeviceDetail")
		return
	}
	err = Store.Register(userDevice)
	if err != nil {
		logrus.Error(err)
		return
	}

	//val, err := dbstore.Connect.GetUserApps(orgUserArray[1], userDevice.OrgID)
	//if err != nil {
	//	logrus.Error(err)
	//}

	//appsArray := val
	//var fcmTokens []string
	//fcmTokens = append(fcmTokens, deviceDetail.FcmToken)
	//mar, err := json.Marshal(appsArray)
	//if err != nil {
	//	logrus.Error("Json marshal error in Givemedevive detail")
	//}
	//deferAppDetailSync(fcmTokens, string(mar))

}

func callServerForDeviceDetail(deviceID string) (models.UserDevice, error) {

	var requestConfig models.UserDevice
	requestConfig.DeviceID = deviceID

	urlPath := global.GetConfig().Trasa.CloudServer + "/api/v1/devicedetailpipe"

	//	inseccure := global.GetConfig().Security.InsecureSkipVerify
	mar, err := json.Marshal(requestConfig)
	if err != nil {
		return models.UserDevice{}, errors.Errorf("failed to marshal request : %v", err)
	}

	//resp, err := utils.CallTrasaAPI(urlPath, requestConfig, inseccure)
	resp, err := http.Post(urlPath, "application/json", bytes.NewBuffer(mar))
	if err != nil {
		return models.UserDevice{}, errors.Errorf("failed to get device detail: %v", err)
	}

	var dev models.UserDevice
	err = json.NewDecoder(resp.Body).Decode(&dev)
	if err != nil {
		return models.UserDevice{}, errors.Errorf("failed to get device detail: %v", err)
	}
	return dev, err
	//
	//logrus.Debug(resp.Data)
	//result, ok := resp.Data.([]models.UserDevice)
	//if !ok || len(result) != 1 {
	//	logrus.Debug(result)
	//	return models.UserDevice{}, errors.Errorf("failed to get device detail")
	//}
	//
	//return result[0], nil
}

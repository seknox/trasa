package tfa

import (
	"fmt"

	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/global"
	"github.com/tstranex/u2f"
)

//TODO @sshah appID??
var appID = global.GetConfig().Trasa.Dashboard

var trustedFacets = []string{appID}

// Normally these state variables would be stored in a database.
// For the purposes of the demo, we just store them in memory.
var challenge *u2f.Challenge

var registrations []u2f.Registration
var regStore = make(map[string][]u2f.Registration)
var counter uint32

// RegisterRequest generates random challenge and sends it to client (browser) that communicates with user hardware device.
// Device needs to sign this challenge and return the signed data along with device details back to server.
//func RegisterRequest(w http.ResponseWriter, r *http.Request) {
//	userContext := r.Context().Value("user").(models.UserContext)
//
//	c, err := u2f.NewChallenge(appID, trustedFacets)
//	if err != nil {
//		logrus.Error(err)
//		utils.TrasaResponse(w, 200, "failed", "Could not start device registration process.", "RegisterRequest")
//		return
//	}
//
//	challval, _ := json.Marshal(c)
//
//	err = redis.Store.Set(
//		fmt.Sprintf("%s:%s", userContext.User.ID, "yubikeyreg"),
//		time.Second*400,
//		"challenge", string(challval),
//		"userMeta", "...")
//
//	if err != nil {
//		logrus.Error(err)
//		utils.TrasaResponse(w, 200, "failed", "Could not set yubico challenge.", "RegisterRequest")
//		return
//	}
//	req := u2f.NewWebRegisterRequest(c, registrations)
//
//	utils.TrasaResponse(w, 200, "success", "challenge presented.", "RegisterRequest", req)
//}
//
//// RegisterResponse gets signed response form client(browser) which contains device details that uniquely identifies user security token(hardware device).
//func RegisterResponse(w http.ResponseWriter, r *http.Request) {
//	userContext := r.Context().Value("user").(models.UserContext)
//
//	var regResp u2f.RegisterResponse
//	if err := json.NewDecoder(r.Body).Decode(&regResp); err != nil {
//		logrus.Error(err)
//		utils.TrasaResponse(w, 200, "failed", "could not finish device registration process.", "RegisterResponse","")
//		return
//	}
//
//	val:= redis.Store.Get(fmt.Sprintf("%s:%s", userContext.User.ID, "yubikeyreg"), "challenge")
//	if val=="" {
//		utils.TrasaResponse(w, 200, "failed", "Could not get yubico challenge.", "RegisterRequest")
//		return
//	}
//
//	config := &u2f.Config{
//		// Chrome 66+ doesn't return the device's attestation
//		// certificate by default.
//		SkipAttestationVerify: true,
//	}
//
//	//byted := byte[](vals[0])
//	var challval u2f.Challenge
//	err:=json.Unmarshal([]byte(val), &challval)
//	if err != nil {
//		logrus.Error(err)
//		utils.TrasaResponse(w, 200, "failed", "could not unmarshall challenge.", "RegisterResponse", "")
//		return
//	}
//
//	reg, err := u2f.Register(regResp, challval, config)
//	if err != nil {
//		logrus.Error(err)
//		utils.TrasaResponse(w, 200, "failed", "could not finish device registration process.", "RegisterResponse",  "")
//		return
//	}
//
//	// store the reg data to database in userdevice table.
//	var device models.UserDevice
//	device.DeviceID = utils.GetUUID()
//	device.DeviceType = "htoken"
//	device.OrgID = userContext.User.OrgID
//	device.UserID = userContext.User.ID
//	device.AddedAt = time.Now().Unix() //.In(nep).String()
//	device.FcmToken = utils.EncodeBase64(reg.KeyHandle)
//	//device.DeviceFinger = string(reg.Raw)
//	finger, _ := json.Marshal(reg.Raw)
//
//	device.DeviceFinger = string(finger)
//
//	pubBytes := utils.GetEcdsaPublicKeyBytes(&reg.PubKey)
//
//	device.PublicKey = utils.EncodeBase64(pubBytes)
//
//	//dat, _ := json.Marshal(device)
//	//fmt.Println(string(dat))
//
//	err = devices.Store.Register(device)
//	if err != nil {
//		logrus.Error(err)
//		utils.TrasaResponse(w, 200, "failed", "could not finish device registration process.", "RegisterResponse", "")
//		return
//	}
//	// skipping counter for now. Yubikeys are tamper resistant anyway....over engineering.@https://developers.yubico.com/U2F/Libraries/Advanced_topics.html
//	//counter = 0
//
//	//log.Printf("Registration success: %+v", reg)
//	utils.TrasaResponse(w, 200, "success", "device registered", "RegisterResponse")
//}
//
//func SignRequest(w http.ResponseWriter, r *http.Request) {
//
//	userContext := r.Context().Value("user").(models.UserContext)
//	//var device utils.UserDevice
//
//	deviceDetailFromDB := dbstore.Connect.GetUserHtokenDevice(userContext.User.ID, userContext.User.OrgID)
//
//	var reg u2f.Registration
//	var regs []u2f.Registration
//	keyval := fmt.Sprintf("%s:%s:%s", userContext.User.ID, userContext.User.OrgID, "yubikesign")
//
//	for _, v := range deviceDetailFromDB {
//		json.Unmarshal([]byte(v.DeviceFinger), &reg.Raw)
//
//		//fmt.Println("reg.fcm: ", v.FcmToken)
//		bit, err := utils.DecodeBase64(v.FcmToken)
//		if err != nil {
//			fmt.Println("cannot retrieve key handle: error decoding base64: ", err)
//		}
//
//		reg.KeyHandle = make([]byte, len(bit))
//		//copy(tmp, arr)
//		copy(reg.KeyHandle, bit)
//		pubB, err := utils.DecodeBase64(v.PublicKey)
//		if err != nil {
//			logrus.Error(err)
//			return
//		}
//		err = utils.GetEcdsaPublicKeyFromBytes(&reg, pubB)
//		if err != nil {
//			logrus.Error(err)
//			//logrus.Error(err)
//		}
//		//json.Unmarshal([]byte(v.PublicKey), &reg.PubKey)
//		regs = append(regs, reg)
//	}
//
//	// if registrations == nil {
//	// 	http.Error(w, "registration missing", http.StatusBadRequest)
//	// 	return
//	// }
//
//	c, err := u2f.NewChallenge(appID, trustedFacets)
//	if err != nil {
//		logrus.Error("u2f.NewChallenge error: %v", err)
//		http.Error(w, "error", http.StatusInternalServerError)
//		return
//	}
//
//	challval, _ := json.Marshal(c)
//
//	err = dbstore.Connect.SetYubicoChall(utils.EncodeBase64(c.Challenge), "challenge", string(challval), "userMeta", keyval)
//	if err != nil {
//		logrus.Error(err)
//		//logrus.Error(err)
//		//utils.Do.Systemlogrus(consts.RedisErr, false, err.Error(), "RegisterRequest:u2f.NewChallenge error:", "error")
//		return
//	}
//
//	// store registration data in memory
//	regStore[keyval] = regs
//
//	req := c.SignRequest(regs)
//
//	utils.TrasaResponse(w, 200, "success", "auth challenge presented.", "SignRequest", nil, req)
//}

// signReq is used by TfaHandler
func SignReq(userID, orgID string) (*u2f.WebSignRequest, error) {

	return nil, nil

	//
	//deviceDetailFromDB := dbstore.Connect.GetUserHtokenDevice(userID, orgID)
	//
	//var reg u2f.Registration
	//var regs []u2f.Registration
	//keyval := fmt.Sprintf("%s:%s:%s", userID, orgID, "yubikesign")
	//
	//for _, v := range deviceDetailFromDB {
	//	json.Unmarshal([]byte(v.DeviceFinger), &reg.Raw)
	//
	//	//fmt.Println("reg.fcm: ", v.FcmToken)
	//	bit, err := utils.DecodeBase64(v.FcmToken)
	//	if err != nil {
	//		logrus.Warning(err)
	//		//fmt.Println("cannot retrieve key handle: error decoding base64: ", err)
	//	}
	//
	//	reg.KeyHandle = make([]byte, len(bit))
	//	//copy(tmp, arr)
	//	copy(reg.KeyHandle, bit)
	//	pubB, err := utils.DecodeBase64(v.PublicKey)
	//	if err != nil {
	//		logrus.Error(err)
	//		return nil, err
	//	}
	//	err = utils.GetEcdsaPublicKeyFromBytes(&reg, pubB)
	//	if err != nil {
	//		logrus.Error(err)
	//		return nil, err
	//	}
	//	//json.Unmarshal([]byte(v.PublicKey), &reg.PubKey)
	//	regs = append(regs, reg)
	//}
	//
	//// if registrations == nil {
	//// 	http.Error(w, "registration missing", http.StatusBadRequest)
	//// 	return
	//// }
	//
	//c, err := u2f.NewChallenge(appID, trustedFacets)
	//if err != nil {
	//	logrus.Errorf("u2f.NewChallenge error: %v", err)
	//	return nil, err
	//}
	//
	//challval, _ := json.Marshal(c)
	//
	//err = dbstore.Connect.SetYubicoChall(utils.EncodeBase64(c.Challenge), "challenge", string(challval), "userMeta", keyval)
	//if err != nil {
	//	logrus.Error(err)
	//	//logrus.Error(err)
	//	//utils.Do.Systemlogrus(consts.RedisErr, false, err.Error(), "RegisterRequest:u2f.NewChallenge error:", "error")
	//	return nil, err
	//}
	//
	//// store registration data in memory
	//regStore[keyval] = regs
	//
	//req := c.SignRequest(regs)
	//
	//return req, err

}

//
//func SignResponseHandler(w http.ResponseWriter, r *http.Request) {
//	userContext := r.Context().Value("user").(models.UserContext)
//	var signResp u2f.SignResponse
//	if err := json.NewDecoder(r.Body).Decode(&signResp); err != nil {
//		http.Error(w, "invalid response: "+err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	keyval := fmt.Sprintf("%s:%s:%s", userContext.User.ID, userContext.User.OrgID, "yubikesign")
//
//	vals, err := dbstore.Connect.GetYubicoChall(keyval, "challenge", "userMeta")
//	if err != nil {
//		logrus.Error(err)
//		//logrus.Error(err)
//		//utils.Do.Systemlogrus(consts.RedisErr, false, err.Error(), "RegisterRequest:u2f.NewChallenge error:", "error")
//		return
//	}
//
//	err = SignResp(signResp, keyval, vals[0])
//
//	if err != nil {
//		utils.TrasaResponse(w, 200, "success", "auth challenge verified.", "SignResponse", nil, nil)
//		return
//	}
//	logrus.Warnf("VerifySignResponse error: %v", err)
//	utils.TrasaResponse(w, 200, "failed", "auth challenge imvalid.", "SignResponse", nil, nil)
//
//	//var challval u2f.Challenge
//	//json.Unmarshal([]byte(vals[0]), &challval)
//	//
//	//registrations := regStore[keyval]
//	//if registrations == nil {
//	//	http.Error(w, "registration missing", http.StatusBadRequest)
//	//	return
//	//} else {
//	//	delete(regStore, keyval)
//	//}
//	//
//	////fmt.Println("registration: ", registrations)
//	//
//	//for _, reg := range registrations {
//	//	//var r u2f.Registration
//	//
//	//	newCounter, authErr := reg.Authenticate(signResp, challval, counter)
//	//	if authErr == nil {
//	//		log.Printf("newCounter: %d", newCounter)
//	//		counter = newCounter
//	//		utils.TrasaResponse(w, 200, "success", "auth challenge verified.", "SignResponse", nil, nil)
//	//		return
//	//	}
//	//}
//	//
//	//log.Printf("VerifySignResponse error: %v", err)
//	//utils.TrasaResponse(w, 200, "failed", "auth challenge imvalid.", "SignResponse", nil, nil)
//}

// SignResp is U2F sign response function that can be used where signing function needs to be called within other functions
func SignResp(signResp *u2f.SignResponse, keyval, challengeJson string) error {

	return nil

	//
	//
	//
	//var challval u2f.Challenge
	//err := json.Unmarshal([]byte(challengeJson), &challval)
	//if err != nil {
	//	logrus.Error(err)
	//	return err
	//}
	//
	//registrations := regStore[keyval]
	//if registrations == nil {
	//	return fmt.Errorf("%s", "could not fetch registration data from cache")
	//} else {
	//	delete(regStore, keyval)
	//}
	//
	////fmt.Println("registration: ", registrations)
	//
	//for _, reg := range registrations {
	//	//var r u2f.Registration
	//
	//	newCounter, authErr := reg.Authenticate(signResp, challval, counter)
	//	if authErr == nil {
	//		logrus.Tracef("newCounter: %d", newCounter)
	//		counter = newCounter
	//		return nil
	//	}
	//}
	//
	//return fmt.Errorf("%s", "failed verifying user data")
}

func checkU2FY(signResponse *u2f.SignResponse, userID, orgID string) (string, error) {
	keyval := fmt.Sprintf("%s:%s:%s", userID, orgID, "yubikesign")

	challenge, _ := redis.Store.Get(keyval, "challenge")

	//TODO @sshah whats the use of userMeta??
	_, _ = redis.Store.Get(keyval, "userMeta")

	err := SignResp(signResponse, keyval, challenge)
	if err != nil {
		//TODO @sshah fix case for valid u2fy
		//get deviceID
		return "", err
	}
	//TODO return device ID
	return "", nil
	//logrus.Tracef("VerifySignResponse error: %v", err)
}

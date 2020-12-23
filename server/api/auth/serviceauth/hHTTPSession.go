package serviceauth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/seknox/trasa/server/api/services"

	"github.com/seknox/trasa/server/api/accesscontrol"
	"github.com/seknox/trasa/server/api/auth/tfa"
	"github.com/seknox/trasa/server/api/devices"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/policies"
	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/nacl/secretbox"
)

type NewSession struct {
	ExtToken  string `json:"extToken"`
	HostName  string `json:"hostName"`
	TfaMethod string `json:"tfaMethod"`
	TotpCode  string `json:"totpCode"`
}

// sessionStore stores session data (logSession Modal) in memory until client logs out.
// during logout sequence, the caller must add logout time and upload it to elasticsearch
var sessionStore = make(map[string]logs.AuthLog)
var sessionStoreMutex sync.Mutex

// a cache session store that tells if 2fa attempt has already been queued for this authentication attempt.
// if inithttprequest has not return but the trasa extension already sends another inithttp request,
// this causes trasa to send multiple 2fa request. we store information regarding if 2fa request matching
// same user and authapp is already in progress. if found, we simply ignore 2fa requirement.
// we should delete the queue from cache store Immediately once 2fa is returned either success or failure.
type http2faCache struct {
	ExtToken string `json:"exToken"`
	Domain   string `json:"domain"`
	Status   bool   `json:"status"`
}

// var http2faCacheStore = make(map[string]bool)

// AuthHTTPAccessProxy initiates http access proxy session. Intent should be 'AUTH_HTTP_ACCESS_PROXY'
func AuthHTTPAccessProxy(w http.ResponseWriter, r *http.Request) {
	logger.Trace("AuthHTTPAccessProxy request")
	var req NewSession

	authlog := logs.NewLog(r, "http")

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "AuthHTTPAccessProxy", nil)
		logs.Store.LogLogin(&authlog, consts.REASON_MALFORMED_REQUEST_RECEIVED, false)
		return
	}

	orgID, deviceID, userID, err := devices.Store.GetDeviceAndOrgIDFromExtID(req.ExtToken)
	if err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid extension id", "AuthHTTPAccessProxy", nil)
		logs.Store.LogLogin(&authlog, consts.REASON_DEVICE_NOT_ENROLLED, false)
		return
	}

	authlog.AccessDeviceID = deviceID

	serviceDetailFromDB, err := services.Store.GetFromHostname(req.HostName, "http", "", orgID)
	if err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not fetch service detail", "AuthHTTPAccessProxy", nil)
		logs.Store.LogLogin(&authlog, consts.REASON_INVALID_SERVICE_HOSTNAME, false)
		return
	}

	authlog.UpdateService(serviceDetailFromDB)

	userDetailFromDB, err := users.Store.GetFromID(userID, orgID)
	if err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not fetch user details", "AuthHTTPAccessProxy", nil)
		logs.Store.LogLogin(&authlog, consts.REASON_USER_NOT_FOUND, false)
		return
	}

	authlog.UpdateUser(&models.UserWithPass{User: *userDetailFromDB})

	orgDetailFromDB, err := orgs.Store.Get(orgID)
	if err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not fetch user details", "AuthHTTPAccessProxy", nil)
		logs.Store.LogLogin(&authlog, consts.REASON_USER_NOT_FOUND, false)
		return
	}

	policy, adhoc, err := policies.Store.GetAccessPolicy(userDetailFromDB.ID, serviceDetailFromDB.ID, "", orgDetailFromDB.ID)
	if err != nil {
		logger.Debug(err)
		utils.TrasaResponse(w, 200, "failed", "no policy assigned", "AuthHTTPAccessProxy", nil)
		logs.Store.LogLogin(&authlog, consts.REASON_NO_POLICY_ASSIGNED, false)
		return
	}

	// TODO @sshahcodes, check if we can assign user's username here
	//authlog.Privilege = privilege
	//authlog.Privilege = userDetailFromDB.UserName
	authlog.SessionRecord = policy.RecordSession

	ok, reason := accesscontrol.CheckPolicy(&models.ConnectionParams{
		ServiceID: serviceDetailFromDB.ID,
		OrgID:     orgDetailFromDB.ID,
		UserID:    userDetailFromDB.ID,
		SessionID: authlog.SessionID,
		UserIP:    utils.GetIp(r),
		Timezone:  orgDetailFromDB.Timezone,
	}, policy, adhoc)

	if !ok {
		logger.Debugf("policy failed with reason: %s", reason)
		utils.TrasaResponse(w, 200, "failed", string(reason), "AuthHTTPAccessProxy", nil)
		logs.Store.LogLogin(&authlog, reason, false)
		return
	}

	var userMobileDevice models.UserDevice

	userMobileDevice.UserID = userDetailFromDB.ID

	if policy.TfaRequired {
		deviceID, reason, ok := tfa.HandleTfaAndGetDeviceID(nil,
			req.TfaMethod,
			req.TotpCode,
			userDetailFromDB.ID,
			utils.GetIp(r),
			serviceDetailFromDB.Name,
			orgDetailFromDB.Timezone,
			orgDetailFromDB.OrgName,
			orgDetailFromDB.ID,
		)
		authlog.TfaDeviceID = deviceID

		if !ok {
			logger.Debugf("tfa failed with reason: %s", reason)
			utils.TrasaResponse(w, 200, "failed", string(reason), "AuthHTTPAccessProxy", nil)
			logs.Store.LogLogin(&authlog, reason, false)
			return
		}
	}

	//////////////////////////////////////////////////////////////////////
	////////////////		Session generation sequence     //////////////
	//////////////////////////////////////////////////////////////////////

	//// Generate session tokens
	encodedSession := authlog.SessionID
	authKey := utils.GetRandomBytes(17)

	// Insert user session values in database
	orgusr := fmt.Sprintf("%s:%s", orgID, userID)

	// sessionStoreMutex.Lock()
	sessionStore[encodedSession] = authlog
	// sessionStoreMutex.Unlock()

	//logger.Trace(authlog.SessionID)
	//logger.Trace(encodedSession)

	// create csrf token
	var encryptionKey [32]byte
	// we use sessionKey as secret key for encryption
	copy(encryptionKey[:], authKey)

	// generate random nonce
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		logger.Error("error creating nonce")
	}

	// we encrypt keyVal with sessionKey and set cipher text as csrf token.
	// To verify csrf token, we simply decrypt csrf token with session key
	// If decryption fails, it suggest csrf token has been tampered and is invalid.
	csrfToken := secretbox.Seal(nonce[:], []byte(orgusr), &nonce, &encryptionKey)

	//	csrfcookie := http.Cookie{Name: "x-csrf", Value: base64.StdEncoding.EncodeToString(csrfToken), Path: "/"}

	// acessProxySessionCookie := http.Cookie{Name: "x-ap-session", Value: encodedSession, Domain: "app.trasa", HttpOnly: true}
	// http.SetCookie(w, &acessProxySessionCookie)

	authData := fmt.Sprintf("%s:%s:%s", base64.StdEncoding.EncodeToString(authKey), deviceID, "")

	//TODO can we use boolean value here?
	sessionRec := "true"
	if !policy.RecordSession {
		sessionRec = "false"
	}
	//logger.Trace(encodedSession)
	err = redis.Store.SetHTTPAccessProxySession(encodedSession, orgusr, authData, sessionRec)
	if err != nil {
		logger.Errorf("setting session in redis: %v", err)
	}

	//set this session in active session.
	err = logs.Store.AddNewActiveSession(&authlog, authlog.SessionID, "http")
	if err != nil {
		logger.Errorf("add active session: %v", err)
	}

	var sessionIdentifiers Session
	sessionIdentifiers.SessionID = encodedSession
	sessionIdentifiers.CsrfToken = base64.StdEncoding.EncodeToString(csrfToken)
	sessionIdentifiers.SessionRecord = policy.RecordSession

	if policy.RecordSession {
		directoryBuilder := filepath.Join(utils.GetTmpDir(), "trasa", "accessproxy", "http", encodedSession)
		logger.Tracef("Logging to : %s", directoryBuilder)
		utils.CreateDirIfNotExist(directoryBuilder)
		logPath := filepath.Join(directoryBuilder, fmt.Sprintf("%s.http-raw", encodedSession))
		file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			logger.Error(err)
		}
		defer file.Close()
	} else {
		logger.Tracef("Not logging ")
	}

	// we need user detail and app detail to store in user sessionStore.
	// This value can be fetched based on exToken received and fetching user detial based on exToken.
	// App details can be fetched from hostname (which will be unique per http app) of http endpoint.
	//sessionStore[encodedSession] =

	utils.TrasaResponse(w, 200, "success", "successfully created session", "InitHttpsSession", sessionIdentifiers)
	//http2faCacheStore[initReqId] = false

	//sessionStoreWithExtokenDomain[initReqId] = true

}

//DestroyHttpSession ends http session and starts logout sequence
func DestroyHttpSession(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error(err)
	}
	for k, v := range r.Form {
		sessionWriter(k, v[0])
		//fmt.Println(k, v)
	}

	sessionID := r.Form.Get("sid")

	logoutSequence(sessionID)
	logger.Trace("session destroyed ", sessionID)
	utils.TrasaResponse(w, http.StatusOK, "success", "session destroyed", "Destroy session", nil)
}

type Session struct {
	SessionID     string `json:"sessionID"`
	CsrfToken     string `json:"csrfToken"`
	SessionRecord bool   `json:"sessionRecord"`
}

// GetHttpSession receives session recording from extension
func GetHttpSession(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error(err)
	}
	for k, v := range r.Form {
		sessionWriter(k, v[0])
		//fmt.Println(k, v)
	}

	utils.TrasaResponse(w, 200, "success", "sr collected", "srcollector", nil, nil)
}

func sessionWriter(sessionID, shots string) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf(` Panic in sessionWriter: %v. stack: %s`, r, string(debug.Stack()))
		}

	}()

	//logger.Tracef("session writer req received: %s", sessionID)

	// separate multiple screenshot with counter		logrus.Trace(allow,reason)
	dataURIWithCounter := strings.Split(shots, "::")

	// loop over dataURIWithCounter
	for _, v := range dataURIWithCounter {
		logStruct, ok := sessionStore[sessionID]
		if !ok {
			logger.Debug("session store not found for sessionID: ", sessionID)
			return
		}

		//logger.Tracef("shot length: %v", len(v))
		// operate only if length is greater than 6
		if len(v) > 6 {
			// split the string again to separate counter and dataURI
			counterAndImage := strings.Split(v, ":")

			//logger.Tracef("counterAndImage length: %v", len(counterAndImage))

			// the length shouold exactly be 2
			if len(counterAndImage) == 2 {
				videoRecord := logStruct.SessionRecord //strings.Split(logStruct.AppID, ":")
				//home, _ := hdir.Dir()
				directoryBuilder := filepath.Join(utils.GetTmpDir(), "trasa", "accessproxy", "http", sessionID)
				if videoRecord == true {

					// convert to png
					img, err := png.Decode(gopherPNG(counterAndImage[1]))
					if err != nil {
						if _, ok := err.(base64.CorruptInputError); ok {
							logger.Debug("base64 input is corrupt, check service Key")
						}
						logger.Debug(err)
						//return
					}

					//directoryBuilder := fmt.Sprintf("../trasahttpgateway/logs/%s", parseMsg[1])

					utils.CreateDirIfNotExist(directoryBuilder)

					filenameBuilder := filepath.Join(directoryBuilder, fmt.Sprintf("%s.png", counterAndImage[0]))
					//fmt.Println("writing to file: ", filenameBuilder)

					outputFile, err := os.Create(filenameBuilder)
					if err != nil {
						logger.Debug(err)
						// Handle error
					}

					// Encode takes a writer interface and an image interface
					// We pass it the File and the RGBA
					err = png.Encode(outputFile, img)
					if err != nil {
						logger.Debugf("png encode: %v", err)
					}

					// Don't forget to close files
					outputFile.Close()
					// return at last
				}
				logPath := filepath.Join(directoryBuilder, fmt.Sprintf("%s.http-raw", sessionID))
				file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
				if err != nil {
					logger.Debug(err)
				}
				defer file.Close()

			}

		}
	}
	return
}

// gopherPNG creates an io.Reader by decoding the base64 encoded image data string in the gopher constant.
func gopherPNG(imsgData string) io.Reader {
	return base64.NewDecoder(base64.StdEncoding, strings.NewReader(imsgData))
}

func logoutSequence(sessionID string) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf(`Panic in logoutSequence: %v:%s`, r, string(debug.Stack()))
		}

	}()

	logger.Debug("HTTP - logout request received for: ", sessionID)
	logStruct := sessionStore[sessionID]
	//sessionCacheKey := fmt.Sprintf("%s:%s", logStruct.RegisteredDeviceID, logStruct.ServerIP)

	//TODO check if session id exists in sessionStore. if not , this means expiry key is not of http session and we should return process
	if _, ok := sessionStore[sessionID]; !ok {
		return
	}

	// 1) log to elasticsearch

	err := logs.Store.LogLogin(&logStruct, "", true) //dbstore.Connect.LogSession(logStruct)
	if err != nil {
		logger.Error(err)
	}

	// 2) deletes session from redis
	err = redis.Store.Delete(sessionID)
	if err != nil {
		logger.Error(err)
	}

	//remove from active session
	err = logs.Store.RemoveActiveSession(logStruct.SessionID)
	if err != nil {
		logger.Error(err)
	}

	// 3) create video file from image file

	deleteDirectory := false
	directory := filepath.Join(utils.GetTmpDir(), "trasa", "accessproxy", "http", sessionID)

	SessionRecord := logStruct.SessionRecord
	if SessionRecord == true {

		//	directory := fmt.Sprintf("../trasahttpgateway/logs/%s", sessionID)
		videoFileName := fmt.Sprintf("%s.mp4", sessionID)
		rawFileName := fmt.Sprintf("%s.http-raw", sessionID)

		_, err := os.Stat(directory)

		// if directory does not exist, it probably means session record was disabled.
		if os.IsNotExist(err) {
			logger.Trace(err)
			return
		}

		err = createVideo(directory, sessionID)
		if err != nil {
			logger.Error(err)
		}
		videoSource := filepath.Join(directory, videoFileName)

		// 4) upload video file to minio
		err = uploadToMinio(videoSource, logStruct)
		if err != nil {
			logger.Error(err)
		}

		rawSource := filepath.Join(directory, rawFileName)

		// 5) upload raw log file to minio
		err = uploadToMinio(rawSource, logStruct)
		if err != nil {
			logger.Error(err)
		}

		deleteDirectory = true

	}

	// 6) delete directory
	if deleteDirectory == true {

		err = os.RemoveAll(directory)
		if err != nil {
			logger.Error(err)
		}
	} else {
		logger.Tracef("Not deleting directory %s as video failed", sessionID)
	}

	// we delete sessionvalur from sessionStore
	//delete(sessionStoreWithExtokenDomain, sessionCacheKey)
	// sessionStoreMutex.Lock()
	delete(sessionStore, sessionID)
	// sessionStoreMutex.Unlock()

}

func uploadToMinio(src string, login logs.AuthLog) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	err = logs.Store.UploadHTTPLogToMinio(source, login)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

// ffmpeg -r 60 -f image2 -s 1920x1080 -i %0d.png -vf pad="width=ceil(iw/2)*2:height=ceil(ih/2)*2" -vcodec libx264 -crf 25  -pix_fmt yuv420p test.mp4

func createVideo(path, sessionID string) error {
	//cmd := exec.Command("ls", "-lah")
	fileName := fmt.Sprintf("%s.mp4", sessionID)
	cmd := exec.Command("ffmpeg", "-r", "60", "-f", "image2", "-s", "1920x1080", "-i", "%d.png", "-vf", "pad='width=ceil(iw/2)*2:height=ceil(ih/2)*2'", "-vcodec", "libx264", "-crf", "25", "-pix_fmt", "yuv420p", fileName)

	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Errorf("createVideo: %s : cmd.Run() failed with %v || %s", sessionID, err, string(output))
		return err
	}
	//fmt.Printf("combined out:\n%s\n", string(out))
	return nil
}

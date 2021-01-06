package my

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"fmt"
	"github.com/seknox/trasa/server/api/devices"
	"github.com/seknox/trasa/server/api/providers/ca"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/api/accesscontrol"

	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/consts"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"

	"net/http"

	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
)

type SingleUserDetailV2 struct {
	User             models.User         `json:"user"`
	AssignedServices []models.Service    `json:"userAccessMaps"`
	UserDevices      []models.UserDevice `json:"userDevices"`
	UserGroups       []models.Group      `json:"userGroups"`
}

func MyAccountDetails(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	var user SingleUserDetailV2

	user.User = *uc.User

	assignedServices, err := users.Store.GetAssignedServices(uc.User.ID, uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to get user apps", "GetSingleUserDetail", user)
		return
	}

	user.AssignedServices = assignedServices

	userDevices, err := users.Store.GetAllDevices(uc.User.ID, uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to fetch user devices.", "SingleUserDevices", user)
		return
	}

	user.UserDevices = userDevices

	groups, err := users.Store.GetGroups(uc.User.ID, uc.Org.ID)
	if err != nil {
		utils.TrasaResponse(w, http.StatusOK, "failed", "could not get groups", "GetGroupsAssignedToUser", user)
		return

	}

	user.UserGroups = groups

	utils.TrasaResponse(w, 200, "success", "success", "GetSingleUserDetail", user)
}

func GenerateKeyPair(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	userID := userContext.User.ID

	//pass userID in context
	privateKeyBytes, publicKeyBytes, certBytes, err := generateTempCertificateForDeviceAgent(userContext.User.Groups, userContext.DeviceID, userContext.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not generate TempCertificate ForDeviceAgent", "Update hygiene", nil)
		return
	}

	err = users.Store.UpdatePublicKey(userID, strings.TrimSpace(string(publicKeyBytes)))
	if err != nil {
		logrus.Errorf(`could not save public key: %v`, err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "could not save public key", "GenerateKeyPair", nil)
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
			utils.TrasaResponse(w, 200, "failed", "could not generate zipWriter.Create", "GenerateKeyPair", nil)
			return
		}
		_, err = zipFile.Write(file.Body)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "could not generate zipWriter.Create", "GenerateKeyPair", nil)
			return
		}

	}

	// Make sure to check the error on Close.
	err = zipWriter.Close()
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not generate zipWriter.Create", "GenerateKeyPair", nil)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=id_rsa.zip")
	w.Write(buff.Bytes())

	return

}

func generateTempCertificateForDeviceAgent(groups []string, deviceID, orgID string) (privateKeyBytes, publicKeyBytes, certBytes []byte, err error) {

	bitSize := 4096
	privateKey, err := utils.GeneratePrivateKey(bitSize)
	if err != nil {
		logrus.Errorf(`could not generate private key: %v`, err)
		return nil, nil, nil, err
	}

	certHolder, err := ca.Store.GetCertHolder(consts.CERT_TYPE_SSH_CA, "system", orgID)
	if err != nil {
		logrus.Debugf(`could not get CA key: %v`, err)
		return nil, nil, nil, err
	}

	publicKeySSH, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		logrus.Errorf(`could not generate public key: %v`, err)
		return nil, nil, nil, err
	}

	publicKeyBytes = ssh.MarshalAuthorizedKey(publicKeySSH)

	caKey, err := ssh.ParsePrivateKey(certHolder.Key)
	if err != nil {
		logrus.Errorf(`Could not parse CA private key: %v`, err)
		return nil, nil, nil, err
	}

	buf := make([]byte, 8)
	_, err = rand.Read(buf)
	if err != nil {
		logrus.Errorf("failed to read random bytes: %v", err)
		return nil, nil, nil, err
	}
	serial := binary.LittleEndian.Uint64(buf)

	//extentions := make(map[string]string)
	extentions := map[string]string{
		"permit-X11-forwarding":   "",
		"permit-agent-forwarding": "",
		"permit-port-forwarding":  "",
		"permit-pty":              "",
		"permit-user-rc":          "",
		"trasa-device-id":         deviceID,
		"trasa-user-groups":       strings.Join(groups, ","),
	}

	//principals := []string{}

	cert := ssh.Certificate{
		Key:             publicKeySSH,
		Serial:          serial,
		CertType:        ssh.UserCert,
		KeyId:           utils.GetRandomString(10),
		ValidPrincipals: nil,
		ValidAfter:      uint64(time.Now().UTC().Unix()),
		ValidBefore:     uint64(time.Now().UTC().Add(time.Hour * 24).Unix()),
		Permissions: ssh.Permissions{
			Extensions: extentions,
		},
	}

	err = cert.SignCert(rand.Reader, caKey)
	if err != nil {
		logrus.Errorf(`could not sign public key: %v`, err)
		return nil, nil, nil, err
	}

	privateKeyBytes = utils.EncodePrivateKeyToPEM(privateKey)
	certBytes = ssh.MarshalAuthorizedKey(&cert)
	if len(certBytes) == 0 {
		logrus.Errorf("failed to marshal signed certificate, empty result")
		err = errors.New("failed to marshal signed certificate, empty result")
		return nil, nil, nil, err
	}

	return privateKeyBytes, publicKeyBytes, certBytes, nil

}

func GetMyEvents(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	events, err := logs.Store.GetLoginEvents("user", userContext.User.ID, userContext.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "could not get logs", "get org logs by page")
		return

	}
	utils.TrasaResponse(w, http.StatusOK, "success", "invalid size or page", "get org logs by page", events)
}

func GetMyEventsByPage(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	//orgID := uc.User.OrgID
	page, err1 := strconv.ParseInt(chi.URLParam(r, "page"), 10, 32)
	size, err2 := strconv.ParseInt(chi.URLParam(r, "size"), 10, 32)

	loc, err := time.LoadLocation(userContext.Org.Timezone)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to load location", "GetLoginEvents: loadlocation", err)
		return

	}

	//Date format 2020-05-18

	dateFromTime, err3 := time.Parse("2006-01-02", chi.URLParam(r, "dateFrom"))
	dateToTime, err4 := time.Parse("2006-01-02", chi.URLParam(r, "dateTo"))

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		utils.TrasaResponse(w, http.StatusBadRequest, "failed", "invalid size or page", "get org logs by page", nil)
		return
	}

	dateFrom := dateFromTime.In(loc).UnixNano()
	dateTo := dateToTime.In(loc).UnixNano()

	events, err := logs.Store.GetLoginEventsByPage("user", userContext.User.ID, userContext.Org.ID, int(page), int(size), dateFrom, dateTo)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "could not get logs", "get org logs by page")
		return

	}
	utils.TrasaResponse(w, http.StatusOK, "success", "", "get org logs by page", events)

}

type AuthMetaResp struct {
	IsPasswordRequired     bool   `json:"isPasswordRequired"`
	IsDeviceHygeneRequired bool   `json:"isDeviceHygeneRequired"`
	TrasaID                string `json:"trasaID"`
}

//Get authentication metada like isPasswordRequired, isDeviceHygeneRequired etc
//They should be verified later when user authenticate
//This API is called just for the UI
func GetAuthMeta(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	username := chi.URLParam(r, "username")
	appID := chi.URLParam(r, "appID")

	var resp AuthMetaResp

	//TODO remove trasaID if not required
	resp.TrasaID = userContext.User.Email
	resp.IsPasswordRequired = true
	resp.IsDeviceHygeneRequired = false

	app1, err := services.Store.GetFromID(appID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to get service", "", resp)
		return
	}

	acccs := strings.Split(app1.ManagedAccounts, ",")
	for _, ac := range acccs {
		if ac == username {
			//Creds already set
			resp.IsPasswordRequired = false
			break
		}

	}

	sett, err := system.Store.GetGlobalSetting(userContext.Org.ID, consts.GLOBAL_DEVICE_HYGIENE_CHECK)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to get global setting", "", resp)
		return
	}

	resp.IsDeviceHygeneRequired = sett.Status
	utils.TrasaResponse(w, 200, "success", "", "", resp)

}

type MyServiceDetail struct {
	User       models.User               `json:"user"`
	MyServices []models.MyServiceDetails `json:"myServices"`
}

// GetMyServicesDetail retrieves services assigned to current user including permission/policy details
func GetMyServicesDetail(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	var response = MyServiceDetail{
		User:       models.User{},
		MyServices: make([]models.MyServiceDetails, 0),
	}
	//logrus.Debug(userContext)
	//email is needed in calculating adhoc permissions because requester id is email
	userApps, err := Store.GetUserAppsDetailsWithPolicyFromUserID(userContext.User.Groups, userContext.User.ID, userContext.Org.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logrus.Error(err)
			utils.TrasaResponse(w, http.StatusOK, "failed", "no services assigned", "GetMyServicesDetail", response)
			return
		}
		utils.TrasaResponse(w, http.StatusOK, "failed", "could not get apps", "GetMyServicesDetail", response)
		return

	}

	now := time.Now()
	loc, err := time.LoadLocation(userContext.Org.Timezone)
	if err == nil {
		now = now.In(loc)
	}
	current := now.Unix()

	var tempApps = make([]models.MyServiceDetails, 0)

	for _, myApp := range userApps {
		//1) check adhoc is adhoc is enabled
		if myApp.Adhoc == true {
			if (myApp.AuthorizedTill / 1000) < current {
				myApp.IsAuthorised = false
				if myApp.RequestedOn == 0 {
					myApp.Reason = "AdHoc permission failed, not requested"
				} else {
					myApp.Reason = "Request Pending"
				}

			} else {
				myApp.IsAuthorised = true
				myApp.Reason = "Adhoc Success"
			}

		} else {
			// 2) we check users regular policy
			checkPermission, reason := accesscontrol.CheckTrasaUAC(userContext.Org.Timezone, utils.GetIp(r), &myApp.Policy)
			myApp.Reason = reason
			if checkPermission == true {
				myApp.IsAuthorised = true
				myApp.Reason = "user authorised by uac check"

			} else {
				// 3) we check if if adhoc permission is set explicitly

				if (myApp.AuthorizedTill / 1000) < current {
					myApp.IsAuthorised = false
					if myApp.RequestedOn == 0 {
						myApp.Reason = reason
					} else {
						myApp.Reason = "Request Pending"
					}

				} else {
					myApp.IsAuthorised = true
					myApp.Reason = "UAC denied but Adhoc Success"
				}
			}
		}

		tempApps = append(tempApps, myApp)

	}

	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "could not get apps", "GetMyServicesDetail", response)
		return
	}
	response.User = *userContext.User
	response.MyServices = tempApps

	utils.TrasaResponse(w, http.StatusOK, "success", "", "GetMyServicesDetail", response)
	return
}

//RemoveUserDevice removes user device
func RemoveMyDevice(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	deviceID := chi.URLParam(r, "deviceID")

	fmt.Println(deviceID, "deviceID")

	dev, err := devices.Store.GetFromID(deviceID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to delete user device.", "failed to remove user device")
		return
	}

	if dev.UserID != userContext.User.ID {
		logrus.Error(err)
		utils.TrasaResponse(w, 403, "failed", "you do not have permission to remove this device", "failed to remove user device")
		return
	}

	err = devices.Store.Deregister(deviceID, userContext.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to delete user device.", "failed to remove user device")
		return
	}

	utils.TrasaResponse(w, 200, "success", "user device removed", "user device removed")

}

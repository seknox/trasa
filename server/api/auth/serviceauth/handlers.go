package serviceauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/seknox/trasa/server/api/accessmap"
	"github.com/seknox/trasa/server/api/auth"
	"github.com/seknox/trasa/server/api/auth/tfa"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/policies"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

// AppLogin mocks request structure which ssh logins and rdp logins generates
//Deprecated
type appLoginLegacy struct {
	AppID           string `json:"appID"`
	DynamicAuthApp  bool   `json:"dynamicAuthApp"`
	AppSecret       string `json:"appSecret"`
	User            string `json:"user"`
	Password        string `json:"password"`
	PublicKey       []byte `json:"publicKey"`
	TfaMethod       string `json:"tfaMethod"`
	TotpCode        string `json:"totpCode"`
	UserIP          string `json:"userIP"`
	UserWorkstation string `json:"workstation"`
	TrasaID         string `json:"trasaID"`
	SessionID       string `json:"sessionID"`
	AppType         string `json:"appType"`
	//RdpProtocol     string           `json:"rdpProto"`
	OrgID    string `json:"orgID"`
	Hostname string `json:"hostname"`
	//SignResponse    u2f.SignResponse `json:"signResponse"`
	//DeviceHygiene   DeviceHygiene    `json:"deviceHygiene"`
}

type serviceAgentLogin struct {
	ServiceID       string `json:"serviceID"`
	DynamicAuthApp  bool   `json:"dynamicAuthApp"`
	ServiceKey      string `json:"serviceKey"`
	User            string `json:"user"`
	Password        string `json:"password"`
	PublicKey       []byte `json:"publicKey"`
	TfaMethod       string `json:"tfaMethod"`
	TotpCode        string `json:"totpCode"`
	UserIP          string `json:"userIP"`
	UserWorkstation string `json:"workstation"`
	TrasaID         string `json:"trasaID"`
	SessionID       string `json:"sessionID"`
	ServiceType     string `json:"serviceType"`
	//RdpProtocol     string           `json:"rdpProto"`
	OrgID    string `json:"orgID"`
	Hostname string `json:"hostname"`
	//SignResponse    u2f.SignResponse `json:"signResponse"`
	//DeviceHygiene   DeviceHygiene    `json:"deviceHygiene"`
}

func AgentLogin(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("Agent request received")

	var remoteLogin serviceAgentLogin

	if err := json.NewDecoder(r.Body).Decode(&remoteLogin); err != nil {
		utils.TrasaResponse(w, 200, "failed", "invalid request", "AgentLogin", nil)
		logrus.Error(err)
		return
	}

	nativeLogEnabled := true

	logLoginFunc := func(authlog *logs.AuthLog, reason consts.FailedReason, status bool) error {
		if nativeLogEnabled {
			return logs.Store.LogLogin(authlog, consts.REASON_INVALID_SERVICE_CREDS, false)
		}
		return nil
	}

	authlog := logs.NewEmptyLog("rdp")

	temp := strings.Split(remoteLogin.User, "\\")
	if len(temp) > 1 && temp[0] == "" {
		remoteLogin.User = temp[1]
	}

	remoteLogin.User = strings.ToLower(remoteLogin.User)

	authlog.Privilege = remoteLogin.User

	// Get app secret from app ID
	service, err := services.Store.GetFromID(remoteLogin.ServiceID)
	if err != nil || service.SecretKey != remoteLogin.ServiceKey {
		logrus.Errorf("invalid service integration details: %v", err)
		err = logLoginFunc(&authlog, consts.REASON_INVALID_SERVICE_CREDS, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", "invalid service integration details", "AgentLogin", nil)
		return
	}
	if service.Type == "rdp" {
		remoteLogin.UserIP = utils.GetIp(r)
	}

	authlog.UpdateIP(remoteLogin.UserIP)

	authlog.UpdateService(service)
	nativeLogEnabled = service.NativeLog

	usernameExits := accessmap.Store.CheckIfPrivilegeExist(remoteLogin.User, service.OrgID, remoteLogin.ServiceID)
	if service.Passthru == true && !usernameExits {
		utils.TrasaResponse(w, 200, "success", "passthru authentication", "agent-login", nil)
		return
	}

	userDetails, err := auth.Store.GetLoginDetails(remoteLogin.TrasaID, "domain")
	if err != nil {

		err = logLoginFunc(&authlog, consts.REASON_USER_NOT_FOUND, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", "user not found", "AgentLogin", nil)
		return
	}

	authlog.UpdateUser(userDetails)

	if !userDetails.Status {
		utils.TrasaResponse(w, http.StatusOK, "failed", "User Disabled", "AgentLogin", nil)

		err = logLoginFunc(&authlog, consts.REASON_USER_DISABLED, false)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	//logrus.Debug(service.OrgID)

	orgDetail, err := orgs.Store.Get(service.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "INVALID_ORG_ID", "agent-login", nil)
		return
	}

	policy, privilege, adhoc, err := policies.Store.GetAccessPolicy(userDetails.ID, service.ID, orgDetail.ID)
	if err != nil {
		logrus.Debug(err)
		err = logLoginFunc(&authlog, consts.REASON_NO_POLICY_ASSIGNED, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", "policy not assigned", "agent-login", nil)
		return
	}

	normalizedPrivilege := utils.NormalizeString(privilege)
	normalizedRemoteUsername := utils.NormalizeString(remoteLogin.User)

	if normalizedPrivilege != normalizedRemoteUsername {
		err = logLoginFunc(&authlog, consts.REASON_INVALID_PRIVILEGE, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", fmt.Sprintf("privilege assigned is (%s) but received %s ", normalizedPrivilege, normalizedRemoteUsername), "agent-login", nil)
		return
	}

	ok, reason := Store.CheckPolicy(service.ID, userDetails.ID, orgDetail.ID, remoteLogin.UserIP, orgDetail.Timezone, policy, adhoc)
	if !ok {
		err = logLoginFunc(&authlog, reason, false)
		if err != nil {
			logrus.Error(err)
		}
		utils.TrasaResponse(w, 200, "failed", string(reason), "agent-login", nil)
		return
	}

	//}
	remoteLogin.TfaMethod = "totp"
	if remoteLogin.TotpCode == "" {
		remoteLogin.TfaMethod = "u2f"
	}

	tfaDeviceID, reason, ok := tfa.HandleTfaAndGetDeviceID(
		nil,
		remoteLogin.TfaMethod,
		remoteLogin.TotpCode,
		userDetails.ID,
		remoteLogin.UserIP,
		service.Name,
		orgDetail.Timezone,
		orgDetail.OrgName,
		orgDetail.ID,
	)

	authlog.TfaDeviceID = tfaDeviceID

	err = logLoginFunc(&authlog, reason, ok)
	if err != nil {
		logrus.Error(err)
	}

	if !ok {
		utils.TrasaResponse(w, 200, "failed", string(reason), "agent-login", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", string(reason), "agent-login", nil)
	logrus.Trace("Agent rlogin response returned")
	return

}

package serviceauth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/seknox/trasa/server/api/accessmap"
	"github.com/seknox/trasa/server/api/auth"
	"github.com/seknox/trasa/server/api/auth/tfa"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/policies"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func DBLogin(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("Agent request received")

	var remoteLogin ServiceAgentLogin

	if err := json.NewDecoder(r.Body).Decode(&remoteLogin); err != nil {
		utils.TrasaResponse(w, 200, "failed", "invalid request", "AgentLogin", nil)
		logrus.Error(err)
		return
	}

	authlog := logs.NewEmptyLog("db")
	authlog.UpdateIP(remoteLogin.UserIP)

	remoteLogin.User = strings.ToLower(remoteLogin.User)

	authlog.Privilege = remoteLogin.User

	userDetails, err := auth.Store.GetLoginDetails(remoteLogin.TrasaID, "domain")
	if err != nil {
		logs.Store.LogLogin(&authlog, consts.REASON_USER_NOT_FOUND, false)
		utils.TrasaResponse(w, 200, "failed", "user not found", "AgentLogin", nil)
		return
	}

	authlog.UpdateUser(userDetails)

	if !userDetails.Status {
		utils.TrasaResponse(w, http.StatusOK, "failed", "User Disabled", "AgentLogin", nil)
		logs.Store.LogLogin(&authlog, consts.REASON_USER_DISABLED, false)
		return
	}

	service, err := services.Store.GetFromHostname(remoteLogin.Hostname, "db", "", userDetails.OrgID)
	if errors.Is(err, sql.ErrNoRows) {
		service, err = accessmap.CreateDynamicService(remoteLogin.Hostname, "db", userDetails.ID, remoteLogin.TrasaID, remoteLogin.User, remoteLogin.OrgID)
		if err != nil {
			logrus.Debug(err)
			authlog.FailedReason = consts.REASON_DYNAMIC_SERVICE_FAILED
			logs.Store.LogLogin(&authlog, consts.REASON_DYNAMIC_SERVICE_FAILED, false)
			utils.TrasaResponse(w, http.StatusOK, "failed", "could not create dynamic service", "DBLogin", nil)
			return

		}
	} else if err != nil {
		logrus.Error(err)
		logs.Store.LogLogin(&authlog, consts.REASON_INVALID_SERVICE_CREDS, false)
		authlog.FailedReason = consts.REASON_INVALID_SERVICE_CREDS
		utils.TrasaResponse(w, http.StatusOK, "failed", "invalid hostname", "DBLogin", nil)
		return
	}

	authlog.UpdateService(service)

	orgDetail, err := orgs.Store.Get(service.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "INVALID_ORG_ID", "db-login", nil)
		return
	}

	policy, adhoc, err := policies.Store.GetAccessPolicy(userDetails.ID, service.ID, utils.NormalizeString(remoteLogin.User), orgDetail.ID)
	if err != nil {
		logrus.Debug(err)
		logs.Store.LogLogin(&authlog, consts.REASON_NO_POLICY_ASSIGNED, false)
		utils.TrasaResponse(w, 200, "failed", "policy not assigned", "db-login", nil)
		return
	}
	authlog.SessionRecord = policy.RecordSession

	//normalizedPrivilege := utils.NormalizeString(privilege)
	//normalizedRemoteUsername := utils.NormalizeString(remoteLogin.User)
	//
	//if normalizedPrivilege != normalizedRemoteUsername {
	//	err = logs.Store.LogLogin(&authlog, consts.REASON_INVALID_PRIVILEGE, false)
	//	if err != nil {
	//		logrus.Error(err)
	//	}
	//	utils.TrasaResponse(w, 200, "failed", fmt.Sprintf("privilege assigned is (%s) but received %s ", normalizedPrivilege, normalizedRemoteUsername), "db-login", nil)
	//	return
	//}

	ok, reason := Store.CheckPolicy(service.ID, userDetails.ID, orgDetail.ID, remoteLogin.UserIP, orgDetail.Timezone, policy, adhoc)
	if !ok {
		logs.Store.LogLogin(&authlog, reason, false)
		utils.TrasaResponse(w, 200, "failed", string(reason), "db-login", nil)
		return
	}

	//TODO uncomment this after implementing authentication in this API

	//creds, err := services.GetUpstreamCreds(remoteLogin.User, service.ID, service.Type, service.OrgID)
	//if err != nil {
	//	logrus.Error(err)
	//	creds = &models.UpstreamCreds{}
	//}
	creds := models.UpstreamCreds{}

	//}
	remoteLogin.TfaMethod = "totp"
	if remoteLogin.TotpCode == "" {
		remoteLogin.TfaMethod = "u2f"
	}

	tfaDeviceID := ""

	if policy.TfaRequired {
		tfaDeviceID, reason, ok = tfa.HandleTfaAndGetDeviceID(
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

	}

	authlog.TfaDeviceID = tfaDeviceID

	logs.Store.LogLogin(&authlog, reason, ok)

	if !ok {
		utils.TrasaResponse(w, 200, "failed", string(reason), "db-login", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", string(reason), "db-login", creds, policy.RecordSession, authlog.SessionID)
	logrus.Trace("DB login response returned")
	return

}

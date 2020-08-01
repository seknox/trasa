package serviceauth

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/seknox/trasa/server/api/accessmap"
	"github.com/seknox/trasa/server/api/auth"
	"github.com/seknox/trasa/server/api/auth/tfa"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/policies"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/sirupsen/logrus"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

// RadiusLogin is a handler for incoming radius auth request
func RadiusLogin(w radius.ResponseWriter, r *radius.Request) {
	logrus.Trace("RadiusLogin login request received")

	remoteUser := rfc2865.UserName_GetString(r.Packet)
	usernameWithTfaMethod := strings.Split(remoteUser, ":")
	trasaID := usernameWithTfaMethod[0]
	password := rfc2865.UserPassword_GetString(r.Packet)

	ipAddr := strings.Split(r.RemoteAddr.String(), ":")

	service, err := services.Store.GetFromHostname(ipAddr[0], "radius", "", global.GetConfig().Trasa.OrgId)
	if err != nil {
		logrus.Error(err)
		w.Write(r.Response(radius.CodeAccessReject))
	}

	nativeLogEnabled := true

	logLoginFunc := func(authlog *logs.AuthLog, reason consts.FailedReason, status bool) error {
		if nativeLogEnabled {
			return logs.Store.LogLogin(authlog, consts.REASON_INVALID_SERVICE_CREDS, false)
		}
		return nil
	}

	authlog := logs.NewEmptyLog("radius")

	authlog.Privilege = trasaID

	authlog.UpdateService(service)

	userDetails, err := auth.Store.GetLoginDetails(trasaID, "domain")
	if err != nil {
		usernameExits := accessmap.Store.CheckIfPrivilegeExist(trasaID, service.OrgID, service.ID)
		logrus.Trace("passthru test ", service.Passthru, usernameExits)
		if service.Passthru == true && !usernameExits {
			w.Write(r.Response(radius.CodeAccessReject))
			return
		}
		err = logLoginFunc(&authlog, consts.REASON_USER_NOT_FOUND, false)
		if err != nil {
			logrus.Error(err)
		}
		w.Write(r.Response(radius.CodeAccessReject))
		return
	}

	reason, err := auth.CheckPassword(userDetails, trasaID, password)
	if err != nil {
		w.Write(r.Response(radius.CodeAccessReject))

		err = logLoginFunc(&authlog, consts.REASON_INVALID_USER_CREDS, false)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	authlog.UpdateUser(userDetails)

	if !userDetails.Status {
		w.Write(r.Response(radius.CodeAccessReject))

		err = logLoginFunc(&authlog, consts.REASON_USER_DISABLED, false)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	logrus.Debug(service.OrgID)

	orgDetail, err := orgs.Store.Get(service.OrgID)
	if err != nil {
		logrus.Error(err)
		w.Write(r.Response(radius.CodeAccessReject))
		return
	}

	policy, privilege, adhoc, err := policies.Store.GetAccessPolicy(userDetails.ID, service.ID, orgDetail.ID)
	if err != nil {
		logrus.Debug(err)
		err = logLoginFunc(&authlog, consts.REASON_NO_POLICY_ASSIGNED, false)
		if err != nil {
			logrus.Error(err)
		}
		w.Write(r.Response(radius.CodeAccessReject))
		return
	}

	if privilege != trasaID {
		err = logLoginFunc(&authlog, consts.REASON_INVALID_PRIVILEGE, false)
		if err != nil {
			logrus.Error(err)
		}
		w.Write(r.Response(radius.CodeAccessReject))
		return
	}

	ok, reason := Store.CheckPolicy(service.ID, userDetails.ID, orgDetail.ID, ipAddr[0], orgDetail.Timezone, policy, adhoc)
	if !ok {
		err = logLoginFunc(&authlog, reason, false)
		if err != nil {
			logrus.Error(err)
		}
		w.Write(r.Response(radius.CodeAccessReject))
		return
	}

	tfaMethod := "u2f"
	totpVal := "0000"
	if len(usernameWithTfaMethod) == 2 {
		tfaMethod = "totp"
		totpVal = usernameWithTfaMethod[1]
	}

	tfaDeviceID, reason, ok := tfa.HandleTfaAndGetDeviceID(
		nil,
		tfaMethod,
		totpVal,
		userDetails.ID,
		ipAddr[0],
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
		w.Write(r.Response(radius.CodeAccessReject))
		return
	}

	w.Write(r.Response(radius.CodeAccessAccept))
	logrus.Trace("Agent rlogin response returned")
	return

}

type printIP struct {
	Framed net.IP
	LOgin  net.IP
	Nas    net.IP
}

func HandleRadiusReq(w radius.ResponseWriter, r *radius.Request) {
	trasaID := rfc2865.UserName_GetString(r.Packet)
	password := rfc2865.UserPassword_GetString(r.Packet)

	var p printIP
	p.Framed = rfc2865.FramedIPAddress_Get(r.Packet)
	p.LOgin = rfc2865.LoginIPHost_Get(r.Packet)
	p.Nas = rfc2865.NASIPAddress_Get(r.Packet)

	v, err := json.Marshal(p)
	if err != nil {
		logrus.Error(err)
	}

	fmt.Println(string(v))
	// verify password.

	fmt.Printf("user: %s, password: %s\n", trasaID, password)

	authlog := logs.NewEmptyLog("radius")

	// get user info from database
	userDetails, err := auth.Store.GetLoginDetails(trasaID, "")
	if err != nil {
		logrus.Error(err)
		err = logs.Store.LogLogin(&authlog, consts.REASON_USER_NOT_FOUND, false)
		if err != nil {
			logrus.Error(err)
		}
		w.Write(r.Response(radius.CodeAccessReject))
		return
	}

	authlog.UpdateUser(userDetails)

	if !userDetails.Status {
		err = logs.Store.LogLogin(&authlog, consts.REASON_USER_DISABLED, false)
		if err != nil {
			logrus.Error(err)
		}

		w.Write(r.Response(radius.CodeAccessReject))
		return
	}

	_, err = auth.CheckPassword(userDetails, trasaID, password)
	if err != nil {
		logrus.Error(err)
		err = logs.Store.LogLogin(&authlog, consts.REASON_INVALID_USER_CREDS, false)
		if err != nil {
			logrus.Error(err)
		}
	}
	w.Write(r.Response(radius.CodeAccessAccept))
	return
}

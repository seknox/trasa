package sshproxy

import (
	"database/sql"
	"encoding/hex"
	"strings"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/api/accessmap"

	"github.com/gorilla/websocket"
	"github.com/seknox/ssh"
	"github.com/seknox/trasa/server/api/accesscontrol"
	"github.com/seknox/trasa/server/api/auth/tfa"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/sirupsen/logrus"
)

func JoinSSHSession(params models.ConnectionParams, uc models.UserContext, conn *websocket.Conn) {

	//Do not defer conn.Close() because it will immediately close the websocket connection and the live session will not work.
	//All the guest joined connections are closed in WrappedTunnel.Close() method
	conn.WriteMessage(1, []byte("\n\rConnecting...\n\r"))

	checkAndInitParams(&uc, &params)
	params.IsSharedSession = true
	params.ServiceType = "ssh"

	service, err := services.Store.GetFromHostname(params.Hostname, "ssh", "", params.OrgID)
	if err != nil {
		logrus.Error(err)
		conn.WriteMessage(1, []byte("\n\rService not created\n\r"))
		return
	}

	_, reason, ok := tfa.HandleTfaAndGetDeviceID(nil, params.TfaMethod, params.TotpCode, uc.User.ID, params.UserIP, service.Name, uc.Org.Timezone, uc.Org.OrgName, uc.Org.ID)
	if !ok {
		conn.WriteMessage(1, []byte(reason))
		conn.Close()
		return
	}

	guest := GuestClient{
		UserID: uc.User.ID,
		Email:  uc.User.Email,
		Conn:   conn,
	}

	sessionID := params.ConnID

	newViewer, err := SSHStore.GetGuestChannel(sessionID)
	if err != nil || newViewer == nil {
		logrus.Error(err)
		conn.WriteMessage(1, []byte("\n\rNo such connection\n\r"))
		//	conn.Write([]byte("No such connection"))
		conn.Close()
		return
	}
	conn.WriteMessage(1, []byte("\n\rConnected to SSH session.\n\r"))
	newViewer <- guest

	//
	//if err := conn.WriteMessage(1, <-tshp.Sshchan); err != nil {
	//	logrus.Debug(err)
	//	return
	//}

	//}

}

// ConnectNewSSH handles new ssh connection from dashboard.
func ConnectNewSSH(params models.ConnectionParams, uc models.UserContext, conn *websocket.Conn) {

	defer conn.Close()
	conn.WriteMessage(1, []byte("\n\rConnecting...\n\r"))

	params.AccessDeviceID = uc.DeviceID
	params.BrowserID = uc.BrowserID

	checkAndInitParams(&uc, &params)
	authlog := logs.NewEmptyLog("ssh")
	authlog.UpdateUser(&models.UserWithPass{User: *uc.User})
	authlog.Privilege = params.Privilege
	authlog.ServerIP = params.Hostname
	params.SessionID = authlog.SessionID

	authlog.AccessDeviceID = uc.DeviceID
	authlog.BrowserID = uc.BrowserID

	authlog.UpdateIP(params.UserIP)
	params.ServiceType = "ssh"

	defer logSession(&authlog)

	service, err := services.Store.GetFromHostname(params.Hostname, "ssh", "", uc.Org.ID)
	if errors.Is(err, sql.ErrNoRows) {
		service, err = accessmap.CreateDynamicService(params.Hostname, "ssh", params.UserID, params.TrasaID, params.Privilege, params.OrgID)
		if err != nil {
			logrus.Errorf("dynamic access: %v", err)
			authlog.FailedReason = consts.REASON_DYNAMIC_SERVICE_FAILED
			conn.WriteMessage(1, []byte("\n\rService not created\n\r"))
			return
		}

	} else if err != nil {
		logrus.Errorf("get service from hostname: %v", err)
		authlog.FailedReason = consts.REASON_INVALID_SERVICE_CREDS
		conn.WriteMessage(1, []byte("\n\rService does  not created\n\r"))
		return
	}

	authlog.UpdateService(service)
	params.ServiceID = service.ID
	params.ServiceName = service.Name

	creds, err := services.GetUpstreamCreds(params.Privilege, service.ID, "ssh", uc.Org.ID)
	if err != nil {
		logrus.Error(err)
		conn.WriteMessage(1, []byte("\n\rSomething is wrong\n\r"))
		return
	}

	//caKey,err:=ssh.ParsePublicKey([]byte(resp.CaCert))

	//
	//passwordStreanght := zxcvbn.PasswordStrength(params.Password, nil)
	//
	//if creds.EnforceStrongPass && creds.ClientKey == "" {
	//	if passwordStreanght.Score < creds.ZxcvbnScore || len(params.Password) < creds.MinimumChar {
	//		conn.WriteMessage(1, []byte("\n\rPassword policy failed, weak password\n\r"))
	//		logrus.Debug("Weak Password")
	//		err := logs.Store.LogLogin(&authlog, consts.REASON_PASSWORD_POLICY_FAILED, false)
	//		if err != nil {
	//			logrus.Error(err)
	//		}
	//		return
	//	}
	//}

	clientConfig := ssh.ClientConfig{
		User: params.Privilege,
		//HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		BannerCallback: func(message string) error {
			//forward banner to client through websocket
			return conn.WriteMessage(1, []byte(message))
		},
	}

	//Checking TRASA policy
	conn.WriteMessage(1, []byte("\n\rChecking policy...\n\r"))

	policy, adhoc, err := accessmap.GetAssignedPolicy(&params)
	if err != nil {
		logrus.Debug(err)
		conn.WriteMessage(1, []byte("\n\rPolicy not assigned "+"\n\r"))
		logs.Store.LogLogin(&authlog, consts.REASON_NO_POLICY_ASSIGNED, false)
		return
	}

	ok, reason := accesscontrol.CheckPolicy(&params, policy, adhoc)
	if !ok {
		logrus.Debug(err)
		conn.WriteMessage(1, []byte("\n\rPolicy check failed: "+reason+"\n\r"))
		logs.Store.LogLogin(&authlog, reason, false)
		return
	}

	//Handle second factor if enabled in policy
	if policy.TfaRequired {
		conn.WriteMessage(1, []byte("\n\rAuthenticating 2nd Factor\n\r"))
		deviceID, reason, ok := tfa.HandleTfaAndGetDeviceID(nil,
			params.TfaMethod,
			params.TotpCode,
			uc.User.ID,
			params.UserIP,
			service.Name,
			uc.Org.Timezone,
			uc.Org.OrgName,
			uc.Org.ID)
		if !ok {
			conn.WriteMessage(1, []byte("\n\r2FA failed: "+reason+"\n\r"))
			err = logs.Store.LogLogin(&authlog, reason, false)
			if err != nil {
				logrus.Error(err)
			}
			return
		}
		authlog.TfaDeviceID = deviceID
	}

	//logrus.Trace(params.AccessDeviceID)
	reason, ok, err = accesscontrol.CheckDevicePolicy(policy.DevicePolicy, params.AccessDeviceID, authlog.TfaDeviceID, uc.Org.ID)
	if err != nil {
		logrus.Error(err)
	}

	if !ok {
		conn.WriteMessage(1, []byte("\n\rDevice policy failed: "+reason+"\n\r"))
		err = logs.Store.LogLogin(&authlog, reason, false)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	authlog.SessionRecord = policy.RecordSession

	//Callback function to check ssh host key of upstream server
	var hostConfirmFunc = func(message string) bool {
		err := conn.WriteMessage(1, []byte("\r\n"+message+"\n\rPress \"y\" to ignore and save it.\n\r"))
		if err != nil {
			logrus.Error(err)
			return false
		}

		_, ans, err := conn.ReadMessage()
		if err != nil {
			logrus.Error(err)
			return false
		}
		if strings.ToLower(string(ans)) == "y" {
			conn.WriteMessage(1, []byte("\r\nSaving new host key...\n\r"))
			return true
		}
		return false

	}
	clientConfig.HostKeyCallback = HandleHostKeyCallback(creds, service.ID, uc.Org.ID, hostConfirmFunc)

	//Add public key auth method
	pkeyAuth := getPublicKeyAuthMethod(creds)
	if pkeyAuth != nil {
		clientConfig.Auth = append(clientConfig.Auth, getPublicKeyAuthMethod(creds))
	}

	upstreamPassword := ""
	if params.Password == "" || creds.Password != "" {
		//password from vault
		conn.WriteMessage(1, []byte("\n\rUsing password from vault.\n\r"))
		upstreamPassword = creds.Password
	} else {
		//password entered by user
		upstreamPassword = params.Password
	}

	//Add keyboard-interactive auth method to handle TRASA PAM module installed in upstream server
	clientConfig.Auth = append(clientConfig.Auth,
		handleUpstreamTrasaPAM(params.TrasaID, upstreamPassword, params.TotpCode),
	)

	//Add password auth method
	clientConfig.Auth = append(clientConfig.Auth, ssh.Password(upstreamPassword))

	if !strings.Contains(params.Hostname, ":") {
		params.Hostname = params.Hostname + ":22"
	}

	conn.WriteMessage(1, []byte("\n\rAuthenticating...\n\r"))

	sshClient, err := ssh.Dial("tcp", params.Hostname, &clientConfig)
	if err != nil {
		//logrus.Debug(err)
		if strings.Contains(err.Error(), "ssh: unable to authenticate") {
			logrus.Debug(err)
			conn.WriteMessage(1, []byte("\r\nInvalid credentials.\n\r"))
			err = logs.Store.LogLogin(&authlog, consts.REASON_INVALID_USER_CREDS, false)

		} else if strings.Contains(err.Error(), "ssh: handshake failed: Could not verify upstream host key") {
			logrus.Debug(err)
			conn.WriteMessage(1, []byte("\n\rInvalid host key. Could not verify upstream server.\n\r"))
			err = logs.Store.LogLogin(&authlog, consts.REASON_INVALID_HOST_KEY, false)
		} else if strings.Contains(err.Error(), ErrVerifyHost.Error()) {
			logrus.Debug(err)
			conn.WriteMessage(1, []byte("\n\rInvalid host key. Could not verify upstream server.\n\r"))
			err = logs.Store.LogLogin(&authlog, consts.REASON_INVALID_HOST_KEY, false)
		} else {
			logrus.Error(err)
			conn.WriteMessage(1, []byte(err.Error()))
			err = logs.Store.LogLogin(&authlog, consts.REASON_UNKNOWN, false)
		}
		//TODO add host not reachable message
		if err != nil {
			logrus.Error(err)
		}

		return
	}

	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		logrus.Error(err)
		return
	}
	defer session.Close()

	guestChan := SSHStore.CreateGuestChannel(hex.EncodeToString(sshClient.SessionID()))
	defer SSHStore.deleteGuestChannel(hex.EncodeToString(sshClient.SessionID()))

	logs.Store.AddNewActiveSession(&authlog, hex.EncodeToString(sshClient.SessionID()), "ssh")
	defer logs.Store.RemoveActiveSession(hex.EncodeToString(sshClient.SessionID()))
	//logrus.Trace("SESSION ID is", hex.EncodeToString(sshClient.SessionID()))

	wsshFrontEndConn := NewWebSSHFrontEndConn(conn)
	wsshBackendConn, err := NewWebSSHBackend(session)
	if err != nil {
		logrus.Error(err)
		conn.WriteMessage(1, []byte("\n\rCould not create std in/out pipe\n\r"))
		wsshFrontEndConn.Close()
		return
	}

	wrappedChannel, err := NewWrappedTunnel(authlog.SessionID, policy.RecordSession, wsshBackendConn, wsshFrontEndConn, guestChan)
	if err != nil {
		logrus.Error(err)
		return
	}

	//
	//modes := ssh.TerminalModes{
	//	ssh.ECHO:          1,
	//	ssh.ECHOCTL:       0,
	//	ssh.TTY_OP_ISPEED: 14400,
	//	ssh.TTY_OP_OSPEED: 14400,
	//}

	// Request pseudo terminal
	err = session.RequestPty("xterm", int(params.OptHeight), int(params.OptWidth), nil)
	if err != nil {
		//if err := session.RequestPty("xterm-256color", 80, 40, modes); err != nil {
		//if err := session.RequestPty("vt100", 80, 40, modes); err != nil {
		//if err := session.RequestPty("vt220", 80, 40, modes); err != nil {
		logrus.Errorf("request for pseudo terminal failed: %s", err)
		return
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		logrus.Errorf("failed to start shell: %s", err)
		return
	}

	wrappedChannel.pipe()

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Subprotocols:    []string{"xterm"},
}

func checkAndInitParams(uc *models.UserContext, params *models.ConnectionParams) {
	//TODO
	params.OrgID = uc.Org.ID
	params.UserID = uc.User.ID
	params.TrasaID = uc.User.Email
	params.Timezone = uc.Org.Timezone
	params.ServiceType = "ssh"
	params.Groups = uc.User.Groups
	//params.UserAgent = r.UserAgent()

	if params.RdpProtocol == "" {
		params.RdpProtocol = "nla"
	}

}

func logSession(authlog *logs.AuthLog) {

	err := logs.Store.LogLogin(authlog, "", true)
	if err != nil {
		logrus.Errorf("failed to log.  trying again: %v", err)
		logs.Store.LogLogin(authlog, "", true)
	}

	if !authlog.SessionRecord {
		return
	}

	err = SSHStore.uploadSessionLog(authlog)
	if err != nil {
		logrus.Error(err)
	}

}

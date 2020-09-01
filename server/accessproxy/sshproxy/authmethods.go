package sshproxy

import (
	"database/sql"
	"github.com/seknox/trasa/server/api/accesscontrol"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/api/accessmap"
	"github.com/seknox/trasa/server/api/auth"
	"github.com/seknox/trasa/server/api/auth/tfa"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/trustelem/zxcvbn"

	//"crypto/tls"
	"fmt"
	"net"

	logrus "github.com/sirupsen/logrus"

	"github.com/seknox/ssh"

	"strings"
	"time"
)

const (
	gotoKeyboardInteractive = "trasa: goto_keyboard_interactive"
	gotoPublicKey           = "trasa: goto_public_key"
	//gotoPublicKeyOrKeyboardInteractive           = "trasa: goto_public_key_or_kb_interactive"
	failNow = "trasa: fail_now"
)

//Decides which auth method to use next from previous error
func nextAuthMethodHandler(conn ssh.ConnMetadata, prevErr error) ([]string, bool, error) {

	switch prevErr.Error() {
	case gotoKeyboardInteractive:
		return []string{"keyboard-interactive"}, true, nil
	case "ssh: no auth passed yet":
		return []string{"publickey", "keyboard-interactive"}, false, nil
	case gotoPublicKey:
		return []string{"publickey", "keyboard-interactive"}, false, nil
	case failNow:
		return []string{}, false, prevErr
	default:
		if prevErr != nil {
			return []string{"keyboard-interactive"}, false, nil
		} else {

			//TODO check this
			//if sessions[conn.RemoteAddr()]["authMeta"] == nil {
			//	return []string{"publickey"}, false, nil
			//}
			return []string{"keyboard-interactive"}, false, nil
		}

	}
}

func handleUpstreamPasswordAndKey(username, serviceID, hostname string, challengeUser ssh.KeyboardInteractiveChallenge) (ssh.Signer, ssh.HostKeyCallback, string, error) {

	var password string

	creds, err := services.GetUpstreamCreds(username, serviceID, "ssh", global.GetConfig().Trasa.OrgId)
	if err != nil {
		logrus.Debug(err)
		return nil, nil, "", err
	}

	publicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(strings.TrimSpace(creds.ClientCert)))
	if err != nil && creds.ClientCert != "" {
		logrus.Error(err)
	}

	privateKey, privateKeyParseErr := ssh.ParsePrivateKey([]byte(creds.ClientKey))
	if privateKeyParseErr != nil && creds.ClientKey != "" {
		//logrus.Debug(creds.ClientKey)
		logrus.Error(privateKeyParseErr)
	}

	var signer ssh.Signer
	cert, ok := publicKey.(*ssh.Certificate)
	if !ok {
		logrus.Debug("Invalid user certificate")
		if privateKeyParseErr == nil {
			signer = privateKey
		}

	} else {
		signer, err = ssh.NewCertSigner(cert, privateKey)
		if err != nil {
			logrus.Debug(err)
		}
	}

	host := ""
	if strings.Contains(hostname, ":") {
		host = hostname
	} else {
		host = hostname + ":22"
	}

	var hostConfirmFunc = func(message string) bool {
		ans, err := challengeUser("user", "Host key verify",
			[]string{message + "\n\rType \"yes\" to ignore and save host key:\n"},
			[]bool{true})
		if err != nil || len(ans) != 1 {
			return false
		}
		if strings.ToLower(ans[0]) == "yes" {
			return true
		}
		return false
	}

	tempConf := ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(creds.Password)},
		HostKeyCallback: HandleHostKeyCallback(creds, serviceID, global.GetConfig().Trasa.OrgId, hostConfirmFunc),
	}

	if signer != nil {
		tempConf.Auth = append(tempConf.Auth, ssh.PublicKeys(signer))
	}

	c, err := ssh.Dial("tcp", host, &tempConf)

	//logrus.Debug(err)

	if err == nil {
		password = creds.Password
		c.Close()
	} else if strings.Contains(err.Error(), "trasa: could not update cert") {
		challengeUser("", "Could not verify host", nil, nil)
		return signer, nil, "", errors.New("Could not verify host")
	} else if strings.Contains(err.Error(), "unable to authenticate") && (strings.Contains(err.Error(), "password")) {
		//password error

		logrus.Debug("password error")

		var ans []string
		if creds.Password == "" {
			ans, err = challengeUser("user", "Upstream password",
				[]string{"\n\rEnter Password(Upstream Server): "},
				[]bool{false})
			if err != nil {
				return signer, nil, password, err

			}
			if len(ans) != 1 {
				return signer, nil, "", ssh.ErrNoAuth
			}

			password = ans[0]
		} else {
			password = creds.Password
		}

	} else {
		logrus.Trace(err)
	}

	passwordStreanght := zxcvbn.PasswordStrength(password, nil)

	if creds.EnforceStrongPass && creds.ClientKey == "" {
		if passwordStreanght.Score < creds.ZxcvbnScore || len(password) < creds.MinimumChar {
			logrus.Debug("Weak Password")
			challengeUser("", "Weak password. Password policy not met", nil, nil)
			return signer, nil, "", errors.New("Weak Password")
		}
	}

	var allowHost = func(msg string) bool { return true }

	return signer, HandleHostKeyCallback(creds, serviceID, global.GetConfig().Trasa.OrgId, allowHost), password, nil

	//return creds,password,nil

	//When logging in from certificate always reject unknown certificate

	//clientConfig.SetDefaults()
	//		clientConfig.Ciphers = append(clientConfig.Ciphers, "aes128-cbc", "blowfish-cbc", "3des-cbc")

}

func authenticateTRASA(conn ssh.ConnMetadata, challengeUser ssh.KeyboardInteractiveChallenge) (models.User, error) {
	user := models.User{}
	creds, err := challengeUser("user",
		"Enter TRASA credentials",
		[]string{"\n\rEnter Email (TRASA): ", "\n\rEnter Password (TRASA): "},
		[]bool{true, false})

	if len(creds) != 2 {
		logrus.Debug("User canceled")
		return user, ssh.ErrNoAuth
	}
	email := creds[0]
	pass := creds[1]

	//TODO enter domain
	userDetail, err := auth.Store.GetLoginDetails(email, "")
	if err != nil {
		logrus.Error(err)
		challengeUser("", "Email not found", nil, nil)
		return user, fmt.Errorf("get login details: %v", err)
	}

	user = models.CopyUserWithoutPass(*userDetail)
	err = SSHStore.UpdateSessionUser(conn.RemoteAddr(), &user)
	if err != nil {
		logrus.Errorf("update session user: %v", err)
		challengeUser("", "Something is wrong", nil, nil)
		return user, fmt.Errorf("update session user: %v", err)
	}

	reason, err := auth.CheckPassword(userDetail, email, pass)
	if err != nil {
		logrus.Debug(err)
		challengeUser("", string(reason), nil, nil)
		return user, fmt.Errorf("check password. %s", err)
	}

	//logrus.Debug(utils.MarshallStruct(respData))

	return user, nil

}

//This function takes username, valid applist and keyboardInteractive callback function to  get user to choose authapp/hostname
//It returns instance of chosen appuser (app with username) and  isDynamicApp boolean if chosen app is not yet created or assgned
func chooseService(privilege, userID, userEmail string, challengeUser ssh.KeyboardInteractiveChallenge) (*models.Service, error) {

	var isHostDown bool = true
	hostname := ""

	//take input(upstream server) from user and validate
	for isHostDown {

		ans, err := challengeUser("user",
			"Choose Service",
			[]string{"\n\r_____________________________________________________________________________________\n\rEnter Service IP : \n\r"}, []bool{true})
		if len(ans) != 1 || err != nil {
			logrus.Debug("User canceled")
			return nil, fmt.Errorf("User canceled")
		}

		hostname = ans[0]

		h := hostname
		if !strings.Contains(h, ":") {
			h = h + ":22"
		}

		//check if upstream server is down
		tempC, errPing := net.DialTimeout("tcp", h, time.Second*7)

		if errPing != nil {
			isHostDown = true
			challengeUser("", "\n\nThe SSH server is down", nil, nil)
			continue
		}
		tempC.Close()
		isHostDown = false

	}

	service, err := services.Store.GetFromHostname(hostname, "ssh", "", global.GetConfig().Trasa.OrgId)
	if errors.Is(err, sql.ErrNoRows) {
		service, err = accessmap.CreateDynamicService(hostname, "ssh", userID, userEmail, privilege, global.GetConfig().Trasa.OrgId)
		if err != nil {
			logrus.Errorf("dynamic access: %v", err)
			challengeUser("", "Service not assigned. you do nor have dynamic access", nil, nil)
			return nil, errors.WithMessage(err, "dynamic access")
		}

	} else if err != nil {
		logrus.Errorf("get service from hostname: %v", err)
		return nil, errors.WithMessage(err, "get service from hostname")
	}

	return service, nil

}

func keyboardInteractiveHandler(conn ssh.ConnMetadata, challengeUser ssh.KeyboardInteractiveChallenge) (*ssh.Permissions, error) {
	/*if authenticateU2F(conn.User())!="success"{
		return nil,errors.New("Failed to authenticate")
	}*/

	//var accessableServiceDetails []models.AccessMapDetail
	var err error

	err = SSHStore.UpdateSessionMeta(conn.RemoteAddr(), conn)
	if err != nil {
		logrus.Errorf("update session meta: %v", err)
		challengeUser("", "Something is wrong. Contact your admin.", nil, nil)
		return nil, err
	}

	sessionMeta, err := SSHStore.GetSession(conn.RemoteAddr())
	if err != nil {
		logrus.Errorf("get session meta: %v", err)
		challengeUser("", "Something is wrong. Contact your admin.", nil, nil)
		return nil, err
	}

	if sessionMeta.AuthType == "" {
		sessionMeta.AuthType = consts.SSH_AUTH_TYPE_PASSWORD
	}

	//possible value of keyType are:
	// "SSH_AUTH_TYPE_DACERT" which is certificate with embedded data. 2FA is already verified, service is also already chosen.
	// "SSH_AUTH_TYPE_PUB"  which is certificate generated from trasa. need to verify 2FA and ask user to choose app
	// "SSH_AUTH_TYPE_PASSWORD" need to ask user his/her email

	if sessionMeta.AuthType == consts.SSH_AUTH_TYPE_PASSWORD {
		userDetails, err := authenticateTRASA(conn, challengeUser)
		if err != nil {
			challengeUser("", "Login Failed", nil, nil)
			return &ssh.Permissions{}, errors.New("Could not verify TRASA user with email")
		}

		SSHStore.UpdateSessionUser(conn.RemoteAddr(), &userDetails)
		sessionMeta.params.UserID = userDetails.ID
		sessionMeta.params.OrgID = userDetails.OrgID
		sessionMeta.params.TrasaID = userDetails.Email

	}

	//logrus.Debug(sessionMeta.params.UserID, sessionMeta.params.OrgID)
	//call api to authenticate and  enumerate accessible servers
	//accessableServiceDetails, err = users.Store.GetAccessMapDetails(sessionMeta.params.UserID, sessionMeta.params.OrgID)
	//if err != nil {
	//	logrus.Debugf("get access map: %v", err)
	//	challengeUser("", "No services are assigned to you", nil, nil)
	//	return nil, fmt.Errorf("get access map: %v", err)
	//}

	service, err := chooseService(conn.User(), sessionMeta.params.UserID, sessionMeta.params.TrasaID, challengeUser)
	if err != nil {
		logrus.Error(err)
		challengeUser("", "Cannot access this service", nil, nil)
		return nil, err
	}
	sessionMeta.UpdateService(service)
	//sessionMeta.params.ServiceName = params.ServiceName
	//sessionMeta.params.ServiceID = params.ServiceID
	//sessionMeta.params.Hostname = params.Hostname
	//sessionMeta.params.Policy=params.Policy

	sessionMeta.params.Privilege = conn.User()
	policy, reason, err := SSHStore.checkPolicy(&models.ConnectionParams{
		ServiceID: sessionMeta.params.ServiceID,
		OrgID:     sessionMeta.params.OrgID,
		UserID:    sessionMeta.params.UserID,
		UserIP:    utils.GetIPFromAddr(conn.RemoteAddr()), //ip policy check
		Privilege: sessionMeta.params.Privilege,
		SessionID: sessionMeta.ID, //to append adhoc sessions

		//TODO change hard coded value
		Timezone: "Asia/Kathmandu", //for day time policy check
	})
	if err != nil {
		logrus.Debug(err)
		challengeUser("", string(reason), nil, nil)
		return nil, err
	}
	sessionMeta.policy = policy
	sessionMeta.log.SessionRecord = policy.RecordSession

	totp := ""

	ans, _ := challengeUser("user",
		"Second factor authentication",
		[]string{"\n\rEnter OTP(Blank for U2F): "},
		[]bool{true})

	if len(ans) != 1 {
		logrus.Debug("User canceled")
		return nil, ssh.ErrNoAuth
	}

	totp = ans[0]
	tfaMethod := "u2f"
	if totp != "" {
		tfaMethod = "totp"
	}

	//logrus.Debug(sessionMeta.params.Hostname)
	//logrus.Debug(sessionMeta.params.ServiceID)
	//logrus.Debug(sessionMeta.params.ServiceName)

	//logrus.Debug(utils.MarshallStruct(sessionMeta.params))

	if sessionMeta.policy.TfaRequired {
		orgDetail, err := orgs.Store.Get(sessionMeta.params.OrgID)
		if err != nil {
			logrus.Error(err)
		}

		deviceID, reason, ok := tfa.HandleTfaAndGetDeviceID(nil,
			tfaMethod,
			totp, sessionMeta.params.UserID,
			sessionMeta.log.UserIP,
			sessionMeta.params.ServiceName,
			orgDetail.Timezone,
			orgDetail.OrgName,
			sessionMeta.params.OrgID)

		if !ok {
			logrus.Trace("tfa failed ", reason)
			sessionMeta.log.FailedReason = reason
			sessionMeta.log.TfaDeviceID = deviceID
			sessionMeta.log.Status = false
			challengeUser("", string(reason), nil, nil)
			return nil, errors.New("tfa failed")
		}
	}

	logrus.Trace(sessionMeta.params.AccessDeviceID)
	reason, ok, err := accesscontrol.CheckDevicePolicy(policy.DevicePolicy, sessionMeta.params.AccessDeviceID, sessionMeta.log.TfaDeviceID, sessionMeta.params.OrgID)
	if err != nil {
		logrus.Error(err)
	}

	if !ok {
		return nil, errors.Errorf("device policy failed: %s", reason)
	}

	challengeUser("", "Checking host key", nil, nil)
	signer, hostkeyCallback, pass, err := handleUpstreamPasswordAndKey(conn.User(),
		sessionMeta.params.ServiceID, sessionMeta.params.Hostname, challengeUser)
	if err != nil {
		if err.Error() == "Weak Password" {
			logs.Store.LogLogin(sessionMeta.log, consts.REASON_PASSWORD_POLICY_FAILED, false)
		}
		challengeUser("", err.Error(), nil, nil)
		//return nil, err
	}

	updateSessionCredentials(sessionMeta, signer, hostkeyCallback, pass, totp)

	err = SSHStore.SetSession(conn.RemoteAddr(), sessionMeta)
	return nil, err

}

func updateSessionCredentials(sess *Session, signer ssh.Signer, hostkeyCallback ssh.HostKeyCallback, password, totp string) {

	sess.clientConfig.HostKeyCallback = hostkeyCallback

	if password != "" {
		sess.clientConfig.Auth = append(sess.clientConfig.Auth, ssh.Password(password))
	}

	if signer != nil {
		sess.clientConfig.Auth = append(sess.clientConfig.Auth, ssh.PublicKeys(signer))
	}

	sess.clientConfig.Auth = append(sess.clientConfig.Auth,
		ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
			answers := make([]string, len(questions))
			if len(questions) == 1 {

				if strings.Contains(questions[0], "Password") {
					answers[0] = password
					return answers, nil
				} else if strings.Contains(questions[0], "email") {
					answers[0] = sess.params.TrasaID
					return answers, nil
				} else if strings.Contains(questions[0], "totp") {
					answers[0] = totp
					return answers, nil
				} else {
					return answers, errors.New("invalid challenges")
				}

			}

			return answers, nil
		}),
	)

}

//searchAppName
// Returns index of app list.
// Returns -1 if dynamic ip is entered.
// Returns -2 if input is invalid.
//func searchAppName(input string, appUsers []models.AccessMapDetail) int {
//	input = strings.Trim(input, "")
//	for i, app := range appUsers {
//		if strings.EqualFold(input, app.ServiceName) {
//			return i
//		} else if input == app.Hostname {
//			return i
//		}
//	}
//
//	splittedInput := strings.Split(input, ":")
//	//If port is included
//	if (len(splittedInput) == 2) && splittedInput[0] != "" && splittedInput[1] != "" {
//		ip := net.ParseIP(splittedInput[0])
//		if ip != nil {
//			return -1
//		}
//	} else if len(splittedInput) == 1 {
//		ip := net.ParseIP(input)
//		if ip != nil {
//			return -1
//		}
//	}
//
//	index, err := strconv.Atoi(input)
//
//	//appNum, errStrConv = strconv.Atoi(ans[0])
//	//Check if choice is out of index
//	if err != nil || index > len(appUsers) || index < 1 {
//		return -2
//	}
//	return index - 1
//
//}

package sshproxy

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/seknox/ssh"
	"github.com/seknox/trasa/server/api/accesscontrol"
	"github.com/seknox/trasa/server/api/auth"
	"github.com/seknox/trasa/server/api/auth/tfa"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"strings"
)

//handleUpstreamPasswordAndKey tries to connect to upstream server to check host key and client key.
// It asks for password if authentication using key fails.
// Finally it returns public key authentication method(if successful) and password
func handleUpstreamPasswordAndKey(username, serviceID, hostname string, challengeUser ssh.KeyboardInteractiveChallenge) (ssh.AuthMethod, string, error) {

	var password string

	//get credentials from vault
	creds, err := services.GetUpstreamCreds(username, serviceID, "ssh", global.GetConfig().Trasa.OrgId)
	if err != nil {
		logrus.Debug(err)
		return nil, "", err
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

	//try connecting to ssh server to check host key and password.
	tempConf := ssh.ClientConfig{
		User:            username,
		HostKeyCallback: HandleHostKeyCallback(creds, serviceID, global.GetConfig().Trasa.OrgId, hostConfirmFunc),
	}

	publicKeyAuth := getPublicKeyAuthMethod(creds)
	if publicKeyAuth != nil {
		tempConf.Auth = append(tempConf.Auth, publicKeyAuth)
	}

	c, err := ssh.Dial("tcp", host, &tempConf)

	//logrus.Debug(err)

	if err == nil {
		c.Close()
	} else if strings.Contains(err.Error(), "trasa: could not update cert") {
		challengeUser("", "could not verify host", nil, nil)
		return nil, "", errors.New("Could not verify host")
	} else if strings.Contains(err.Error(), "unable to authenticate") {
		//password error

		var ans []string
		if creds.Password == "" {
			ans, err = challengeUser("user", "Upstream password",
				[]string{"\n\rEnter Password(Upstream Server): "},
				[]bool{false})
			if err != nil {
				return nil, "", err

			}
			if len(ans) != 1 {
				return nil, "", ssh.ErrNoAuth
			}

			password = ans[0]
		} else {
			password = creds.Password
		}

	} else {
		logrus.Trace(err)
	}

	//passwordStreanght := zxcvbn.PasswordStrength(password, nil)
	//
	//if creds.EnforceStrongPass && creds.ClientKey == "" {
	//	if passwordStreanght.Score < creds.ZxcvbnScore || len(password) < creds.MinimumChar {
	//		logrus.Debug("Weak Password")
	//		challengeUser("", "Weak password. Password policy not met", nil, nil)
	//		return signer, nil, "", errors.New("Weak Password")
	//	}
	//}

	return publicKeyAuth, password, nil

	//clientConfig.SetDefaults()
	//		clientConfig.Ciphers = append(clientConfig.Ciphers, "aes128-cbc", "blowfish-cbc", "3des-cbc")

}

//authenticateTRASA takes email/password from user, authenticates and returns user details
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

//keyboardInteractiveHandler is a callback function to handle keyboard-interactive authentication
func keyboardInteractiveHandler(conn ssh.ConnMetadata, challengeUser ssh.KeyboardInteractiveChallenge) (*ssh.Permissions, error) {

	var err error

	//Update map which holds info about ssh sessions
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

	//make user choose a service to connect to and get the service details
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

	orgDetail, err := orgs.Store.Get(sessionMeta.params.OrgID)
	if err != nil {
		logrus.Error(err)
	}

	sessionMeta.params.Privilege = conn.User()

	//Check trasa policies
	policy, reason, err := SSHStore.checkPolicy(&models.ConnectionParams{
		ServiceID: sessionMeta.params.ServiceID,
		OrgID:     sessionMeta.params.OrgID,
		UserID:    sessionMeta.params.UserID,
		UserIP:    utils.GetIPFromAddr(conn.RemoteAddr()), //ip policy check
		Privilege: sessionMeta.params.Privilege,
		SessionID: sessionMeta.ID, //to append adhoc sessions

		Timezone: orgDetail.Timezone, //for day time policy check
	})
	if err != nil {
		logrus.Debug(err)
		challengeUser("", string(reason), nil, nil)
		return nil, err
	}
	sessionMeta.policy = policy
	sessionMeta.log.SessionRecord = policy.RecordSession

	totp := ""

	//If TFA is required, ask user to verify the second factor
	if sessionMeta.policy.TfaRequired {

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

		deviceID, reason, ok := tfa.HandleTfaAndGetDeviceID(nil,
			tfaMethod,
			totp, sessionMeta.params.UserID,
			sessionMeta.log.UserIP,
			sessionMeta.params.ServiceName,
			orgDetail.Timezone,
			orgDetail.OrgName,
			sessionMeta.params.OrgID)

		if !ok {
			//logrus.Trace("tfa failed ", reason)
			sessionMeta.log.FailedReason = reason
			sessionMeta.log.TfaDeviceID = deviceID
			sessionMeta.log.Status = false
			challengeUser("", string(reason), nil, nil)
			return nil, errors.New("tfa failed")
		}
	}

	//Check device policy
	reason, ok, err := accesscontrol.CheckDevicePolicy(policy.DevicePolicy, sessionMeta.params.AccessDeviceID, sessionMeta.log.TfaDeviceID, sessionMeta.params.OrgID)
	if err != nil {
		logrus.Error(err)
	}

	if !ok {
		return nil, errors.Errorf("device policy failed: %s", reason)
	}

	//Try connecting to upstream server to verify password and host key
	challengeUser("", "Checking host key", nil, nil)
	pubKeyAuthMethod, pass, err := handleUpstreamPasswordAndKey(conn.User(), sessionMeta.params.ServiceID, sessionMeta.params.Hostname, challengeUser)
	if err != nil {
		if err.Error() == "Weak Password" {
			logs.Store.LogLogin(sessionMeta.log, consts.REASON_PASSWORD_POLICY_FAILED, false)
		}
		challengeUser("", err.Error(), nil, nil)
		return nil, err
	}

	//Since host is already verified, we can ignore it now
	sessionMeta.clientConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	//Add public key auth method
	if pubKeyAuthMethod != nil {
		sessionMeta.clientConfig.Auth = append(sessionMeta.clientConfig.Auth, pubKeyAuthMethod)
	}

	//Add keyboard-interactive auth method for handling upstream trasa PAM
	sessionMeta.clientConfig.Auth = append(sessionMeta.clientConfig.Auth, handleUpstreamTrasaPAM(sessionMeta.params.TrasaID, pass, totp))

	//Add password auth method
	if pass != "" {
		sessionMeta.clientConfig.Auth = append(sessionMeta.clientConfig.Auth, ssh.Password(pass))
	}

	err = SSHStore.SetSession(conn.RemoteAddr(), sessionMeta)
	return nil, err

}

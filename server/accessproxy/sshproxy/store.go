package sshproxy

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/seknox/ssh"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/providers/ca"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

func (s Store) GetUserFromPublicKey(publicKey ssh.PublicKey, orgID string) (*models.User, error) {
	var user models.User

	//If it's a certificate, extract public key only
	cert, ok := publicKey.(*ssh.Certificate)
	if ok {
		publicKey = cert.Key
	}

	publicKeyStr := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(publicKey)))

	//TODO use orgID??
	err := s.DB.QueryRow(`SELECT org_id, id,first_name, email, user_role, status FROM users WHERE public_key = $1 ;`, publicKeyStr).
		Scan(&user.OrgID, &user.ID, &user.FirstName, &user.Email, &user.UserRole, &user.Status)

	return &user, err

}

//tfaCert
//Is tfa already done from device agent
func (s Store) parseSSHCert(addr net.Addr, publicKey ssh.PublicKey) error {
	//TODO
	cert, ok := publicKey.(*ssh.Certificate)
	if !ok {
		return errors.New("invalid ssh certificate")
	}
	deviceID, ok := cert.Extensions["trasa-device-id"]
	if !ok {
		logrus.Error("device ID not found in ssh certificate")
		return errors.New("device ID not found in ssh certificate")
	}

	if s.sessions == nil {
		return errors.New("session map not initialised")
	}

	sess, ok := s.sessions[addr]
	if !ok {
		return errors.New("session not found")
	}

	sess.log.AccessDeviceID = deviceID
	sess.params.AccessDeviceID = deviceID

	groups, ok := cert.Extensions["trasa-user-groups"]

	logrus.Trace(groups, ok)
	if ok {
		sess.params.Groups = strings.Split(groups, ",")
	}

	return nil
}

//validateTempCert
func (s Store) validateTempCert(cert *ssh.Certificate, privilege string, orgID string) error {

	caKey, err := ca.Store.GetCertDetail(orgID, "system", consts.CERT_TYPE_SSH_CA)
	if err != nil {
		logrus.Error(err)
		return err
	}

	checker := ssh.CertChecker{
		IsUserAuthority: func(auth ssh.PublicKey) bool {
			return bytes.Compare(auth.Marshal(), caKey.Cert) == 0

		},
	}

	err = checker.CheckCert(privilege, cert)
	if err != nil {
		return errors.WithMessage(err, "could not verify certificate")
	}

	return nil
}

func (s Store) SetSession(addr net.Addr, session *Session) error {
	if s.sessions == nil {
		return errors.New("session map not initialised")
	}
	s.sessions[addr] = session
	return nil
}

func (s Store) GetSession(addr net.Addr) (*Session, error) {
	if s.sessions == nil {
		return nil, errors.New("session map not initialised")
	}
	sess, ok := s.sessions[addr]
	if !ok {
		return nil, errors.New("session not found")
	}
	return sess, nil
}

func (s Store) DeleteSession(addr net.Addr) error {
	if s.sessions == nil {
		return errors.New("session map not initialised")
	}

	//TODO check if session is actually deleted
	delete(s.sessions, addr)
	return nil
}

func (s Store) UpdateSessionMeta(addr net.Addr, connMeta ssh.ConnMetadata) error {
	if s.sessions == nil {
		return errors.New("session map not initialised")
	}

	sess, ok := s.sessions[addr]
	if !ok {
		return errors.New("session not found")
	}

	sess.ID = hex.EncodeToString(connMeta.SessionID())
	sess.log.EventID = sess.ID
	sess.clientConfig.User = connMeta.User()
	sess.params.Privilege = connMeta.User()
	sess.log.Privilege = connMeta.User()
	sess.log.ServiceType = "ssh"

	sess.log.LoginTime = time.Now().UnixNano()
	sess.log.UserAgent = string(connMeta.ClientVersion())

	sess.log.UserIP = utils.GetIPFromAddr(connMeta.RemoteAddr())

	s.sessions[addr] = sess

	return nil
}

//func (s Store) UpdateSessionParams(addr net.Addr, params *models.AccessMapDetail) error {
//	if s.sessions == nil {
//		return errors.New("session map not initialised")
//	}
//
//	if params == nil {
//		return errors.New("params is nil")
//	}
//
//	sess, ok := s.sessions[addr]
//	if !ok {
//		return errors.New("session not found")
//	}
//
//	sess.params = params
//
//	sess.log.OrgID = params.OrgID
//
//	sess.log.ServiceType = "ssh"
//	sess.log.Email = params.Email
//	sess.log.LoginTime = time.Now().UnixNano()
//	s.sessions[addr] = sess
//	return nil
//}

func (s Store) UpdateSessionUser(addr net.Addr, user *models.User) error {
	if s.sessions == nil {
		return errors.New("session map not initialised")
	}

	if user == nil {
		return errors.New("user is nil")
	}

	sess, ok := s.sessions[addr]
	if !ok {
		return errors.New("session not found")
	}

	if sess.params == nil {
		return errors.New("params is nil")
	}

	sess.params.UserID = user.ID
	sess.params.TrasaID = user.Email
	sess.params.OrgID = user.OrgID

	sess.log.OrgID = user.OrgID
	sess.log.UserID = user.ID
	sess.log.Email = user.Email
	s.sessions[addr] = sess
	return nil
}

func (s Store) SetAuthType(addr net.Addr, authType consts.SSH_AUTH_TYPE) error {
	if s.sessions == nil {
		return errors.New("session map not initialised")
	}

	sess, ok := s.sessions[addr]
	if !ok {
		return errors.New("session not found")
	}

	sess.AuthType = authType

	s.sessions[addr] = sess
	return nil
}

func (s Store) CreateGuestChannel(sessionID string) chan GuestClient {
	guestChan := make(chan GuestClient)
	s.guestChannels[sessionID] = guestChan
	return guestChan
}
func (s Store) GetGuestChannel(sessionID string) (chan GuestClient, error) {
	if s.guestChannels == nil {
		return nil, errors.New("channel is nil")
	}

	guestChan, ok := s.guestChannels[sessionID]
	if !ok {
		return nil, errors.New("channel not found")

	}
	return guestChan, nil
}

func (s Store) closeSession(addr net.Addr) {
	session, err := s.GetSession(addr)
	if err != nil {
		logrus.Error()
		return
	}

	delete(s.sessions, addr)
	delete(s.guestChannels, session.ID)

	session.tempSessionFile.Close()

	err = logs.Store.LogLogin(session.log, "", true)
	if err != nil {
		logrus.Error(err)
	}
	err = s.uploadSessionLog(session.log)
	if err != nil {
		logrus.Errorf("minio upload fail, trying again: %v", err)
		s.uploadSessionLog(session.log)
	}
	return
}

func (s Store) deleteGuestChannel(sessionID string) {
	delete(s.guestChannels, sessionID)
}

func (s Store) uploadSessionLog(authlog *logs.AuthLog) error {

	tempFileDir := filepath.Join(utils.GetTmpDir(), "trasa", "accessproxy", "ssh")
	bucketName := "trasa-ssh-logs"
	sessionID := authlog.SessionID

	loginTime := time.Unix(0, authlog.LoginTime).In(time.UTC)
	authlog.LogoutTime = time.Now().UnixNano()

	objectName := filepath.Join(authlog.OrgID, fmt.Sprintf("%d", loginTime.Year()), fmt.Sprintf("%d", int(loginTime.Month())), fmt.Sprintf("%d", loginTime.Day()), fmt.Sprintf("%s.session", sessionID))
	filePath := filepath.Join(tempFileDir, fmt.Sprintf("%s.session", sessionID))

	// Upload log file to minio
	uploadErr := logs.Store.PutIntoMinio(objectName, filePath, bucketName)
	if uploadErr != nil {
		logrus.Errorf("could not upload to minio, trying again: %v", uploadErr)
		uploadErr = logs.Store.PutIntoMinio(objectName, filePath, bucketName)
	}

	if uploadErr == nil {
		err := os.Remove(filePath)
		if err != nil {
			logrus.Errorf("could not delete session file: %v", err)
		}

	}

	return uploadErr
}

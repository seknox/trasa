package crypt

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"time"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func GetSSHCerts(serviceID, orgID string) (hostc, caHost, caUser string) {

	hostCert, err := Store.GetCertHolder(consts.CERT_TYPE_SSH_HOST_KEY, serviceID, orgID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logrus.Error(err)
		}
		hostc = ``
	}
	hostc = string(hostCert.Cert)

	caU, err := Store.GetCertHolder(consts.CERT_TYPE_SSH_CA, "user", orgID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logrus.Error(err)
		}
		caUser = ``
		return
	}

	caH, err := Store.GetCertHolder(consts.CERT_TYPE_SSH_CA, "host", orgID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logrus.Error(err)
		}
		caHost = ``
		return
	}

	return string(hostCert.Cert), string(caH.Cert), string(caU.Cert)
}

func GenerateTempSSHCert(username, orgID string) (string, string, error) {

	bitSize := 4096
	privateKey, err := utils.GeneratePrivateKey(bitSize)
	if err != nil {
		logrus.Errorf(`could not generate private key: %v`, err)
		return "", "", err
	}

	//caKeyStr, err := dbstore.Connect.GetCAkey(userContext.Org.ID)
	//TODO @sshah get ssh CA key and decrypt it
	//
	sshUserCA, err := Store.GetCertHolder(consts.CERT_TYPE_SSH_CA, "user", orgID)

	if err != nil {
		logrus.Debugf(`could not get CA key: %v`, err)
		return "", "", err
	}

	publicKeySSH, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		logrus.Errorf(`could not generate public key: %v`, err)
		return "", "", err
	}

	//publicKeyBytes := ssh.MarshalAuthorizedKey(publicKeySSH)

	caKey, err := ssh.ParsePrivateKey(sshUserCA.Key)
	if err != nil {
		logrus.Errorf(`Could not parse CA private key: %v`, err)
		return "", "", err
	}

	buf := make([]byte, 8)
	_, err = rand.Read(buf)
	if err != nil {
		logrus.Errorf("failed to read random bytes: %v", err)
		return "", "", err
	}
	serial := binary.LittleEndian.Uint64(buf)

	extentions := make(map[string]string)
	extentions = map[string]string{
		"permit-X11-forwarding":   "",
		"permit-agent-forwarding": "",
		"permit-port-forwarding":  "",
		"permit-pty":              "",
		"permit-user-rc":          "",
	}

	principals := []string{username}

	cert := ssh.Certificate{
		Key:             publicKeySSH,
		Serial:          serial,
		CertType:        ssh.UserCert,
		KeyId:           "some_id",
		ValidPrincipals: principals,
		ValidAfter:      uint64(time.Now().Unix()),
		ValidBefore:     uint64(time.Now().Add(time.Minute * 5).Unix()),
		Permissions: ssh.Permissions{
			Extensions: extentions,
		},
	}

	err = cert.SignCert(rand.Reader, caKey)
	if err != nil {
		logrus.Errorf(`could not sign public key: %v`, err)
		return "", "", err
	}

	privateKeyBytes := utils.EncodePrivateKeyToPEM(privateKey)
	certBytes := ssh.MarshalAuthorizedKey(&cert)
	if len(certBytes) == 0 {
		logrus.Errorf("failed to marshal signed certificate, empty result")
		return "", "", err
	}

	return string(privateKeyBytes), string(certBytes), nil
}

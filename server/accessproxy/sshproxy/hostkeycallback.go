package sshproxy

import (
	"bytes"
	"net"

	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/models"

	"github.com/pkg/errors"
	"github.com/seknox/ssh"
	logger "github.com/sirupsen/logrus"
)

var ErrVerifyHost = errors.New("trasa: could not update cert")

func HandleHostKeyCallback(creds *models.UpstreamCreds, serviceID, orgID string, confirmSkipVerify func(message string) bool) ssh.HostKeyCallback {

	caKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(creds.HostCaCert))
	if err != nil && creds.HostCaCert != "" {
		//If error is not nil, it means it could not get
		logger.Error(err)
		return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return err
		}
	}
	hostKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(creds.HostCert))
	if err != nil && creds.HostCert != "" {
		logger.Error(err)
		return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return err
		}
	}

	certChecker := ssh.CertChecker{
		SupportedCriticalOptions: nil,

		//If host presents a certificate
		IsHostAuthority: func(k ssh.PublicKey, address string) bool {

			//If caKey from trasacore is invalid
			if caKey == nil {
				if confirmSkipVerify("Upstream Host provided a certificate which could not be verified. Do you want to skip the verification and save the key?  \n") {
					//logger.Debug(string(k.Marshal()))
					//TODO verify cert
					// put orgID
					err := services.Store.UpdateHostCert(string(ssh.MarshalAuthorizedKey(k)), serviceID, orgID)
					if err != nil {
						logger.Error(err)
					}
					return err == nil
				} else {
					return false
				}
			}
			return bytes.Equal(k.Marshal(), caKey.Marshal())
		},

		//If host presents a key that is not a certificate
		HostKeyFallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			//if host certs is not set update it
			//logger.Debug(creds.HostCert)
			if creds.HostCert == "" {
				//logger.Debug(ssh.MarshalAuthorizedKey(key))
				if confirmSkipVerify("Upstream host key cannot be verified. Do you want to skip verification and save the key?\n") {
					//TODO put orgID
					return services.Store.UpdateHostCert(string(ssh.MarshalAuthorizedKey(key)), serviceID, orgID)
				} else {
					return ErrVerifyHost
				}

			}
			if bytes.Equal(key.Marshal(), hostKey.Marshal()) {
				return nil
			}
			return ErrVerifyHost
		},
	}

	return certChecker.CheckHostKey

}

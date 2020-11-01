package services

import (
	"encoding/json"

	"github.com/seknox/trasa/server/api/providers/ca"
	"github.com/seknox/trasa/server/api/providers/vault/tsxvault"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/sirupsen/logrus"
)

func GetUpstreamCreds(user, serviceID, serviceType, orgID string) (*models.UpstreamCreds, error) {

	//TODO wrap errors

	var resp models.UpstreamCreds

	switch serviceType {
	case "db":
		certHolder, err := ca.Store.GetCertHolder(consts.CERT_TYPE_TLS_KEYS, serviceID, orgID)
		if err == nil {
			resp.ClientCert = string(certHolder.Cert)
			//resp.ClientKey = string(certHolder.Key)
			resp.UserCaCert = string(certHolder.Csr)
		}

	case "ssh":
		hostCert, caHost, caUser := ca.GetSSHCerts(serviceID, orgID)
		resp.HostCert = hostCert
		resp.UserCaCert = caUser
		resp.HostCaCert = caHost

		userKey, userCrt, err := ca.GenerateTempSSHCert(user, orgID)
		if err != nil {
			logrus.Error(err)
			resp.ClientKey = ""
			resp.ClientCert = ""
		} else {
			resp.ClientKey = userKey
			resp.ClientCert = userCrt
		}

		//TODO handle rdp host cert
	}

	passPolicy, err := system.Store.GetGlobalSetting(orgID, consts.GLOBAL_PASSWORD_CONFIG)
	if err != nil {
		logrus.Error(err)
		//return nil, errors.WithMessage(err, "error fetching global settings")
	} else {
		var policy models.PasswordPolicy
		err = json.Unmarshal([]byte(passPolicy.SettingValue), &policy)
		if err != nil {
			logrus.Error(err)
		} else {
			resp.MinimumChar = policy.MinimumChars
			resp.ZxcvbnScore = policy.ZxcvbnScore
			resp.EnforceStrongPass = policy.EnforceStrongPass
		}

	}

	//First try to get key from vault
	//
	clientKey, err := tsxvault.Store.GetSecret(orgID, serviceID, "key", user)
	if err != nil {
		logrus.Debug(err)
	} else {
		resp.ClientKey = clientKey
		resp.ClientCert = ""
	}

	//userDetails := userContext.User
	//logrus.Debug(orgID, serviceID, "password", user)
	pass, err := tsxvault.Store.GetSecret(orgID, serviceID, "password", user)
	if err != nil {
		logrus.Debug(err)
		resp.Password = ""
	} else {
		resp.Password = pass
	}

	return &resp, nil
}

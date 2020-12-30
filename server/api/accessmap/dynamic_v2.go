package accessmap

import (
	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"time"
)

func GetDynamicPolicyV2(groups []string, userID, orgID string) (*models.Policy, error) {
	err := checkDynamicAccessSetting(orgID)
	if err != nil {
		return nil, err
	}
	dPolicy, err := Store.GetDynamicAccessPolicy(groups, userID, orgID)
	if err != nil {
		return nil, errors.Errorf(`failed to get dynamic access policy: %v`, err)
	}

	return dPolicy, nil
}

func CreateDynamicServiceV2(hostname, serviceType, userEmail, orgID string) (*models.Service, error) {
	err := checkDynamicAccessSetting(orgID)
	if err != nil {
		return nil, err
	}

	newService := NewService(hostname, serviceType, orgID)

	err = checkService(newService)
	if err != nil {
		return nil, errors.WithMessage(err, "could not connect to service")
	}

	err = services.Store.Create(&newService)
	if err != nil {
		return nil, errors.WithMessage(err, "create new dynamic service")
	}

	notifyAdmins(orgID, hostname, serviceType, userEmail)

	return &newService, err
}

func checkService(s models.Service) error {
	host := s.Hostname
	if !strings.Contains(s.Hostname, ":") {
		switch s.Type {
		case "ssh":
			host = host + ":22"
		case "rdp":
			host = host + ":3389"
		}
	}

	c, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}
	return c.Close()
}

func checkDynamicAccessSetting(orgID string) error {
	//Check global settings
	sett, err := system.Store.GetGlobalSetting(orgID, consts.GLOBAL_DYNAMIC_ACCESS)
	if err != nil {
		logrus.Error(err)
		return errors.Errorf("trasa: could not get dynamic access settings: %v", err)

	}
	if !sett.Status {
		return errors.New("trasa: dynamic access disabled")
	}
	return nil
}

//NewService returns empty struct of service
func NewService(hostname, stype, orgID string) models.Service {
	return models.Service{
		ID:              utils.GetUUID(),
		OrgID:           orgID,
		Name:            hostname,
		SecretKey:       utils.GetRandomString(17),
		Passthru:        false,
		Hostname:        hostname,
		Type:            stype,
		ManagedAccounts: "",
		RemoteAppName:   "",
		Adhoc:           false,
		NativeLog:       false,
		RdpProtocol:     "nla",
		ProxyConfig: models.ReverseProxy{
			RouteRule:           "",
			PassHostheader:      false,
			UpstreamServer:      "",
			StrictTLSValidation: true,
		},
		PublicKey:             "",
		ExternalProviderName:  "",
		ExternalID:            "",
		ExternalSecurityGroup: "{}",
		DistroName:            "",
		DistroVersion:         "",
		IPDetails: models.IPDetails{
			IpAddress:      "0.0.0.0",
			NetMask:        "",
			DefaultGateway: "0.0.0.0",
		},
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		DeletedAt: 0,
	}
}

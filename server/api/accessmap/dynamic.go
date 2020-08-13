package accessmap

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/api/groups"
	"github.com/seknox/trasa/server/api/notif"
	"github.com/seknox/trasa/server/api/policies"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	logger "github.com/sirupsen/logrus"
)

func CreateDynamicService(hostname, serviceType, userID, userEmail, privilege, orgID string) (*models.Service, error) {
	sett, err := checkDynamicSettings(userID, orgID)
	if err != nil {
		return nil, err
	}
	service, err := createDynamicServiceAndAccessMap(hostname, serviceType, userID, userEmail, privilege, orgID, sett)
	if err != nil {
		return nil, err
	}
	return service, err
}

func CreateDynamicAccessMap(serviceID, userID, userEmail, privilege, orgID string) (*models.Policy, error) {
	sett, err := checkDynamicSettings(userID, orgID)
	if err != nil {
		return nil, err
	}
	err = Store.CreateServiceUserMap(&models.ServiceUserMap{
		MapID:     utils.GetUUID(),
		ServiceID: serviceID,
		OrgID:     orgID,
		UserID:    userID,
		PolicyID:  sett.PolicyID,
		Privilege: privilege,
		AddedAt:   time.Now().Unix(),
	})

	if err != nil {
		return nil, errors.WithMessage(err, "create new dynamic access map")
	}

	policy, err := policies.Store.GetPolicy(sett.PolicyID, orgID)
	if err != nil {
		logger.Errorf("get dynamic access policy: %v", err)
		return nil, err
	}

	service, err := services.Store.GetFromID(serviceID)
	if err != nil {
		return &policy, err
	}

	notifyAdmins(orgID, fmt.Sprintf("%s (%s)", service.Name, service.Hostname), service.Type, userEmail)

	return &policy, nil
}

func checkDynamicSettings(userID, orgID string) (*models.GlobalDynamicAccessSettings, error) {

	//Check global settings
	sett, err := system.Store.GetGlobalSetting(orgID, consts.GLOBAL_DYNAMIC_ACCESS)
	if err != nil || !sett.Status {
		logger.Error(err)
		return nil, errors.Errorf("trasa: dynamic auth app not set: %v", err)
	}

	var settingVal models.GlobalDynamicAccessSettings
	err = json.Unmarshal([]byte(sett.SettingValue), &settingVal)
	if err != nil {
		logger.Error(err)
		return nil, errors.WithMessage(err, "trasa: invalid global settings")
	}

	check, err := groups.Store.CheckIfUserInGroup(userID, orgID, settingVal.UserGroups)
	if err != nil || !check {
		logger.Error(err)
		return &settingVal, errors.Errorf("trasa: user cannot access dynamic auth app: %v", err)
	}

	return &settingVal, nil
}

func createDynamicServiceAndAccessMap(hostname, serviceType, userID, userEmail, privilege, orgID string, sett *models.GlobalDynamicAccessSettings) (*models.Service, error) {

	newservice := models.NewEmptyServiceStruct()

	newservice.Name = hostname
	newservice.ID = "<Dynamic AuthApp>"
	newservice.Hostname = hostname
	newservice.OrgID = orgID
	newservice.Type = serviceType
	//newservice.SessionRecord = settings.VideoRecord
	newservice.RdpProtocol = "nla"
	newservice.CreatedAt = time.Now().Unix() //.In(nep).String() // .UTC().Format(time.RFC3339)
	newservice.UpdatedAt = time.Now().Unix() //.In(nep).String()

	//models.SendDynamicAccessEmail(userEmail, hostname, appType, primaryContact, newservice.CreatedAt, adminsEmails)

	newservice.ID = utils.GetUUID()

	newservice.Adhoc = false
	newservice.SecretKey = utils.GetRandomString(17)
	newservice.ManagedAccounts = ""
	newservice.CreatedAt = time.Now().Unix() //.In(nep).String() // .UTC().Format(time.RFC3339)
	newservice.UpdatedAt = time.Now().Unix() //.In(nep).String()
	err := services.Store.Create(&newservice)
	if err != nil {
		return nil, errors.WithMessage(err, "create new dynamic service")
	}

	err = Store.CreateServiceUserMap(&models.ServiceUserMap{
		MapID:     utils.GetUUID(),
		ServiceID: newservice.ID,
		OrgID:     orgID,
		UserID:    userID,
		PolicyID:  sett.PolicyID,
		Privilege: privilege,
		AddedAt:   time.Now().Unix(),
	})

	if err != nil {
		return nil, errors.WithMessage(err, "create new dynamic access map")
	}

	notifyAdmins(orgID, hostname, serviceType, userEmail)

	return &newservice, err
}

func notifyAdmins(orgID, hostname, appType, userEmail string) {

	admins, err := users.Store.GetAdminEmails(orgID)
	if err != nil {
		logger.Errorf("get admin email: %v", err)
	}

	if len(admins) < 1 {
		logger.Error("no admins found")
		return
	}

	var tmplt models.EmailDynamicAccess
	tmplt.User = userEmail
	tmplt.Hostname = hostname
	tmplt.AppType = appType
	tmplt.ReceiverEmail = admins[0]
	tmplt.CC = admins
	tmplt.TimeInt = time.Now().Unix()
	err = notif.Store.SendEmail(orgID, consts.EMAIL_DYNAMIC_ACCESS, tmplt)
	if err != nil {
		logger.Errorf("send dynamic access email", err)
	}

}

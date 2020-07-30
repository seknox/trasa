package system

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/core/crypt"
	logger "github.com/sirupsen/logrus"

	"github.com/seknox/trasa/consts"
	"github.com/seknox/trasa/global"
	"github.com/seknox/trasa/models"
	"github.com/seknox/trasa/utils"
	"github.com/sirupsen/logrus"
)

func InitStore(state *global.State) {
	Store = SystemStore{state}
	storeGlobalDynamicServiceSetting()
	storeGlobalPasswordPolicy()
	storeGlobalTRASASshAuth()
	storeGlobalEmailSettings()
	storeDefaultSecRules()
	storeDeviceHygieneCheck()
}

func InitStoreMock() *SystemMock {
	lmock := new(SystemMock)
	Store = lmock
	return lmock
}

var Store SystemAdapter

type SystemStore struct {
	*global.State
}

type SystemAdapter interface {
	// global settings
	SetGlobalSetting(setting models.GlobalSettings) error
	GetGlobalSetting(orgID, settingType string) (models.GlobalSettings, error)
	UpdateGlobalSetting(setting models.GlobalSettings) error

	// security rules
	CreateSecurityRule(rule models.SecurityRule) error
	GetSecurityRuleByName(orgID, constName string) (models.SecurityRule, error)
	getSecurityRules(orgID string) ([]models.SecurityRule, error)
	updateSecurityRule(orgID, statusstr, ruleID string) error

	// backups
	storeBackupMeta(backup models.Backup) error
	getBackupMeta(backup models.Backup) (models.Backup, error)
	getBackupMetas(orgID string) ([]models.Backup, error)
}

func storeGlobalPasswordPolicy() {
	var setting models.GlobalSettings
	var passPolicy models.PasswordPolicy
	passPolicy.Expiry = "never"
	passPolicy.MinimumChars = 8
	passPolicy.ZxcvbnScore = 2

	passJson, _ := json.Marshal(passPolicy)

	setting.OrgID = global.GetConfig().Trasa.OrgId
	setting.SettingID = utils.GetRandomID(7)
	setting.SettingType = consts.GLOBAL_PASSWORD_CONFIG
	setting.SettingValue = string(passJson)
	setting.UpdatedBy = "SYSTEM"
	setting.UpdatedOn = time.Now().Unix()
	err := Store.UpdateGlobalSetting(setting)
	if err != nil {
		logrus.Trace(err)
		return
	}

	logrus.Trace("Global Password Policy stored")
}
func storeDeviceHygieneCheck() {
	var setting models.GlobalSettings

	setting.OrgID = global.GetConfig().Trasa.OrgId
	setting.SettingID = utils.GetRandomID(7)
	setting.SettingType = consts.GLOBAL_DEVICE_HYGIENE_CHECK
	setting.SettingValue = "{}"
	setting.UpdatedBy = "SYSTEM"
	setting.UpdatedOn = time.Now().Unix()
	setting.Status = false
	err := Store.SetGlobalSetting(setting)
	if err != nil {
		logrus.Trace(err)
		return
	}

	logrus.Trace("Global Password Policy stored")
}

func storeGlobalEmailSettings() {
	var setting models.GlobalSettings
	setting.SettingValue = "{}"
	setting.Status = false
	setting.OrgID = global.GetConfig().Trasa.OrgId
	setting.SettingType = consts.GLOBAL_EMAIL_CONFIG
	setting.UpdatedBy = "SYSTEM"
	setting.SettingID = utils.GetRandomID(7)
	setting.SettingType = consts.GLOBAL_EMAIL_CONFIG
	setting.UpdatedOn = time.Now().Unix()
	err := Store.SetGlobalSetting(setting)
	if err != nil {
		//logger.Error(err)
		return
	}

	logrus.Trace("Global Email Setting stored")
}

func storeGlobalTRASASshAuth() {

	var gssha models.GlobalTrasaSshAuth
	gssha.MandatoryCertAuth = false

	v, _ := json.Marshal(gssha)
	var setting models.GlobalSettings
	setting.SettingValue = string(v)
	setting.Status = true
	setting.OrgID = global.GetConfig().Trasa.OrgId
	setting.UpdatedBy = "SYSTEM"
	setting.SettingID = utils.GetRandomID(7)
	setting.SettingType = consts.GLOBAL_TRASA_SSH_CERT_ENFORCE
	setting.UpdatedOn = time.Now().Unix()
	err := Store.SetGlobalSetting(setting)
	if err != nil {
		logrus.Trace(err)
		return
	}

	logrus.Trace("Global TRASA SSH auth")
}

func storeGlobalDynamicServiceSetting() {

	var sett models.GlobalDynamicAccessSettings
	sett.UserGroups = []string{}
	sett.PolicyID = ""

	v, _ := json.Marshal(sett)
	var setting models.GlobalSettings
	setting.SettingValue = string(v)
	setting.Status = true
	setting.OrgID = global.GetConfig().Trasa.OrgId
	setting.UpdatedBy = "SYSTEM"
	setting.SettingID = utils.GetRandomID(7)
	setting.SettingType = consts.GLOBAL_DYNAMIC_ACCESS
	setting.UpdatedOn = time.Now().Unix()
	err := Store.SetGlobalSetting(setting)
	if err != nil {
		logrus.Trace(err)
		return
	}

	logrus.Trace("Global dynamic access")
}

func initSystemCA() {
	_, err := crypt.Store.GetCertHolder(consts.CERT_TYPE_SSH_CA, "user", global.GetConfig().Trasa.OrgId)
	if !errors.Is(err, sql.ErrNoRows) {
		logrus.Debug("ssh CA already initialised")
		return
	}

	privateKey, err := utils.GeneratePrivateKey(4096)
	if err != nil {
		panic(err)
		return
	}
	pubKey, err := utils.GeneratePublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
		return
	}

	privateKeyBytes := utils.EncodePrivateKeyToPEM(privateKey)

	ca := models.CertHolder{
		CertID:      utils.GetUUID(),
		OrgID:       global.GetConfig().Trasa.OrgId,
		EntityID:    "system",
		Cert:        pubKey,
		Key:         privateKeyBytes,
		Csr:         nil,
		CertType:    consts.CERT_TYPE_SSH_CA,
		CreatedAt:   time.Now().Unix(),
		CertMeta:    "",
		LastUpdated: time.Now().Unix(),
	}
	err = crypt.Store.StoreCert(ca)
	if err != nil {
		logger.Error(err)
		return
	}
}

var secRules = `[
	{"name": "User Deleted", "description": "A user has been deleted from TRASA.", "constName":"DELETE_USER", "scope":"allusers", "source": "Admin activity", "action":"emailadmins"},
	{"name": "Admin User Deleted", "description": "A user with administrative privilege has been deleted from TRASA.", "constName":"DELETE_ADMIN_USER", "scope":"allusers", "source": "Admin activity", "action":"emailadmins"},
	{"name": "New User Added", "description": "A new user has been added to TRASA.", "constName":"CREATE_USER", "scope":"allusers", "source": "Admin activity", "action":"emailadmins"},
	{"name": "Admin User Added", "description": "A new administrator has been added to the TRASA.", "constName":"CREATE_ADMIN_USER", "scope":"orgAdmin", "source": "Admin activity", "action":"emailadmins"},
	{"name": "User Granted Admin Privilege", "description": "A user is granted an admin privilege.", "constName":"GRANT_ADMIN_PRIVILEGE", "scope":"orgAdmin", "source": "Admin activity", "action":"emailadmins"},
	{"name": "Admin User Profie Edited", "description": "A user profile with administrative privelege has been edited.", "constName":"ADMIN_PROFILE_EDITED", "scope":"orgAdmin", "source": "Admin activity", "action":"emailadmins"},
	{"name": "User Admin Privilege Revoked", "description": "An administrator is revoked of their admin privilege.", "constName":"REVOKE_ADMIN_PRIVILEGE", "scope":"orgAdmin", "source": "Admin activity", "action":"emailadmins"  },
	{"name": "User Admin Forgot Password", "description": "A user with administrative privilege has request forgot password process.", "constName":"ADMIN_FORGOT_PASSWORD", "scope":"orgAdmin", "source": "Admin activity", "action":"emailadmins"  },
	{"name": "Suspicious Login", "description": "TRASA detected a sign-in attempt that doesn't match a user's normal behavior, such as a sign-in from an unusual location..", "constName":"SUSPICIOUS_LOGIN", "scope":"users,orgAdmins,allusers", "condition": "ZERO_TRUST" ,"source": "user login", "action":"emailadmins" },
	{"name": "Multiple Failed Login", "description": "Multiple failed login attempt for user..", "constName":"SUSPICIOUS_LOGIN", "scope":"users,orgAdmins,allusers", "condition": "2" ,"source": "user login", "action":"emailadmins" },
	{"name": "Session Recording Disabled", "description": "Session recording for authapp has been disabled..", "constName":"SESSION_REC_DISABLED", "scope":"authapps" ,"source": "authapp setting", "action":"emailadmins" },
	{"name": "Low System Resource", "description": "Low system resource in TRASA server", "constName":"LOW_SYSTEM_RESOURCE", "scope":"SYSTEM", "source": "SYSTEM", "action":"emailadmins"  },
	{"name": "Unusual traffic received", "description": "TRASA server is getting unusual traffic", "constName":"UNUSUAL_TRAFFIC", "scope":"SYSTEM", "source": "request traffic", "action":"emailadmins"  }
]
	`

func storeDefaultSecRules() {
	var rule []models.SecurityRule

	err := json.Unmarshal([]byte(secRules), &rule)
	if err != nil {
		logrus.Debug("error storeDefaultSecRules: ", err)
	}

	for _, v := range rule {
		v.RuleID = utils.GetRandomID((10))
		v.OrgID = global.GetConfig().Trasa.OrgId
		v.CreatedBy = "SYSTEM"
		v.Action = `{"alertType":"email", "alertTo":"orgAdmins"}`
		v.CreatedAt = time.Now().Unix()
		v.LastModified = time.Now().Unix()

		err := Store.CreateSecurityRule(v)
		if err != nil {
			logrus.Trace(err)
			//fmt.Println("error storeDefaultSecRules: ", err)
		}
	}
}

package initdb

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/policies"
	"github.com/seknox/trasa/server/api/providers/vault"
	"github.com/seknox/trasa/server/api/providers/vault/tsxvault"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func InitDB() {
	//migrations

	//
	//init org
	initOrg()

	//init system settings
	storeGlobalDynamicServiceSetting()
	storeGlobalPasswordPolicy()
	storeGlobalEmailSettings()
	storeDefaultSecRules()
	storeDeviceHygieneCheck()

	//init CA

	initDefaultPolicies()

	setVaultConfigInGlobalVar()

}

// initOrg will create org and user if it doesn't exists
func initOrg() {
	orgID, err := orgs.Store.CheckOrgExists()
	if err != nil {
		panic(err)
	}

	if orgID != "" {
		global.SetOrgID(orgID)
		return
	}

	org := models.Org{
		ID:        utils.GetUUID(),
		OrgName:   "Trasa",
		Domain:    "trasa.io",
		Timezone:  "UTC",
		CreatedAt: time.Now().Unix(),
		License:   models.License{},
	}
	user := models.UserWithPass{
		User: models.User{
			ID:         utils.GetUUID(),
			OrgID:      org.ID,
			UserName:   "root",
			FirstName:  "Admin",
			MiddleName: "",
			LastName:   "Admin",
			Email:      "",
			UserRole:   "orgAdmin",
			Status:     true,
			IdpName:    "trasa",
			CreatedAt:  time.Now().Unix(),
			UpdatedAt:  time.Now().Unix(),
		},
		Password: "",
	}

	// Create Organization with db handle
	err = orgs.Store.CreateOrg(&org)
	if err != nil {
		logrus.Error(err)
		panic(err)
		return
	}
	global.SetOrgID(org.ID)

	hashedpass, err := bcrypt.GenerateFromPassword([]byte(consts.DEFAULT_ROOT_PASSWORD), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error(err)
		return
	}
	user.Password = string(hashedpass)

	err = users.Store.Create(&user)
	if err != nil {
		panic(err)
		return
	}
}

func storeGlobalPasswordPolicy() {
	var setting models.GlobalSettings
	var passPolicy models.PasswordPolicy
	passPolicy.Expiry = "never"
	passPolicy.MinimumChars = 8
	passPolicy.ZxcvbnScore = 2

	passJson, _ := json.Marshal(passPolicy)

	setting.OrgID = global.GetConfig().Trasa.OrgId
	setting.SettingID = utils.GetRandomString(7)
	setting.SettingType = consts.GLOBAL_PASSWORD_CONFIG
	setting.SettingValue = string(passJson)
	setting.UpdatedBy = "SYSTEM"
	setting.UpdatedOn = time.Now().Unix()

	_, err := system.Store.GetGlobalSetting(setting.OrgID, setting.SettingType)
	if errors.Is(err, sql.ErrNoRows) {
		errr := system.Store.UpdateGlobalSetting(setting)
		if errr != nil {
			logrus.Error(errr)
			return
		}
		logrus.Trace("global password policy initialised")

	}

}
func storeDeviceHygieneCheck() {
	var setting models.GlobalSettings

	setting.OrgID = global.GetConfig().Trasa.OrgId
	setting.SettingID = utils.GetRandomString(7)
	setting.SettingType = consts.GLOBAL_DEVICE_HYGIENE_CHECK
	setting.SettingValue = "{}"
	setting.UpdatedBy = "SYSTEM"
	setting.UpdatedOn = time.Now().Unix()
	setting.Status = false

	_, err := system.Store.GetGlobalSetting(setting.OrgID, setting.SettingType)
	if errors.Is(err, sql.ErrNoRows) {
		errr := system.Store.SetGlobalSetting(setting)
		if errr != nil {
			logrus.Error(errr)
			return
		}
		logrus.Trace("global device hygiene setting initialised")

	}

}

func initDefaultPolicies() {
	_, err := policies.Store.GetPolicy("f022d753-5f5f-4035-b3d4-59db0079d634", global.GetConfig().Trasa.OrgId)
	if err == nil {
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		logrus.Error(err)
		return
	}

	var fullPolicy = models.Policy{
		PolicyID:   "f022d753-5f5f-4035-b3d4-59db0079d634",
		OrgID:      global.GetConfig().Trasa.OrgId,
		PolicyName: "full-access",
		DayAndTime: []models.DayAndTimePolicy{{
			Days:     []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			FromTime: "00:00",
			ToTime:   "23:59",
		}},
		TfaRequired:      true,
		RecordSession:    true,
		FileTransfer:     true,
		IPSource:         "0.0.0.0/0",
		AllowedCountries: "",
		DevicePolicy: models.DevicePolicy{
			BlockUntrustedDevices:           false,
			BlockAutologinEnabled:           false,
			BlockTfaNotConfigured:           false,
			BlockJailBroken:                 false,
			BlockDebuggingEnabled:           false,
			BlockEmulated:                   false,
			BlockOpenWifiConn:               false,
			BlockIdleScreenLockDisabled:     false,
			BlockRemoteLoginEnabled:         false,
			BlockEncryptionNotSet:           false,
			BlockFirewallDisabled:           false,
			BlockCriticalAutoUpdateDisabled: false,
			BlockAntivirusDisabled:          false,
		},
		RiskThreshold: 5.0,
		CreatedAt:     0,
		UpdatedAt:     0,
		Expiry:        "2069-04-20",
		IsExpired:     false,
	}

	fullPolicy.OrgID = global.GetConfig().Trasa.OrgId

	err = policies.Store.CreatePolicy(fullPolicy)
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Trace("default policy initialised")
}

func storeGlobalEmailSettings() {
	var setting models.GlobalSettings
	setting.SettingValue = "{}"
	setting.Status = false
	setting.OrgID = global.GetConfig().Trasa.OrgId
	setting.SettingType = consts.GLOBAL_EMAIL_CONFIG
	setting.UpdatedBy = "SYSTEM"
	setting.SettingID = utils.GetRandomString(7)
	setting.SettingType = consts.GLOBAL_EMAIL_CONFIG
	setting.UpdatedOn = time.Now().Unix()

	_, err := system.Store.GetGlobalSetting(setting.OrgID, setting.SettingType)
	if errors.Is(err, sql.ErrNoRows) {
		errr := system.Store.SetGlobalSetting(setting)
		if errr != nil {
			logrus.Error(errr)
			return
		}
		logrus.Trace("global email setting initialised")

	}

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
	setting.SettingID = utils.GetRandomString(7)
	setting.SettingType = consts.GLOBAL_DYNAMIC_ACCESS
	setting.UpdatedOn = time.Now().Unix()

	_, err := system.Store.GetGlobalSetting(setting.OrgID, setting.SettingType)
	if errors.Is(err, sql.ErrNoRows) {
		errr := system.Store.SetGlobalSetting(setting)
		if errr != nil {
			logrus.Error(errr)
			return
		}
		logrus.Trace("global dynamic access setting initialised")
	}

}

var secRules = `[
	{"name": "User Deleted", "description": "A user has been deleted from TRASA.", "constName":"DELETE_USER", "scope":"allusers", "source": "Admin activity", "action":"emailadmins"},
	{"name": "Admin User Deleted", "description": "A user with administrative privilege has been deleted from TRASA.", "constName":"DELETE_ADMIN_USER", "scope":"allusers", "source": "Admin activity", "action":"emailadmins"},
	{"name": "New User Added", "description": "A new user has been added to TRASA.", "constName":"CREATE_USER", "scope":"allusers", "source": "Admin activity", "action":"emailadmins"},
	{"name": "Admin User Added", "description": "A new administrator has been added to the TRASA.", "constName":"CREATE_ADMIN_USER", "scope":"orgAdmin", "source": "Admin activity", "action":"emailadmins"},
	{"name": "User Granted Admin Privilege", "description": "A user is granted an admin privilege.", "constName":"GRANT_ADMIN_PRIVILEGE", "scope":"orgAdmin", "source": "Admin activity", "action":"emailadmins"},
	{"name": "Admin User Profie Edited", "description": "A user profile with administrative privilege has been edited.", "constName":"ADMIN_PROFILE_EDITED", "scope":"orgAdmin", "source": "Admin activity", "action":"emailadmins"},
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
		v.RuleID = utils.GetRandomString((10))
		v.OrgID = global.GetConfig().Trasa.OrgId
		v.CreatedBy = "SYSTEM"
		v.Action = `{"alertType":"email", "alertTo":"orgAdmins"}`
		v.CreatedAt = time.Now().Unix()
		v.LastModified = time.Now().Unix()

		_, err := system.Store.GetSecurityRuleByName(v.OrgID, v.ConstName)
		if errors.Is(err, sql.ErrNoRows) {
			errr := system.Store.CreateSecurityRule(v)
			if errr != nil {
				logrus.Error(errr)
				continue
			}
			logrus.Tracef("security rule %v initialised", v.ConstName)
		}

	}
}


// stores key and vault setting in global var
func setVaultConfigInGlobalVar() {
	if global.GetConfig().Vault.SaveMasterKey == true {
    logrus.Debug("SaveMasterKey is true, retrieving key from config file.")
	vaultsetting, err := system.Store.GetGlobalSetting(global.GetConfig().Trasa.OrgId, consts.GLOBAL_TSXVAULT)
	if err != nil {
		logrus.Debug("TsxVault is not initialized yet.")
		return
	}

	// store cred prov setting in global tsxvkey struct
	var cred models.CredProvProps
	err = json.Unmarshal([]byte(vaultsetting.SettingValue), &cred)
	if err != nil {
		logrus.Error(err)
		
		panic(err)
	}

	keyFromConfigFile := global.GetConfig().Vault.Key
	masterkey, err := hex.DecodeString(keyFromConfigFile)
	if err != nil {
		panic(err)
	}
	nkey := new([32]byte)
	copy(nkey[:], masterkey)

	// get access token from keyholder if credprov is hashicorp vault
	if cred.ProviderName == consts.CREDPROV_HCVAULT {
		ct, err := vault.Store.GetKeyOrTokenWithKeyval(global.GetConfig().Trasa.OrgId, string(consts.CREDPROV_HCVAULT_TOKEN))
		if err != nil {
			logrus.Error(err)
			panic(err)
		}
		
		if vaultsetting.Status == true {
			pt, err := utils.AESDecrypt(nkey[:], ct.KeyVal)
			if err != nil {
				logrus.Error(err)
				panic(err)
			}
			
			cred.ProviderAccessToken = string(pt)
		}
	
	}
		
	
	
	tsxvault.Store.SetTsxVaultKey(nkey, true, cred)
	}
}
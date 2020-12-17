package system

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/seknox/trasa/server/api/providers/vault/tsxvault"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

type GlobalSettingsResp struct {
	DynamicAccess  models.GlobalSettings `json:"dynamicAccess"`
	PasswordPolicy models.GlobalSettings `json:"passPolicy"`
	EmailSettings  models.GlobalSettings `json:"emailSettings"`
	DeviceHygiene  models.GlobalSettings `json:"deviceHygiene"`
}

//GlobalSettings returns all global settings
func GlobalSettings(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("request received")
	userContext := r.Context().Value("user").(models.UserContext)
	var resp GlobalSettingsResp

	passPolicy, err := Store.GetGlobalSetting(userContext.User.OrgID, consts.GLOBAL_PASSWORD_CONFIG)
	if err != nil {
		logrus.Error(err)

	}

	emailConfigs, err := Store.GetGlobalSetting(userContext.User.OrgID, consts.GLOBAL_EMAIL_CONFIG)
	if err != nil {
		logrus.Error(err)
	}

	resp.DeviceHygiene, err = Store.GetGlobalSetting(userContext.User.OrgID, consts.GLOBAL_DEVICE_HYGIENE_CHECK)
	if err != nil {
		logrus.Error(err)
	}
	resp.DynamicAccess, err = Store.GetGlobalSetting(userContext.Org.ID, consts.GLOBAL_DYNAMIC_ACCESS)
	if err != nil {
		logrus.Error(err)
	}

	resp.PasswordPolicy = passPolicy
	resp.EmailSettings = emailConfigs

	utils.TrasaResponse(w, 200, "success", "global settings fetched", "", resp)

}

type welcomeNoteResp struct {
	Intent string      `json:"intent"`
	Show   bool        `json:"show"`
	Data   interface{} `json:"data"`
}

// WelcomeNote processes any init events that needs to be presented to the admin after successfull dashboard login.
// For now, we use this handler to auto init vault and respond with sharded keys if it is not already initialized.
//  This handler can be used for any other similar operations.
func WelcomeNote(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	var vaultKeys []string = make([]string, 0)

	var resp []welcomeNoteResp = make([]welcomeNoteResp, 0)
	var note welcomeNoteResp
	note.Intent = consts.SHOW_VAULT_KEYS
	note.Data = 0
	note.Show = false

	if uc.User.UserName == "root" {
		// check and fire pending global settings (first time setup jobs like init vault)

		// check vault status. Auto init if its not already initialized
		vaultInitStatus, err := Store.GetGlobalSetting(uc.Org.ID, consts.GLOBAL_TSXVAULT)
		if err != nil || vaultInitStatus.Status == false {
			logrus.Trace("First time root login: Proceeding auto init TsxVault")
			// auto init
			vaultKeys, err = InitTsxvault(uc.User.OrgID, uc.User.ID)
			if err != nil {
				// what can we do "extra" here?
				logrus.Error(err)
			}
			note.Show = true
			note.Data = vaultKeys
		}

	}

	resp = append(resp, note)
	utils.TrasaResponse(w, 200, "success", "global settings triggered.", "WelcomeNote", resp)

}

//UpdateDeviceHygieneSetting updates device hygiene enforce settings
func UpdateDeviceHygieneSetting(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("device hygeiene req")
	uc := r.Context().Value("user").(models.UserContext)
	var req struct {
		EnableDeviceHygieneCheck bool `json:"enableDeviceHygieneCheck"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error parsing request", "device hygiene setting not updated", nil, nil)
		return
	}

	err := Store.UpdateGlobalSetting(models.GlobalSettings{
		SettingID:    utils.GetUUID(),
		OrgID:        uc.Org.ID,
		Status:       req.EnableDeviceHygieneCheck,
		SettingType:  consts.GLOBAL_DEVICE_HYGIENE_CHECK,
		SettingValue: "{}",
		UpdatedBy:    uc.User.ID,
		UpdatedOn:    time.Now().Unix(),
	})

	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not update  setting", "device hygiene setting not updated", nil, nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "Setting updated", "device hygiene setting updated", nil, nil)
}

//UpdateDynamicAccessSetting updates dynamic access settings
func UpdateDynamicAccessSetting(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	var req models.GlobalDynamicAccessSettings

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error parsing request", "dynamic access setting not updated", nil, nil)
		return
	}

	settVal, err := json.Marshal(req)
	if err != nil {
		logrus.Error(err)
	}

	err = Store.UpdateGlobalSetting(models.GlobalSettings{
		SettingID:    utils.GetUUID(),
		OrgID:        uc.Org.ID,
		Status:       req.Status,
		SettingType:  consts.GLOBAL_DYNAMIC_ACCESS,
		SettingValue: string(settVal),
		UpdatedBy:    uc.User.ID,
		UpdatedOn:    time.Now().Unix(),
	})

	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not update  setting", "dynamic access setting not updated", nil, nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "Setting updated", "dynamic access setting updated", nil, nil)
}

//UpdateEmailSetting updates email settings
func UpdateEmailSetting(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("UpdateEmailSetting request received")
	uc := r.Context().Value("user").(models.UserContext)
	var req models.EmailIntegrationConfig

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error parsing request", "email setting not updated", nil, nil)
		return
	}

	// switch req.IntegrationType {
	// case "EMAIL_SMTP":
	// 	req.IntegrationType = string(consts.EMAIL_SMTP)
	// case "mailgun":
	// 	req.IntegrationType = string(consts.EMAIL_MAILGUN)

	// }
	var key models.KeysHolder

	start := ""
	if len(req.AuthPass) > 4 {
		start = req.AuthPass[0:4]
	}

	key.OrgID = uc.User.OrgID
	key.KeyID = utils.GetRandomString(5)
	key.KeyTag = fmt.Sprintf("%sxxxx-xxxx...", start)
	key.AddedBy = uc.User.ID
	key.AddedAt = time.Now().Unix()
	key.KeyName = consts.GLOBAL_EMAIL_CONFIG_SECRET
	key.KeyVal = []byte(req.AuthPass)
	_, err := EncryptAndStoreKeyOrToken(key)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to update email settings", "email settings not updated")
		return
	}

	req.AuthPass = key.KeyTag
	j, err := json.Marshal(req)
	if err != nil {

	}

	var store models.GlobalSettings
	store.SettingValue = string(j)
	store.Status = true
	store.OrgID = uc.User.OrgID
	store.SettingType = consts.GLOBAL_EMAIL_CONFIG
	store.UpdatedBy = uc.User.ID
	store.UpdatedOn = time.Now().Unix()

	err = Store.UpdateGlobalSetting(store)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error updating global settings", "email settings not updated", nil, nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "successfully updated email setting", "password policy updated", nil, nil)

}

// StoreCloudProxyKey handles signed TRASA cloud proxy access key storage.
func StoreCloudProxyKey(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	type storekey struct {
		TsxCPxyAddr string `json:"tsxCPxyAddr"`
		APIKey      string `json:"apiKey"`
	}

	var req storekey

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error parsing request", "cloud proxy key", nil, nil)
		return
	}

	var key models.KeysHolder

	start := ""
	if len(req.APIKey) > 4 {
		start = req.APIKey[0:4]
	}

	key.OrgID = uc.User.OrgID
	key.KeyID = utils.GetRandomString(5)
	key.KeyTag = fmt.Sprintf("%sxxxx-xxxx...", start)
	key.AddedBy = uc.User.ID
	key.AddedAt = time.Now().Unix()
	key.KeyName = consts.GLOBAL_CLOUDPROXY_APIKEY
	key.KeyVal = []byte(req.APIKey)
	_, err := EncryptAndStoreKeyOrToken(key)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to encrypt key", "cloud proxy key not updated")
		return
	}

	// update config
	global.UpdateTRASACPxyAddr(req.TsxCPxyAddr)
	tsxvault.Store.SetTsxCPxyKey(req.APIKey)

	// proceed updating global setting
	req.APIKey = key.KeyTag
	j, err := json.Marshal(req)
	if err != nil {
		logrus.Error(err)
	}

	var store models.GlobalSettings
	store.SettingValue = string(j)
	store.Status = true
	store.OrgID = uc.User.OrgID
	store.SettingType = consts.GLOBAL_CLOUDPROXY_APIKEY
	store.UpdatedBy = uc.User.ID
	store.UpdatedOn = time.Now().Unix()

	err = Store.UpdateGlobalSetting(store)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to store key", "cloud proxy key not updated", nil, nil)
		return
	}

	logrus.Trace("CLOUDPROXYADDR: ", global.GetConfig().Trasa.CloudServer)

	utils.TrasaResponse(w, 200, "success", "key stored", "cloud proxy key updated", nil, nil)

}

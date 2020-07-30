package system

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/seknox/trasa/consts"
	"github.com/seknox/trasa/models"
	"github.com/seknox/trasa/utils"
	"github.com/sirupsen/logrus"
)

type globalSettings struct {
	DynamicAccess  models.GlobalSettings `json:"dynamicAccess"`
	SSHCertSetting models.GlobalSettings `json:"sshCertSetting"`
	PasswordPolicy models.GlobalSettings `json:"passPolicy"`
	EmailSettings  models.GlobalSettings `json:"emailSettings"`
	DeviceHygiene  models.GlobalSettings `json:"deviceHygiene"`
}

func GlobalSettings(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("request received")
	userContext := r.Context().Value("user").(models.UserContext)
	var resp globalSettings

	passPolicy, err := Store.GetGlobalSetting(userContext.User.OrgID, consts.GLOBAL_PASSWORD_CONFIG)
	if err != nil {
		logrus.Error(err)

	}

	emailConfigs, err := Store.GetGlobalSetting(userContext.User.OrgID, consts.GLOBAL_EMAIL_CONFIG)
	if err != nil {
		logrus.Error(err)
	}

	sshCertConf, err := Store.GetGlobalSetting(userContext.User.OrgID, consts.GLOBAL_TRASA_SSH_CERT_ENFORCE)
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
	resp.SSHCertSetting = sshCertConf

	utils.TrasaResponse(w, 200, "success", "global settings fetched", "GlobalSettings", resp)

}

type updatePassPolicyReq struct {
	Policy models.PasswordPolicy `json:"policy"`
	Enable bool                  `json:"enable"`
}

func UpdatePasswordPolicy(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("request received")
	userContext := r.Context().Value("user").(models.UserContext)
	var req updatePassPolicyReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error parsing request", "Password policy not updated", nil, nil)
		return
	}

	var store models.GlobalSettings
	encodePolicy, _ := json.Marshal(req.Policy)
	store.SettingValue = string(encodePolicy)
	store.Status = req.Enable
	store.OrgID = userContext.User.OrgID
	store.SettingType = consts.GLOBAL_PASSWORD_CONFIG
	store.UpdatedBy = userContext.User.ID
	store.UpdatedOn = time.Now().Unix()

	err := Store.UpdateGlobalSetting(store)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error updating global settings", "Password policy not updated", nil, nil)
		return
	}

	reason := "password policy enabled"
	if req.Enable == false {
		reason = "password policy disabled"
	}

	utils.TrasaResponse(w, 200, "success", reason, "Password policy updated", nil, nil)

}

func UpdateSSHCertSetting(w http.ResponseWriter, r *http.Request) {
	//mandatoryCert

	uc := r.Context().Value("user").(models.UserContext)
	var req struct {
		MandatoryCert bool `json:"mandatoryCert"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error parsing request", "Password policy not updated", nil, nil)
		return
	}

	err := Store.UpdateGlobalSetting(models.GlobalSettings{
		SettingID:    utils.GetUUID(),
		OrgID:        uc.Org.ID,
		Status:       req.MandatoryCert,
		SettingType:  consts.GLOBAL_TRASA_SSH_CERT_ENFORCE,
		SettingValue: "{}",
		UpdatedBy:    "",
		UpdatedOn:    time.Now().Unix(),
	})

	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not update  setting", "SSH certificate policy not updated", nil, nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "Setting updated", "SSH certificate policy updated", nil, nil)
}

func UpdateDeviceHygieneSetting(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("device hygeiene req")
	uc := r.Context().Value("user").(models.UserContext)
	var req struct {
		EnableDeviceHygieneCheck bool `json:"enableDeviceHygieneCheck"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error parsing request", "Password policy not updated", nil, nil)
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
		utils.TrasaResponse(w, 200, "failed", "Could not update  setting", "UpdateDeviceHygieneSetting", nil, nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "Setting updated", "UpdateDeviceHygieneSetting", nil, nil)
}

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

func UpdateEmailSetting(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("UpdateEmailSetting request received")
	uc := r.Context().Value("user").(models.UserContext)
	var req models.EmailIntegrationConfig

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error parsing request", "Password policy not updated", nil, nil)
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
	key.KeyID = utils.GetRandomID(5)
	key.KeyTag = fmt.Sprintf("%sxxxx-xxxx...", start)
	key.AddedBy = uc.User.ID
	key.AddedAt = time.Now().Unix()
	key.KeyName = consts.GLOBAL_EMAIL_CONFIG_SECRET
	key.KeyVal = []byte(req.AuthPass)
	_, err := EncryptAndStoreKeyOrToken(key)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to update email settings", "Email settings not updated")
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
		utils.TrasaResponse(w, 200, "failed", "error updating global settings", "Email settings not updated", nil, nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "successfully updated email setting", "Password policy updated", nil, nil)

}

package models

//Deprecated
type GlobalServiceSettings struct {
	Status              bool    `json:"status"`
	VideoRecord         bool    `json:"videoRecord"`
	AutocreateApp       bool    `json:"autocreateApp"`
	PolicyID            string  `json:"policyID"`
	PolicyName          string  `json:"policyName"`
	UserGroups          []Group `json:"userGroups"`
	UserGroupName       string  `json:"userGroupName"`
	DynamicServiceGroup string  `json:"dynamicServiceGroup"` //map_id of appgroup_usergroup_mapv1
}

// GlobalSettings holds model for global settings that can be applied to users in
// TRASA (not appusers. appusers can be managed from policy or compliance settings)
// multiple settings can be applied to global user settings.
// These settings should dictate user login behaviours, lock outs, password rotations etc...
// A minimun threshould settings should be auto generated and stored in database. Administrators should be
// able to modify these settings later on.
type GlobalSettings struct {
	SettingID string `json:"settingID"`
	OrgID     string `json:"orgID"`
	// Status is either active or disabled based on boolean value
	Status bool `json:"status"`
	// SettingType is name of setting
	SettingType string `json:"settingType"`
	// SettingValue holds json object of settings
	SettingValue string `json:"settingValue"`
	// UpdatedBy should be userID of user that updated this setting
	UpdatedBy string `json:"updatedBy"`
	UpdatedOn int64  `json:"updatedOn"`
}

type GlobalDynamicAccessSettings struct {
	Status     bool     `json:"status"`
	PolicyID   string   `json:"policyID"`
	Admins     []string `json:"admins"`
	UserGroups []string `json:"userGroups"`
}

type Config struct {
	Backup struct {
		Backupdir string `toml:"backupdir"`
	} `toml:"backup"`
	Database struct {
		Dbname     string `toml:"dbname"`
		Dbuser     string `toml:"dbuser"`
		Port       string `toml:"port"`
		Server     string `toml:"server"`
		Sslenabled bool   `toml:"sslenabled"`
		Usercert   string `toml:"usercert"`
		Userkey    string `toml:"userkey"`
		Cacert     string `toml:"cacert"`
	} `toml:"database"`

	//Deprecated
	Etcd struct {
		Server   string `toml:"server"`
		Usercert string `toml:"usercert"`
		Userkey  string `toml:"userkey"`
		Cacert   string `toml:"cacert"`
	} `toml:"etcd"`

	//Deprecated
	Logging struct {
		Env string `toml:"env"`
	} `toml:"logging"`
	Minio struct {
		Key    string `toml:"key"`
		Secret string `toml:"secret"`
		Server string `toml:"server"`
		Usessl bool   `toml:"usessl"`
	} `toml:"minio"`
	Platform struct {
		Base string `toml:"base"`
	} `toml:"platform"`
	Redis struct {
		Port       string   `toml:"port"`
		Server     []string `toml:"server"`
		Sslenabled bool     `toml:"sslenabled"`
		Usercert   string   `toml:"usercert"`
		Userkey    string   `toml:"userkey"`
		Cacert     string   `toml:"cacert"`
	} `toml:"redis"`
	Timezone struct {
		Location string `toml:"location"`
	} `toml:"timezone"`
	Security struct {
		InsecureSkipVerify bool `toml:"insecureSkipVerify"`
	} `toml:"security"`
	Trasa struct {
		Dashboard   string `toml:"dashboard"`
		Rootdomain  string `toml:"rootdomain"`
		CloudServer string `toml:"cloudServer"`
		Ssodomain   string `toml:"ssodomain"`

		//Deprecated
		Trasacore string `toml:"trasacore"`
		Rootdir   string `toml:"rootdir"`
		OrgId     string `toml:"orgID"`
	} `toml:"trasa"`
	Vault struct {
		Tsxvault bool   `toml:"tsxvault"`
		Port     string `toml:"port"`
		Server   string `toml:"server"`
		Token    string `toml:"token"`
	} `toml:"vault"`
	InternalHosts struct {
		Hosts string `toml:"hosts"`
	} `toml:"internalHosts"`
}

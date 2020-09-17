package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// UserDevice models devices registered to users.
type UserDevice struct {
	UserID     string `json:"userID"`
	OrgID      string `json:"orgID"`
	DeviceID   string `json:"deviceID"`
	MachineID  string `json:"machineID"`
	DeviceType string `json:"deviceType"`
	FcmToken   string `json:"fcmToken"`
	TotpSec    string `json:"-"`
	PublicKey  string `json:"publicKey"`
	//Deprecated
	DeviceFinger  string        `json:"deviceFinger"`
	Trusted       bool          `json:"trusted"`
	DeviceHygiene DeviceHygiene `json:"deviceHygiene"`
	AddedAt       int64         `json:"addedAt"`
}

//Deprecated
type DeviceFinger struct {
	Browser        string `json:"browser"`
	BrowserVersion string `json:"browserVersion"`
	Engine         string `json:"engine"`
	Os             string `json:"os"`
	Device         string `json:"device"`
	IP             string `json:"ip"`
	Blob           string `json:"blob"`
}

////////////////////////////////////////////////////
////////////////////////////////////////////////////
/////// 		Device Hygiene
////////////////////////////////////////////////////
////////////////////////////////////////////////////

// WorkstationHygiene stores health of user workstation.
type DeviceHygiene struct {
	//DeviceID         string           `json:"deviceID"`   // TRASA unique identifier for this device.
	//DeviceType       string           `json:"deviceType"` // can be mobile or workstation
	DeviceInfo    DeviceInfo    `json:"deviceInfo"`
	DeviceOS      DeviceOS      `json:"deviceOS"`
	LoginSecurity LoginSecurity `json:"loginSecurity"`
	//DeviceBrowser    DeviceBrowser    `json:"deviceBrowser"`
	NetworkInfo      NetworkInfo      `json:"networkInfo"`
	EndpointSecurity EndpointSecurity `json:"endpointSecurity"`
	LastCheckedTime  int64            `json:"lastCheckedTime"`
}

func (a DeviceHygiene) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *DeviceHygiene) Scan(value interface{}) error {
	if value == nil {
		*a = DeviceHygiene{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}
	return nil
}

//DeviceInfo
type DeviceInfo struct {
	DeviceName    string `json:"deviceName"`
	DeviceVersion string `json:"deviceVersion"`
	MachineID     string `json:"machineID"`
	Brand         string `json:"brand"`        // iOS: "Apple" // Android: "xiaomi"
	Manufacturer  string `json:"manufacturer"` // iOS: "Apple"  // Android: "Google"
	DeviceModel   string `json:"deviceModel"`  // iOS: "iPhone7,2"  //
}

//DeviceOS
type DeviceOS struct {
	OSName              string   `json:"osName"`    //(OS Name) iOS: "iOS" on newer iOS devices "iPhone OS" on older devices, including older iPad's. // Android: "Android"
	OSVersion           string   `json:"osVersion"` //(OS version) iOS: "11.0" // Android: "7.1.1"
	KernelType          string   `json:"kernelType"`
	KernelVersion       string   `json:"kernelVersion"`
	ReadableVersion     string   `json:"readableVersion"`
	LatestSecurityPatch string   `json:"latestSecurityPatch"` //// "2018-07-05"
	AutoUpdate          bool     `json:"autoUpdate"`
	PendingUpdates      []string `json:"pendingUpdates"`
	JailBroken          bool     `json:"jailBroken"`
	DebugModeEnabled    bool     `json:"debugModeEnabled"` // only for mobile device
	IsEmulator          bool     `json:"isEmulator"`       // only for mobile device
}

//LoginSecurity is device hygiene related to login
type LoginSecurity struct {
	// checks if device requires login before console.
	AutologinEnabled bool `json:"autologinEnabled" `
	// value can be password/pin/pattern/faceID/fingerprint
	LoginMethod         string `json:"loginMethod"`
	PasswordLastUpdated string `json:"passwordLastUpdated"`
	TfaConfigured       bool   `json:"tfaConfigured"`
	// IdleDeviceScreenLockTime stores if device screen lock time. value can be "never","5 minute idle",
	IdleDeviceScreenLockTime string `json:"idleDeviceScreenLockTime"`
	IdleDeviceScreenLock     bool   `json:"idleDeviceScreenLock"`
	RemoteLoginEnabled       bool   `json:"remoteLoginEnabled"`
}

// AntiVirus collects data of installed antivirus or antimalware  or any endpoint protection agent available in user device. eg window defender, crowdstrike, kaspersky.
type EndpointSecurity struct {
	EpsConfigured           bool   `json:"epsConfigured"` // endpoint security enabled
	EpsVendorName           string `json:"epsVendorName"` // endpoint security vendor name. eg. win defender, avira
	EpsVersion              string `json:"epsVersion"`    // version of endpoint security sute
	EpsMeta                 string `json:"epsMeta"`
	FirewallEnabled         bool   `json:"firewallEnabled"`
	FirewallPolicy          string `json:"firewallPolicy"` // FirewallPolicy stores firewall config policy found on user device
	DeviceEncryptionEnabled bool   `json:"deviceEncryptionEnabled"`
	DeviceEncryptionMeta    string `json:"deviceEncryptionMeta"` // DeviceEncryptionMeta stores metadata related to disk encryption (only if enabled)
}

// NetworkInfo collects network information of currently active connection during time of access
type NetworkInfo struct {
	Hostname         string `json:"hostname"`
	DomainControlled bool   `json:"domainControl"`
	DomainName       string `json:"domainName"`
	InterfaceName    string `json:"interfaceName"` // Name of outgoing net interface. eg. eth0, wlaan11
	IPAddress        string `json:"ipAddress"`
	MacAddress       string `json:"macAddress"`
	WirelessNetwork  bool   `json:"wirelessNetwork"`
	OpenWifiConn     bool   `json:"openWifiConn"`
	NetworkName      string `json:"networkName"`     // Name of connected network. eg. OfficeWIFI, marketingLAN
	NetworkSecurity  string `json:"networkSecurity"` // detail about current active connection. Eg. if using wifi, then is it open wifi? or wpa2psk wifi?
}

type DevicePolicy struct {
	BlockUntrustedDevices bool `json:"blockUntrustedDevices"`

	//May not/ does not work
	BlockAutologinEnabled bool `json:"blockAutologinEnabled"`
	BlockTfaNotConfigured bool `json:"blockTfaNotConfigured"`
	BlockJailBroken       bool `json:"blockJailBroken"`
	BlockDebuggingEnabled bool `json:"blockDebuggingEnabled"`
	BlockEmulated         bool `json:"blockEmulated"`
	BlockOpenWifiConn     bool `json:"blockOpenWifiConn"`

	//Works
	BlockIdleScreenLockDisabled bool `json:"blockIdleScreenLockDisabled"`
	BlockRemoteLoginEnabled     bool `json:"blockRemoteLoginEnabled"`
	BlockEncryptionNotSet       bool `json:"blockEncryptionNotSet"`
	BlockFirewallDisabled       bool `json:"blockFirewallDisabled"`
	//BlockPendingUpdates             bool `json:"blockPendingUpdates"`
	BlockCriticalAutoUpdateDisabled bool `json:"blockCriticalAutoUpdateDisabled"`
	BlockAntivirusDisabled          bool `json:"blockAntivirusDisabled"`
}

type DevicePolicyMaker struct {
	RuleID         string `json:"ruleID"`
	OrgID          string `json:"orgID"`
	Name           string `json:"name"`
	ConstName      string `json:"constName"`
	Description    string `json:"description"`
	Scope          string `json:"scope"`          // eg. ALL_DEVICE, MOBILE, WORKSTATION
	Constraint     string `json:"constraint"`     // Constraint can be version, name or certain value
	ConstraintType string `json:"constraintType"` // Type can be EQ, LT, GT (equal to, less than, greater than, boolean)
	ConsraintValue string `json:"contraintValue"` // eg. 10, windows xp,
	Status         bool   `json:"status"`         // enabled or disabled
	Source         string `json:"source"`         // source of event
	Action         string `json:"action"`         // action to take. eg. BLOCK, ALERT
	CreatedBy      string `json:"createdBy"`
	CreatedAt      int64  `json:"createdAt"`
	LastModified   int64  `json:"lastModified"`
}

// DeviceBrowser hygiene should be always updated along with workstation.
type DeviceBrowser struct {
	ID    string `json:"ID"`
	OrgID string `json:"orgID"`
	// DeviceID should be deviceID of workstation which this browser in context belongs to.
	DeviceID   string              `json:"deviceID"`
	Name       string              `json:"name"`
	Version    string              `json:"version"`
	Build      string              `json:"build"`
	IsBot      bool                `json:"isBot"`
	UserAgent  string              `json:"userAgent"`
	Extensions []BrowserExtensions `json:"extensions"`
}

type BrowserExtensions struct {
	// device id is id of device which maps to deviceID of userdevices
	DeviceID string `json:"deviceID"`
	// userID maps to userID from users
	UserID string `json:"userID"`
	OrgID  string `json:"orgID"`
	// ExtensionID is unique identifier of extension that is provided by extensions to browser vendors.
	ExtensionID     string   `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Version         string   `json:"version"`
	MayDisable      bool     `json:"mayDisable"`
	Enabled         bool     `json:"enabled"`
	InstallType     string   `json:"installType"`
	Type            string   `json:"type"`
	Permissions     []string `json:"permissions"`
	HostPermissions []string `json:"hostPermissions"`
	IsVulnerable    bool     `json:"isVulnerable"`
	VulnReason      string   `json:"vulnReason"`
	// LastChecked stores date of when this extension was last uddated in trasa database
	LastChecked int64 `json:"lastChecked"`
}

// func (a BrowserExtensions) Value() (driver.Value, error) {
// 	return json.Marshal(a)
// }

// // Make the Attrs struct implement the sql.Scanner interface. This method
// // simply decodes a JSON-encoded value into the struct fields.
// func (a *BrowserExtensions) Scan(value interface{}) error {
// 	if value == nil {
// 		*a = BrowserExtensions{}
// 		return nil
// 	}
// 	b, ok := value.([]byte)
// 	if !ok {
// 		return errors.New("type assertion to []byte failed")
// 	}

// 	err := json.Unmarshal(b, &a)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

type DeviceAuthFinger struct {
	MachineID      string `json:"machineID"`
	Hostname       string `json:"hostname"`
	KernelType     string `json:"kernelType"`
	OsName         string `json:"osName"`
	OsVersion      string `json:"osVersion"`
	DeviceName     string `json:"deviceName"`
	SecurityStatus struct {
		IsPasswordSet            bool     `json:"isPasswordSet"`
		PasswordLastUpdated      string   `json:"passwordLastUpdated"`
		IsFirewallSet            bool     `json:"isFirewallSet"`
		IsDeviceEncryptionSet    bool     `json:"isDeviceEncryptionSet"`
		IsRemoteLoginEnabled     bool     `json:"isRemoteLoginEnabled"`
		IsScreenLockEnabled      bool     `json:"isScreenLockEnabled"`
		CriticalAutoUpdateStatus bool     `json:"criticalAutoUpdateStatus"`
		PendingUpdates           []string `json:"pendingUpdates"`
		//InstalledApplications []string
	} `json:"securityStatus"`
}

// Mobile
type MobileDeviceHygiene1 struct {
	InstalledApps string `json:"installedApps"`
	DeviceName    string `json:"name"`         // iOS: "Becca's iPhone 6" // Android: ?
	Brand         string `json:"brand"`        // iOS: "Apple" // Android: "xiaomi"
	Manufacturer  string `json:"manufacturer"` // iOS: "Apple"  // Android: "Google"
	OSName        string `json:"osName"`       //(OS Name) iOS: "iOS" on newer iOS devices "iPhone OS" on older devices, including older iPad's. // Android: "Android"
	OSVersion     string `json:"osVersion"`    //(OS version) iOS: "11.0" // Android: "7.1.1"

	DeviceModel      string `json:"deviceModel"` // iOS: "iPhone7,2"  // Android: "goldfish"
	UserAgent        string `json:"userAgent"`   // iOS: "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143" // Android: ?
	IsJailBroken     bool   `json:"isJailBroken"`
	HooksDetected    bool   `json:"hooksDetected"`
	DebugModeEnabled bool   `json:"debugModeEnabled"`
	DeviceID         string `json:"deviceID"` // iOS: "FCDBD8EF-62FC-4ECB-B2F5-92C9E79AC7F9" // Android: "dd96dec43fb81c97"
	IpAddress        string `json:"ipAddress"`
	MacAddress       string `json:"macAddress"`
	ReadableVersion  string `json:"readableVersion"` //(application version+build number) iOS: 1.0.1.32  // Android: 1.0.1.234
	SecurityPatch    string `json:"securityPatch"`   //// "2018-07-05"
	AppVersion       string `json:"appVersion"`      //Gets the application version.
	IsEmulator       bool   `json:"isEmulator"`
	DeviceLockSet    bool   `json:"deviceLockSet"`
	DeviceLockType   string `json:"deviceLockType"` // pin/pattern/faceID/fingerprint
}

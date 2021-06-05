package hygiene

/**
 * Copyright (C) 2020 Seknox Pte Ltd.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// DeviceHygiene
type DeviceHygiene struct {
	DeviceInfo       DeviceInfo       `json:"deviceInfo"`
	DeviceOS         DeviceOS         `json:"deviceOS"`
	LoginSecurity    LoginSecurity    `json:"loginSecurity"`
	NetworkInfo      NetworkInfo      `json:"networkInfo"`
	EndpointSecurity EndpointSecurity `json:"endpointSecurity"`
	LastCheckedTime  int64            `json:"lastCheckedTime"`
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

// DeviceBrowser hygiene should be always updated along with workstation.
type DeviceBrowser struct {
	ID         string              `json:"ID"`
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
	LastChecked     int64    `json:"lastChecked"`
}

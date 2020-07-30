package models

import (
	"github.com/seknox/trasa/consts"
	"github.com/tstranex/u2f"
)

type UpstreamCreds struct {
	Password          string `json:"password"`
	HostCert          string `json:"hostCert"`
	HostCaCert        string `json:"hostCaCert"`
	UserCaCert        string `json:"UserCaCert"`
	ClientCert        string `json:"clientCert"`
	ClientKey         string `json:"clientKey"`
	SkipHostVerify    bool   `json:"skipHostVerify"`
	MinimumChar       int    `json:"minimumChar"`
	ZxcvbnScore       int    `json:"zxcvbnScore"`
	EnforceStrongPass bool   `json:"enforceStrongPass"`
}

//TODO remove useless fields
// also omit unnecessary fields in json
type ConnectionParams struct {
	ServiceID       string           `json:"serviceID"`
	ServiceName     string           `json:"-"`
	AppSecret       string           `json:"appSecret"`
	TfaMethod       string           `json:"tfaMethod"`
	TotpCode        string           `json:"totpCode"`
	TrasaID         string           `json:"trasaID"`
	OrgID           string           `json:"orgID"`
	Privilege       string           `json:"privilege"`
	Password        string           `json:"password"`
	UserID          string           `json:"userID"`
	SessionID       string           `json:"sessionID"`
	UserIP          string           `json:"userIP"`
	Skip2FA         bool             `json:"skip2FA"`
	SignResponse    u2f.SignResponse `json:"signResponse"`
	CSRF            string           `json:"csrf"`
	SESSION         string           `json:"session"`
	OptHeight       int64            `json:"optHeight"`
	OptWidth        int64            `json:"optWidth"`
	IsSharedSession bool             `json:"isSharedSession"`
	ConnID          string           `json:"connID"`
	Token           string           `json:"token"`
	ServiceType     string           `json:"serviceType"`
	RdpProtocol     string           `json:"rdpProto"`
	SessionRecord   bool             `json:"-"`
	CanTransferFile bool             `json:"-"`
	DeviceHygiene   DeviceHygiene    `json:"deviceHygiene"`
	AccessDeviceID  string           `json:"-"`
	TfaDeviceID     string           `json:"-"`
	BrowserID       string           `json:"-"`
	Hostname        string           `json:"hostname"`
	Timezone        string           `json:"-"`
	OrgName         string           `json:"-"`
	//UserAgent       string
}

type CheckPolicyFunc func(params *ConnectionParams, policy *Policy, adhoc bool) (bool, consts.FailedReason)

package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Here starts structures related to Services. Services are basic component of trasa.
// Everything is connected as piece of Oauth Services and trasa itself is 1st Service.
// Service structure. this holds typical Service structure also known as Oauth clients
type Service struct {
	ID              string       `json:"ID"`
	OrgID           string       `json:"orgID"`
	Name            string       `json:"serviceName" validate:"printascii"`
	SecretKey       string       `json:"secretKey"`
	Passthru        bool         `json:"passthru"`
	Hostname        string       `json:"hostname" validate:"printascii,required"`
	Type            string       `json:"serviceType" validate:"printascii,required"`
	ManagedAccounts string       `json:"managedAccounts"`
	RemoteAppName   string       `json:"remoteAppName"`
	Adhoc           bool         `json:"adhoc"`
	NativeLog       bool         `json:"nativeLog"`
	RdpProtocol     string       `json:"rdpProtocol"`
	ProxyConfig     ReverseProxy `json:"proxyConfig"`
	PublicKey       string       `json:"publicKey"`
	// ExternalProviderName is name of provider from which this Services details was fetched(eg, digital ocean, aws)
	ExternalProviderName string `json:"externalProviderName"`
	// ExternalID is ID of service that exists outside of trasa. (eg, digital ocean, aws)
	ExternalID            string `json:"externalID"`
	ExternalSecurityGroup string `json:"externalSecurityGroup"`
	// DistroName can be any specefic distribution version. eg ubuntu, debian, windows.
	DistroName    string    `json:"distroName"`
	DistroVersion string    `json:"distroVersion"`
	IPDetails     IPDetails `json:"ipDetails"`
	CreatedAt     int64
	UpdatedAt     int64
	DeletedAt     int64
}

func NewEmptyServiceStruct() Service {
	return Service{
		ID:              "",
		OrgID:           "",
		Name:            "",
		SecretKey:       "",
		Passthru:        false,
		Hostname:        "",
		Type:            "",
		ManagedAccounts: "",
		RemoteAppName:   "",
		//SessionRecord:   false,
		Adhoc:       false,
		NativeLog:   false,
		RdpProtocol: "",
		ProxyConfig: ReverseProxy{
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
		IPDetails: IPDetails{
			IpAddress:      "0.0.0.0",
			NetMask:        "",
			DefaultGateway: "0.0.0.0",
		},
		CreatedAt: 0,
		UpdatedAt: 0,
		DeletedAt: 0,
	}
}

// ReverseProxy defines proxy config for http access proxy
type ReverseProxy struct {
	RouteRule           string `json:"routeRule"`
	PassHostheader      bool   `json:"passHostHeader"`
	UpstreamServer      string `json:"upstreamServer"`
	StrictTLSValidation bool   `json:"strictTLSValidation"`
}

func (r ReverseProxy) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *ReverseProxy) Scan(value interface{}) error {

	if value == nil {
		*r = ReverseProxy{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	return nil
}

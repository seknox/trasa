package models

import (
	"github.com/seknox/trasa/server/consts"
)

//AuthLog is a log struct for all authentication events in trasa including dashboard login
type AuthLog struct {
	EventID        string                `json:"eventID"`
	Endpoint       consts.ConstEndpoints `json:"endpoint"`
	SessionID      string                `json:"sessionID"`
	OrgID          string                `json:"orgID"`
	ServiceName    string                `json:"ServiceName"`
	ServiceID      string                `json:"ServiceID"`
	ServiceType    string                `json:"ServiceType"`
	ServerIP       string                `json:"serverIP"`
	ServerName     string                `json:"serverName"`
	UserName       string                `json:"userName"`
	Email          string                `json:"email"`
	UserID         string                `json:"userID"`
	UserAgent      string                `json:"userAgent"`
	AccessDeviceID string                `json:"accessDeviceID"`
	TfaDeviceID    string                `json:"tfaDeviceID"`
	DeviceType     string                `json:"deviceType"`
	Commands       []string              `json:"commands"`
	UserIP         string                `json:"userIP"`
	GeoLocation    struct {
		IsoCountryCode string    `json:"isoCountryCode"`
		City           string    `json:"city"`
		TimeZone       string    `json:"timeZone"`
		Location       []float64 `json:"location"`
	} `json:"geoLocation"`
	LoginMethod     string              `json:"loginMethod"`
	Status          bool                `json:"status"`
	MarkedAs        string              `json:"markedAs"`
	LoginTime       int64               `json:"loginTime"`
	LogoutTime      int64               `json:"logoutTime"`
	SessionDuration string              `json:"sessionDuration"`
	SessionRecord   bool                `json:"sessionRecord"`
	FailedReason    consts.FailedReason `json:"failedReason"`
	Guests          []string            `json:"guests"`
}

//InAppTrail is struct of inapp audit log of trasa
type InAppTrail struct {
	EventID      string      `json:"eventID"`
	Status       bool        `json:"status"`
	OrgID        string      `json:"orgID"`
	UserID       string      `json:"userID"`
	Email        string      `json:"email"`
	Description  string      `json:"description"`
	UserAgent    string      `json:"userAgent"`
	RequestDump  interface{} `json:"requestDump"`
	ResponseDump interface{} `json:"responseDump"`
	EventType    string      `json:"eventType"`
	EventTime    int64       `json:"eventTime"`
	ClientIP     string      `json:"clientIP"`
}

//SignupLog
type SignupLog struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Company     string `json:"company"`
	JobTitle    string `json:"jobTitle"`
	Country     string `json:"country"`
	PhoneNumber string `json:"phoneNumber"`
	Reference   string `json:"reference"`
	SignupTime  int64  `json:"signupTime"`
}

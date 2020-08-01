package logs

import (
	"encoding/json"
	"net/http"
	"time"

	"net"

	"github.com/seknox/trasa/server/api/misc"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

//TODO Move this to models
type AuthLog struct {
	EventID         string              `json:"eventID"`
	SessionID       string              `json:"sessionID"`
	OrgID           string              `json:"orgID"`
	ServiceName     string              `json:"serviceName"`
	ServiceID       string              `json:"serviceID"`
	ServiceType     string              `json:"serviceType"`
	ServerIP        string              `json:"serverIP"`
	Privilege       string              `json:"privilege"`
	Email           string              `json:"email"`
	UserID          string              `json:"userID"`
	UserAgent       string              `json:"userAgent"`
	AccessDeviceID  string              `json:"accessDeviceID"`
	TfaDeviceID     string              `json:"tfaDeviceID"`
	BrowserID       string              `json:"browserID"`
	DeviceType      string              `json:"deviceType"`
	UserIP          string              `json:"userIP"`
	GeoLocation     models.GeoLocation  `json:"geoLocation"`
	Status          bool                `json:"status"`
	LoginTime       int64               `json:"loginTime"`
	LogoutTime      int64               `json:"logoutTime"`
	SessionDuration string              `json:"sessionDuration"`
	SessionRecord   bool                `json:"sessionRecord"`
	FailedReason    consts.FailedReason `json:"failedReason"`
	Guests          []string            `json:"guests"`
}

//ActiveSession is a struct which represents data stored in redis as active session
type ActiveSession struct {
	ConnID      string `json:"connID"`
	ServiceName string `json:"serviceName"`
	ServiceType string `json:"serviceType"`
	ServiceID   string `json:"serviceID"`
	ServerIP    string `json:"serverIP"`
	Email       string `json:"email"`
	LoginTime   int64  `json:"loginTime"`
	UserAgent   string `json:"userAgent"`
	Privilege   string `json:"privilege"`
	UserIP      string `json:"userIP"`
}

//newActiveSession creates new ActiveSession object
func newActiveSession(session *AuthLog, connID, appType string) (string, error) {

	newSession := ActiveSession{
		ServiceName: session.ServiceName,
		ServiceID:   session.ServiceID,
		ServiceType: appType,
		ConnID:      connID,
		Email:       session.Email,
		ServerIP:    session.ServerIP,
		LoginTime:   session.LoginTime,
		UserAgent:   session.UserAgent,
		Privilege:   session.Privilege,
		UserIP:      session.UserIP,
	}

	mars, err := json.Marshal(newSession)
	return string(mars), err
}

func NewEmptyLog(sType string) AuthLog {
	eventID := utils.GetUUID()
	authLog := AuthLog{
		EventID:       eventID,
		ServiceType:   string(sType),
		SessionID:     eventID,
		Status:        false,
		LoginTime:     time.Now().UnixNano(),
		LogoutTime:    time.Now().UnixNano(),
		SessionRecord: false,
	}

	if global.GetConfig().Platform.Base == "private" {
		authLog.OrgID = global.GetConfig().Trasa.OrgId
	}

	return authLog
}

func NewLog(r *http.Request, sType string) AuthLog {
	ip := utils.GetIp(r)

	geo, err := misc.Store.GetGeoLocation(ip)
	if err != nil {
		logrus.Trace(err)
	}

	eventID := utils.GetUUID()

	authLog := AuthLog{
		EventID:       eventID,
		ServiceType:   string(sType),
		SessionID:     eventID,
		UserIP:        ip,
		GeoLocation:   geo,
		UserAgent:     r.UserAgent(),
		Status:        false,
		LoginTime:     time.Now().UnixNano(),
		LogoutTime:    time.Now().UnixNano(),
		SessionRecord: false,
	}

	if global.GetConfig().Platform.Base == "private" {
		authLog.OrgID = global.GetConfig().Trasa.OrgId
	}

	return authLog
}

// NewHTTPAuthLog returns initialized Auth log struct
func NewHTTPAuthLog(endpoint consts.ConstEndpoints, ip, sessionID, userAgent string, reason consts.FailedReason, privilege string, status bool, user *models.User, service *models.Service, deviceID, domainName string) (AuthLog, error) {
	geo, err := misc.Store.GetGeoLocation(ip)
	if err != nil {
		logrus.Trace(err)
	}

	var session AuthLog
	session.EventID = utils.GetRandomID(5)
	session.SessionID = sessionID
	session.OrgID = user.OrgID
	session.UserID = user.ID
	session.Email = user.Email
	session.Privilege = privilege
	session.ServiceID = service.ID
	session.ServiceName = service.Name
	session.ServiceType = service.Type
	session.UserAgent = userAgent
	session.DeviceType = ""
	session.ServerIP = domainName
	session.TfaDeviceID = deviceID
	session.UserIP = ip
	//session.ServerIP = "0.0.0.0"
	session.GeoLocation.IsoCountryCode = geo.IsoCountryCode
	session.GeoLocation.City = geo.City
	session.GeoLocation.TimeZone = geo.TimeZone
	session.GeoLocation.Location = geo.Location
	//log.GeoLocation.Longitude =
	session.Status = status
	session.LoginTime = time.Now().UnixNano()
	session.FailedReason = reason
	return session, err
}

func (l *AuthLog) UpdateUser(user *models.UserWithPass) {
	l.Email = user.Email
	l.UserID = user.ID
	if l.ServiceType == "dashboard" {
		l.Privilege = user.UserRole
	}

	if user.OrgID != "" {
		l.OrgID = user.OrgID
	}
}

func (l *AuthLog) UpdateService(service *models.Service) {
	l.ServiceID = service.ID
	l.ServiceName = service.Name
	l.ServiceType = service.Type
	l.ServerIP = service.Hostname

	if service.OrgID != "" {
		l.OrgID = service.OrgID
	}
}

func (l *AuthLog) UpdateAddr(addr net.Addr) {
	ip := utils.GetIPFromAddr(addr)

	geo, err := misc.Store.GetGeoLocation(ip)
	if err != nil {
		logrus.Trace(err)
	}

	l.GeoLocation = geo
	l.UserIP = ip
}
func (l *AuthLog) UpdateIP(ip string) {

	geo, err := misc.Store.GetGeoLocation(ip)
	if err != nil {
		logrus.Trace(err)
	}

	l.GeoLocation = geo
	l.UserIP = ip
}

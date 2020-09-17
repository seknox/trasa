package notif

import (
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	Store = notifStore{State: state}
}

//InitStoreMock will init mock state of this package
func InitStoreMock() *notifMock {
	m := new(notifMock)
	Store = m
	return m
}

//Store is the package state variable which contains database connections
var Store adapter

type notifStore struct {
	*global.State
}

type adapter interface {
	GetPendingNotif(userID, orgID string) ([]models.InAppNotification, error)
	StoreNotif(notif models.InAppNotification) (err error)
	UpdateNotif(notif models.InAppNotification) error
	UpdateNotifFromNotifID(notif models.InAppNotification) error

	SendEmail(orgID string, emailType consts.EmailType, emailTemplate interface{}) error

	SendPushNotification(fcmToken, orgName, appName, ipAddr, time, challenge string) error
	CallTrasaCloudProxy(path string, reqBody interface{}, insecure bool) (resp models.TrasaResponseStruct, err error)
}

package notif

import (
	"github.com/seknox/trasa/consts"
	"github.com/seknox/trasa/global"
	"github.com/seknox/trasa/models"
)

func InitStore(state *global.State) {
	Store = NotifStore{State: state}
}

func InitStoreMock() *NotifMock {
	m := new(NotifMock)
	Store = m
	return m
}

var Store NotifAdapter

type NotifStore struct {
	*global.State
}

type NotifAdapter interface {
	GetPendingNotif(userID, orgID string) ([]models.InAppNotification, error)
	StoreNotif(notif models.InAppNotification) (err error)
	UpdateNotif(notif models.InAppNotification) error
	UpdateNotifFromNotifID(notif models.InAppNotification) error

	SendEmail(orgID string, emailType consts.EmailType, emailTemplate interface{}) error

	SendPushNotification(fcmToken, orgName, appName, ipAddr, time, challenge string) error
}

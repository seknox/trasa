package notif

import (
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

type notifMock struct {
	mock.Mock
}

func (n *notifMock) GetPendingNotif(userID, orgID string) ([]models.InAppNotification, error) {
	panic("implement me")
}

func (n *notifMock) StoreNotif(notif models.InAppNotification) (err error) {
	panic("implement me")
}

func (n *notifMock) UpdateNotif(notif models.InAppNotification) error {
	panic("implement me")
}

func (n *notifMock) UpdateNotifFromNotifID(notif models.InAppNotification) error {
	panic("implement me")
}

func (n *notifMock) SendEmail(orgID string, emailType consts.EmailType, emailTemplate interface{}) error {
	args := n.Called(orgID, emailType)
	if emailTemplate == nil {
		panic(emailTemplate)
	}

	return args.Error(0)
}

func (n *notifMock) SendPushNotification(fcmToken, orgName, appName, ipAddr, time, challenge string) error {
	panic("implement me")
}

package notif

import (
	"github.com/seknox/trasa/consts"
	"github.com/seknox/trasa/models"
	"github.com/stretchr/testify/mock"
)

type NotifMock struct {
	mock.Mock
}

func (n NotifMock) GetPendingNotif(userID, orgID string) ([]models.InAppNotification, error) {
	panic("implement me")
}

func (n NotifMock) StoreNotif(notif models.InAppNotification) (err error) {
	panic("implement me")
}

func (n NotifMock) UpdateNotif(notif models.InAppNotification) error {
	panic("implement me")
}

func (n NotifMock) UpdateNotifFromNotifID(notif models.InAppNotification) error {
	panic("implement me")
}

func (n NotifMock) SendEmail(orgID string, emailType consts.EmailType, emailTemplate interface{}) error {
	args := n.Called(orgID, emailType)
	if emailTemplate == nil {
		panic(emailTemplate)
	}

	return args.Error(0)
}

func (n NotifMock) SendPushNotification(fcmToken, orgName, appName, ipAddr, time, challenge string) error {
	panic("implement me")
}

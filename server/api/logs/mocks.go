package logs

import (
	"io"
	"os"

	"github.com/gorilla/websocket"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

type logsMock struct {
	mock.Mock
}

func (l *logsMock) LogSignup(signup *models.InitSignup) error {
	return l.Called(signup).Error(0)
}

func (l *logsMock) LogLogin(log *AuthLog, reason consts.FailedReason, status bool) error {
	if log == nil {
		panic(log)

	}

	return nil
	//return l.Called(log,reason,status).Error(0)
}

func (l *logsMock) GetLoginEvents(entityType, entityID, orgID string) (logEvents []AuthLog, err error) {
	args := l.Called(entityType, entityID, orgID)
	return args.Get(0).([]AuthLog), args.Error(0)

}

func (l *logsMock) GetLoginEventsByPage(entityType, entityID, orgID string, page int, size int, dateFrom, dateTo int64) ([]AuthLog, error) {
	args := l.Called(entityType, entityID, orgID, page, size, dateFrom, dateTo)
	return args.Get(0).([]AuthLog), args.Error(0)

}

//TODO implement methods

func (l *logsMock) AddNewActiveSession(session *AuthLog, connID, appType string) error {
	panic("implement me")
}

func (l *logsMock) RemoveActiveSession(connID string) error {
	panic("implement me")
}

func (l *logsMock) RemoveAllActiveSessions() {
	panic("implement me")
}

func (l *logsMock) ServeLiveSessions(ws *websocket.Conn) {
	panic("implement me")
}

func (l *logsMock) LogInAppTrail(ip, userAgent, description string, user *models.User, status bool) error {
	panic("implement me")
}

func (l *logsMock) GetOrgInAppTrails(orgID string, page int, size int, dateFrom, dateTo int64) ([]models.InAppTrail, error) {
	panic("implement me")
}

func (l *logsMock) GetFromMinio(path, bucketName string) (io.ReadSeeker, error) {
	panic("implement me")
}

func (l *logsMock) PutIntoMinio(path, filepath, bucketName string) error {
	panic("implement me")
}

func (l *logsMock) UploadHTTPLogToMinio(file *os.File, login AuthLog) error {
	panic("implement me")
}

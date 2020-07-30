package logs

import (
	"os"

	"github.com/gorilla/websocket"
	"github.com/minio/minio-go"
	"github.com/seknox/trasa/consts"
	"github.com/seknox/trasa/models"
	"github.com/stretchr/testify/mock"
)

type LogsMock struct {
	mock.Mock
}

func (l LogsMock) LogSignup(signup *models.InitSignup) error {
	return l.Called(signup).Error(0)
}

func (l LogsMock) LogLogin(log *AuthLog, reason consts.FailedReason, status bool) error {
	if log == nil {
		panic(log)

	}

	return nil
	//return l.Called(log,reason,status).Error(0)
}

func (l LogsMock) GetLoginEvents(entityType, entityID, orgID string) (logEvents []AuthLog, err error) {
	args := l.Called(entityType, entityID, orgID)
	return args.Get(0).([]AuthLog), args.Error(0)

}

func (l LogsMock) GetLoginEventsByPage(entityType, entityID, orgID string, page int, size int, dateFrom, dateTo int64) ([]AuthLog, error) {
	args := l.Called(entityType, entityID, orgID, page, size, dateFrom, dateTo)
	return args.Get(0).([]AuthLog), args.Error(0)

}

//TODO implement methods

func (l LogsMock) AddNewActiveSession(session *AuthLog, connID, appType string) error {
	panic("implement me")
}

func (l LogsMock) RemoveActiveSession(connID string) error {
	panic("implement me")
}

func (l LogsMock) RemoveAllActiveSessions() {
	panic("implement me")
}

func (l LogsMock) ServeLiveSessions(ws *websocket.Conn) {
	panic("implement me")
}

func (l LogsMock) LogInAppTrail(ip, userAgent, description string, user *models.User, status bool) error {
	panic("implement me")
}

func (l LogsMock) GetOrgInAppTrails(orgID string, page int, size int, dateFrom, dateTo int64) ([]models.InAppTrail, error) {
	panic("implement me")
}

func (l LogsMock) GetFromMinio(path, bucketName string) (*minio.Object, error) {
	panic("implement me")
}

func (l LogsMock) PutIntoMinio(path, filepath, bucketName string) error {
	panic("implement me")
}

func (l LogsMock) UploadHTTPLogToMinio(file *os.File, login AuthLog) error {
	panic("implement me")
}

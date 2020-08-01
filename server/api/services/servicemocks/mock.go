package servicemocks

import (
	"errors"

	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

type ServiceStoreMock struct {
	mock.Mock
}

func (as ServiceStoreMock) GetByID(serviceID, orgID string) (*models.Service, error) {
	args := as.Called(serviceID, orgID)
	return args.Get(0).(*models.Service), args.Error(1)
}

func (as ServiceStoreMock) GetTotalForOrg(orgID string) (int64, error) {
	args := as.Called(orgID)
	return int64(args.Int(0)), args.Error(1)
}

func (as ServiceStoreMock) CreateApp(app *models.Service) error {
	args := as.Called(app)
	if app.ID == "" || app.OrgID == "" || app.Hostname == "" {
		return errors.New("")
	}
	return args.Error(1)
}

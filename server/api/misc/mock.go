package misc

import (
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

type miscMock struct {
	mock.Mock
}

func (m *miscMock) GetGeoLocation(ip string) (geo models.GeoLocation, err error) {
	return models.GeoLocation{}, err
	//panic("implement me")
}

func (m *miscMock) GetEntityDescription(entityID string, entityType consts.EntityConst, orgID string) (string, string, error) {
	panic("implement me")
}

//GetAdmins returns all users of a organizations
func (s *miscMock) GetAdmins(orgID string) ([]models.User, error) {
	panic("implement me")

}

package misc

import (
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/stretchr/testify/mock"
)

type MiscMock struct {
	mock.Mock
}

func (m MiscMock) GetGeoLocation(ip string) (geo models.GeoLocation, err error) {
	return models.GeoLocation{}, err
	//panic("implement me")
}

func (m MiscMock) GetEntityDescription(entityID string, entityType consts.EntityConst, orgID string) (string, string, error) {
	panic("implement me")
}

//GetAdmins returns all users of a organizations
func (s MiscMock) GetAdmins(orgID string) ([]models.User, error) {
	panic("implement me")

}

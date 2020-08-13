package misc

import (
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	Store = MiscStore{State: state}
}

func InitMock() *MiscMock {
	m := new(MiscMock)
	Store = m
	return m
}

//Store is the package state variable which contains database connections
var Store MiscAdapter

type MiscStore struct {
	*global.State
}

type MiscAdapter interface {
	GetAdmins(orgID string) ([]models.User, error)

	GetGeoLocation(ip string) (geo models.GeoLocation, err error)

	GetEntityDescription(entityID string, entityType consts.EntityConst, orgID string) (string, string, error)
}

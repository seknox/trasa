package uidp

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	Store = idpStore{state}
}

//Store is the package state variable which contains database connections
var Store adapter

type idpStore struct {
	*global.State
}

type adapter interface {
	// user idp
	GetAllIdps(orgID string) ([]models.IdentityProvider, error)
	GetByID(orgID, idpID string) (models.IdentityProvider, error)
	GetByName(orgID, idpName string) (models.IdentityProvider, error)
	CreateIDP(idp *models.IdentityProvider) error
	UpdateIDP(idp *models.IdentityProvider) error
	UpdateLDAPIDP(idp *models.IdentityProvider) error
	activateOrDisableIdp(orgID, idpID string, updateTime int64, updateVal bool) error
}

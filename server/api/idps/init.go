package idps

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

func InitStore(state *global.State) {
	Store = IDPStore{state}
}

var Store IDPAdapter

type IDPStore struct {
	*global.State
}

type IDPAdapter interface {
	// user idp
	GetAllIdps(orgID string) ([]models.IdentityProvider, error)
	GetByID(orgID, idpID string) (models.IdentityProvider, error)
	CreateIDP(idp *models.IdentityProvider) error
	UpdateIDP(idp *models.IdentityProvider) error
	UpdateLDAPIDP(idp *models.IdentityProvider) error
	activateOrDisableIdp(orgID, idpID string, updateTime int64, updateVal bool) error

	// cloudIaas idp
	GetCloudSyncState(orgID, cName string) (*models.CloudIaaSSync, error)
}

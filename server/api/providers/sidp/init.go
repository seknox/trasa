package sidp

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

	// cloudIaas idp
	GetCloudSyncState(orgID, cName string) (*models.CloudIaaSSync, error)
}

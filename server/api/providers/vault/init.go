package vault

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	Store = cryptStore{State: state}
}

//Store is the package state variable which contains database connections
var Store adapter

type adapter interface {
	StoreKey(k models.KeysHolder) error
	GetKeyOrTokenWithTag(orgID string, keyName string) (*models.KeysHolder, error)
	GetKeyOrTokenWithKeyval(orgID, keyName string) (*models.KeysHolder, error)
	GetKeyOrTokenWithKeyvalAndID(orgID, keyName, keyID string) (*models.KeysHolder, error)
}

type cryptStore struct {
	*global.State
}

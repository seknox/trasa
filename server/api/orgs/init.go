package orgs

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	Store = orgStore{State: state}
}

func InitStoreMock() *OrgMock {
	lmock := new(OrgMock)
	Store = lmock
	return lmock
}

//Store is the package state variable which contains database connections
var Store orgAdapter

type orgStore struct {
	*global.State
}

type orgAdapter interface {
	Get(orgID string) (models.Org, error)
	CheckOrgExists() (orgID string, err error)
	CreateOrg(org *models.Org) error
	GetIDP(orgID, idpName string) (models.IdentityProvider, error)

	RemoveAllManagedAccounts(orgID string) error
}

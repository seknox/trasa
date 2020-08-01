package orgs

import (
	"github.com/seknox/trasa/global"
	"github.com/seknox/trasa/models"
)

func InitStore(state *global.State) {
	Store = OrgStore{State: state}
}

func InitStoreMock() *OrgMock {
	lmock := new(OrgMock)
	Store = lmock
	return lmock
}

var Store OrgAdapter

type OrgStore struct {
	*global.State
}

type OrgAdapter interface {
	Get(orgID string) (models.Org, error)
	CheckOrgExists() (orgID string, err error)
	CreateOrg(org *models.Org) error
	GetIDP(orgID, idpName string) (models.IdentityProvider, error)

	RemoveAllManagedAccounts(orgID string) error
}

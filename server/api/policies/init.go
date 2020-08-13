package policies

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	Store = policyStore{State: state}
}

//Store is the package state variable which contains database connections
var Store adapter

type policyStore struct {
	*global.State
}

type adapter interface {
	GetPolicy(policyID, orgID string) (models.Policy, error)
	GetAllPolicies(orgID string) ([]models.Policy, error)
	CreatePolicy(policy models.Policy) error
	UpdatePolicy(policy models.Policy) error
	DeletePolicy(policyID, orgID string) error

	GetAccessPolicy(userID, serviceID, privilege, orgID string) (policy *models.Policy, adhoc bool, err error)
}

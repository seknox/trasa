package policies

import (
	"github.com/seknox/trasa/global"
	"github.com/seknox/trasa/models"
)

func InitStore(state *global.State) {
	Store = PolicyStore{State: state}
}

var Store PolicyAdapter

type PolicyStore struct {
	*global.State
}

type PolicyAdapter interface {
	GetPolicy(policyID, orgID string) (models.Policy, error)
	GetAllPolicies(orgID string) ([]models.Policy, error)
	CreatePolicy(policy models.Policy) error
	UpdatePolicy(policy models.Policy) error
	DeletePolicy(policyID, orgID string) error

	GetAccessPolicy(userID, serviceID, orgID string) (policy *models.Policy, privilege string, adhoc bool, err error)
}

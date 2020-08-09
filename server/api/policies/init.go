package policies

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
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

	GetAccessPolicy(userID, serviceID, privilege, orgID string) (policy *models.Policy, adhoc bool, err error)
}

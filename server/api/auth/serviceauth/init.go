package serviceauth

import (
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

func InitStore(state *global.State, policyFunc models.CheckPolicyFunc) {
	Store = SStore{
		State:           state,
		checkPolicyFunc: policyFunc,
	}
}

var Store SAdapter

type SStore struct {
	*global.State
	checkPolicyFunc models.CheckPolicyFunc
}

type SAdapter interface {
	CheckPolicy(serviceID, userID, orgID, userIP, timezone string, policy *models.Policy, adhoc bool) (bool, consts.FailedReason)
}

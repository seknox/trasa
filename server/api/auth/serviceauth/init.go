package serviceauth

import (
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	Store = sStore{
		State: state,
	}
}

//Store is the package state variable which contains database connections
var Store adapter

type sStore struct {
	*global.State
}

type adapter interface {
	CheckPolicy(serviceID, userID, orgID, userIP, timezone string, policy *models.Policy, adhoc bool) (bool, consts.FailedReason)
}

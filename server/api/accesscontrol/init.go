package accesscontrol

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

//InitStore initialises package state
func InitStore(state *global.State, checkPolicyFunc models.CheckPolicyFunc) {
	ACStore = store{
		State:           state,
		CheckPolicyFunc: checkPolicyFunc,
	}
}

//TODO use this store for ssh and rdp also

//TODO use interface
var ACStore store

type store struct {
	*global.State
	CheckPolicyFunc models.CheckPolicyFunc
}

type adapter interface {
}

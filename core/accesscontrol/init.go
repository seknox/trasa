package accesscontrol

import (
	"github.com/seknox/trasa/global"
	"github.com/seknox/trasa/models"
)

func InitStore(state *global.State, checkPolicyFunc models.CheckPolicyFunc) {
	ACStore = Store{
		State:           state,
		CheckPolicyFunc: checkPolicyFunc,
	}
}

//TODO use this store for ssh and rdp also

//TODO use interface
var ACStore Store

type Store struct {
	*global.State
	CheckPolicyFunc models.CheckPolicyFunc
}

type Adapter interface {
}

package accesscontrol

import (
	"github.com/seknox/trasa/server/global"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	ACStore = store{
		State: state,
	}
}

//TODO use interface
var ACStore store

type store struct {
	*global.State
}

type adapter interface {
}

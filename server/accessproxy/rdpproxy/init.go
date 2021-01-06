package rdpproxy

import (
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/global"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	Store = GWStore{
		proxy: NewProxy(),
		State: state,
	}
}

//Store is the package state variable which contains database connections
var Store GWAdapter

type GWStore struct {
	*global.State
	proxy *Proxy
}

type GWAdapter interface {
	uploadSessionLog(authlog *logs.AuthLog) error
}

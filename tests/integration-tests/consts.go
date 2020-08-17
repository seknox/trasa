package integration_tests

import (
	"github.com/gorilla/websocket"
	"net/http"
)

const (
	trasaEmail   = "root"
	trasaPass    = "changeme"
	totpSEC      = "AV2COXZHVG4OAFSF"
	upstreamUser = "bhrg3se"
	upstreamPass = "root"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		//TODO
		return true
	},
	Subprotocols: []string{"trasa", "guacamole", "livesessions", "xterm"},
}

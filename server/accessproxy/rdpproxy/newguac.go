package rdpproxy

import (
	"time"

	protocol2 "github.com/seknox/guacamole"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/models"

	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Proxy struct {
	mu       sync.Mutex
	upgrader *websocket.Upgrader
}

const (
	SocketTimeout  = 15 * time.Second
	MaxGuacMessage = 8192
)

//TODO remove proxy struct
func NewProxy() *Proxy {

	return &Proxy{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  websocketReadBufferSize,
			WriteBufferSize: websocketWriteBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true // TODO
			},
			Subprotocols: []string{"guacamole"}, // fixed by guacamole-client
		},
	}

}

// ServeWS serves rdp connection through guacamole proxy
func (p *Proxy) ServeWS(params models.ConnectionParams, uc models.UserContext, ws *websocket.Conn) {

	checkAndInitParams(&uc, &params)
	authlog := logs.NewEmptyLog("rdp")
	authlog.UpdateIP(params.UserIP)
	authlog.UpdateUser(&models.UserWithPass{User: *uc.User})
	authlog.AccessDeviceID = uc.DeviceID
	authlog.BrowserID = uc.BrowserID
	authlog.Privilege = params.Privilege
	authlog.UserIP = params.UserIP
	authlog.OrgID = uc.Org.ID
	params.SessionID = authlog.SessionID

	//If connection is not "joined" one, write authlog
	//Joined connections have connID, new connection don't
	if params.ConnID == "" {
		defer logSession(&authlog)
	}

	newSession, err := NewSession(&params, &authlog)
	if err != nil {
		logrus.Debug(err)
		// code 3339 is for trasa related errors
		p.sendError(ws, err.Error(), "3339")
		ws.WriteMessage(websocket.CloseMessage, nil)
		return
	}

	errcode, err := newSession.Start(ws)
	if err != nil {
		logrus.Debug(err)
		p.sendError(ws, err.Error(), errcode)
	}

	ws.Close()
}

//sendError Sends  error to client through guacamole protocol
func (p *Proxy) sendError(ws *websocket.Conn, msg, code string) error {
	inst := protocol2.NewInstruction("error", msg, code)
	logrus.Debug(inst.String())
	return ws.WriteMessage(websocket.TextMessage, inst.Byte())
}

func checkAndInitParams(uc *models.UserContext, params *models.ConnectionParams) {
	//TODO
	params.OrgID = uc.Org.ID
	params.UserID = uc.User.ID
	params.TrasaID = uc.User.Email
	params.Timezone = uc.Org.Timezone
	params.ServiceType = "rdp"
	params.OrgName = uc.Org.OrgName
	params.AccessDeviceID = uc.DeviceID
	params.BrowserID = uc.BrowserID
	params.Groups = uc.User.Groups
	//params.UserAgent = r.UserAgent()

	if params.RdpProtocol == "" {
		params.RdpProtocol = "nla"
	}

}

func logSession(authlog *logs.AuthLog) {

	err := logs.Store.LogLogin(authlog, authlog.FailedReason, authlog.Status)
	if err != nil {
		logrus.Errorf("failed to log.  trying again: %v", err)
		logs.Store.LogLogin(authlog, authlog.FailedReason, authlog.Status)
	}

	if !authlog.SessionRecord {
		return
	}

	err = Store.uploadSessionLog(authlog)
	if err != nil {
		logrus.Error(err)
	}

}

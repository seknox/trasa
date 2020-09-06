package rdpproxy

import (
	"bytes"
	"database/sql"
	"github.com/seknox/trasa/server/global"
	"io"
	"net"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/seknox/guacamole"
	"github.com/seknox/trasa/server/api/accesscontrol"
	"github.com/seknox/trasa/server/api/accessmap"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/policies"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

const (
	websocketReadBufferSize  = guacamole.MaxGuacMessage
	websocketWriteBufferSize = guacamole.MaxGuacMessage * 2

	//TODO calibrate the right number of instructions before assuming connection is successful

	initialInstructions int = 40
)

type Session struct {
	isOwner bool
	log     *logs.AuthLog
	params  *models.ConnectionParams
	policy  *models.Policy
	service *models.Service
	tunnel  guacamole.Tunnel
}

// NewSession Create new guacamole session
func NewSession(params *models.ConnectionParams, authlog *logs.AuthLog) (*Session, error) {

	service, err := services.Store.GetFromHostname(params.Hostname, "rdp", "", params.OrgID)
	if errors.Is(err, sql.ErrNoRows) {
		service, err = accessmap.CreateDynamicService(params.Hostname, "rdp", params.UserID, params.TrasaID, params.Privilege, params.OrgID)
		if err != nil {
			authlog.FailedReason = consts.REASON_DYNAMIC_SERVICE_FAILED
			return nil, err
		}
	} else if err != nil {
		authlog.FailedReason = consts.REASON_INVALID_SERVICE_CREDS
		return nil, err
	}

	authlog.UpdateService(service)
	params.ServiceName = service.Name

	session := Session{service: service}

	if params.ConnID == "" {
		session.isOwner = true
	}

	if session.isOwner {

		policy, _, err := policies.Store.GetAccessPolicy(params.UserID, service.ID, params.Privilege, params.OrgID)
		if errors.Is(err, sql.ErrNoRows) {
			//if service is not assigned to user, create one (only if dynamic access is enabled)
			policy, err = accessmap.CreateDynamicAccessMap(params.SessionID, params.UserID, params.TrasaID, params.Privilege, params.OrgID)
			if err != nil {
				logrus.Errorf("dynamic access map: %v", err)
				return nil, err
			}

		} else if err != nil {
			logrus.Errorf("get service from hostname: %v", err)
			return nil, err
		}

		params.CanTransferFile = policy.FileTransfer
		params.SessionRecord = policy.RecordSession
		params.Skip2FA = !policy.TfaRequired

		session.policy = policy

		//ok, reason := handlePolicy(params, service.ID)
		//if !ok {
		//	authlog.FailedReason = reason
		//	return nil, errors.Errorf("policy failed: %v", reason)
		//}
		authlog.SessionRecord = params.SessionRecord

	} else {
		params.Skip2FA = false
		params.CanTransferFile = false
		params.SessionRecord = false
	}

	creds := &models.UpstreamCreds{}

	if session.isOwner {
		creds, err = handlePass(params)
		if err != nil {
			logrus.Debug(err)
			return nil, err
		}

		params.ConnID = utils.GetUUID()
		session.log = authlog

	}

	config, err := makeConfig(params, creds)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if !session.isOwner {
		config.ConnectionID = params.ConnID
	}

	guacdAddr := global.GetConfig().Proxy.GuacdAddr
	if guacdAddr == "" {
		guacdAddr = "127.0.0.1:4822"
	}

	addr, err := net.ResolveTCPAddr("tcp", guacdAddr)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		logrus.Errorln("error while connecting to guacd", err)
		return nil, err
	}

	stream := guacamole.NewStream(conn, guacamole.SocketTimeout)

	logrus.Debugf("Starting handshake with ")
	err = stream.Handshake(config)
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}
	logrus.Debug("Socket configured")

	tunnel := guacamole.NewSimpleTunnel(stream, params.ConnID)

	session.tunnel = tunnel
	session.params = params
	return &session, nil

}

func (s *Session) Start(ws *websocket.Conn) (errcode string, err error) {

	if s.tunnel == nil {
		return "", errors.New("Session not initialised")
	}

	writer := s.tunnel.AcquireWriter()
	defer s.tunnel.ReleaseWriter()
	reader := s.tunnel.AcquireReader()
	defer s.tunnel.ReleaseReader()

	errcode, err = catchInitialErrors(*ws, writer, reader)
	if err != nil {
		logrus.Error(errcode, err)
		if errcode == "519" || errcode == "769" {
			s.log.FailedReason = consts.REASON_INVALID_USER_CREDS
			return "3339", errors.New("Invalid credentials")
		}
		return errcode, err
	}

	//prevent timeout

	//logrus.Debug("auth success,  tfa")

	tfaSuccess := true
	var tfaDeviceID string
	tfaSuccess, tfaDeviceID = s.handleTfa(ws, writer)

	//logrus.Debug("tfa success")

	if !tfaSuccess {
		return "3339", errors.New("TFA Failed")
	}

	if s.isOwner {

		reason, ok, err := accesscontrol.CheckDevicePolicy(s.policy.DevicePolicy, s.params.AccessDeviceID, tfaDeviceID, s.params.OrgID)

		if err != nil {
			logrus.Error(err)
		}

		if !ok {
			return "3339", errors.Errorf("Device policy failed: %v", reason)
		}

		ok, reason = Store.CheckPolicy(s.params, s.policy, s.service.Adhoc)
		if !ok {
			return "3339", errors.Errorf("policy failed: %v", reason)
		}

		s.log.Status = true
		err = logs.Store.AddNewActiveSession(s.log, s.tunnel.ConnectionID(), "rdp")
		if err != nil {
			logrus.Error(err)
		}
		defer logs.Store.RemoveActiveSession(s.tunnel.ConnectionID())

	}

	sessionDone := make(chan bool, 2)

	go func() {
		<-sessionDone
		s.tunnel.Close()
		ws.Close()
	}()
	go WSToGuacd(ws, writer, sessionDone)
	GuacdToWS(ws, reader, sessionDone)

	logrus.Trace("Session ended ", s.isOwner)
	return "", nil
}

//It will listen for error within first few instructions
//If everything seem fine continue to serveIO
func catchInitialErrors(ws websocket.Conn, guacdWriter io.Writer, guacdReader guacamole.InstructionReader) (errcode string, err error) {
	wg := sync.WaitGroup{}
	exit := make(chan error, 2)
	wg.Add(2)
	var done = false

	go func(conn guacamole.InstructionReader, ws websocket.Conn) {
		var err error
		var raw []byte
		var inst *guacamole.Instruction
		for i := 0; i < initialInstructions; i++ {
			raw, err = conn.ReadSome()
			if err != nil {
				break
			}

			//logrus.Debug(string(raw))

			inst, err = guacamole.Parse(raw)
			if err != nil {
				break
			}

			//If opcode is error return
			if inst.Opcode == "error" {
				if len(inst.Args) != 2 {
					err = errors.New("trasa: invalid number of error parameters.")
					break
				}
				errcode = inst.Args[1]
				err = errors.New(inst.Args[0])
				break
			}
			err = ws.WriteMessage(websocket.TextMessage, raw)
			if err != nil {
				break
			}
		}
		done = true
		exit <- err
		logrus.Info("trasa: reading from guacd terminated.", err)
		wg.Done()
	}(guacdReader, ws)

	go func(conn io.Writer, ws websocket.Conn) {
		var err error
		var buf []byte
		for !done {
			_, buf, err = ws.ReadMessage()
			if err != nil {
				break
			}
			_, err = conn.Write(buf)
			if err != nil {
				break
			}
		}

		//If err==nil loop was closed by done=false
		// which is caused by catching error opcode from guacd
		if err != nil {
			exit <- err
		}
		logrus.Info("trasa: reading from client terminated.", err)
		wg.Done()
	}(guacdWriter, ws)
	err = <-exit
	done = true

	//conn.Close()
	wg.Wait()
	logrus.Info("trasa: IO goroutines are terminated.", err)
	return
}

// MessageReader wraps a websocket connection and only permits Reading
type MessageReader interface {
	// ReadMessage should return a single complete message to send to guac
	ReadMessage() (int, []byte, error)
}

func WSToGuacd(ws MessageReader, guacd io.Writer, done chan bool) {
	defer func() { done <- true }()
	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			logrus.Traceln("Error reading message from ws", err)
			return
		}

		if bytes.HasPrefix(data, guacamole.InternalOpcodeIns) {
			// messages starting with the InternalDataOpcode are never sent to guacd
			continue
		}

		if _, err = guacd.Write(data); err != nil {
			logrus.Traceln("Failed writing to guacd", err)
			return
		}
	}
}

// MessageWriter wraps a websocket connection and only permits Writing
type MessageWriter interface {
	// WriteMessage writes one or more complete guac commands to the websocket
	WriteMessage(int, []byte) error
}

func GuacdToWS(ws MessageWriter, guacd guacamole.InstructionReader, done chan bool) {
	defer func() { done <- true }()
	buf := bytes.NewBuffer(make([]byte, 0, guacamole.MaxGuacMessage*2))

	for {
		ins, err := guacd.ReadSome()
		if err != nil {
			logrus.Traceln("Error reading from guacd", err)
			return
		}

		if bytes.HasPrefix(ins, guacamole.InternalOpcodeIns) {
			// messages starting with the InternalDataOpcode are never sent to the websocket
			continue
		}

		if _, err = buf.Write(ins); err != nil {
			logrus.Traceln("Failed to buffer guacd to ws", err)
			return
		}

		// if the buffer has more data in it or we've reached the max buffer size, send the data and reset
		if !guacd.Available() || buf.Len() >= guacamole.MaxGuacMessage {
			if err = ws.WriteMessage(1, buf.Bytes()); err != nil {
				if err == websocket.ErrCloseSent {
					return
				}
				logrus.Traceln("Failed sending message to ws", err)
				return
			}
			buf.Reset()
		}
	}
}

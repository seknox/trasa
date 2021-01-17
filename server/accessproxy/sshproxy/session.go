package sshproxy

import (
	"encoding/hex"
	"net"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/seknox/ssh"
	"github.com/seknox/ssh/agent"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

type Session struct {
	ID              string
	AuthType        consts.SSH_AUTH_TYPE
	log             *logs.AuthLog
	tempSessionFile *os.File
	params          *models.ConnectionParams
	policy          *models.Policy
	Conn            *net.Conn
	clientConfig    *ssh.ClientConfig
	sshClient       *ssh.Client
	guestChan       chan GuestClient
	guests          []GuestClient
}

type GuestClient struct {
	UserID string
	Email  string
	Conn   *websocket.Conn
}

func (s *Session) appendGuest() {

}

func (s *Session) UpdateConMeta(connMeta ssh.ConnMetadata) {
	s.ID = hex.EncodeToString(connMeta.SessionID())

	s.params.Privilege = connMeta.User()
	//s.log.EventID = s.ID
	//s.log.SessionID = s.ID
	s.log.Privilege = connMeta.User()
	s.log.ServiceType = "ssh"

	s.log.LoginTime = time.Now().UnixNano()
	s.log.UserAgent = string(connMeta.ClientVersion())

	s.log.UserIP = utils.GetIPFromAddr(connMeta.RemoteAddr())

}

func (s *Session) UpdateService(params *models.Service) {

	s.params.ServiceName = params.Name
	s.params.ServiceID = params.ID
	s.params.Hostname = params.Hostname
	//	s.params.Policy = params.Policy

	//s.params = params
	s.log.ServiceName = params.Name
	s.log.ServiceID = params.ID
	s.log.OrgID = params.OrgID
	s.log.ServerIP = params.Hostname
	///s.log.SessionRecord = params.Policy.RecordSession

	s.log.ServiceType = "ssh"
	//s.log.Email = params.Email
	s.log.LoginTime = time.Now().UnixNano()

}

func NewSession(serverConn *net.Conn) *Session {

	authlog := logs.NewEmptyLog("ssh")
	authlog.UpdateAddr((*serverConn).RemoteAddr())

	//	logrus.Debug((*serverConn).RemoteAddr())
	params := models.ConnectionParams{
		//TODO
		OrgID: global.GetConfig().Trasa.OrgId,
		//User:            serverConn.User(),
		//SessionID:       string(serverConn.SessionID()),
		//UserIP:          utils.GetIPFromAddr(serverConn.RemoteAddr()),
		ServiceType: "ssh",
	}

	clientConfig := &ssh.ClientConfig{}

	clientConfig.SetDefaults()
	clientConfig.Ciphers = append(clientConfig.Ciphers, "aes128-cbc", "blowfish-cbc", "3des-cbc")

	s := Session{
		log:             &authlog,
		tempSessionFile: nil,
		params:          &params,
		Conn:            serverConn,
		clientConfig:    clientConfig,
		sshClient:       nil,
	}

	return &s

}

func start(conn net.Conn, serverConf *ssh.ServerConfig) error {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("Recovered in serve()", r, string(debug.Stack()))
		}
	}()

	frontEndConn, frontEndChans, frontEndReqs, err := ssh.NewServerConn(conn, serverConf)
	if err != nil {
		logrus.Debug("failed to handshake", err)
		return (err)
	}

	defer frontEndConn.Close()

	session, err := SSHStore.GetSession(conn.RemoteAddr())
	if err != nil {
		logrus.Error("session not found")
		return err
	}

	logs.Store.AddNewActiveSession(session.log, session.ID, "ssh")
	defer logs.Store.RemoveActiveSession(session.ID)

	var signers []ssh.Signer
	agentChan, agentReqs, err := frontEndConn.OpenChannel("auth-agent@openssh.com", nil)
	if err != nil {
		logrus.Debug(err)
		go ssh.DiscardRequests(agentReqs)
	} else {
		go ssh.DiscardRequests(agentReqs)
		signers, err = agent.NewClient(agentChan).Signers()
	}

	session.guestChan = SSHStore.CreateGuestChannel(hex.EncodeToString(frontEndConn.SessionID()))

	//	logrus.Trace("SESSION ID is", hex.EncodeToString(serverConn.SessionID()))

	clientConn, err := getClient(frontEndConn, signers)
	if err != nil {
		logrus.Error(err)
		return (err)
	}

	defer clientConn.Close()
	defer SSHStore.closeSession(conn.RemoteAddr())

	go ssh.DiscardRequests(frontEndReqs)
	//_ = reqs
	for frontEndChan := range frontEndChans {
		//logrus.Debug("||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||\n", frontEndChan.ExtraData(), frontEndChan.ChannelType())

		backEndChannel, backEndRequests, err2 := clientConn.OpenChannel(frontEndChan.ChannelType(), frontEndChan.ExtraData())
		if err2 != nil {
			logrus.Debugf("Could not accept client channel: %s", err2.Error())
			return err
		}

		acceptedFrontEndChannel, requests, err := frontEndChan.Accept()
		if err != nil {
			logrus.Debugf("Could not accept server channel: %s", err.Error())
			return err
		}
		defer acceptedFrontEndChannel.Close()
		var isSubSystem = make(chan bool, 1)

		// connect requests
		go func() {
			//log.Printf("Waiting for request")

		r:
			for {
				var req *ssh.Request
				var dst ssh.Channel

				select {
				case req = <-requests:
					dst = backEndChannel
				case req = <-backEndRequests:
					dst = acceptedFrontEndChannel
				}

				//	logrus.Debug( dst, req.WantReply, req.Payload)
				if req == nil {
					break r
				}

				b, err := dst.SendRequest(req.Type, req.WantReply, req.Payload)
				//	logrus.Debug("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n", string(req.Payload), req.Type, req.WantReply)
				if err != nil {
					logrus.Errorf("%s", err)
				}

				if req.WantReply {
					req.Reply(b, nil)
				}
				//logrus.Debug(req.Type, req.WantReply)
				switch req.Type {
				case "exit-status":
					close(isSubSystem)
					break r
				case "shell":
					isSubSystem <- false
				case "exec":
					isSubSystem <- false
					//	logrus.Debug("exec Type---------")
					req.Reply(false, nil)
					//break r
					// not supported (yet)
				case "subsystem":
					//If file transfer is not allowed close the connection
					//TODO skip session log for sftp
					if !session.policy.FileTransfer {
						logs.Store.LogLogin(session.log, consts.REASON_FILE_TRANSFER_NOT_ALLOWED, false)
						close(isSubSystem)
						break r
					}
					isSubSystem <- true
					logs.Store.LogLogin(session.log, "", true)
					//break r

				default:
					//logrus.Debug(req.Type)
				}
			}

			acceptedFrontEndChannel.Close()
			backEndChannel.Close()
		}()

		// connect channels
		//	log.Printf("Connecting channels.")

		//	var wrappedFrontEndChannel io.ReadCloser = acceptedFrontEndChannel

		//var wrappedChannel1  = NewTypeWriterReadCloser(acceptedFrontEndChannel)
		var wrappedBackendChannel *WrappedTunnel

		//wrappedFrontEndChannel, err = p.wrapFn(serverConn,acceptedFrontEndChannel)
		wrappedBackendChannel, err = NewWrappedTunnel(
			session.log.SessionID,
			session.policy.RecordSession && !<-isSubSystem, //skip session record for subsystem like sftp
			backEndChannel,
			acceptedFrontEndChannel,
			session.guestChan)
		if err != nil {
			logrus.Error(err)
			return err
		}
		defer wrappedBackendChannel.Close()

		wrappedBackendChannel.pipe()
	}

	return nil

}

//func pipe(frontEndchannel ssh.Channel, wrappedTunnel *WrappedTunnel)  {
//	go func() {
//		for !wrappedTunnel.closed{
//			buff := make([]byte, 100)
//			n, err := wrappedTunnel.Read(buff)
//			if err != nil {
//				logrus.Debug(err)
//				return
//			}
//			if n > 0 {
//				_,err := frontEndchannel.Write( buff)
//				if err != nil {
//					logrus.Debug(err)
//					return
//				}
//			}
//		}
//	}()
//
//
//	func (){
//		for !wrappedTunnel.closed {
//			buff := make([]byte, 100)
//			n, err := frontEndchannel.Read(buff)
//			if err != nil {
//				logrus.Debug(err)
//				return
//			}
//			if n > 0 {
//				_,err := wrappedTunnel.Write( buff)
//				if err != nil {
//					logrus.Debug(err)
//					return
//				}
//			}
//		}
//	}()
//
//	logrus.Debug("Session Ended")
//}

func getClient(c ssh.ConnMetadata, signers []ssh.Signer) (*ssh.Client, error) {
	sess, err := SSHStore.GetSession(c.RemoteAddr())
	if err != nil {
		return nil, err
	}

	if len(signers) != 0 {
		//logrus.Debug(signers)
		sess.clientConfig.Auth = append(sess.clientConfig.Auth, ssh.PublicKeys(signers...))

	}

	if !strings.Contains(sess.params.Hostname, ":") {
		sess.params.Hostname = sess.params.Hostname + ":22"
	}

	client, err := ssh.Dial("tcp", sess.params.Hostname, sess.clientConfig)
	if err != nil {
		if strings.Contains(err.Error(), "ssh: unable to authenticate") {
			logrus.Debug(err)
			logs.Store.LogLogin(sess.log, consts.REASON_INVALID_USER_CREDS, false)
		} else if strings.Contains(err.Error(), "ssh: handshake failed: Could not verify upstream host key") {
			logrus.Debug(err)
			logs.Store.LogLogin(sess.log, consts.REASON_INVALID_HOST_KEY, false)
		} else {
			logrus.Error(err)
			logs.Store.LogLogin(sess.log, consts.REASON_UNKNOWN, false)
		}

		return nil, err
	}

	logrus.Tracef("Connection accepted from: %s", c.RemoteAddr())

	return client, nil
}

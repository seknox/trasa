package sshproxy

import (
	"github.com/gorilla/websocket"
	"github.com/seknox/ssh"
	"github.com/sirupsen/logrus"
	"io"
)

type webSSHBackendConn struct {
	stdIn   io.Writer
	stdOut  io.Reader
	session *ssh.Session
}

func NewWebSSHBackend(session *ssh.Session) (*webSSHBackendConn, error) {
	stdInPipe, err := session.StdinPipe()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	stdOutPipe, err := session.StdoutPipe()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &webSSHBackendConn{
		stdIn:   stdInPipe,
		stdOut:  stdOutPipe,
		session: session,
	}, nil

}

func (wssb webSSHBackendConn) Read(p []byte) (int, error) {
	return wssb.stdOut.Read(p)
}
func (wssb *webSSHBackendConn) Write(p []byte) (int, error) {
	return wssb.stdIn.Write(p)
}
func (wssb *webSSHBackendConn) Close() error {
	return wssb.session.Close()
}

type webSSHFrontEndConn struct {
	wsConn *websocket.Conn
}

func NewWebSSHFrontEndConn(ws *websocket.Conn) *webSSHFrontEndConn {
	return &webSSHFrontEndConn{
		wsConn: ws,
	}
}

func (websshConn *webSSHFrontEndConn) Read(p []byte) (int, error) {
	_, byt, err := websshConn.wsConn.ReadMessage()
	if err != nil {
		return 0, err
	}
	copy(p, byt)
	return len(byt), nil
}

func (websshConn *webSSHFrontEndConn) Write(p []byte) (int, error) {
	err := websshConn.wsConn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func (websshConn *webSSHFrontEndConn) Close() error {
	return websshConn.wsConn.Close()
}

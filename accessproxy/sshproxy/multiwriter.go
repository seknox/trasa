package sshproxy

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime/debug"
)

//WrappedTunnel wraps upstream(backend) ssh connection and writes data to session file,guests
//It also writes data coming from guests to upstream(backend) ssh connection
type WrappedTunnel struct {
	io.WriteCloser
	io.Reader
	sessionRecord bool
	tempLogFile   *os.File
	guests        []*websocket.Conn
}

func NewWrappedTunnel(sessionID string, sessionRecord bool, backendReader io.Reader, backendWriter io.WriteCloser, guestChan chan GuestClient) (*WrappedTunnel, error) {

	err := os.MkdirAll("/tmp/trasa/trasagw/", 0644)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	tunn := &WrappedTunnel{
		WriteCloser:   backendWriter,
		Reader:        backendReader,
		sessionRecord: sessionRecord,
		guests:        nil,
	}

	if sessionRecord {
		tunn.tempLogFile, err = os.OpenFile("/tmp/trasa/trasagw/"+sessionID+".session", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

	}

	go tunn.ListenToNewGuests(guestChan)
	return tunn, nil
}

func (lr *WrappedTunnel) ListenToNewGuests(guestChan chan GuestClient) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error(r, string(debug.Stack()))
		}
	}()
	for v := range guestChan {
		//Append wesocket connection to list of viewers
		lr.guests = append(lr.guests, v.Conn)

		//TODO
		//Append guest email to  log
		//lr.proxyMeta.Log.Guests = append(lr.proxyMeta.Log.Guests, v.Email)
		go lr.readFromGuests(v.Conn)
		logrus.Debug("New viewer joined")
	}
}

func (lr *WrappedTunnel) readFromGuests(guest *websocket.Conn) {
	defer func() {
		r := recover()
		if r != nil {
			logrus.Error(r, string(debug.Stack()))
		}
	}()
	for {
		_, data, err := guest.ReadMessage()
		if err != nil {
			logrus.Debugf(`Could not read from viewer, disconnected: `, err)
			return
		}
		lr.WriteCloser.Write(data)
	}

}

func (lr *WrappedTunnel) Read(p []byte) (n int, err error) {
	n, err = lr.Reader.Read(p)

	//Write to file
	if lr.sessionRecord {
		lr.tempLogFile.WriteString(fmt.Sprintf("%s", string(p[:n])))
	}

	///Write to guest
	for _, guest := range lr.guests {
		if guest != nil {
			guest.WriteMessage(1, p[:n])
		}
	}

	return n, err
}

func (lr *WrappedTunnel) Close() error {
	lr.tempLogFile.Close()

	for _, v := range lr.guests {
		if v != nil {
			v.Close()
		}
	}
	return lr.WriteCloser.Close()

}

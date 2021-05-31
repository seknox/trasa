package sshproxy

import (
	"net"
	"sync"

	"github.com/seknox/ssh"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

//InitStore initialises package state
func InitStore(state *global.State) {
	SSHStore = Store{
		State:         state,
		sessions:      make(map[net.Addr]*Session),
		guestChannels: map[string]chan GuestClient{},
		lock:          &sync.Mutex{},
	}
}

var SSHStore Adapter

type Store struct {
	lock *sync.Mutex
	*global.State
	sessions      map[net.Addr]*Session
	guestChannels map[string]chan GuestClient
}

type Adapter interface {
	GetUserFromPublicKey(publicKey ssh.PublicKey, orgID string) (*models.User, error)
	validateTempCert(publicKey *ssh.Certificate, privilege string, orgID string) error
	parseSSHCert(addr net.Addr, publicKey ssh.PublicKey) error

	SetSession(addr net.Addr, session *Session) error
	GetSession(addr net.Addr) (*Session, error)
	DeleteSession(addr net.Addr) error
	UpdateSessionMeta(addr net.Addr, connMeta ssh.ConnMetadata) error
	//	UpdateSessionParams(addr net.Addr, params *models.AccessMapDetail) error
	UpdateSessionUser(addr net.Addr, user *models.User) error
	SetAuthType(addr net.Addr, authType consts.SSH_AUTH_TYPE) error

	CreateGuestChannel(sessionID string) chan GuestClient
	GetGuestChannel(sessionID string) (chan GuestClient, error)
	deleteGuestChannel(sessionID string)

	uploadSessionLog(authlog *logs.AuthLog) error
	closeSession(addr net.Addr)
}

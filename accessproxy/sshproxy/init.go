package sshproxy

import (
	"net"

	"github.com/seknox/ssh"
	"github.com/seknox/trasa/consts"
	"github.com/seknox/trasa/core/logs"
	"github.com/seknox/trasa/global"
	"github.com/seknox/trasa/models"
)

func InitStore(state *global.State, checkPolicyFunc models.CheckPolicyFunc) {
	SSHStore = Store{
		State:           state,
		checkPolicyFunc: checkPolicyFunc,
		sessions:        make(map[net.Addr]*Session),
		guestChannels:   map[string]chan GuestClient{},
	}
}

var SSHStore Adapter

type Store struct {
	*global.State
	sessions        map[net.Addr]*Session
	guestChannels   map[string]chan GuestClient
	checkPolicyFunc models.CheckPolicyFunc
}

type Adapter interface {
	checkPolicy(params *models.ConnectionParams) (*models.Policy, consts.FailedReason, error)
	getUserFromPublicKey(publicKey ssh.PublicKey, orgID string) (*models.User, error)
	validateTempCert(publicKey ssh.PublicKey, privilege string, orgID string) error
	tfaCert(publicKey ssh.PublicKey) (*models.AccessMapDetail, error)

	SetSession(addr net.Addr, session *Session) error
	GetSession(addr net.Addr) (*Session, error)
	DeleteSession(addr net.Addr) error
	UpdateSessionMeta(addr net.Addr, connMeta ssh.ConnMetadata) error
	UpdateSessionParams(addr net.Addr, params *models.AccessMapDetail) error
	UpdateSessionUser(addr net.Addr, user *models.User) error
	SetAuthType(addr net.Addr, authType consts.SSH_AUTH_TYPE) error

	CreateGuestChannel(sessionID string) chan GuestClient
	GetGuestChannel(sessionID string) (chan GuestClient, error)
	deleteGuestChannel(sessionID string)

	uploadSessionLog(authlog *logs.AuthLog) error
	closeSession(addr net.Addr)
}

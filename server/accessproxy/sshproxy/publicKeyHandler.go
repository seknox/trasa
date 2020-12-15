package sshproxy

import (
	"github.com/pkg/errors"
	"github.com/seknox/ssh"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/sirupsen/logrus"
)

func publicKeyCallbackHandler(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {

	session, err := SSHStore.GetSession(conn.RemoteAddr())
	if err != nil {
		logrus.Trace(err)
		return nil, errors.New(failNow)
	}
	session.UpdateConMeta(conn)
	SSHStore.UpdateSessionMeta(conn.RemoteAddr(), conn)

	publicKey, ok := key.(ssh.PublicKey)
	if !ok {
		return nil, errors.New(failNow)
	}
	user, err := SSHStore.GetUserFromPublicKey(publicKey, global.GetConfig().Trasa.OrgId)
	if err != nil {
		logrus.Debug(err)
		//Try again with next key
		return nil, errors.New(gotoPublicKey)
	}

	SSHStore.SetAuthType(conn.RemoteAddr(), consts.SSH_AUTH_TYPE_PUB)

	err = SSHStore.UpdateSessionUser(conn.RemoteAddr(), user)
	if err != nil {
		logrus.Errorf("update session user: %v", err)
		return nil, errors.New(failNow)
	}

	err = SSHStore.validateTempCert(key, conn.User(), global.GetConfig().Trasa.OrgId)
	if err != nil {
		logrus.Trace(err)
		return nil, errors.New(gotoPublicKey)
	}

	//parse and validate connection deviceID embedded in ssh certificate
	err = SSHStore.parseSSHCert(conn.RemoteAddr(), key)
	if err != nil {
		return nil, errors.New(gotoPublicKey)
	}

	return nil, errors.New(gotoKeyboardInteractive)
}

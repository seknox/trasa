package sshproxy

import (
	"github.com/pkg/errors"
	"github.com/seknox/ssh"
	"github.com/seknox/trasa/consts"
	"github.com/seknox/trasa/global"
	"github.com/sirupsen/logrus"
)

func publicKeyCallbackHandler(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {

	session, err := SSHStore.GetSession(conn.RemoteAddr())
	if err != nil {
		logrus.Trace(err)
		return nil, errors.New(fail_now)
	}
	session.UpdateConMeta(conn)
	SSHStore.UpdateSessionMeta(conn.RemoteAddr(), conn)
	//TODO check enforce cert flag

	//	if resp.EnforceCert {
	//		return nil, errors.New(fail_now)
	//	}

	//If its just a notmarlk certificate without tfa
	user, err := SSHStore.getUserFromPublicKey(key, global.GetConfig().Trasa.OrgId)
	if err != nil {
		logrus.Debug(err)
		//Try again with next key
		return nil, errors.New(goto_public_key)
	}

	SSHStore.SetAuthType(conn.RemoteAddr(), consts.SSH_AUTH_TYPE_PUB)

	err = SSHStore.UpdateSessionUser(conn.RemoteAddr(), user)
	if err != nil {
		logrus.Errorf("update session user: %v", err)
		return nil, errors.New(fail_now)
	}

	err = SSHStore.validateTempCert(key, conn.User(), global.GetConfig().Trasa.OrgId)
	if err != nil {
		logrus.Trace(err)
		return nil, errors.New(goto_public_key)
	}

	//tfa already done from device agent
	//parse and validate connection params embedded in ssh certificate
	params, err := SSHStore.tfaCert(key)
	if err != nil {
		return nil, errors.New(goto_keyboard_interactive)
	}

	SSHStore.SetAuthType(conn.RemoteAddr(), consts.SSH_AUTH_TYPE_DACERT)
	err = SSHStore.UpdateSessionParams(conn.RemoteAddr(), params)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New(fail_now)
	}

	return nil, errors.New(goto_keyboard_interactive)
}

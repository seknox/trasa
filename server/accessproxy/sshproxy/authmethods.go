package sshproxy

import (
	"github.com/seknox/ssh"
)

const (
	gotoKeyboardInteractive = "trasa: goto_keyboard_interactive"
	gotoPublicKey           = "trasa: goto_public_key"
	//gotoPublicKeyOrKeyboardInteractive           = "trasa: goto_public_key_or_kb_interactive"
	failNow = "trasa: fail_now"
)

//Decides which auth method to use next from previous error
func nextAuthMethodHandler(conn ssh.ConnMetadata, prevErr error) ([]string, bool, error) {

	switch prevErr.Error() {
	case gotoKeyboardInteractive:
		return []string{"keyboard-interactive"}, true, nil
	case "ssh: no auth passed yet":
		return []string{"publickey", "keyboard-interactive"}, false, nil
	case gotoPublicKey:
		return []string{"publickey", "keyboard-interactive"}, false, nil
	case failNow:
		return []string{}, false, prevErr
	default:
		if prevErr != nil {
			return []string{"keyboard-interactive"}, false, nil
		} else {

			//TODO check this
			//if sessions[conn.RemoteAddr()]["authMeta"] == nil {
			//	return []string{"publickey"}, false, nil
			//}
			return []string{"keyboard-interactive"}, false, nil
		}

	}
}

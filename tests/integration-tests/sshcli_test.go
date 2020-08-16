package integration_tests

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/seknox/trasa/server/accessproxy/sshproxy"
	"github.com/seknox/trasa/server/utils"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"testing"
	"time"
)

const (
	trasaEmail   = "root"
	trasaPass    = "changeme"
	totpSEC      = "AV2COXZHVG4OAFSF"
	upstreamUser = "bhrg3se"
	upstreamPass = "downwiththe"
)

// tryAuthBothSides runs the handshake and returns the resulting errors from both sides of the connection.
func tryAuthBothSides(t *testing.T, config *ssh.ClientConfig) (clientError error, serverAuthErrors []error) {
	c1, c2, err := netPipe()
	if err != nil {
		t.Fatalf("netPipe: %v", err)
	}
	defer c1.Close()
	defer c2.Close()

	certChecker := ssh.CertChecker{
		IsUserAuthority: func(k ssh.PublicKey) bool {
			return bytes.Equal(k.Marshal(), testPublicKeys["ecdsa"].Marshal())
		},
		UserKeyFallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if conn.User() == "testuser" && bytes.Equal(key.Marshal(), testPublicKeys["rsa"].Marshal()) {
				return nil, nil
			}

			return nil, fmt.Errorf("pubkey for %q not acceptable", conn.User())
		},
		IsRevoked: func(c *ssh.Certificate) bool {
			return c.Serial == 666
		},
	}

	serverConfig := &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if conn.User() == "testuser" && string(pass) == clientPassword {
				return nil, nil
			}
			return nil, errors.New("password auth failed")
		},
		PublicKeyCallback: certChecker.Authenticate,
		KeyboardInteractiveCallback: func(conn ssh.ConnMetadata, challenge ssh.KeyboardInteractiveChallenge) (*ssh.Permissions, error) {
			ans, err := challenge("user",
				"instruction",
				[]string{"question1", "question2"},
				[]bool{true, true})
			if err != nil {
				return nil, err
			}
			ok := conn.User() == "testuser" && ans[0] == "answer1" && ans[1] == "answer2"
			if ok {
				challenge("user", "motd", nil, nil)
				return nil, nil
			}
			return nil, errors.New("keyboard-interactive failed")
		},
	}
	serverConfig.AddHostKey(testSigners["rsa"])

	serverConfig.AuthLogCallback = func(conn ssh.ConnMetadata, method string, err error) {
		serverAuthErrors = append(serverAuthErrors, err)
	}

	go newServer(c1, serverConfig)
	_, _, _, err = ssh.NewClientConn(c2, "", config)
	return err, serverAuthErrors
}

func TestSSHAuthWithoutPublicKey(t *testing.T) {
	done := make(chan bool, 1)
	go sshproxy.ListenSSH(done)

	time.Sleep(time.Second * 2)

	cconf := ssh.ClientConfig{
		User: upstreamUser,
		Auth: []ssh.AuthMethod{
			handleKBAuth(t),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		BannerCallback: nil,
	}

	client, err := ssh.Dial("tcp", "127.0.0.1:8022", &cconf)
	if err != nil {
		t.Fatal(err)
	}

	s, err := client.NewSession()
	if err != nil {
		t.Fatalf(`could not start session: %v`, err)
	}

	//s.Close()
	//t.Log("closed++++++++++++++++++++++++++++++++++++++")

	err = s.Run("ls")
	if err != nil {
		t.Fatalf(`could not run command: %v`, err)
	}

	done <- true
}

func handleKBAuth(t *testing.T) ssh.AuthMethod {
	return ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {

		switch true {
		case strings.Contains(instruction, "Enter TRASA credentials"):
			if len(questions) != 2 {
				t.Fatalf(`incorrect number of question, want: %d got: %d`, 2, len(questions))
			}
			if !strings.Contains(questions[0], "Email") {
				t.Fatalf(`incorrect  question, want: %s got: %s`, "Enter Email (TRASA)", questions[0])
			}
			t.Log("Enter TRASA credentials")

			return []string{trasaEmail, trasaPass}, nil

		case strings.Contains(instruction, "Choose Service"):
			if len(questions) != 1 {
				t.Fatalf(`incorrect number of question, want: %d got: %d`, 1, len(questions))
			}
			t.Log("Choose service")

			return []string{`127.0.0.1:2222`}, nil

		case strings.Contains(instruction, "Second factor authentication"):
			if len(questions) != 1 {
				t.Fatalf(`incorrect number of question, want: %d got: %d`, 1, len(questions))
			}
			_, totp, _ := utils.CalculateTotp(totpSEC)
			t.Log("Second factor authentication " + totp)
			return []string{totp}, nil

		case strings.Contains(instruction, "Host key verify"):
			if len(questions) != 1 {
				t.Fatalf(`incorrect number of question, want: %d got: %d`, 1, len(questions))
			}
			return []string{"yes"}, nil

		case strings.Contains(instruction, "Upstream password"):
			if len(questions) != 1 {
				t.Fatalf(`incorrect number of question, want: %d got: %d`, 1, len(questions))
			}
			return []string{upstreamPass}, nil

		default:
			if len(questions) != 0 {
				t.Fatalf(`incorrect number of question, %v want: %d got: %d`, questions, 0, len(questions))
			}
			return nil, nil
		}

	})
}

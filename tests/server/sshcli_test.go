package server_test

import (
	"github.com/seknox/trasa/server/accessproxy/sshproxy"
	"github.com/seknox/trasa/server/utils"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"testing"
	"time"
)

//TODO make tabular tests

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

func TestSSHAuthWithPublicKey(t *testing.T) {
	done := make(chan bool, 1)
	go sshproxy.ListenSSH(done)

	time.Sleep(time.Second * 2)

	pk, err := ssh.ParsePrivateKey([]byte(testPrivateKey))
	if err != nil {
		t.Fatal(err)
	}

	cconf := ssh.ClientConfig{
		User: upstreamUser,
		Auth: []ssh.AuthMethod{
			handleKBAuth(t),
			ssh.PublicKeys(pk),
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

func TestSSHAuthWithAuthorisedPublicKey(t *testing.T) {
	done := make(chan bool, 1)
	go sshproxy.ListenSSH(done)

	time.Sleep(time.Second * 2)

	pk, err := ssh.ParsePrivateKey([]byte(testPrivateKey2))
	if err != nil {
		t.Fatal(err)
	}

	cconf := ssh.ClientConfig{
		User: upstreamUser,
		Auth: []ssh.AuthMethod{
			handleKBAuth(t),
			ssh.PublicKeys(pk),
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
			//t.Log("Enter TRASA credentials")

			return []string{trasaEmail, trasaPass}, nil

		case strings.Contains(instruction, "Choose Service"):
			if len(questions) != 1 {
				t.Fatalf(`incorrect number of question, want: %d got: %d`, 1, len(questions))
			}
			//t.Log("Choose service")

			return []string{`127.0.0.1:2222`}, nil

		case strings.Contains(instruction, "Second factor authentication"):
			if len(questions) != 1 {
				t.Fatalf(`incorrect number of question, want: %d got: %d`, 1, len(questions))
			}
			_, totp, _ := utils.CalculateTotp(totpSEC)
			//t.Log("Second factor authentication " + totp)
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

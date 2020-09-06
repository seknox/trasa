package accessproxytest

import (
	"github.com/seknox/ssh"
	"github.com/seknox/trasa/server/accessproxy/sshproxy"
	"github.com/seknox/trasa/server/api/my"
	"github.com/seknox/trasa/server/utils"
	"github.com/seknox/trasa/tests/server/testutils"
	"net"
	"net/http"
	"net/http/httptest"
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
		User: testutils.MockupstreamUser,
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

	pk, err := ssh.ParsePrivateKey([]byte(testutils.MockPrivateKey))
	if err != nil {
		t.Fatal(err)
	}

	cconf := ssh.ClientConfig{
		User: testutils.MockupstreamUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(pk),
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

func TestSSHAuthWithAuthorisedPublicKey(t *testing.T) {
	done := make(chan bool, 1)
	go sshproxy.ListenSSH(done)

	time.Sleep(time.Second * 2)

	key := downloadKey(t)

	pk, err := ssh.ParsePrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}

	cconf := ssh.ClientConfig{
		User: testutils.MockupstreamUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(pk),
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
			//t.Log("Enter TRASA credentials")

			return []string{testutils.MockTrasaID, testutils.MocktrasaPass}, nil

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
			_, totp, _ := utils.CalculateTotp(testutils.MocktotpSEC)
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
			return []string{testutils.MockupstreamPass}, nil

		default:
			if len(questions) != 0 {
				t.Fatalf(`incorrect number of question, %v want: %d got: %d`, questions, 0, len(questions))
			}
			t.Logf("default %s  %v", instruction, questions)

			return nil, nil
		}

	})
}

func downloadKey(t *testing.T) []byte {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testutils.AddTestUserContext(my.GenerateKeyPair))

	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	k, err := ssh.ParsePrivateKey(rr.Body.Bytes())
	if err != nil {
		t.Errorf(`invalid user key`)
	}

	user, err := sshproxy.SSHStore.GetUserFromPublicKey(k.PublicKey(), testutils.MockOrgID)
	if err != nil {
		t.Errorf(`incorrect user key`)
	}

	if user.ID != testutils.MockUserID {
		t.Errorf(`incorrect user ID, want=%v got=%v`, testutils.MockUserID, user.ID)
	}
	return rr.Body.Bytes()

}

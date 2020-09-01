package accessproxytest

import (
	"github.com/google/goexpect"
	"github.com/seknox/trasa/server/accessproxy/sshproxy"
	"github.com/seknox/trasa/tests/server/testutils"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestSSHCLI(t *testing.T) {
	done := make(chan bool, 1)
	go sshproxy.ListenSSH(done)

	time.Sleep(time.Second * 2)

	//c:=exec.Command("ssh",testutils.MockupstreamUser+"@localhost","-p","8022")
	c := exec.Command("/usr/local/bin/trasacli", "-u", testutils.MockupstreamUser, "ssh")
	in, err := c.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}

	out, err := c.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	c.Stdout = os.Stdout

	resCh := make(chan error)

	err = c.Start()
	if err != nil {
		t.Fatal(err)
	}
	e, _, err := expect.SpawnGeneric(&expect.GenOptions{
		In:  in,
		Out: out,
		Wait: func() error {
			return <-resCh
		},
		Close: func() error {
			close(resCh)
			return nil
		},
		Check: func() bool { return true },
	}, time.Second)

	//e, _, err := expect.Spawn(fmt.Sprintf("/usr/local/bin/trasacli -u %s ssh",testutils.MockupstreamUser), -1)
	if err != nil {
		t.Fatal(err)
	}
	defer e.Close()

	//res,err:=e.ExpectBatch([]expect.Batcher{
	//	&expect.BExp{R: `Enter TRASA URL:`},
	//	&expect.BSnd{S: "https://localhost" + "\n"},
	//
	//
	//
	//}, time.Second*5)

	//TODO fix the regex

	res, err := e.ExpectBatch([]expect.Batcher{
		&expect.BExp{R: `Enter TRASA URL:`},
		&expect.BSnd{S: "https://localhost" + "\n"},
		&expect.BExp{R: `Enter TRASA email/username:`},
		&expect.BSnd{S: testutils.MockTrasaID + "\n"},
		&expect.BExp{R: "Enter TRASA password:"},
		&expect.BSnd{S: testutils.MocktrasaPass + "\n"},
		&expect.BExp{R: "Enter Service IP :"},
		&expect.BSnd{S: "127.0.0.1:2222" + "\n"},
		&expect.BExp{R: `Enter OTP(Blank for U2F):`},
		&expect.BSnd{S: testutils.GetTotpCode(testutils.MocktotpSEC) + "\n"},
		&expect.BExp{R: `Enter Password(Upstream Server):`},
		&expect.BSnd{S: testutils.GetTotpCode(testutils.MocktrasaPass) + "\n"},
		&expect.BExp{R: `#`},
		&expect.BSnd{S: `echo "Hello"\n`},
		&expect.BExp{R: `Hello`},
	}, time.Second*5)

	if err != nil {
		t.Fatalf(`%v  %v`, res, err)
	}

	//err = s.Run("ls")
	//if err != nil {
	//	t.Fatalf(`could not run command: %v`, err)
	//}

	done <- true
}

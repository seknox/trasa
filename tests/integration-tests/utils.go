package integration_tests

import (
	"crypto/rand"
	"fmt"
	"github.com/seknox/trasa/server/utils"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/testdata"
	"net"
)

var clientPassword = "tiger"

// netPipe is analogous to net.Pipe, but it uses a real net.Conn, and
// therefore is buffered (net.Pipe deadlocks if both sides start with
// a write.)
func netPipe() (net.Conn, net.Conn, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		listener, err = net.Listen("tcp", "[::1]:0")
		if err != nil {
			return nil, nil, err
		}
	}
	defer listener.Close()
	c1, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		return nil, nil, err
	}

	c2, err := listener.Accept()
	if err != nil {
		c1.Close()
		return nil, nil, err
	}

	return c1, c2, nil
}

var (
	testPrivateKeys map[string]interface{}
	testSigners     map[string]ssh.Signer
	testPublicKeys  map[string]ssh.PublicKey
)

func init() {
	var err error

	n := len(testdata.PEMBytes)
	testPrivateKeys = make(map[string]interface{}, n)
	testSigners = make(map[string]ssh.Signer, n)
	testPublicKeys = make(map[string]ssh.PublicKey, n)
	for t, k := range testdata.PEMBytes {
		testPrivateKeys[t], err = ssh.ParseRawPrivateKey(k)
		if err != nil {
			panic(fmt.Sprintf("Unable to parse test key %s: %v", t, err))
		}
		testSigners[t], err = ssh.NewSignerFromKey(testPrivateKeys[t])
		if err != nil {
			panic(fmt.Sprintf("Unable to create signer for test key %s: %v", t, err))
		}
		testPublicKeys[t] = testSigners[t].PublicKey()
	}

	// Create a cert and sign it for use in tests.
	testCert := &ssh.Certificate{
		Nonce:           []byte{},                       // To pass reflect.DeepEqual after marshal & parse, this must be non-nil
		ValidPrincipals: []string{"gopher1", "gopher2"}, // increases test coverage
		ValidAfter:      0,                              // unix epoch
		ValidBefore:     ssh.CertTimeInfinity,           // The end of currently representable time.
		Reserved:        []byte{},                       // To pass reflect.DeepEqual after marshal & parse, this must be non-nil
		Key:             testPublicKeys["ecdsa"],
		SignatureKey:    testPublicKeys["rsa"],
		Permissions: ssh.Permissions{
			CriticalOptions: map[string]string{},
			Extensions:      map[string]string{},
		},
	}
	testCert.SignCert(rand.Reader, testSigners["rsa"])
	testPrivateKeys["cert"] = testPrivateKeys["ecdsa"]
	testSigners["cert"], err = ssh.NewCertSigner(testCert, testSigners["ecdsa"])
	if err != nil {
		panic(fmt.Sprintf("Unable to create certificate signer: %v", err))
	}
}

type server struct {
	*ssh.ServerConn
	chans <-chan ssh.NewChannel
}

func newServer(c net.Conn, conf *ssh.ServerConfig) (*server, error) {

	sconn, chans, reqs, err := ssh.NewServerConn(c, conf)
	if err != nil {
		return nil, err
	}
	go ssh.DiscardRequests(reqs)
	return &server{sconn, chans}, nil
}

func getTotpCode(secret string) string {
	_, t, _ := utils.CalculateTotp(secret)
	return t
}

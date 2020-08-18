package server_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/testdata"
	"net"
	"net/http"
)

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

// AddTestUserContext is a middleware that adds  mock userContext
func AddTestUserContext(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userContext := models.UserContext{
			User: &models.User{
				ID:         "13c45cfb-72ca-4177-b968-03604cab6a27",
				OrgID:      "153f7582-5ae2-46ba-8c1c-79ef73fe296e",
				UserName:   "root",
				FirstName:  "Bhargab",
				MiddleName: "",
				LastName:   "Acharya",
				Email:      "bhargab@seknox.com",
				Groups:     nil,
				UserRole:   "orgAdmin",
				Status:     true,
				IdpName:    "trasa",
			},
			Org: models.Org{
				ID:             "153f7582-5ae2-46ba-8c1c-79ef73fe296e",
				OrgName:        "Trasa",
				Domain:         "trasa.io",
				PrimaryContact: "",
				Timezone:       "Asia/Kathmandu",
				PhoneNumber:    "",
				CreatedAt:      0,
				PlatformBase:   "",
				License:        models.License{},
			},
			DeviceID:  "",
			BrowserID: "",
		}
		ctx := context.WithValue(r.Context(), "user", userContext)
		next(w, r.WithContext(ctx))

	})

}

// AddTestUserContextWS is a middleware that adds  mock userContext to ws handlers
func AddTestUserContextWS(next func(params models.ConnectionParams, uc models.UserContext, ws *websocket.Conn)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logrus.Error(err)
			return
		}
		//defer conn.Close()

		//TODO use different generic model for session validation
		var params models.ConnectionParams
		err = conn.ReadJSON(&params)
		if err != nil {
			logrus.Error(err)
			conn.WriteMessage(1, []byte(err.Error()))
			conn.Close()
			return
		}

		userContext := models.UserContext{
			User: &models.User{
				ID:         "13c45cfb-72ca-4177-b968-03604cab6a27",
				OrgID:      "153f7582-5ae2-46ba-8c1c-79ef73fe296e",
				UserName:   "root",
				FirstName:  "Bhargab",
				MiddleName: "",
				LastName:   "Acharya",
				Email:      "bhargab@seknox.com",
				Groups:     nil,
				UserRole:   "orgAdmin",
				Status:     true,
				IdpName:    "trasa",
			},
			Org: models.Org{
				ID:             "153f7582-5ae2-46ba-8c1c-79ef73fe296e",
				OrgName:        "Trasa",
				Domain:         "trasa.io",
				PrimaryContact: "",
				Timezone:       "Asia/Kathmandu",
				PhoneNumber:    "",
				CreatedAt:      0,
				PlatformBase:   "",
				License:        models.License{},
			},
			DeviceID:  "",
			BrowserID: "",
		}

		params.UserIP = utils.GetIp(r)

		next(params, userContext, conn)

	})

}

package main

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/seknox/trasa/server/api/auth/serviceauth"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/global"
	"github.com/sirupsen/logrus"
	"layeh.com/radius"
)

func StartRadiusServer(done chan bool) {
	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(serviceauth.RadiusLogin),
		SecretSource: DynamicSecretSource(),
	}
	go func() {
		<-done
		server.Shutdown(context.Background())
	}()

	logrus.Info("Radius server started on port 1812/udp")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}

// SecretSource is our custom secret source
type SecretSource interface {
	RADIUSSecret(ctx context.Context, remoteAddr net.Addr) ([]byte, error)
}

// DynamicSecretSource returns a secret by quering service db
func DynamicSecretSource() SecretSource {
	return &dynamicSecretSource{}
}

type dynamicSecretSource struct {
	secret []byte
}

func (s *dynamicSecretSource) RADIUSSecret(ctx context.Context, remoteAddr net.Addr) ([]byte, error) {
	logrus.Trace("request received from: ", remoteAddr.String())

	ipAddr := strings.Split(remoteAddr.String(), ":")

	if len(ipAddr) == 2 {
		srv, err := services.Store.GetFromHostname(ipAddr[0], "radius", "", global.GetConfig().Trasa.OrgId)
		if err != nil {
			logrus.Error(err)
			return []byte(""), fmt.Errorf("invalid remote host")
		}

		return []byte(srv.SecretKey), nil
	}

	//return s.secret, nil
	return []byte(""), nil
}

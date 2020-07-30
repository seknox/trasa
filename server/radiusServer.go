package server

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/seknox/trasa/core/auth/serviceauth"
	"github.com/seknox/trasa/core/services"
	"github.com/seknox/trasa/global"
	"github.com/sirupsen/logrus"
	"layeh.com/radius"
)

func startRadiusServer() {
	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(serviceauth.HandleRadiusReq),
		SecretSource: DynamicSecretSource(),
	}

	logrus.Debug("Starting radius server on UDP :1812")
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

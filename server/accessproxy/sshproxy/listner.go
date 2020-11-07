package sshproxy

import (
	"github.com/pkg/errors"
	"github.com/seknox/ssh"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"path/filepath"
)

func ListenSSH(closeChan chan bool) error {

	//defer func() {
	//	if r := recover(); r != nil {
	//		logrus.Errorf("Recovered in ListenSSH(),: %v : %s", r, string(debug.Stack()))
	//	}
	//}()

	privateBytes, err := ioutil.ReadFile(filepath.Join(utils.GetETCDir(), "trasa", "certs", "id_rsa"))
	if err != nil {
		pkey, err := utils.GeneratePrivateKey(4082)
		if err != nil {
			return errors.WithMessage(err, "Failed to load private key")
		}
		privateBytes = utils.EncodePrivateKeyToPEM(pkey)
	}

	//TODO IMPORTANT key pass
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return errors.WithMessage(err, "Failed to parse private key")
	}

	/////////// SERVER CONFIG 	///////////////
	config := &ssh.ServerConfig{
		BannerCallback:              bannerCallBackHandler,
		KeyboardInteractiveCallback: keyboardInteractiveHandler,
		PublicKeyCallback:           publicKeyCallbackHandler,
		NextAuthMethodsCallback:     nextAuthMethodHandler,
		MaxAuthTries:                -1,
	}
	config.SetDefaults()
	config.Ciphers = append(config.Ciphers, "aes128-cbc", "blowfish-cbc", "3des-cbc")

	config.AddHostKey(private)
	logrus.Info("TRASA SSH access proxy started")

	listenAddr := global.GetConfig().Proxy.SSHListenAddr
	if listenAddr == "" {
		listenAddr = ":8022"
	}
	sshListener, err := net.Listen("tcp", listenAddr)
	if err != nil {

		return errors.Errorf("net.Listen failed: %v", err)
	}

	defer sshListener.Close()

	var done = false
	go func() {
		<-closeChan
		done = true
	}()
	for !done {
		conn, err := sshListener.Accept()
		if err != nil {
			logrus.Errorf("listen.Accept failed: %v", err)
			continue
			//return err
		}

		session := NewSession(&conn)
		err = SSHStore.SetSession(conn.RemoteAddr(), session)
		if err != nil {
			logrus.Error(err)
			continue
		}
		go start(conn, config)
	}

	return nil
}

func bannerCallBackHandler(conn ssh.ConnMetadata) string {
	return "Welcome to TRASA SSH access proxy\n\rIF YOU ARE NOT AUTHORISED LEAVE IMMEDIATELY!\n\r\n\r" +
		"Copyright @ 2020, Seknox Cybersecurity, All rights reserved.\n\r\n\r "
}

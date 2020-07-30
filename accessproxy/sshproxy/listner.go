package sshproxy

import (
	"io/ioutil"
	"net"
	"runtime/debug"

	"github.com/seknox/ssh"
	"github.com/seknox/trasa/global"
	"github.com/sirupsen/logrus"
)

func ListenSSH() {

	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in ListenSSH(),: %v : %s", r, string(debug.Stack()))
		}
	}()

	privateBytes, err := ioutil.ReadFile("/etc/trasa/certs/id_rsa")

	if err != nil {
		panic("Failed to load private key" + err.Error())
	}

	//TODO IMPORTANT key pass
	private, err := ssh.ParsePrivateKeyWithPassphrase(privateBytes, []byte("testpass"))
	if err != nil {
		panic("Failed to parse private key")
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
	logrus.Debug("TRASA gateway started")

	listenAddr := global.GetConfig().SSHProxy.ListenAddr
	sshListener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logrus.Debugf("net.Listen failed: %v", err)
		return
	}

	defer sshListener.Close()

	for {
		conn, err := sshListener.Accept()
		if err != nil {
			logrus.Error("listen.Accept failed: %v", err)
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

}

func bannerCallBackHandler(conn ssh.ConnMetadata) string {
	return "Welcome to TRASA ssh Gateway\n\rIF YOU ARE NOT AUTHORISED LEAVE IMMEDIATELY!\n\r\n\r" +
		"Copyright @ 2020, Seknox Cybersecurity, All rights reserved.\n\r\n\r "
}

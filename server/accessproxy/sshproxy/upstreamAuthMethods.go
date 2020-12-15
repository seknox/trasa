package sshproxy

import (
	"errors"
	"github.com/seknox/ssh"
	"github.com/seknox/trasa/server/models"
	"github.com/sirupsen/logrus"
	"strings"
)

//handleUpstreamTrasaPAM returns keyboard-interactive auth method to handle TRASA PAM module installed in upstream server
func handleUpstreamTrasaPAM(trasaID, password, totp string) ssh.AuthMethod {
	return ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
		//logrus.Tracef("user: %s \n instruction: %v \n questions: %v",user,instruction,questions)
		answers := make([]string, len(questions))
		if len(questions) == 1 {

			if strings.Contains(questions[0], "Password") {
				answers[0] = password
				//logrus.Tracef("pass: %v ",upstreamPassword)
				//logrus.Tracef("ans: %v ",answers)
				return answers, nil
			} else if strings.Contains(questions[0], "email") {
				answers[0] = trasaID
				return answers, nil
			} else if strings.Contains(questions[0], "totp") {
				answers[0] = totp
				return answers, nil
			} else {
				return answers, errors.New("unexpected challenges")
			}

		}

		return answers, nil
	})
}

func getPublicKeyAuthMethod(creds *models.UpstreamCreds) ssh.AuthMethod {
	publicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(creds.ClientCert))
	if err != nil && creds.ClientCert != "" {
		logrus.Error(err)
	}

	privateKey, privateKeyParseErr := ssh.ParsePrivateKey([]byte(creds.ClientKey))
	if privateKeyParseErr != nil && creds.ClientKey != "" {
		//logrus.Debug(creds.ClientKey)
		logrus.Error(privateKeyParseErr)
	}

	var signer ssh.Signer
	cert, ok := publicKey.(*ssh.Certificate)
	if !ok {
		logrus.Debug("\n\rInvalid user certificate\n\r")
		if privateKeyParseErr == nil {
			signer = privateKey
		}

	} else {
		signer, err = ssh.NewCertSigner(cert, privateKey)
		if err != nil {
			logrus.Debug(err)
		}
	}

	if signer != nil {
		return ssh.PublicKeys(signer)
	}
	return nil
}

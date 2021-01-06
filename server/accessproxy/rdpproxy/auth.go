package rdpproxy

import (
	"io"
	"strings"

	"github.com/seknox/trasa/server/api/auth/tfa"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/seknox/guacamole"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/models"
	"github.com/sirupsen/logrus"
	"github.com/trustelem/zxcvbn"
)

func makeConfig(params *models.ConnectionParams, creds *models.UpstreamCreds) (*guacamole.Config, error) {

	config := guacamole.NewGuacamoleConfiguration()

	config.Protocol = "rdp"
	config.OptimalScreenHeight = int(params.OptHeight)
	config.OptimalScreenWidth = int(params.OptWidth)

	config.Parameters = make(map[string]string)

	port := "3389"
	host := params.Hostname
	if strings.Contains(params.Hostname, ":") {
		splitted := strings.Split(params.Hostname, ":")
		port = splitted[1]
		host = splitted[0]

	}
	config.Parameters["hostname"] = host
	config.Parameters["port"] = port
	splitted := strings.Split(params.Privilege, `\`)
	if strings.Contains(params.Privilege, `\`) && len(splitted) == 2 {
		config.Parameters["username"] = splitted[1]
		config.Parameters["domain"] = splitted[0]
	} else {
		config.Parameters["username"] = params.Privilege
	}

	config.Parameters["password"] = creds.Password

	if params.SessionRecord {
		config.Parameters["recording-path"] = "/tmp/trasa/accessproxy/guac"
		config.Parameters["create-recording-path"] = "true"
		config.Parameters["recording-name"] = params.SessionID + ".guac"

	}

	if params.CanTransferFile {
		config.Parameters["enable-drive"] = "true"
		config.Parameters["create-drive-path"] = "true"
		config.Parameters["drive-path"] = "/tmp/trasa/accessproxy/guac/shared/" + params.UserID
		config.Parameters["drive-name"] = "TRASA shared drive"
	}

	config.Parameters["client-name"] = "TRASA"
	config.Parameters["security"] = params.RdpProtocol
	config.Parameters["ignore-cert"] = "true"

	//logrus.Debug(config.Parameters)
	return config, nil
}

func handlePass(params *models.ConnectionParams) (*models.UpstreamCreds, error) {
	creds, err := services.GetUpstreamCreds(params.Privilege, params.ServiceID, params.ServiceType, params.OrgID)
	if err != nil {
		return nil, err
	}

	password := ""
	if params.Password == "" {
		password = creds.Password
	} else {
		password = params.Password
		creds.Password = params.Password
	}

	passwordStrength := zxcvbn.PasswordStrength(password, nil)

	if creds.EnforceStrongPass && !params.IsSharedSession {
		if passwordStrength.Score < creds.ZxcvbnScore || len(password) < creds.MinimumChar {
			logrus.Debug("weak password")
			return nil, errors.New("weak password")
		}

	}

	return creds, err

}

func (s Session) handleTfa(ws *websocket.Conn, guacdWriter io.Writer) (ok bool, tfaDeviceID string) {

	logrus.Trace("sending tfa instruction")

	var tfaInstruction *guacamole.Instruction
	if s.params.Skip2FA {
		//we need to send tfa instruction with skip argument to tell guacamole client (js) to skip TFA and start capturing input
		tfaInstruction = guacamole.NewInstruction("tfa", "skip")
	} else {
		tfaInstruction = guacamole.NewInstruction("tfa")
	}

	err := ws.WriteMessage(websocket.TextMessage, tfaInstruction.Byte())
	if err != nil {
		logrus.Error(err)
		return false, ""
	}
	//logrus.Debug("tfa instruction  sentt")

	totp := ""

	var pendingInstructions []guacamole.Instruction

	//Wait for tfa response
	for i := 0; i < 100; i++ {
		_, b, err := ws.ReadMessage()
		if err != nil {
			logrus.Error(err)
			return false, ""
		}

		inst, err := guacamole.Parse(b)
		if err != nil {
			logrus.Error(err)
			return false, ""
		}

		//	logrus.Debug(inst.String())

		if inst.Opcode == guacamole.TfaOpcode {
			if len(inst.Args) > 0 {
				totp = inst.Args[0]
			}

			totp = inst.Args[0]
			break
		} else {
			pendingInstructions = append(pendingInstructions, *inst)
		}

	}

	//TODO handle u2fy
	s.params.TotpCode = totp
	if totp == "" {
		s.params.TfaMethod = "u2f"
	} else {
		s.params.TfaMethod = "totp"
	}

	if !s.params.Skip2FA {
		deviceID, reason, ok := tfa.HandleTfaAndGetDeviceID(nil,
			s.params.TfaMethod,
			totp,
			s.params.UserID,
			s.params.UserIP,
			s.params.ServiceName,
			s.params.Timezone, s.params.OrgName, s.params.OrgID)

		if !ok {
			s.log.FailedReason = reason
			s.log.TfaDeviceID = deviceID
			s.log.Status = false
			return false, deviceID
		}

	}

	//forward pending instructions
	for _, i := range pendingInstructions {
		_, err := guacdWriter.Write(i.Byte())
		if err != nil {
			logrus.Error(err)
		}

	}

	return true, tfaDeviceID

}

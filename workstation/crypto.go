package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

// keyHolder holds secret key for encryption.
// It is a map of type [purpose]secretKey
var keyHolder = make(map[string][]byte)

type kexRequest struct {
	// Intent of kex, example - KEX_EXPORT_DH, KEX_ENROL_DEVICE, KEX_HTTP_SR
	Intent string `json:"intent"`
	// IntentID is unique identifier for key exchange process. Should be trasaID in KEX_ENROL_DEVICE and sessionID in case of other intents
	IntentID string `json:"intentID"`
	// Unique id of device which is making kex request. Machine ID can be used here?
	DeviceID string `json:"deviceID"`
	// Public key of client
	PublicKey string `json:"publicKey"`
}

type trasaResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

// kex handles key exchange with trasa server.
func kex(userID, deviceID, intent, hostName string) error {

	// gen keypair for trasaWrkstnAgent.
	priv, pub, err := utils.ECDHGenKeyPair()
	if err != nil {
		logrus.Error(err)
	}

	// send our pub key to server
	client := getHTTPClient(true)
	url := fmt.Sprintf("%s/auth/crypto/kex", hostName)

	var req kexRequest
	req.Intent = intent
	req.PublicKey = hex.EncodeToString(pub[:])
	req.DeviceID = deviceID
	req.IntentID = userID

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	hresp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	defer hresp.Body.Close()

	var resp trasaResponse
	err = json.NewDecoder(hresp.Body).Decode(&resp)
	if err != nil {
		return errors.Errorf("failed to read kex response: %v", err)
	}

	// update secret key
	pubBytes, err := hex.DecodeString(resp.Data)
	if err != nil {
		return err
	}

	var serverPublicKey [32]byte

	copy(serverPublicKey[:], pubBytes)

	sec := utils.ECDHComputeSecret(priv, &serverPublicKey)

	logrus.Debug("our secret key: ", hex.EncodeToString(sec))
	keyHolder[intent] = sec

	return nil

}

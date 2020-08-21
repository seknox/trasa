package crypt

import (
	"encoding/hex"
	"net/http"
	"time"

	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

type KexRequest struct {
	// Intent of kex, example - KEX_EXPORT_DH, KEX_ENROL_DEVICE, KEX_HTTP_SR
	Intent string `json:"intent"`
	// IntentID is unique identifier for key exchange process. Should be trasaID in KEX_ENROL_DEVICE and sessionID in case of other intents
	IntentID string `json:"intentID"`
	// Unique id of device which is making kex request. Machine ID can be used here?
	DeviceID string `json:"deviceID"`
	// Public key of client
	PublicKey string `json:"publicKey"`
}

// Kex expects  (1) intent, intentID and public key of client, (2) generates server keypair (3) generates secret key (4) store secret key in ecdhKexDerivedKey map and (5) resopods with public key received in step 2.
//  This handler is only to be used be client side flow i.e. kex request initiated by client. eg. device register process.
func Kex(w http.ResponseWriter, r *http.Request) {

	logrus.Trace("kex request received")
	var req KexRequest
	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "success", "invalid request", "Kex", nil)
		return
	}

	//  TODO imp!! @bhrg3se, if intent is not KEX_ENROL_DEVICE and deviceID is not present in device table, exit this handler.

	// generate our keypair
	priv, pub, err := utils.ECDHGenKeyPair()
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to generate key pair", "Kex", nil)
		return
	}

	pubKeyBytes, err := hex.DecodeString(req.PublicKey)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to decode public key", "Kex", nil)
		return
	}

	var clientPublicKey [32]byte

	copy(clientPublicKey[:], pubKeyBytes)

	// gen secret
	sec := utils.ECDHComputeSecret(priv, &clientPublicKey)

	logrus.Trace("our secret: ", hex.EncodeToString(sec))

	var kexData = global.KexDerivedKey{
		DeviceID:  req.DeviceID,
		Secretkey: sec,
		Timestamp: time.Now().Unix(),
	}

	// store secret in ECDHKexDerivedKey
	global.ECDHKexDerivedKey[req.IntentID] = kexData

	// respond with our public key
	utils.TrasaResponseWithDataString(w, 200, "success", "kex process completed on server", "Kex", hex.EncodeToString(pub[:]))

}

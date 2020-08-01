package tfa

import (
	"net/http"

	"github.com/seknox/trasa/server/utils"
)

//todo check if fields are correct
type U2f struct {
	serviceID  string `json:"serviceID"`
	AppName    string `json:"appSecret"`
	Answer     string `json:"answer"`
	Challange  string `json:"challenge"`
	DeviceID   string `json:"deviceId"`
	Signature  string `json:"signature"`
	DeviceInfo string `json:"deviceInfo"`
}

// U2f handler
func U2fHandler(w http.ResponseWriter, r *http.Request) {
	var u2fRequest U2f

	if err := utils.ParseAndValidateRequest(r, &u2fRequest); err != nil {
		utils.TrasaResponse(w, 200, "failed", "invalid request", "U2FHandle")
		return
	}

	c, ok := u2fChanMap[u2fRequest.Challange]
	if !ok {
		utils.TrasaResponse(w, 200, "failed", "invalid request", "", nil)
		return
	}

	c <- u2fRequest
	utils.TrasaResponse(w, 200, "success", "", "", nil)

}

//TODO u2f sign and verify

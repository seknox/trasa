package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	curl "github.com/andelf/go-curl"
)

type trasaPAMConfig struct {
	TrasaServerURL     string `json:"trasaServerURL"`
	ServiceID          string `json:"serviceID"`
	ServiceKey         string `json:"serviceKey"`
	OfflineUsers       string `json:"offlineUsers"`
	InsecureSkipVerify bool   `json:"insecureSkipVerify"`
	Debug              bool   `json:"debug"`
}

type configFile struct {
	TrasaPAMConfig trasaPAMConfig `toml:"trasaPAM"`
}

func readConfigFromFile() (configFile, error) {
	var config configFile

	if _, err := toml.DecodeFile("/etc/trasa/config/trasapam.toml", &config); err != nil {
		return config, fmt.Errorf("%v", err)
	}

	return config, nil

}

// AppLogin mocks request structure which ssh logins and rdp logins generates
type serviceAgentLogin struct {
	ServiceID       string `json:"serviceID"`
	DynamicAuthApp  bool   `json:"dynamicAuthApp"`
	ServiceKey      string `json:"serviceKey"`
	User            string `json:"user"`
	TfaMethod       string `json:"tfaMethod"`
	TotpCode        string `json:"totpCode"`
	UserIP          string `json:"userIP"`
	UserWorkstation string `json:"workstation"`
	TrasaID         string `json:"trasaID"`
	OrgID           string `json:"orgID"`
}

type trasaResponseStruct struct {
	Status string      `json:"status"`
	Error  error       `json:"error,omitempty"`
	Reason string      `json:"reason,omitempty"`
	Intent string      `json:"intent,omitempty"`
	Data   interface{} `json:"data"`
}

func sendTfaReq(user, trasaID, totpCode, userIP string) bool {

	var reqData serviceAgentLogin
	reqData.ServiceID = ServiceConfig.TrasaPAMConfig.ServiceID
	reqData.ServiceKey = ServiceConfig.TrasaPAMConfig.ServiceKey
	reqData.User = user
	reqData.TrasaID = trasaID
	if totpCode == "" || len(totpCode) != 6 {
		reqData.TfaMethod = "U2F"
	} else {
		reqData.TfaMethod = "totp"
	}

	reqData.TotpCode = totpCode
	reqData.UserIP = userIP

	requestBody, _ := json.Marshal(&reqData)

	urlPath := fmt.Sprintf("%s/auth/agent/nix", ServiceConfig.TrasaPAMConfig.TrasaServerURL)

	// begin curl operation
	easy := curl.EasyInit()
	defer easy.Cleanup()

	easy.Setopt(curl.OPT_URL, urlPath)
	if ServiceConfig.TrasaPAMConfig.InsecureSkipVerify == true {
		easy.Setopt(curl.OPT_SSL_VERIFYPEER, false)
	} else {
		easy.Setopt(curl.OPT_SSL_VERIFYPEER, true)
	}

	easy.Setopt(curl.OPT_POSTFIELDS, string(requestBody))

	authStatus := false

	// handleResponse is a callbackfunction to parse and process trasa response
	handleResponse := func(buf []byte, userdata interface{}) bool {

		var result trasaResponseStruct
		err := json.Unmarshal(buf, &result)
		if err != nil {
			writeLog(fmt.Sprintf("[SendTfaReq] failed to parse response body"))
			return false
		}

		if ServiceConfig.TrasaPAMConfig.Debug {

			mar, err := json.Marshal(result)
			if err != nil {
				writeLog("Failed to parse http json response")
			}

			writeLog(fmt.Sprintf("response data: %s", string(mar)))
		}

		// return true if result is true
		if result.Status == "success" {
			authStatus = true

		}

		return true
	}

	if ServiceConfig.TrasaPAMConfig.Debug {

		writeLog(fmt.Sprintf("sending tfa request with  url: %s |  request data: %s", urlPath, string(requestBody)))
	}

	easy.Setopt(curl.OPT_WRITEFUNCTION, handleResponse)

	if err := easy.Perform(); err != nil {
		writeLog(fmt.Sprintf("[Perform()] failed connect to trasa server %s.", err.Error()))
		// If we fail here, it means no contact was made to trasa server.
		// we will check and return for offline user access.
		offU := offlineUsers(user, ServiceConfig.TrasaPAMConfig.OfflineUsers)
		if offU == true {
			writeLog(fmt.Sprintf("[Offline Access] Allowing offline access for user %s.", user))
			return true
		}
	}

	return authStatus

}

// offlineUsers splits csv offline users retreived from confif file and returns boolean value based on username match
func offlineUsers(username string, usernames string) bool {
	users := strings.Split(usernames, ",")

	resp := false
	for _, v := range users {
		str := strings.TrimSpace(v)
		if str == username {
			resp = true
		}

	}

	return resp
}

// writeLog is error or info log writer for trasapam. writes log in /var/log/trasapam.log
func writeLog(errval string) {
	logPath := "/var/log/trasapam.log"
	//logPath := "/root/trasapam.log"
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var buf []byte
	data := fmt.Sprintf("\n%s - %s\n", time.Now().String(), errval)

	buf = append(buf, []byte(data)...)

	_, err = file.Write(buf)
}

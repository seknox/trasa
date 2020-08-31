package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/seknox/trasa/cli/config"
	"github.com/seknox/trasa/server/api/auth"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
)

// TrasaAuth gets apps assigned to user
func SendHygiene(email, pass string, hygiene []byte, pubKey string) ([]byte, error) {

	var hygieneWithKey struct {
		ClientPubKey string `json:"clientPubKey"`
		EncryptedDH  string `json:"encryptedDH"`
	}

	err := json.Unmarshal(hygiene, &hygieneWithKey)
	if err != nil {
		logger.Debug(err)
		fmt.Println("Invalid device hygiene from trasaExtComm")
		return nil, err
	}

	var reqBody auth.UpdateHygienereq
	reqBody.TrasaID = email
	reqBody.DeviceHygiene = hygieneWithKey.EncryptedDH
	reqBody.ClientKey = hygieneWithKey.ClientPubKey
	reqBody.Token = pubKey

	url := config.Context.TRASA_URL + "/auth/device/cli/updatehygiene"
	//logger.Debug((reqBody.PublicKey))
	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}

	//fmt.Printf("request sent was: %s\n", req.RequestURI)
	var respData struct {
		models.TrasaResponseStruct
		Data [][]byte `json:"data"`
	}

	//TODO add --insecure flag
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//
	//fmt.Printf("resp body was: %s\n", string(body))

	err = json.Unmarshal([]byte(body), &respData)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("status was: %s\n", respData.AppUsers)

	if (respData.Status != "success") || len(respData.Data) < 1 {
		return nil, errors.New(respData.Reason)
	}

	return respData.Data[0], nil

}

func Auth(trasaID, pass string) (pubKey string, dhrequired bool, err1 error) {
	var reqBody auth.LoginRequest
	reqBody.Email = trasaID
	reqBody.Password = pass
	reqBody.IdpName = "trasa"

	url := config.Context.TRASA_URL + "/auth/identity"
	//logger.Debug((reqBody.PublicKey))
	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		return "", false, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	if err != nil {
		logger.Fatal(err)
		return "", false, err
	}

	//fmt.Printf("request sent was: %s\n", req.RequestURI)
	var respData struct {
		models.TrasaResponseStruct
		Data []string `json:"data"`
	}

	//TODO add --insecure flag
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return "", false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", false, err
	}

	//
	//fmt.Printf("resp body was: %s\n", string(body))

	err = json.Unmarshal([]byte(body), &respData)
	if err != nil {
		return "", false, err
	}
	//fmt.Printf("status was: %s\n", respData.AppUsers)

	if (respData.Status != "success") || len(respData.Data) < 1 {
		return "", respData.Intent == consts.AUTH_RESP_TFA_DH_REQUIRED, errors.New(respData.Reason)
	}

	return respData.Data[0], respData.Intent == consts.AUTH_RESP_TFA_DH_REQUIRED, nil

}

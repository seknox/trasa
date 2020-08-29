package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/seknox/trasa/cli/config"
	"github.com/seknox/trasa/server/api/auth"
	"github.com/seknox/trasa/server/models"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
)

// TrasaAuth gets apps assigned to user
func SendHygiene(email, pass, hygiene string) ([]byte, error) {
	var reqBody auth.UpdateHygienereq
	reqBody.TrasaID = email
	reqBody.Password = pass
	reqBody.DeviceHygiene = hygiene

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

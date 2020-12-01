package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/seknox/trasa/server/models"
)

var Validator *validator.Validate

func init() {
	Validator = validator.New()
}

// TrasaResponseWithDataString expects string value in data
func TrasaResponseWithDataString(w http.ResponseWriter, httpRespCode int, status, reason string, intent, data string) {
	var resp models.TrasaResponseStructWIthDataString
	w.Header().Set("Content-Type", "application/json")

	//TODO remove this after all inapp trail related apis are ported to TrasaResponseWithTrail
	w.Header().Set("trailDescription", intent)
	w.Header().Set("trailStatus", status)

	w.WriteHeader(httpRespCode)

	resp.Status = status
	resp.Reason = reason
	resp.Intent = intent
	resp.Data = data
	write, err := json.Marshal(resp)
	if err != nil {
		logrus.Error(err)
	}
	w.Write(write)
}

// TrasaResponse is generic response function for http api
//  Use TrasaResponseWithTrail if in apptrail is needed. TrasaResponseWithTrail has separate trailDecription for inapp trail description
func TrasaResponse(w http.ResponseWriter, httpRespCode int, status, reason string, intent string, data ...interface{}) {
	var resp models.TrasaResponseStruct
	w.Header().Set("Content-Type", "application/json")

	//TODO remove this after all inapp trail related apis are ported to TrasaResponseWithTrail
	w.Header().Set("trailDescription", intent)
	w.Header().Set("trailStatus", status)
	w.Header().Set("Cache-Control", "no-store, pragma")
	w.WriteHeader(httpRespCode)

	resp.Status = status
	resp.Reason = reason
	resp.Intent = intent
	resp.Data = data
	write, err := json.Marshal(resp)
	if err != nil {
		logrus.Error(err)
	}
	w.Write(write)
}

// TrasaResponseWithTrail is generic response function for http api.
func TrasaResponseWithTrail(w http.ResponseWriter, httpRespCode int, status, reason, intent, trailDecription string, data ...interface{}) {

	var resp models.TrasaResponseStruct
	w.Header().Set("Content-Type", "application/json")

	if trailDecription != "" {
		w.Header().Set("trailDescription", trailDecription)
		w.Header().Set("trailStatus", status)
	}

	w.WriteHeader(httpRespCode)

	resp.Status = status
	resp.Reason = reason
	resp.Intent = intent
	resp.Data = data
	write, err := json.Marshal(resp)
	if err != nil {
		logrus.Error(err)
	}
	w.Write(write)
}

// GetIp returns user's origin IP address
func GetIp(r *http.Request) string {

	ip := r.Header.Get("X-Real-IP")
	if strings.Compare(ip, "") == 0 {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			logrus.Error(err)
			return "0.0.0.0"
		}
		return ip
	}

	return ip
}

//ParseAndValidateRequest unmarshalls request body into given struct and also verify json fields
func ParseAndValidateRequest(r *http.Request, reqStruct interface{}) error {

	if err := json.NewDecoder(r.Body).Decode(&reqStruct); err != nil {
		return err
	}
	if err := Validator.Struct(reqStruct); err != nil {
		return err
	}
	return nil
}

//GetHttpClient return a http client
func GetHttpClient(insecure bool) *http.Client {
	var client *http.Client

	if insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = &http.Client{}
	}

	return client
}

//CallTrasaAPI is generic function to call TRASA API
func CallTrasaAPI(path string, reqBody interface{}, insecure bool) (resp models.TrasaResponseStruct, err error) {
	var client *http.Client

	if insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = &http.Client{}
	}
	defer client.CloseIdleConnections()

	var reqBodyBytes []byte
	reqBodyBytes, err = json.Marshal(reqBody)
	if err != nil {
		return
	}

	var req *http.Request

	req, err = http.NewRequest("POST", path, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return
	}

	var esp *http.Response
	esp, err = client.Do(req)
	if err != nil {
		return
	}
	defer esp.Body.Close()

	var respBody []byte
	respBody, err = ioutil.ReadAll(esp.Body)
	if err != nil {
		return
	}

	//logrus.Debug(string(respBody))
	err = json.Unmarshal(respBody, &resp)

	return
}

func ParseTrasaResponse(data []byte) (resp models.TrasaResponseStruct, err error) {
	err = json.Unmarshal(data, &resp)
	return
}

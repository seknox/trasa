package notif

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"firebase.google.com/go/messaging"
	"github.com/seknox/trasa/server/models"
)

// SendPushNotification sends push notification to TRASA mobile app
func (s notifStore) SendPushNotification(fcmToken, orgName, appName, ipAddr, time, challenge string) error {
	ctx := context.Background()

	// Obtain a messaging.Client from the App.
	//client, err := app.Messaging(ctx)

	client, err := s.FirebaseClient.Messaging(ctx)
	if err != nil {
		return err
	}

	var notif messaging.Notification

	notif.Title = "TRASA"
	notif.Body = "authorize login request"
	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"orgName":   orgName,
			"appName":   appName,
			"ipAddr":    ipAddr,
			"time":      time,
			"challenge": challenge,
			"type":      "u2f",
		},
		Notification: &notif,
		Token:        fcmToken,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	_, err = client.Send(ctx, message)

	return err

	// Response is a message ID string.
	//logger.Trace("successfully sent message:", response)

}

//CallTrasaCloudProxy is generic function to call TRASA Cloud Proxy with api key
func (s notifStore) CallTrasaCloudProxy(path string, reqBody interface{}, insecure bool) (resp models.TrasaResponseStruct, err error) {
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

	// retreive
	tsxcpxyKey := s.TsxCPxyKey
	req.Header.Set("TSXCPXY-KEY", tsxcpxyKey)

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

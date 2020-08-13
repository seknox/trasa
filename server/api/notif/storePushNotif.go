package notif

import (
	"context"

	"firebase.google.com/go/messaging"
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

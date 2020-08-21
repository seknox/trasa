package server

import (
	"github.com/seknox/trasa/tests/server/notiftest"
	"testing"
)

func TestServer(t *testing.T) {
	t.Run("test inapp notifications", func(t *testing.T) {
		n := notiftest.AddNotif(t)
		notiftest.GetPendingNotif(t, n)
		notiftest.ResolvNotif(t, n.NotificationID)
	})

}

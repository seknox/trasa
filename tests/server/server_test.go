package server

import (
	"github.com/seknox/trasa/tests/server/crudtest"
	"github.com/seknox/trasa/tests/server/notiftest"
	"testing"
)

func TestServer(t *testing.T) {
	t.Run("test inapp notifications", func(t *testing.T) {
		n := notiftest.AddNotif(t)
		notiftest.GetPendingNotif(t, n)
		notiftest.ResolvNotif(t, n.NotificationID)
	})

	t.Run("crud and accessmap ", func(t *testing.T) {
		crudtest.CreateService(t)
		serviceID := crudtest.GetService(t)
		crudtest.UpdateService(t, serviceID)
		crudtest.DeleteService(t, serviceID)

		p := crudtest.CreatePolicy(t)
		crudtest.GetPolicy(t, p)
		crudtest.GetPolicies(t, p)
		crudtest.UpdatePolicy(t, p.PolicyID)
		crudtest.DeletePolicy(t, p.PolicyID)

	})

}

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

		org := crudtest.UpdateOrg(t)
		crudtest.GetOrg(t, org)

		crudtest.CreateService(t)
		serviceID := crudtest.GetService(t)
		crudtest.UpdateService(t, serviceID)

		p := crudtest.CreatePolicy(t)
		crudtest.GetPolicy(t, p)
		crudtest.GetPolicies(t, p)
		crudtest.UpdatePolicy(t, p.PolicyID)

		user := crudtest.CreateUser(t)
		crudtest.GetUser(t, user)
		crudtest.GetAllUsers(t, user)
		crudtest.UpdateUser(t, user.ID)

		userGroupID := crudtest.CreateGroup(t, "usergroup")
		serviceGroupID := crudtest.CreateGroup(t, "servicegroup")

		crudtest.GetAllGroups(t, "user", userGroupID)
		crudtest.GetAllGroups(t, "service", serviceGroupID)

		crudtest.UpdateGroup(t, userGroupID)
		crudtest.UpdateServiceGroup(t, serviceID, serviceGroupID, "add")
		crudtest.GetServiceGroupDetail(t, serviceID, serviceGroupID)
		crudtest.UpdateServiceGroup(t, serviceID, serviceGroupID, "remove")

		crudtest.UpdateUserGroup(t, user.ID, userGroupID, "add")
		crudtest.GetUserGroupDetail(t, user.ID, userGroupID)
		crudtest.UpdateUserGroup(t, user.ID, userGroupID, "remove")

		//Delete
		crudtest.DeleteService(t, serviceID)
		crudtest.DeletePolicy(t, p.PolicyID)
		crudtest.DeleteGroup(t, userGroupID)
		crudtest.DeleteGroup(t, serviceGroupID)
		crudtest.DeleteUser(t, user.ID)

	})

}

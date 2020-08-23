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

		//create service
		crudtest.CreateService(t)
		serviceID := crudtest.GetService(t)
		crudtest.UpdateService(t, serviceID)

		//create policy
		policy := crudtest.CreatePolicy(t)
		crudtest.GetPolicy(t, policy)
		crudtest.GetPolicies(t, policy)
		crudtest.UpdatePolicy(t, policy.PolicyID)

		//create user
		user := crudtest.CreateUser(t)
		crudtest.GetUser(t, user)
		crudtest.GetAllUsers(t, user)
		crudtest.UpdateUser(t, user.ID)

		//create group
		userGroupID := crudtest.CreateGroup(t, "usergroup")
		serviceGroupID := crudtest.CreateGroup(t, "servicegroup")

		crudtest.GetAllGroups(t, "user", userGroupID)
		crudtest.GetAllGroups(t, "service", serviceGroupID)

		//add service to group
		crudtest.UpdateGroup(t, userGroupID)
		crudtest.UpdateServiceGroup(t, serviceID, serviceGroupID, "add")
		crudtest.GetServiceGroupDetail(t, serviceID, serviceGroupID)

		//add user to group
		crudtest.UpdateUserGroup(t, user.ID, userGroupID, "add")
		crudtest.GetUserGroupDetail(t, user.ID, userGroupID)

		//accessmaps
		crudtest.CreateUserAccessMap(t, serviceID, user.ID, policy.PolicyID, "testpriv")
		crudtest.CreateUserGroupAccessMap(t, serviceID, userGroupID, policy.PolicyID, "testpriv")
		crudtest.CreateUserGroupServiceGroupAccessMap(t, serviceGroupID, userGroupID, policy.PolicyID, "testpriv")

		userAccessMapID := crudtest.GetUserAccessMap(t, serviceID, user.ID, policy.PolicyID, "testpriv")
		userGroupAccesMapID := crudtest.GetUserGroupsAssignedToServiceGroups(t, serviceID, userGroupID, policy.PolicyID, "testpriv")
		usergroupServiegroupMapID := crudtest.GetUserGroupServiceGroupAccessMaps(t, serviceGroupID, userGroupID, policy.PolicyID, "testpriv")

		crudtest.UpdateUserAccessMap(t, userAccessMapID)
		crudtest.UpdateUserGroupAccessMap(t, userGroupAccesMapID)
		crudtest.UpdateUserGroupAccessMap(t, usergroupServiegroupMapID)

		//Delete

		crudtest.DeleteUserAccessMap(t, userAccessMapID)
		crudtest.DeleteUserGroupServiceGroupAccessMap(t, userGroupAccesMapID)
		crudtest.DeleteUserGroupServiceGroupAccessMap(t, usergroupServiegroupMapID)

		crudtest.UpdateUserGroup(t, user.ID, userGroupID, "remove")
		crudtest.UpdateServiceGroup(t, serviceID, serviceGroupID, "remove")

		crudtest.DeleteService(t, serviceID)
		crudtest.DeletePolicy(t, policy.PolicyID)
		crudtest.DeleteGroup(t, userGroupID)
		crudtest.DeleteGroup(t, serviceGroupID)
		crudtest.DeleteUser(t, user.ID)

	})

}

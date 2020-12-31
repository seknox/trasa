package accessmap

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/groups"
	"github.com/seknox/trasa/server/api/misc"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

func GetUserAccessMaps(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	serviceID := chi.URLParam(r, "serviceID")

	appusers, err := Store.GetServiceUserMaps(serviceID, userContext.Org.ID)
	if err != nil {
		logrus.Error(err)
		return
	}

	utils.TrasaResponse(w, 200, "success", "app details fetched", "GetAppDetailFromserviceIDV2", appusers)
}

func GetUserGroupServiceGroupAccessMaps(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("request received")
	userContext := r.Context().Value("user").(models.UserContext)
	groupID := chi.URLParam(r, "serviceGroupID")

	userGroups, err := Store.GetAssignedUserGroupsWithPolicies(groupID, userContext.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not get access maps", "AllUsergroupAndPoliciesToAdd")
		return
	}

	utils.TrasaResponse(w, 200, "success", "", "AllAddedUserGroups", userGroups)
}

type AssignUserToApp struct {
	ServiceID string   `json:"serviceID"`
	OrgID     string   `json:"orgID"`
	Privilege string   `json:"privilege"`
	Users     []string `json:"users"`
	PolicyID  []string `json:"policyID"`
}

//AssignUser To Service
//User Can be assigned toservicein different ways
// User Group/Array to Service
// User to Service Group/Array
// Single User To Single Service
//There are different APIs to do these

func CreateServiceUserMap(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	var request AssignUserToApp

	if err := utils.ParseAndValidateRequest(r, &request); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error parsing request", "could not assign users to app")
		return
	}

	//Removed this validation to support multiple username in same authserviceby same user
	//But we still need to validate same user is assigned to sameservicewith same username
	//which can be done by database constraint

	// Add User to App

	successNum := 0
	for _, v := range request.Users {
		var serviceUserMap models.ServiceUserMap

		serviceUserMap.MapID = utils.GetUUID()
		serviceUserMap.ServiceID = request.ServiceID
		serviceUserMap.OrgID = userContext.Org.ID
		serviceUserMap.UserID = v
		serviceUserMap.Privilege = strings.ToLower(request.Privilege)
		serviceUserMap.PolicyID = request.PolicyID[0]
		serviceUserMap.AddedAt = time.Now().Unix() //.In(nep).String()
		err := Store.CreateServiceUserMap(&serviceUserMap)

		if err != nil {
			//TODO check constraint
			logrus.Error(err)
			continue
		}
		successNum = successNum + 1
	}
	appName, err := services.Store.GetServiceNameFromID(request.ServiceID, userContext.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "", "could not assign users to app")
		return
	}
	intent := fmt.Sprintf("asigned %d users to %s", successNum, appName)

	utils.TrasaResponse(w, http.StatusOK, "success", "successfully assigned users", intent)

}

type DeleteServiceUserMapReq struct {
	MapIDs []string `json:"mapIDs"`
}

func DeleteServiceUserMap(w http.ResponseWriter, r *http.Request) {
	var req DeleteServiceUserMapReq

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		utils.TrasaResponse(w, 200, "failed", "invalid request", "could not remove user from Service")
		return
	}

	userContext := r.Context().Value("user").(models.UserContext)

	var serviceName string
	deletedCount := 0
	for _, v := range req.MapIDs {
		serviceName, err := Store.DeleteServiceUserMap(v, userContext.Org.ID)
		if err != nil {
			logrus.Debug(err)
			utils.TrasaResponse(w, 200, "failed", "failed removing user", "could not remove user from service")
			return
		}
		serviceName = serviceName
		deletedCount = deletedCount + 1
	}

	utils.TrasaResponse(w, 200, "success", "", fmt.Sprintf(`%d user removed from "%s" service`, deletedCount, serviceName), nil)
	return
}

//Groups

type ServiceGroupUserGroupMapRequest struct {
	MapID          string   `json:"mapID"`
	ServiceGroupID string   `json:"serviceGroupID"`
	MapType        string   `json:"mapType"`
	UserGroupID    []string `json:"userGroupID"`
	Privilege      string   `json:"privilege"`
	OrgID          string   `json:"orgID"`
	PolicyID       []string `json:"policyID"`
	CreatedAt      int64    `json:"createdAt"`
}

func CreateServiceGroupUserGroupMap(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("request received")
	userContext := r.Context().Value("user").(models.UserContext)
	var req ServiceGroupUserGroupMapRequest

	var store models.ServiceGroupUserGroupMap

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "json parse error", "user group not assigned to servicegroup")
		return
	}

	store.OrgID = userContext.User.OrgID
	store.PolicyID = req.PolicyID[0]
	store.ServiceGroupID = req.ServiceGroupID
	store.MapType = req.MapType
	store.Privilege = strings.ToLower(req.Privilege)

	var service *models.Service
	var serviceGroup models.Group
	var err error
	if store.MapType == "service" {
		service, err = services.Store.GetFromID(req.ServiceGroupID)
		if err != nil {
			logrus.Debug(err)
			utils.TrasaResponse(w, 200, "failed", "invalid service id", "user group not assigned to servicegroup")
			return
		}
	} else if store.MapType == "servicegroup" {
		serviceGroup, err = groups.Store.Get(req.ServiceGroupID, userContext.Org.ID)
		if err != nil {
			logrus.Debug(err)
			utils.TrasaResponse(w, 200, "failed", "invalid group id", "user group not assigned to servicegroup")
			return
		}
	} else {
		logrus.Debug("invalid map type")
		utils.TrasaResponse(w, 200, "failed", "invalid map type", "user group not assigned to servicegroup")
		return
	}

	var addedGroups = make([]string, 0)
	for _, v := range req.UserGroupID {
		store.MapID = utils.GetUUID()
		store.CreatedAt = time.Now().Unix()
		store.UserGroupID = v
		err := Store.CreateServiceGroupUserGroupMap(&store)
		if err != nil {
			logrus.Error(err)
			if store.MapType == "service" {
				utils.TrasaResponse(w, 200, "failed", "user group not assigned", fmt.Sprintf(`user group not assigned to serviceg "%s"`, service.Name), nil, nil)
				return
			}
			utils.TrasaResponse(w, 200, "failed", "user group not assigned", fmt.Sprintf(`user group not assigned to servicegroup "%s"`, serviceGroup.GroupName), nil, nil)
			return
		}

		//collect group name for inapp trail
		groupName, _, err := misc.Store.GetEntityDescription(v, consts.ENTITY_GROUP, userContext.Org.ID)
		if err != nil {
			logrus.Error(err)
			continue
		}
		addedGroups = append(addedGroups, groupName)
	}

	if store.MapType == "service" {
		utils.TrasaResponse(w, 200, "success", "group mapping successful", fmt.Sprintf(`user groups %s assigned to service "%s"`, strings.Join(addedGroups, ","), service.Name), addedGroups)
		return
	}

	utils.TrasaResponse(w, 200, "success", "group mapping successful", fmt.Sprintf(`user groups %s assigned to servicegroup "%s"`, strings.Join(addedGroups, ","), serviceGroup.GroupName), addedGroups)
}

type RmGroupMap struct {
	MapID []string `json:"mapID"`
}

func DeleteServiceGroupUserGroupMap(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("request received")
	userContext := r.Context().Value("user").(models.UserContext)

	var req RmGroupMap

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "access map not deleted", "RemoveUsergroupsFromGroupMap")
		return
	}

	var appGroupName, userGroupName string
	appGroupName = ""
	var err error
	removedUserGroups := make([]string, 0)
	for _, v := range req.MapID {
		appGroupName, userGroupName, err = Store.DeleteServiceGroupUserGroupMap(v, userContext.User.OrgID)
		if err != nil {
			logrus.Debug(err)
		}
		removedUserGroups = append(removedUserGroups, userGroupName)
		// TODO if any deletion failed, it must be reported back to user.

	}

	utils.TrasaResponse(w, 200, "success", "access map removed", fmt.Sprintf(`user groups "%s" removed from appgroup %s`, strings.Join(removedUserGroups, ","), appGroupName))
}

func GetUserGroupsAssignedToServiceGroups(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("request received")
	userContext := r.Context().Value("user").(models.UserContext)
	groupID := chi.URLParam(r, "groupID")

	userGroups, err := Store.GetAssignedUserGroupsWithPolicies(groupID, userContext.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not get access map", "AllUsergroupAndPoliciesToAdd")
		return
	}

	utils.TrasaResponse(w, 200, "success", "access map fetched", "AllAddedUserGroups", userGroups)
}

type UpdatePrivilege struct {
	MapID     string `json:"mapID"`
	Privilege string `json:"privilege"`
}

func UpdateServiceUserMap(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	var req UpdatePrivilege
	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "Privilege not updated")
		return
	}

	err := Store.UpdateServiceUserMap(req.MapID, userContext.User.OrgID, req.Privilege)
	if err != nil {
		logrus.Debug(err)
		utils.TrasaResponse(w, 200, "failed", "could not update", "Privilege not updated")
		return
	}

	utils.TrasaResponse(w, 200, "success", "username update", fmt.Sprintf(`Privilege updated to "%s"`, req.Privilege))

}

func UpdateServiceGroupUserGroup(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	var req UpdatePrivilege
	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "privilege not updated", nil, nil)
		return
	}

	err := Store.UpdateServiceGroupUserGroupMap(req.MapID, userContext.User.OrgID, req.Privilege)
	if err != nil {
		//TODO check constraint
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not update privilege", "privilege not updated", nil, nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "username update", fmt.Sprintf(`Privilege updated to "%s"`, req.Privilege))

}

func UserGroupsToAdd(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("request received")
	userContext := r.Context().Value("user").(models.UserContext)

	userGroups, err := Store.GetUserGroupsToAddInServiceGroup(userContext.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "group not created", "AllUsergroupAndPoliciesToAdd")
		return
	}

	utils.TrasaResponse(w, 200, "success", "group details fetched", "AllUsergroupAndPoliciesToAdd", userGroups)
}

type UserGroupOfServiceGroup struct {
	MapID         string `json:"mapID"`
	UsergroupID   string `json:"usergroupID"`
	UsergroupName string `json:"userGroupName"`
	Privilege     string `json:"privilege"`
	PolicyName    string `json:"policyName"`
	PolicyID      string `json:"policyID"`
	AddedAt       int64  `json:"addedAt"`
}

//GetAllDynamicAccessRules returns all dynamic access rules
func GetAllDynamicAccessRules(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	dps, err := Store.GetAllDynamicAccessRules(uc.Org.ID)

	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not get rules", "get dynamic access rules", nil, nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "", "get dynamic access rules", dps)
}

//GetAllUserGroupsWithIDPs returns all user groups
func GetAllUserGroupsWithIDPs(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	gps, err := Store.getAllUserGroupsWithIDPs(uc.Org.ID)

	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not get groups", "get dynamic access user groups")
		return
	}

	utils.TrasaResponse(w, 200, "success", "", "get dynamic access user groups", gps)
}

//CreateDynamicAccessRule creates a dynamic access rule
func CreateDynamicAccessRule(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	var req models.DynamicAccessRule

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error parsing request", "dynamic access rule not created", nil, nil)
		return
	}

	req.RuleID = utils.GetUUID()
	req.OrgID = uc.Org.ID
	req.CreatedAt = time.Now().Unix()

	err := Store.CreateDynamicAccessRule(req)

	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not create  rule", "dynamic access rule not created", nil, nil)
		return
	}

	allrules, err := Store.GetAllDynamicAccessRules(uc.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not get rules", "dynamic access rule not created", nil, nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "Rule added", "dynamic access rule created", allrules)
}

//DeleteDynamicAccessRule deletes a dynamic access rule
func DeleteDynamicAccessRule(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	var req models.DynamicAccessRule

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error parsing request", "dynamic access rule not deleted", nil, nil)
		return
	}

	err := Store.DeleteDynamicAccessRule(req.RuleID, uc.Org.ID)

	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not delete rule", "dynamic access rule not deleted", nil, nil)
		return
	}

	allrules, err := Store.GetAllDynamicAccessRules(uc.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not get rules", "dynamic access rule not deleted", nil, nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "Rule deleted", "dynamic access rule deleted", allrules)

}

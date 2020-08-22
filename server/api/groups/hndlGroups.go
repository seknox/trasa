package groups

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/misc"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

// CreateGroup handles grouping of any entity - users, services or policies.
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	var req models.Group

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		http.Error(w, http.StatusText(400), 400)
		logrus.Error(err)
		return
	}

	req.GroupID = utils.GetUUID()
	req.GroupName = strings.ToLower(req.GroupName)
	req.OrgID = userContext.User.OrgID
	req.Status = true
	req.CreatedAt = time.Now().Unix()
	req.UpdatedAt = req.CreatedAt
	err := Store.Create(&req)
	if err != nil {
		reason := utils.GetConstraintErrorMessage(err)
		utils.TrasaResponse(w, 200, "failed", reason, "CreateGroup")
		return
	}
	utils.TrasaResponse(w, 200, "success", "", fmt.Sprintf(`group %s created of type %s`, req.GroupName, req.GroupType), req.GroupID)
}

// UpdateGroup handles grouping of any entity - users, services or policies.
func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	var req models.Group

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		http.Error(w, http.StatusText(400), 400)
		logrus.Error(err)
		return
	}
	req.UpdatedAt = time.Now().Unix()
	req.OrgID = userContext.User.OrgID
	req.GroupName = strings.ToLower(req.GroupName)
	err := Store.Update(&req)
	if err != nil {
		reason := utils.GetConstraintErrorMessage(err)
		utils.TrasaResponse(w, 200, "failed", reason, "group not updated")
		return
	}
	utils.TrasaResponse(w, 200, "success", "group updated", fmt.Sprintf(`%s group updated`, req.GroupName), req.GroupID)
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	groupID := chi.URLParam(r, "groupID")
	groupName, err := Store.Delete(groupID, userContext.User.OrgID)
	if err != nil {
		logrus.Errorf("could not delete group: %v", err)
		utils.TrasaResponse(w, 200, "failed", "group not deleted", "group not deleted")
		return
	}
	utils.TrasaResponse(w, 200, "success", "group deleted", fmt.Sprintf(`group "%s" deleted`, groupName))

}

func GetAllGroups(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	groupType := chi.URLParam(r, "groupType")

	var groups []models.Group
	var err error

	//TODO make const for group type
	if groupType == "user" {
		groups, err = Store.GetAllUserGroups(uc.User.OrgID)
	} else if groupType == "service" {
		groups, err = Store.GetAllServiceGroups(uc.User.OrgID)
	}
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error fetching groups", "GetGroups")
		return
	}
	utils.TrasaResponse(w, 200, "success", "group fetched", "GetGroups", groups)
}

func GetUserGroup(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	groupID := chi.URLParam(r, "groupid")

	var resp GroupUsers

	group, err := Store.Get(groupID, uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error fetching group detail", "GetUsersGroupDetail")
		return
	}

	resp.GroupMeta = group

	resp.AddedUsers, err = Store.GetUsersInGroup(groupID, uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error fetching group detail", "GetUsersInGroup")
		return
	}

	// TODO : Bug @bhrg3se - should only return user's that are not already in the group.
	resp.UnaddedUsers, err = Store.GetUsersNotInGroup(groupID, uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error fetching group detail", "GetUsersNotInGroup")
		return
	}

	utils.TrasaResponse(w, 200, "success", "apps fetched", "GetUsersGroupDetail", resp)

}

type GroupUsers struct {
	GroupMeta    models.Group  `json:"groupMeta"`
	AddedUsers   []models.User `json:"addedUsers"`
	UnaddedUsers []models.User `json:"unaddedUsers"`
}

func GetServiceGroup(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	groupID := chi.URLParam(r, "groupID")

	var resp GroupApps

	group, err := Store.Get(groupID, uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error fetching group detail", "UpdateServiceGroup")
		return
	}
	resp.GroupMeta = group

	resp.AddedServices, err = Store.GetServicesInGroup(groupID, uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error fetching group detail", "UpdateServiceGroup")
		return
	}

	resp.UnaddedServices, err = Store.GetServicesNotInGroup(groupID, uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error fetching group detail", "UpdateServiceGroup")
		return
	}

	utils.TrasaResponse(w, 200, "success", "apps fetched", "GetServiceGroupDetail", resp)

}

type GroupApps struct {
	GroupMeta       models.Group     `json:"groupMeta"`
	AddedServices   []models.Service `json:"addedServices"`
	UnaddedServices []models.Service `json:"unaddedServices"`
}

func UpdateServiceGroup(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	var req UpdateServiceGroupReq

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		http.Error(w, http.StatusText(400), 400)
		logrus.Error(err)
		return
	}

	group, err := Store.Get(req.GroupID, userContext.User.OrgID)
	if err != nil || group.GroupID == "" {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error fetching group detail", "UpdateServiceGroup")
		return
	}
	logrus.Debug("update ssss ", req.ServiceIDs)
	if req.UpdateType == "add" {
		logrus.Debug("update type is add ", req.ServiceIDs)
		var sgMap models.ServiceGroupMap
		sgMap.OrgID = userContext.User.OrgID
		sgMap.GroupID = req.GroupID

		err := Store.AddServicesToGroup(group, req.ServiceIDs)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "could not add services to group", fmt.Sprintf(`could not add services to group "%s"`, group.GroupName))
			return
		}

		var addedServices []string = make([]string, 0)
		for _, serviceID := range req.ServiceIDs {
			added, _, err := misc.Store.GetEntityDescription(serviceID, consts.ENTITY_APP, userContext.Org.ID)
			if err != nil {
				logrus.Error(err)
				continue
			}
			addedServices = append(addedServices, added)
		}

		utils.TrasaResponse(w, 200, "success", "successfully added", fmt.Sprintf(`services %s added to group "%s"`, strings.Join(addedServices, ","), group.GroupName), addedServices)
		return
	}

	if req.UpdateType == "remove" {
		// remove users
		err := Store.RemoveServicesFromGroup(req.GroupID, userContext.User.OrgID, req.ServiceIDs)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "could not remove services", "update service group")
			return
		}

		var removedServices []string = make([]string, 0)
		for _, serviceID := range req.ServiceIDs {
			name, _, err := misc.Store.GetEntityDescription(serviceID, consts.ENTITY_APP, userContext.Org.ID)
			if err != nil {
				logrus.Error(err)
				continue
			}
			removedServices = append(removedServices, name)
		}

		utils.TrasaResponse(w, 200, "success", "", fmt.Sprintf(`services %s removed from group "%s"`, strings.Join(removedServices, ","), group.GroupName), removedServices)
		return
	}
}

type UpdateServiceGroupReq struct {
	GroupID    string   `json:"groupID"`
	UpdateType string   `json:"updateType"` // add or delete
	ServiceIDs []string `json:"serviceIDs"`
}

type UpdateUsersGroupReq struct {
	GroupID    string   `json:"groupID"`
	UpdateType string   `json:"updateType"` // add or delete
	UserIDs    []string `json:"userIDs",validate:"min=1"`
}

func UpdateUsersGroup(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("request received")
	userContext := r.Context().Value("user").(models.UserContext)

	var req UpdateUsersGroupReq

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 400, "failed", "error parsing request", "group not updated")
		return
	}

	group, err := Store.Get(req.GroupID, userContext.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error fetching group detail", "group not updated")
		return
	}

	if req.UpdateType == "add" {
		var userGroup models.UserGroupMap
		userGroup.OrgID = userContext.User.OrgID
		userGroup.GroupID = req.GroupID

		// add services to userGroupv1 table
		err := Store.AddUsersToGroup(group, req.UserIDs)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "could not add users.", fmt.Sprintf(`could not add users to group "%s"`, group.GroupName))
			return
		}

		var addedUsers []string = make([]string, 0)
		for _, userID := range req.UserIDs {
			email, _, err := misc.Store.GetEntityDescription(userID, consts.ENTITY_USER, userContext.Org.ID)
			if err != nil {
				logrus.Error(err)
				continue
			}
			addedUsers = append(addedUsers, email)
		}
		utils.TrasaResponse(w, 200, "success", "users added", fmt.Sprintf(`users %s added to group "%s"`, strings.Join(addedUsers, ","), group.GroupName), addedUsers)
		return
	}

	if req.UpdateType == "remove" {
		// remove users
		err := Store.RemoveUsersFromGroup(req.GroupID, userContext.User.OrgID, req.UserIDs)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "could not add users.", fmt.Sprintf(`could not remove users to group "%s"`, group.GroupName), nil, nil)
			return
		}

		var removedUsers []string = make([]string, 0)
		for _, userID := range req.UserIDs {
			email, _, err := misc.Store.GetEntityDescription(userID, consts.ENTITY_USER, userContext.Org.ID)
			if err != nil {
				logrus.Error(err)
				continue
			}
			removedUsers = append(removedUsers, email)
		}

		utils.TrasaResponse(w, 200, "success", "users removed", fmt.Sprintf(`removed users %s from group "%s"`, strings.Join(removedUsers, ","), group.GroupName), removedUsers)
		return
	}
}

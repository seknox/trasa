package uidp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/notif"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

/////////////////////////////////////////////////////////////////////////
/////////////////////      SCIM Reference     ///////////////////////////
/////////////////////////////////////////////////////////////////////////
// 	HTTP   SCIM Usage
// 	Method
// 	------ --------------------------------------------------------------
// 	GET    Retrieves one or more complete or partial resources.

// 	POST   Depending on the endpoint, creates new resources, creates a
// 		   search request, or MAY be used to bulk-modify resources.

// 	PUT    Modifies a resource by replacing existing attributes with a
// 		   specified set of replacement attributes (replace).  PUT
// 		   MUST NOT be used to create new resources.

// 	PATCH  Modifies a resource with a set of client-specified changes
// 		   (partial update).

// 	DELETE Deletes a resource.

// 						 Table 1: SCIM HTTP Methods

//  Hunt, et al.                 Standards Track                    [Page 9]

//  RFC 7644               SCIM Protocol Specification        September 2015

// 	Resource Endpoint         Operations             Description
// 	-------- ---------------- ---------------------- --------------------
// 	User     /Users           GET (Section 3.4.1),   Retrieve, add,
// 							  POST (Section 3.3),    modify Users.
// 							  PUT (Section 3.5.1),
// 							  PATCH (Section 3.5.2),
// 							  DELETE (Section 3.6)

// SCIMCreateUser creates user in TRASA based on data from SCIM request.
func SCIMCreateUser(w http.ResponseWriter, r *http.Request) {

	var req models.ScimUser

	// parse json value into struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "CreateUser", nil, nil)
		return
	}

	// get ScimContext
	uc := r.Context().Value("scimprov").(models.ScimContext)

	primaryEmail := getPrimaryEmail(req.Emails)
	// normalized. i.e change to lower case.
	normalizedEmail := utils.NormalizeString(primaryEmail)
	normalizedUserName := utils.NormalizeString(req.UserName)

	var user models.UserWithPass
	user.ID = utils.GetUUID()
	user.OrgID = uc.OrgID
	user.ExternalID = req.ExternalID
	user.OrgName = uc.Orgname
	user.Email = normalizedEmail
	user.UserName = normalizedUserName
	user.FirstName = req.Name.GivenName
	user.MiddleName = req.Name.MiddleName
	user.LastName = req.Name.FamilyName
	user.UserRole = "selfUser"
	user.IdpName = utils.NormalizeString(uc.IdpName)
	user.Password = ""
	user.Status = true

	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	err := users.Store.Create(&user)
	if err != nil {
		logrus.Error(err)
		scimNotFoundOrConflictResp(w, 409, "user already exists", consts.SCIM_USER_SCHEMA)
		return
	}

	req.ID = user.ID
	req.X509Certificates = make([]models.ScimUserX509Certificates, 0)
	req.Groups = make([]models.ScimUserGroups, 0)


	scimUserResp(w, 201, req)

}

// SCIMGetSingleUser get detail of single user based on userID supplied in SCIM request.
func SCIMGetSingleUser(w http.ResponseWriter, r *http.Request) {
	
	uc := r.Context().Value("scimprov").(models.ScimContext)

	userID := chi.URLParam(r, "userID")
	userDetailFromDb, err := users.Store.GetFromID(userID, uc.OrgID)
	if err != nil {
		logrus.Error(err)
		scimNotFoundOrConflictResp(w, 404, "user not found", consts.SCIM_ERR)
		return
	}

	s := transformTuserToSuser(userDetailFromDb)

	scimUserResp(w, 200, s)
}

// SCIMGetSingleUsersWithFilter is filter based SCIM query. Currently we only support "eq" filter, meaning query with specefic trasaID or return all.
func SCIMGetSingleUsersWithFilter(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("scimprov").(models.ScimContext)

	
	filter := r.URL.Query().Get("filter")

	if filter != "" {
		// perform filter query

		// first we parse query string
		str, err := url.PathUnescape(filter)
		if err != nil {
			logrus.Error(err)
			scimNotFoundOrConflictResp(w, 403, "invalid query string", consts.SCIM_USER_SCHEMA)
			return
		}
		// when urlquery parsed, query string is of format userName eq "Runscope175Vmqxbvhrj999@atko.com"
		// we only support eq filter for now
		querystrs := strings.Split(str, " ")
		userName := strings.Trim(querystrs[2], "\"")

		userDetailFromDb, err := users.Store.GetFromTRASAID(userName, uc.OrgID)
		if err != nil {
			logrus.Debug(err)
			var ss = make([]models.ScimUser, 0)
			//s := transformTuserToSuser(userDetailFromDb)
			//ss = append(ss, s)
			scimUserListResp(w, 200, ss)
			return
		}

		//s := transformTuserToSuser(userDetailFromDb)

		var ss = make([]models.ScimUser, 0)
	
		s := transformTuserToSuser(userDetailFromDb)
		s.X509Certificates = make([]models.ScimUserX509Certificates, 0)
		s.Groups = make([]models.ScimUserGroups, 0)
		ss = append(ss, s)

		// scimUserResp(w, 200, s)
		scimUserListResp(w, 200, ss)
		return

	}

	// We will also handle "count" request. This implementation is clumsy handle of count and limit but will work
	//  since this will be mostly be used for SCIM API connection tests only.
	limit := r.URL.Query().Get("count")
	i, _ := strconv.Atoi(limit)
	userDetailFromDb, err := users.Store.GetFromWithLimit(uc.OrgID, i)
	if err != nil {
		logrus.Error(err)
		var ss = make([]models.ScimUser, 0)
		scimUserListResp(w, 200, ss)
		return
	}

	var ss = make([]models.ScimUser, 0)
	s := transformTuserToSuser(userDetailFromDb)
	ss = append(ss, s)

	scimUserListResp(w, 200, ss)

}

// SCIMPutSingleUser updates user profile (all details supplied by request). For single element update, use patch.
func SCIMPutSingleUser(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("scimprov").(models.ScimContext)
	userID := chi.URLParam(r, "userID")

	var req models.ScimUser

	//email := chi.URLParam(r, "userID")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		scimNotFoundOrConflictResp(w, 409, "user already exists", consts.SCIM_USER_SCHEMA)
		return
	}

	trasaUser := transformSuserToTuser(req, uc)
	trasaUser.ID = userID
	trasaUser.UserRole = "selfUser"

	userName := req.UserName
	if userName == "" {
		primaryEmail := getPrimaryEmail(req.Emails)
		userName = primaryEmail
	}
	trasaUser.UserName = userName

	if req.Active == false {
		err := users.Store.UpdateStatus(false, trasaUser.ID, uc.OrgID) //createUser(&user)
		if err != nil {
			logrus.Error(err)
			var s models.ScimUser
			scimUserResp(w, 200, s)
			return
		}
	}

	err := users.Store.Update(trasaUser) //createUser(&user)
	if err != nil {
		logrus.Error(err)
		var s models.ScimUser
		scimUserResp(w, 200, s)
		return
	}

	suser := transformTuserToSuser(&trasaUser)
	scimUserResp(w, 200, suser)

}

// SCIMPatchSingleUser patch single element. Currently only update user status is implemented.
func SCIMPatchSingleUser(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("scimprov").(models.ScimContext)

	var req models.ScimUser

	userID := chi.URLParam(r, "userID")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		scimNotFoundOrConflictResp(w, 409, "user already exists", consts.SCIM_USER_SCHEMA)
		return
	}

	err := users.Store.UpdateStatus(req.Active, userID, uc.OrgID) //createUser(&user)
	if err != nil {
		logrus.Error(err)
		var s models.ScimUser
		scimUserResp(w, 200, s)
		return
	}

	userDetailFromDb, err := users.Store.GetFromID(userID, uc.OrgID)
	if err != nil {
		logrus.Error(err)
		var s models.ScimUser
		scimUserResp(w, 200, s)
		return
	}

	suser := transformTuserToSuser(userDetailFromDb)
	scimUserResp(w, 200, suser)

}

// DeleteUser should be atomic transaction - a single user delete call should delete users detail from
// every database tabble connected with user.
func SCIMDeleteUser(w http.ResponseWriter, r *http.Request) {

	userID := chi.URLParam(r, "userID")

	sc := r.Context().Value("scimprov").(models.ScimContext)
	email, userRole, err := users.Store.Delete(userID, sc.OrgID) //createUser(&user)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed deleting user", "user not deleted", nil)
		return
	}

	err = users.Store.DeleteAllUserAccessMaps(userID, sc.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed deleting user.", "user not deleted", nil)
		return
	}

	err = users.Store.DeregisterUserDevices(userID, sc.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed deleting user.", "user not deleted", nil)
		return
	}

	utils.TrasaResponse(w, http.StatusOK, "success", "successfully deleted user", fmt.Sprintf(`user "%s" deleted`, email), nil)

	//TODO @sshah
	_ = userRole
	if userRole == "orgAdmin" {
		go notif.CheckAndFireSecurityRule(sc.OrgID, consts.DELETE_ADMIN_USER, email)

	} else {
		go notif.CheckAndFireSecurityRule(sc.OrgID, consts.DELETE_USER, email)
	}

}

/*

////////////////////////////////////////////////
///////// 	SCIM Utility functions for TRASA
////////////////////////////////////////////////


*/

// scimUserResp is a generic scim response handler
func scimUserResp(w http.ResponseWriter, respVal int, u models.ScimUser) {

	w.WriteHeader(respVal)
	w.Header().Set("Content-Type", "application/json")

	write, err := json.MarshalIndent(u, "	", "	")
	logrus.Trace(string(write))
	if err != nil {
		logrus.Error(err)
	}
	w.Write(write)
}

func scimUserListResp(w http.ResponseWriter, respVal int, u []models.ScimUser) {

	var l models.ScimListUser
	l.Schemas = []string{consts.SCIM_LISTRESP}
	l.Resources = u
	l.ItemsPerPage = 1
	l.TotalResults = len(u)
	l.StartIndex = 1
	w.WriteHeader(respVal)
	w.Header().Set("Content-Type", "application/json")

	write, _ := json.MarshalIndent(l, "	", "	")
	logrus.Trace(string(write))
	w.Write(write)
}

// scimNotFoundOrConflictResp is generic SCIM response for Not Found and Conflict.
func scimNotFoundOrConflictResp(w http.ResponseWriter, respVal int, detail, schema string) {

	var c models.ScimConflict
	c.Schemas = []string{schema}
	c.Detail = detail
	c.Status = respVal

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(respVal)

	write, err := json.Marshal(c)
	if err != nil {
		logrus.Error(err)
	}

	w.Write(write)
}

// transformTuserToSuser transforms trasa user to scim user
func transformTuserToSuser(u *models.User) models.ScimUser {
	var s models.ScimUser
	s.Schemas = []string{consts.SCIM_USER_SCHEMA}
	s.ID = u.ID
	s.ExternalID = u.ExternalID

	s.Emails = []models.ScimUserEmails{{
		Primary: true,
		Value:   u.Email,
		//Type:    "trasa",
	}}
	s.UserName = u.UserName
	s.Name = models.ScimUserName{
		GivenName:  u.FirstName,
		MiddleName: u.MiddleName,
		FamilyName: u.LastName,
	}

	var sgroups models.ScimUserGroups
	for _, v := range u.Groups {
		sgroups.Value = v
		s.Groups = append(s.Groups, sgroups)
	}

	s.Active = u.Status
	s.X509Certificates = []models.ScimUserX509Certificates{}

	s.Meta = models.ScimMeta{
		Created:      time.Unix(u.CreatedAt, 0).String(),
		LastModified: time.Unix(u.CreatedAt, 0).String(),
	}

	return s
}

// transformTuserToSuser transforms trasa user to scim user
func transformSuserToTuser(s models.ScimUser, uc models.ScimContext) models.User {
	var u models.User
	primaryEmail := ""
	for _, v := range s.Emails {
		if v.Primary == true {
			primaryEmail = v.Value
		}
	}
	// normalized. i.e change to lower case.
	normalizedEmail := utils.NormalizeString(primaryEmail)
	normalizedUserName := utils.NormalizeString(s.UserName)

	u.Email = normalizedEmail
	u.UserName = normalizedUserName
	u.UserRole = s.UserRole
	u.UpdatedAt = time.Now().Unix() //.In(timezone).String()
	u.FirstName = s.Name.FamilyName
	u.MiddleName = s.Name.MiddleName
	u.LastName = s.Name.FamilyName
	for _, v := range s.Groups {
		u.Groups = append(u.Groups, v.Value)
	}
	u.Status = s.Active
	return u
}

func getPrimaryEmail(emails []models.ScimUserEmails) string {
	primaryEmail := ""
	for _, v := range emails {
		if v.Primary == true {
			primaryEmail = v.Value
		}
	}
	return primaryEmail
}

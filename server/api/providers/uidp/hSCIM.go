package uidp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	logger "github.com/sirupsen/logrus"
)

// SCIMCreateUser creates user in TRASA based on data from SCIM request.
func SCIMCreateUser(w http.ResponseWriter, r *http.Request) {

	var req models.ScimUser

	// parse json value into struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "CreateUser", nil, nil)
		return
	}

	// get ScimContext
	uc := r.Context().Value("scimprov").(models.ScimContext)

	primaryEmail := ""
	for _, v := range req.Emails {
		if v.Primary == true {
			primaryEmail = v.Value
		}
	}
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

	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	err := users.Store.Create(&user)
	if err != nil {
		logger.Error(err)
		scimNotFoundOrConflictResp(w, 409, "user already exists", consts.SCIM_USER_SCHEMA)
		return
	}

	req.ID = user.ID

	scimUserResp(w, 201, req)

}

// scimUserResp is a generic scim response handler
func scimUserResp(w http.ResponseWriter, respVal int, u models.ScimUser) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(respVal)

	write, err := json.Marshal(u)
	if err != nil {
		logger.Error(err)
	}
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
		logger.Error(err)
	}

	w.Write(write)
}

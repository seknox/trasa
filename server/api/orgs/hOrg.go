package orgs

import (
	"net/http"

	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

// Get returns org account detail
func Get(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("get org detail request received")
	uc := r.Context().Value("user").(models.UserContext)

	org, err := Store.Get(uc.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to get org account details", "Get", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "", "Get", org)

}

// Update updates org accont details
func Update(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("Update org detail request received")

	var req models.Org

	// parse json value into struct
	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "Update", nil, nil)
		return
	}

	uc := r.Context().Value("user").(models.UserContext)
	req.ID = uc.User.OrgID

	err := Store.update((req))
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to update account details", "Update", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "", "Update", nil)

}

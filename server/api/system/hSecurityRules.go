package system

import (
	"net/http"

	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

// SecurityRules retrieves orgWide security rules
func SecurityRules(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("SecurityRules request received")
	uc := r.Context().Value("user").(models.UserContext)

	rules, err := Store.getSecurityRules(uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", err.Error(), "SecurityRules", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "SecurityRules fetched", "SecurityRules", rules)

}

type UpdateSecurityRulesReq struct {
	Status string `json:"status"`
	RuleID string `json:"ruleID"`
}

// UpdateSecurityRule updates system security rules
func UpdateSecurityRule(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("request received UpdateSecurityRule")
	uc := r.Context().Value("user").(models.UserContext)

	var req UpdateSecurityRulesReq

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request format", "UpdateSecurityRule", nil)
		return
	}

	err := Store.updateSecurityRule(uc.User.OrgID, req.Status, req.RuleID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", err.Error(), "UpdateSecurityRule", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "status updated", "UpdateSecurityRule", nil)

}

package policies

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

// CreatePolicy expects permissions from client and creates new policy group in database.
func CreatePolicy(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("CreatePolicy request received")
	uc := r.Context().Value("user").(models.UserContext)

	var req models.Policy

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "json unmarshal error", "Policy not created")
		return
	}

	req.PolicyName = strings.ToLower(req.PolicyName)
	req.PolicyID = utils.GetUUID()
	req.OrgID = uc.User.OrgID
	req.RiskThreshold = 5
	req.CreatedAt = time.Now().Unix()
	req.UpdatedAt = time.Now().Unix()
	req.AllowedCountries = ""
	err := Store.CreatePolicy(req)
	if err != nil {
		reason := utils.GetConstraintErrorMessage(err)
		utils.TrasaResponse(w, 200, "failed", reason, "Policy not created")
		return
	}

	utils.TrasaResponse(w, 200, "success", "successfully created policy", fmt.Sprintf(`Policy "%s" created`, req.PolicyName), req)

}

//UpdatePolicy updates given policy of an organization
func UpdatePolicy(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("UpdatePolicy request received")
	uc := r.Context().Value("user").(models.UserContext)

	var req models.Policy

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "json unmarshall error", "Policy not updated")
		return
	}

	req.PolicyName = strings.ToLower(req.PolicyName)
	req.OrgID = uc.User.OrgID
	req.RiskThreshold = 5
	req.UpdatedAt = time.Now().Unix()
	err := Store.UpdatePolicy(req)
	if err != nil {
		reason := utils.GetConstraintErrorMessage(err)
		utils.TrasaResponse(w, 200, "failed", reason, "Policy not updated")
		return
	}

	utils.TrasaResponse(w, 200, "success", "successfully updated policy", fmt.Sprintf(`Policy "%s" updated`, req.PolicyName))

}

//GetPolicies returns all policies of an organization
func GetPolicies(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("GetPolicies request received")
	uc := r.Context().Value("user").(models.UserContext)

	policies, err := Store.GetAllPolicies(uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to get policy", "Get policies")
		return
	}

	loc, err := time.LoadLocation(uc.Org.Timezone)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to get timezone", "Get policies")
		return
	}
	for i, policy := range policies {

		expiry, err := time.Parse("2006-01-02", policy.Expiry) //)
		if err != nil {
			logrus.Error(err)
		}
		expiry = expiry.In(loc)
		policy.IsExpired = expiry.Before(time.Now().In(loc))
		policies[i] = policy
	}

	utils.TrasaResponse(w, 200, "success", "", "get policies", policies)

}

//GetPolicy returns policy details from policyID
func GetPolicy(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("CreatePolicy request received")
	uc := r.Context().Value("user").(models.UserContext)

	policyID := chi.URLParam(r, "policyID")
	policy, err := Store.GetPolicy(policyID, uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to get policy", "Get policy")
		return
	}

	utils.TrasaResponse(w, 200, "success", "", "get policy", policy)

}

type policyIDs struct {
	PolicyID []string `json:"policyID"`
}

//DeletePolicies deletes given policies
func DeletePolicies(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("request received")
	userContext := r.Context().Value("user").(models.UserContext)

	var req policyIDs
	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "policy not created", "Policy not deleted")
		return
	}

	deletedNum := 0
	for _, v := range req.PolicyID {
		// delete policy
		err := Store.DeletePolicy(v, userContext.User.OrgID)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "Could not delete policy", "Policy not deleted")
			return
		}
		deletedNum = deletedNum + 1
	}

	// send response
	utils.TrasaResponse(w, 200, "success", "successfully removed policy", fmt.Sprintf(`%d Policies deleted`, deletedNum))
}

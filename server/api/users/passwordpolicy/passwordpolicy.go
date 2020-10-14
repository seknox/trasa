package passwordpolicy

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

//EnforcePasswordPolicyNow will enforce password policy immediately
func EnforcePasswordPolicyNow(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("request received")
	userContext := r.Context().Value("user").(models.UserContext)

	allUsers, err := users.Store.GetAll(userContext.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error fetching users", "failed to enforce password policy", nil, nil)
		return
	}

	var policy models.PolicyEnforcer

	policy.OrgID = userContext.User.OrgID

	policy.EnforceType = consts.ChangePassword
	policy.Pending = true
	policy.AssignedBy = userContext.User.ID
	policy.AssignedOn = time.Now().Unix()
	policy.ResolvedOn = time.Now().Unix()

	for _, v := range allUsers {
		policy.EnforceID = utils.GetRandomString(7)
		policy.UserID = v.ID
		err := users.Store.EnforcePolicy(policy)
		if err != nil {
			logrus.Error(err)
		}
	}

	utils.TrasaResponse(w, 200, "success", "Password policy enforced", "password policy enforced", nil, nil)

}

// CheckPendingPasswordRotationForUser checks if a user has pending password expiry that needs to be enforced.
func CheckPendingPasswordRotationForUser(userID, orgID string) (bool, error) {

	orgDetail, err := orgs.Store.Get(orgID)
	if err != nil {
		return false, err
	}

	gsetting, err := system.Store.GetGlobalSetting(orgID, consts.GLOBAL_PASSWORD_CONFIG)
	if err != nil || gsetting.Status == false {
		// if there is error getting row or global setting is not enabled, we return
		return false, err
	}

	var policy models.PasswordPolicy
	err = json.Unmarshal([]byte(gsetting.SettingValue), &policy)
	if err != nil {
		return false, err
	}

	passState, err := users.Store.GetPasswordState(userID, orgID)
	if err != nil || len(passState.LastPasswords) <= 0 {
		// If passstate is empty, this means user has never changed the password.
		// We check if gsetting specifies policy change
		if gsetting.Status == true && policy.Expiry == "never" {
			return false, err
		}
		if gsetting.Status == true && policy.Expiry != "never" {
			// we check if password is expired but based on date of gsetting
			check := checkisExpired(gsetting.UpdatedOn, orgDetail.Timezone, policy.Expiry)
			return check, err
		}

	}

	// if GetPasswordStateForUser does not return error or is non empty,
	// check password expiry based on time of passState.lastUpdated
	check := checkisExpired(passState.LastUpdated, orgDetail.Timezone, policy.Expiry)
	return check, nil

}

// if passState returns row, check the date of lastUpdate exceeds expiry days time.
// If passState returns error or nil, check date of gsetting ang check it exceeds expiry days time.
// return boolean value
func checkisExpired(thenTime int64, timeZone, days string) bool {
	// parse unix timestamp
	then := time.Unix(thenTime, 0)

	// parse time based on org current timezone
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		logrus.Errorf("load location: %v", err)
		loc, err = time.LoadLocation("UTC")
		if err != nil {
			logrus.Errorf("load location: %v", err)
			return true
		}
	}
	current := time.Now().In(loc)

	// parse elapsed duration since last update based on current timezone
	dur := current.Sub(then)

	// get diference in int
	diff := int(dur.Hours())
	// get days
	diffDays := diff / 24

	// expiry date are stored as string in format example "30 days".
	// get integere value of day. Consider storing expiry days in int rather than string to avoid conversion.
	day := strings.Split(days, " ")
	i, err := strconv.Atoi(day[0])
	if err != nil {
		return false
	}

	// if diff day is greated expiry day, return true
	if diffDays >= i {
		return true
	}

	// return false by default
	return false
}

// EnforceChangePassword will enforce users to change password
func EnforceChangePassword(userID, orgID string) error {
	var policy models.PolicyEnforcer

	policy.OrgID = orgID

	policy.EnforceType = consts.ChangePassword
	policy.Pending = true
	policy.AssignedBy = consts.SYSTEM
	policy.AssignedOn = time.Now().Unix()
	policy.ResolvedOn = time.Now().Unix()
	policy.EnforceID = utils.GetUUID()
	policy.UserID = userID
	err := users.Store.EnforcePolicy(policy)

	return err

}

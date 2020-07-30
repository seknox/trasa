package serviceauth

import (
	"github.com/seknox/trasa/consts"
	"github.com/seknox/trasa/models"
)

func (s SStore) CheckPolicy(serviceID, userID, orgID, userIP, timezone string, policy *models.Policy, adhoc bool) (bool, consts.FailedReason) {

	params := models.ConnectionParams{
		ServiceID: serviceID,
		UserID:    userID,
		OrgID:     orgID,
		UserIP:    userIP,
		Timezone:  timezone,
	}
	return s.checkPolicyFunc(&params, policy, adhoc)
}

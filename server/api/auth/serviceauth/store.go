package serviceauth

import (
	"github.com/seknox/trasa/server/api/accesscontrol"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
)

func (s sStore) CheckPolicy(serviceID, userID, orgID, userIP, timezone string, policy *models.Policy, adhoc bool) (bool, consts.FailedReason) {

	params := models.ConnectionParams{
		ServiceID: serviceID,
		UserID:    userID,
		OrgID:     orgID,
		UserIP:    userIP,
		Timezone:  timezone,
	}
	return accesscontrol.CheckPolicy(&params, policy, adhoc)
}

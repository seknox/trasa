package notif

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/seknox/trasa/server/api/misc"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

// CheckAndFireSecurityRule, when called based on constName queries database to see if there is a rule
// that has been enabled and matches situation. If true, it should also perform actions according to rule.
// Currently only hardcoded and email and dashboard notifications are supported.
func CheckAndFireSecurityRule(orgID, constName, entityValue string) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error(fmt.Errorf(fmt.Sprintf(`%v:%s`, r, string(debug.Stack()))), "Panic in CheckAndFireSecurityRule")
		}

	}()
	// get rule from database
	rule, err := system.Store.GetSecurityRuleByName(orgID, constName)
	if err != nil {
		logrus.Error(err)
		return
	}

	// get all admins emails
	getAdmins, err := misc.Store.GetAdmins(orgID) //dbstore.Connect.CRDBGetUsersBasedOnRole("orgAdmin", orgID)
	if err != nil {
		logrus.Error(err)
		return
	}
	// send alert email
	var tmplt models.EmailSecurityAlert
	tmplt.EntityName = entityValue
	tmplt.SecurityRuleTitle = rule.Name
	tmplt.SecurityRuleText = rule.Description

	for _, v := range getAdmins {
		tmplt.ReceiverEmail = v.Email
		err = Store.SendEmail(orgID, consts.EMAIL_SECURITY_ALERT, tmplt)
		if err != nil {
			logrus.Error(err)
		}
		// sendMail, id, err := utils.SecurityAlertMail(v.Email, rule.Name, rule.Description, entityValue)
		// if err != nil {
		// 	logrus.Error(err)
		// }
		// _, _ = id, sendMail
	}

	// store dashboard alert
	var notif models.InAppNotification
	notif.OrgID = orgID
	notif.NotificationLabel = "security-alert"
	notif.NotificationText = rule.Name
	notif.CreatedOn = time.Now().Unix() //time.Now().In(nep).Unix()
	notif.EmitterID = constName
	notif.IsResolved = false
	notif.ResolvedOn = 0

	for _, v := range getAdmins {
		notif.NotificationID = utils.GetRandomID(10)
		notif.UserID = v.ID
		// we call notification store to store notification for this event in user id who is to be notified.
		err = Store.StoreNotif(notif)
		if err != nil {
			logrus.Error(err)
		}
	}

}

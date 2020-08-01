package notif

import (
	"net/http"
	"time"

	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

func GetPendingNotif(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	notifs, err := Store.GetPendingNotif(uc.User.ID, uc.User.OrgID)
	if err != nil {
		logrus.Debug(err)
		utils.TrasaResponse(w, 200, "failed", "failed to get notification", "GetPendingNotif")
		return
	}

	utils.TrasaResponse(w, 200, "success", "notification fetched", "GetPendingNotif", notifs)
}

type resolveNotifReq struct {
	NotifID string `json:"notifID"`
}

func ResolveNotif(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	var req resolveNotifReq

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request format", "ResolveNotif")
		return
	}

	var notif models.InAppNotification
	notif.NotificationID = req.NotifID
	notif.OrgID = userContext.User.OrgID
	notif.IsResolved = true
	notif.ResolvedOn = time.Now().Unix()

	// we call notification store to store notification for this event in user id who is to be notified.
	err := Store.UpdateNotifFromNotifID(notif)
	if err != nil {
		logrus.Debug(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request format", "ResolveNotif")
		return
	}

	utils.TrasaResponse(w, 200, "success", "notification fetched", "ResolveNotif")
}

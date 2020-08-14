package notif

import (
	"github.com/seknox/trasa/server/models"
)

// CRDBStoreNotif stores notification that is to be notified to user.
func (s notifStore) StoreNotif(notif models.InAppNotification) (err error) {
	_, err = s.DB.Exec(`INSERT into inapp_notifs (id, user_id, emitter_id, org_id, label,text, created_on, is_resolved, resolved_on)
						 values($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		notif.NotificationID, notif.UserID, notif.EmitterID, notif.OrgID, notif.NotificationLabel, notif.NotificationText, notif.CreatedOn, notif.IsResolved, notif.ResolvedOn)
	return
}

func (s notifStore) UpdateNotif(notif models.InAppNotification) error {
	_, err := s.DB.Exec(`UPDATE inapp_notifs set is_resolved=$3, resolved_on=$4 WHERE emitter_id=$1 AND org_id=$2;`,
		notif.EmitterID, notif.OrgID, notif.IsResolved, notif.ResolvedOn)
	return err

}

func (s notifStore) UpdateNotifFromNotifID(notif models.InAppNotification) error {

	_, err := s.DB.Exec(`UPDATE inapp_notifs set is_resolved=$1, resolved_on=$2 WHERE id=$3 AND org_id=$4;`,
		notif.IsResolved, notif.ResolvedOn, notif.NotificationID, notif.OrgID)
	return err

}

func (s notifStore) GetPendingNotif(userID, orgID string) ([]models.InAppNotification, error) {

	/////////////////
	var notif models.InAppNotification
	var notifs = make([]models.InAppNotification, 0)

	rows, err := s.DB.Query("SELECT id, user_id, emitter_id, org_id, label,text, created_on FROM inapp_notifs WHERE user_id=$1 AND org_id=$2 AND is_resolved=$3 ORDER BY created_on DESC;",
		userID, orgID, false)

	if err != nil {
		return notifs, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&notif.NotificationID, &notif.UserID, &notif.EmitterID, &notif.OrgID, &notif.NotificationLabel, &notif.NotificationText, &notif.CreatedOn)
		if err != nil {
			return notifs, err

		}
		notifs = append(notifs, notif)
	}

	return notifs, nil

}

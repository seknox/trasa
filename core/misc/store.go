package misc

import (
	"net"

	"github.com/seknox/trasa/consts"
	"github.com/seknox/trasa/models"
)

//GetAdminEmails returns all users of a organizations
func (s MiscStore) GetAdmins(orgID string) ([]models.User, error) {
	var users = make([]models.User, 0)

	rows, err := s.DB.Query(`SELECT users.org_id, users.id,  username, first_name, middle_name,
								   last_name, email, user_role,
								   users.status AND COALESCE(idp.is_enabled,true) as status,
								   created_at, updated_at, users.idp_name
							FROM users
							LEFT JOIN idp  on users.idp_name = idp.name WHERE users.user_role = $1 AND users.org_id = $2`, "orgAdmin", orgID)

	if err != nil {
		return users, err
	}

	defer rows.Close()
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.OrgID, &user.ID,
			&user.UserName, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.UserRole, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.IdpName)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, err

}

// CRDBStoreNotif stores notification that is to be notified to user.
func (s MiscStore) StoreNotif(notif models.InAppNotification) (err error) {
	_, err = s.DB.Exec(`INSERT into inapp_notifs (id, user_id, emitter_id, org_id, label,text, created_on, is_resolved, resolved_on)
						 values($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		notif.NotificationID, notif.UserID, notif.EmitterID, notif.OrgID, notif.NotificationLabel, notif.NotificationText, notif.CreatedOn, notif.IsResolved, notif.ResolvedOn)
	return
}

func (s MiscStore) UpdateNotif(notif models.InAppNotification) error {
	_, err := s.DB.Exec(`UPDATE inapp_notifs set is_resolved=$3, resolved_on=$4 WHERE emitter_id=$1 AND org_id=$2;`,
		notif.EmitterID, notif.OrgID, notif.IsResolved, notif.ResolvedOn)
	return err

}

func (s MiscStore) UpdateNotifFromNotifID(notif models.InAppNotification) error {

	_, err := s.DB.Exec(`UPDATE inapp_notifs set is_resolved=$1, resolved_on=$2 WHERE id=$3 AND org_id=$4;`,
		notif.IsResolved, notif.ResolvedOn, notif.NotificationID, notif.OrgID)
	return err

}

func (s MiscStore) GetGeoLocation(ip string) (geo models.GeoLocation, err error) {
	locations, err := s.Geoip.City(net.ParseIP(ip))
	if err != nil {
		return geo, err

	}

	geo.Country = locations.Country.Names["en"]
	geo.City = locations.City.Names["en"]
	geo.TimeZone = locations.Location.TimeZone
	geo.Location = []float64{locations.Location.Longitude, locations.Location.Latitude}

	return geo, nil
}

func (s MiscStore) GetEntityDescription(entityID string, entityType consts.EntityConst, orgID string) (string, string, error) {
	var query = ``
	var val1, val2 string
	switch entityType {
	case consts.ENTITY_APP:
		query = `SELECT name,type FROM services where id=$1 AND org_id=$2`
	case consts.ENTITY_USER:
		query = `SELECT email,username FROM users where id=$1 AND org_id=$2`
	case consts.ENTITY_APP_USER_MAP:
		query = `SELECT u.email,a.name from user_accessmaps au
					JOIN services a ON a.id=au.service_id
					JOIN users u ON u.id= au.user_id
					JOIN policies p ON au.policy_id=p.id
					WHERE au.id=$1 AND a.org_id=$2`
	case consts.ENTITY_GROUP:
		query = `SELECT name,'' from groups WHERE id=$1 AND org_id=$2`
	case consts.ENTITY_APP_GROUP_MAP:
		query = `SELECT g.name,s.name FROM service_group_maps sg 
					JOIN groups g on g.id = sg.group_id
					JOIN services s on sg.service_id = s.id
					WHERE sg.id=$1 AND sg.org_id=$2`
	case consts.ENTITY_USER_GROUP_MAP:
		query = `SELECT g.name,email from user_group_maps ug
					JOIN users u ON u.id=ug.user_id
					JOIN groups g ON g.id=ug.group_id
					WHERE ug.id=$1 AND ug.org_id=$2`
	case consts.ENTITY_USERGROUP_APPGROUP_MAP:
		query = `SELECT ug.name,ag.name from usergroup_accessmaps ag_ug
					JOIN groups ug on ag_ug.usergroup_id = ug.id
					JOIN groups ag on ag_ug.servicegroup_id= ag.id
					WHERE ag_ug.id=$1 AND ag_ug.org_id=$2`
	case consts.ENTITY_USER_DEVICE:
		query = `SELECT device_hygiene FROM devices WHERE id=$1 AND org_id=$2`

	}
	err := s.DB.QueryRow(query, entityID, orgID).Scan(&val1, &val2)
	if err != nil {
		return "", "", err
	}
	return val1, val2, nil

}

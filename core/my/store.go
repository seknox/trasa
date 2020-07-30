package my

import (
	"encoding/json"

	"github.com/seknox/trasa/models"
)

//This is new api to get my apps which retrives authapps based on group policies

func (s MyStore) GetUserAppsDetailsWithPolicyFromUserID(userID, orgID string) ([]models.MyServiceDetails, error) {
	var appusers = make([]models.MyServiceDetails, 0)

	//New Query which returns adhoc permissions too
	//Joins user_accessmaps,authappsv1 and adhoc_perms
	//adhoc_perms is sub-queried to get only unexpired
	rows, err := s.DB.Query(`SELECT DISTINCT services.name, services.id,gappusersv1.privilege, user_id,
                policy_id,
                policy_name,
                device_policy::bytes ,
                record_session,
                file_transfer,
                COALESCE(ip_source,''),
                risk_threshold,
                gappusersv1.tfa_enabled,
                day_time::bytes,
                expiry,
                COALESCE(allowed_countries,''),
                services.created_at,adhoc,services.type,hostname,
                COALESCE(authorized_period,0) AS authorized_till,
                COALESCE(authorized_on,0) AS authorized_on,
                COALESCE(requested_on,0) AS requested_on
FROM (
         --user assigned to app
         SELECT service_id,user_id,
                policies.id as policy_id,
                policies.name as policy_name,
                device_policy::varchar,
                record_session,
                file_transfer,
                ip_source,
                risk_threshold,
                policies.tfa_enabled,
                day_time::varchar,
                expiry,
                allowed_countries,
                privilege,tfa_enabled from
             user_accessmaps
                 JOIN policies on user_accessmaps.policy_id=policies.id
         where user_id=$1

               --usergroup assigned to app
         UNION SELECT servicegroup_id as authapp_id,user_id,
                      p.id as policy_id,
                      p.name as policy_name,
                      device_policy::varchar,
                      record_session,
                      file_transfer,
                      ip_source,
                      risk_threshold,
                      p.tfa_enabled,
                      day_time::varchar,
                      expiry,
                      allowed_countries,
                      privilege,tfa_enabled FROM
             (
              usergroup_accessmaps ag_ug
                 join policies p on ag_ug.policy_id=p.id
                 join user_group_maps ug on ug.group_id=ag_ug.usergroup_id

                 ) WHERE ag_ug.map_type='service'

                         --usergroup assigned to appgroup
         UNION 	SELECT ag.service_id,user_id,
                         p.id as policy_id,
                         p.name as policy_name,
                         device_policy::varchar,
                         record_session,
                         file_transfer,
                         ip_source,
                         risk_threshold,
                         p.tfa_enabled,
                         day_time::varchar,
                         expiry,
                         allowed_countries,
                         ag_ug.privilege,tfa_enabled FROM
             (
              usergroup_accessmaps ag_ug
                 join policies p on ag_ug.policy_id=p.id
                 join user_group_maps ug on ug.group_id=ag_ug.usergroup_id
                 join service_group_maps ag on ag.group_id=ag_ug.servicegroup_id
                 ) WHERE ag_ug.map_type='servicegroup'


     ) as gappusersv1
         JOIN services ON  gappusersv1.service_id=services.id
         LEFT JOIN (
    SELECT requester_id,service_id,authorized_period,authorized_on,requested_on from adhoc_perms
    WHERE NOT adhoc_perms.is_expired
) as adhoc_perms
                   ON adhoc_perms.requester_id=gappusersv1.user_id AND adhoc_perms.service_id=services.id
WHERE gappusersv1.user_id=$1 AND services.org_id=$2;`,
		userID, orgID)
	if err != nil {
		return appusers, err
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return appusers, err
	}

	defer rows.Close()
	for rows.Next() {
		var appuser models.MyServiceDetails

		var dayTime string
		err := rows.Scan(&appuser.ServiceName, &appuser.ServiceID, &appuser.Privilege, &appuser.UserID,
			&appuser.Policy.PolicyID,
			&appuser.Policy.PolicyName,
			&appuser.Policy.DevicePolicy,
			&appuser.Policy.RecordSession,
			&appuser.Policy.FileTransfer,
			&appuser.Policy.IPSource,
			&appuser.Policy.RiskThreshold,
			&appuser.Policy.TfaRequired,
			&dayTime,
			&appuser.Policy.Expiry,
			&appuser.Policy.AllowedCountries,

			&appuser.UserAddedAt, &appuser.Adhoc, &appuser.ServiceType, &appuser.Hostname, &appuser.AuthorizedTill, &appuser.AuthorizedOn, &appuser.RequestedOn)
		if err != nil {
			return appusers, err

		}
		err = json.Unmarshal([]byte(dayTime), &appuser.Policy.DayAndTime)
		if err != nil {
			return appusers, err
		}

		appuser.Usernames = make([]string, 0)

		var newAppuser bool = true
		for i, tempuser := range appusers {
			if tempuser.ServiceID == appuser.ServiceID {
				appusers[i].Usernames = append(appusers[i].Usernames, appuser.Privilege)
				newAppuser = false
			}
		}
		if newAppuser {
			appuser.Usernames = append(appuser.Usernames, appuser.Privilege)
			appusers = append(appusers, appuser)
		}

	}

	//	mar, _ := json.Marshal(appusers)
	return appusers, err
}

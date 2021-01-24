package accessmap

import (
	"encoding/json"
	"github.com/lib/pq"
	"github.com/seknox/trasa/server/models"
	"github.com/sirupsen/logrus"
)

//TODO handle error
func (s accessMapStore) CheckIfPrivilegeExist(privilege, orgID, serviceID string) bool {
	var check bool = true
	s.DB.QueryRow(`
SELECT EXISTS 
(
	(Select privilege FROM 
	user_accessmaps 
	WHERE privilege=$1 AND org_id=$2 AND service_id=$3)
	UNION (
	SELECT ag_ug.privilege from usergroup_accessmaps ag_ug
		JOIN user_group_maps ug ON ag_ug.usergroup_id=ug.group_id AND map_type='service'
	WHERE ag_ug.privilege=$1 AND ag_ug.org_id=$2 AND servicegroup_id=$3 
	)
	UNION (
	SELECT ag_ug.privilege FROM usergroup_accessmaps ag_ug
		JOIN user_group_maps ug ON ag_ug.usergroup_id=ug.group_id AND map_type='servicegroup'
		JOIN service_group_maps ag ON ag_ug.servicegroup_id=ag.group_id
	WHERE ag_ug.privilege=$1 AND ag_ug.org_id=$2 AND ag.service_id=$3
	)
)	
	
	
	
	`, privilege, orgID, serviceID).Scan(&check)
	return check

}

func (s accessMapStore) GetServiceUserMaps(serviceID, orgID string) (appusers []models.AccessMapDetail, err error) {
	appusers = make([]models.AccessMapDetail, 0)
	rows, err := s.DB.Query(`
SELECT map.id, service_id, map.user_id,
       COALESCE(NULLIF(users.email,''),users.username),
       map.privilege, map.policy_id, policies.name, map.added_at 
	FROM user_accessmaps map
	    JOIN users ON map.user_id=users.id
	    join policies on policies.id=map.policy_id
WHERE service_id= $1 AND map.org_id=$2`, serviceID, orgID)
	if err != nil {
		return appusers, err
	}

	defer rows.Close()
	for rows.Next() {
		appuser := models.AccessMapDetail{
			Policy: models.Policy{},
		}
		err = rows.Scan(&appuser.MapID, &appuser.ServiceID, &appuser.UserID, &appuser.TrasaID, &appuser.Privilege, &appuser.Policy.PolicyID, &appuser.Policy.PolicyName, &appuser.UserAddedAt)
		if err != nil {
			return
		}

		appusers = append(appusers, appuser)
	}

	return appusers, err
}

func (s accessMapStore) UpdateServiceUserMap(mapID, orgID, privilege string) error {
	_, err := s.DB.Exec(`UPDATE user_accessmaps SET privilege = $1 WHERE id = $2 AND org_id = $3;`,
		privilege, mapID, orgID)
	return err
}

func (s accessMapStore) CreateServiceUserMap(appUser *models.ServiceUserMap) error {
	_, err := s.DB.Exec(`INSERT into user_accessmaps (id, service_id,  org_id, user_id,privilege, policy_id, added_at)
						values($1, $2, $3, $4, $5, $6, $7);`, &appUser.MapID, &appUser.ServiceID, &appUser.OrgID, &appUser.UserID, &appUser.Privilege, &appUser.PolicyID, &appUser.AddedAt)

	return err
}

func (s accessMapStore) DeleteServiceUserMap(mapID, orgID string) (string, error) {
	var name string
	//TODO check query
	err := s.DB.QueryRow(`
	With deleted AS (DELETE FROM user_accessmaps WHERE id = $1 AND org_id=$2 RETURNING *)
	select services.name from deleted d join services ON d.service_id=services.id
		`, mapID, orgID).
		Scan(&name)
	if err != nil {
		return "", err
	}

	return name, err
}

////////////////////////////////////////////////////////////////////////////////////
////////////////   	AppGroup UserGroup MAP 	////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////
func (s accessMapStore) CreateServiceGroupUserGroupMap(data *models.ServiceGroupUserGroupMap) error {
	_, err := s.DB.Exec(`INSERT into usergroup_accessmaps (id, org_id, servicegroup_id, map_type, usergroup_id,privilege, policy_id, created_at)
						values($1, $2, $3, $4, $5, $6, $7, $8);`,
		&data.MapID, &data.OrgID, &data.ServiceGroupID, &data.MapType, &data.UserGroupID, &data.
			Privilege, &data.PolicyID, &data.CreatedAt)

	return err
}

func (s accessMapStore) UpdateServiceGroupUserGroupMap(mapID, orgID, privilege string) error {
	_, err := s.DB.Exec(`UPDATE usergroup_accessmaps SET privilege = $1 WHERE id = $2 AND org_id = $3;`,
		privilege, mapID, orgID)

	return err
}

func (s accessMapStore) DeleteServiceGroupUserGroupMap(mapID, orgID string) (string, string, error) {
	var appGroupName, userGroupName string
	//TODO check query
	err := s.DB.QueryRow(`WITH deleted as (DELETE FROM usergroup_accessmaps WHERE id = $1 AND org_id=$2 RETURNING *)
						SELECT ug.name,ag.name from deleted d 
						JOIN groups ug ON d.usergroup_id=ug.id
						JOIN groups ag ON d.servicegroup_id=ag.id
						
`, mapID, orgID).Scan(&userGroupName, &appGroupName)

	return appGroupName, userGroupName, err
}

func (s accessMapStore) GetUserGroupsToAddInServiceGroup(orgID string) ([]models.Group, error) {
	var userGroups = make([]models.Group, 0)
	rows, err := s.DB.Query(`select id, name FROM groups where org_id=$1 AND type=$2;`, orgID, "usergroup")
	if err != nil {
		return userGroups, err
	}
	for rows.Next() {
		var userGroup models.Group
		err := rows.Scan(&userGroup.GroupID, &userGroup.GroupName)
		if err != nil {
			return userGroups, err
		}

		userGroups = append(userGroups, userGroup)
	}
	return userGroups, err
}

func (s accessMapStore) GetAssignedUserGroupsWithPolicies(groupID, orgID string) ([]UserGroupOfServiceGroup, error) {
	var userGroups = make([]UserGroupOfServiceGroup, 0)

	rows, err := s.DB.Query(`
							select ag_ug.id, usergroup_id, g.name , ag_ug.policy_id ,p.name, ag_ug.privilege, ag_ug.created_at FROM 
							usergroup_accessmaps ag_ug
									 join groups g on ag_ug.usergroup_id=g.id
									 join policies p on ag_ug.policy_id=p.id 
							where ag_ug.org_id=$1 AND ag_ug.servicegroup_id=$2;`,
		orgID, groupID)
	if err != nil {
		return userGroups, err
	}
	for rows.Next() {
		var userGroup UserGroupOfServiceGroup
		err := rows.Scan(&userGroup.MapID, &userGroup.UsergroupID, &userGroup.UsergroupName, &userGroup.PolicyID, &userGroup.PolicyName, &userGroup.Privilege, &userGroup.AddedAt)
		if err != nil {
			return userGroups, err
		}

		userGroups = append(userGroups, userGroup)
	}
	return userGroups, err
}

func (s accessMapStore) GetDynamicAccessPolicy(groups []string, userID, orgID string) (*models.Policy, error) {
	policy := &models.Policy{}
	var tempDayTimeStr string

	err := s.DB.QueryRow(`										
				SELECT p.id, p.name, p.org_id, day_time,risk_threshold, tfa_enabled,file_transfer,record_session,ip_source,device_policy,expiry,allowed_countries , p.created_at, p.updated_at
from dynamic_access
         JOIN policies p on dynamic_access.policy_id = p.id
         LEFT JOIN groups g on dynamic_access.group_name = g.name
         LEFT JOIN user_group_maps ugm on ugm.group_id=g.id
WHERE (dynamic_access.group_name = ANY ($1) OR ugm.user_id=$2) AND  dynamic_access.org_id=$3 ORDER BY dynamic_access.created_at DESC LIMIT 1;`,
		pq.Array(groups), userID, orgID).
		Scan(&policy.PolicyID, &policy.PolicyName, &policy.OrgID, &tempDayTimeStr, &policy.RiskThreshold, &policy.TfaRequired, &policy.FileTransfer, &policy.RecordSession, &policy.IPSource, &policy.DevicePolicy, &policy.Expiry, &policy.AllowedCountries, &policy.CreatedAt, &policy.UpdatedAt)
	if err != nil {
		return policy, err
	}

	err = json.Unmarshal([]byte(tempDayTimeStr), &policy.DayAndTime)

	return policy, err

}

//CreateDynamicAccessRule creates a dynamic access policy for a particular group or idp
func (s accessMapStore) CreateDynamicAccessRule(setting models.DynamicAccessRule) error {
	_, err := s.DB.Exec(`INSERT INTO dynamic_access (id,org_id, group_name, policy_id, created_at) VALUES ($1,$2,$3,$4,$5)`,
		setting.RuleID, setting.OrgID, setting.GroupName, setting.PolicyID, setting.CreatedAt)
	return err
}

//DeleteDynamicAccessRule delete  a dynamic access policy for a particular group or idp
func (s accessMapStore) DeleteDynamicAccessRule(id, orgID string) error {
	_, err := s.DB.Exec(`DELETE FROM dynamic_access WHERE id=$1 AND org_id=$2`,
		id, orgID)
	return err
}

//GetAllDynamicAccessRules returns dynamic access policies for all particular groups or idps
func (s accessMapStore) GetAllDynamicAccessRules(orgID string) ([]models.DynamicAccessRule, error) {
	var das []models.DynamicAccessRule = make([]models.DynamicAccessRule, 0)
	rows, err := s.DB.Query(`
				SELECT dynamic_access.id,dynamic_access.org_id, group_name, policy_id,p.name, dynamic_access.created_at 
				FROM dynamic_access
				JOIN policies p on dynamic_access.policy_id = p.id
				WHERE dynamic_access.org_id=$1`,
		orgID)
	if err != nil {
		return das, err
	}

	for rows.Next() {
		var da models.DynamicAccessRule
		err = rows.Scan(&da.RuleID, &da.OrgID, &da.GroupName, &da.PolicyID, &da.PolicyName, &da.CreatedAt)
		if err != nil {
			logrus.Error(err)
			continue
		}
		das = append(das, da)
	}
	return das, nil
}

//getAllUserGroupsWithIDPs returns all user groups and idp names
func (s accessMapStore) getAllUserGroupsWithIDPs(orgID string) ([]string, error) {

	var gps []string = make([]string, 0)

	rows, err := s.DB.Query(`
select DISTINCT name from (
	select name,org_id from groups where type='usergroup'
	UNION
	select name,org_id from idp
                              ) gps
WHERE gps.org_id=$1;`,
		orgID)
	if err != nil {
		return gps, err
	}

	for rows.Next() {
		var gp string
		err = rows.Scan(&gp)
		if err != nil {
			logrus.Error(err)
			continue
		}
		gps = append(gps, gp)
	}
	return gps, nil
}

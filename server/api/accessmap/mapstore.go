package accessmap

import (
	"github.com/seknox/trasa/server/models"
)

//TODO handle error
func (s AccessMapStore) CheckIfPrivilegeExist(privilege, orgID, serviceID string) bool {
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

func (s AccessMapStore) GetServiceUserMaps(serviceID, orgID string) (appusers []models.AccessMapDetail, err error) {
	appusers = make([]models.AccessMapDetail, 0)
	rows, err := s.DB.Query(`
SELECT map.id, service_id, map.user_id,users.email, map.privilege, map.policy_id, policies.name, map.added_at 
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
		err = rows.Scan(&appuser.MapID, &appuser.ServiceID, &appuser.UserID, &appuser.Email, &appuser.Privilege, &appuser.Policy.PolicyID, &appuser.Policy.PolicyName, &appuser.UserAddedAt)
		if err != nil {
			return
		}

		appusers = append(appusers, appuser)
	}

	return appusers, err
}

func (s AccessMapStore) UpdateServiceUserMap(mapID, orgID, privilege string) error {
	_, err := s.DB.Exec(`UPDATE user_accessmaps SET privilege = $1 WHERE id = $2 AND org_id = $3;`,
		privilege, mapID, orgID)
	return err
}

func (s AccessMapStore) CreateServiceUserMap(appUser *models.ServiceUserMap) error {
	_, err := s.DB.Exec(`INSERT into user_accessmaps (id, service_id,  org_id, user_id,privilege, policy_id, added_at)
						values($1, $2, $3, $4, $5, $6, $7);`, &appUser.MapID, &appUser.ServiceID, &appUser.OrgID, &appUser.UserID, &appUser.Privilege, &appUser.PolicyID, &appUser.AddedAt)

	return err
}

func (s AccessMapStore) DeleteServiceUserMap(mapID, orgID string) (string, error) {
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
func (s AccessMapStore) CreateServiceGroupUserGroupMap(data *models.ServiceGroupUserGroupMap) error {
	_, err := s.DB.Exec(`INSERT into usergroup_accessmaps (id, org_id, servicegroup_id, map_type, usergroup_id,privilege, policy_id, created_at)
						values($1, $2, $3, $4, $5, $6, $7, $8);`,
		&data.MapID, &data.OrgID, &data.ServiceGroupID, &data.MapType, &data.UserGroupID, &data.
			Privilege, &data.PolicyID, &data.CreatedAt)

	return err
}

func (s AccessMapStore) UpdateServiceGroupUserGroupMap(mapID, orgID, privilege string) error {
	_, err := s.DB.Exec(`UPDATE usergroup_accessmaps SET privilege = $1 WHERE id = $2 AND org_id = $3;`,
		privilege, mapID, orgID)

	return err
}

func (s AccessMapStore) DeleteServiceGroupUserGroupMap(mapID, orgID string) (string, string, error) {
	var appGroupName, userGroupName string
	//TODO check query
	err := s.DB.QueryRow(`WITH deleted as (DELETE FROM usergroup_accessmaps WHERE id = $1 AND org_id=$2 RETURNING *)
						SELECT ug.name,ag.name from deleted d 
						JOIN groups ug ON d.usergroup_id=ug.id
						JOIN groups ag ON d.servicegroup_id=ag.id
						
`, mapID, orgID).Scan(&userGroupName, &appGroupName)

	return appGroupName, userGroupName, err
}

func (s AccessMapStore) GetUserGroupsToAddInServiceGroup(orgID string) ([]models.Group, error) {
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

func (s AccessMapStore) GetAssignedUserGroupsWithPolicies(groupID, orgID string) ([]userGroupOfServiceGroup, error) {
	var userGroups = make([]userGroupOfServiceGroup, 0)

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
		var userGroup userGroupOfServiceGroup
		err := rows.Scan(&userGroup.MapID, &userGroup.UsergroupID, &userGroup.UsergroupName, &userGroup.PolicyID, &userGroup.PolicyName, &userGroup.Privilege, &userGroup.AddedAt)
		if err != nil {
			return userGroups, err
		}

		userGroups = append(userGroups, userGroup)
	}
	return userGroups, err
}

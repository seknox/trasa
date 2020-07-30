package groups

import (
	"time"

	"github.com/seknox/trasa/models"
	"github.com/seknox/trasa/utils"
)

func (s GroupStore) CheckIfUserInGroup(userID, orgID string, groupIDs []string) (bool, error) {
	var exists bool
	for _, groupID := range groupIDs {
		err := s.DB.QueryRow(`SELECT EXISTS (select user_id from user_group_maps where user_id=$1 AND group_id=$2 AND org_id=$3)`,
			userID,
			groupID,
			orgID,
		).Scan(&exists)
		if err != nil {
			return false, err
		}
		if exists {
			return true, nil
		}

	}
	return false, nil

}

// GGet users that are already assigned to this group.
func (s GroupStore) GetUsersInGroup(groupID, org string) ([]models.User, error) {
	var users = make([]models.User, 0)
	var user models.User
	rows, err := s.DB.Query(`SELECT users.id, username, first_name, last_name, email, idp_name, user_role,users.status, users.created_at, users.updated_at
 										FROM users
 										 join user_group_maps map on users.id=map.user_id 
 										 WHERE  group_id=$1 AND map.org_id=$2`, groupID, org)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.IdpName, &user.UserRole, &user.Status, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, err

}

func (s GroupStore) GetUsersNotInGroup(groupID, orgID string) ([]models.User, error) {
	var users = make([]models.User, 0)
	rows, err := s.DB.Query(`select id, username, email,first_name, last_name FROM users  
										where org_id=$1 AND id not in
										 	(select id from user_group_maps where group_id=$2 AND org_id=$3);`,
		orgID, groupID, orgID)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.UserName, &user.Email, &user.FirstName, &user.LastName)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}
	return users, err
}

// GetApps returns apps available for organizations.
func (s GroupStore) GetServicesInGroup(groupID, org string) ([]models.Service, error) {
	var services = make([]models.Service, 0)

	rows, err := s.DB.Query(`SELECT services.id, name, adhoc,hostname, type
 								FROM services 
 								join service_group_maps map on services.id=map.service_id 
 								WHERE  group_id=$1 AND services.org_id=$2`,
		groupID, org)
	if err != nil {
		return services, err
	}
	defer rows.Close()
	for rows.Next() {
		var service models.Service
		err := rows.Scan(&service.ID, &service.Name, &service.Adhoc, &service.Hostname, &service.Type)
		if err != nil {
			return services, err
		}
		services = append(services, service)
	}
	return services, err
}

func (s GroupStore) GetServicesNotInGroup(groupID, orgID string) ([]models.Service, error) {
	var services = make([]models.Service, 0)
	rows, err := s.DB.Query(`select id, name 
										from services where id not in (
													select service_id from service_group_maps where group_id=$1 AND org_id=$2
													) AND org_id=$2;`,
		groupID, orgID)
	if err != nil {
		return services, err
	}
	for rows.Next() {
		var service models.Service
		err = rows.Scan(&service.ID, &service.Name)
		if err != nil {
			return services, err
		}
		services = append(services, service)
	}
	return services, nil
}

func (s GroupStore) AddServicesToGroup(group models.Group, serviceIDs []string) (err error) {
	for i := 0; i < len(serviceIDs); i++ {
		mapId := utils.GetUUID()
		createdAt := time.Now().Unix()
		updatedAt := createdAt
		_, err = s.DB.Exec(`INSERT into service_group_maps (id, group_id, org_id, service_id, status, created_at,updated_at)
		values($1, $2, $3, $4, $5, $6, $7);`,
			mapId, &group.GroupID, &group.OrgID, serviceIDs[i], &group.Status, createdAt, updatedAt)
		if err != nil {
			return err
		}
	}

	return err
}

func (s GroupStore) RemoveServicesFromGroup(groupID, orgID string, serviceIDs []string) error {
	var err error
	for i := 0; i < len(serviceIDs); i++ {
		// delete the user from group
		_, err = s.DB.Exec(`DELETE FROM service_group_maps WHERE service_id = $1 AND group_id=$2 AND org_id=$3;`,
			serviceIDs[i], groupID, orgID)
		if err != nil {
			return err
		}

	}

	return err
}

//AddUsersToGroup adds users to user-group
func (s GroupStore) AddUsersToGroup(group models.Group, userIDs []string) (err error) {

	for i := 0; i < len(userIDs); i++ {
		mapId := utils.GetUUID()
		createdAt := time.Now().Unix()
		updateAt := createdAt
		_, err = s.DB.Exec(`INSERT into user_group_maps (id, group_id, org_id, user_id, status, created_at,updated_at)
		values($1, $2, $3, $4, $5, $6, $7);`, mapId, &group.GroupID, &group.OrgID, userIDs[i], &group.Status, createdAt, updateAt)
		if err != nil {
			return err
		}
	}

	return err
}

// RemoveUsersFromGroup should remove users from group and also corresponding authorization from appusers table.
// service usertable should store user group id if the assignment was initiated from within group.
// group_id in appusersv1 can be either group id or "no-group" value.
func (s GroupStore) RemoveUsersFromGroup(groupID, orgID string, userIDs []string) (err error) {
	for i := 0; i < len(userIDs); i++ {
		// delete the user from group
		_, err = s.DB.Exec(`DELETE FROM user_group_maps WHERE user_id = $1 AND group_id=$2 AND org_id=$3;`,
			userIDs[i], groupID, orgID)
		if err != nil {
			return err
		}

		// Commenting out these line. why do we need this? @sshahcodes 27/09/2019
		// // also delete user from appusersv1 table which references this groupid
		// _, err = conn.db.Exec(`DELETE FROM appusersv1 WHERE user_id = $1 AND group_id=$2 AND org_id=$3 RETURNING *;`, users[i], groupID, orgID)
		// if err != nil {
		// 	return err
		// 	//logger.Debug(err)
		//utils.Do.SystemLogger(consts.DatabaseError, false, err.Error(), "RemoveUsersFromGroup-appusersv1", "error")
		// }
	}

	return err
}

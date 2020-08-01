package groups

import (
	"github.com/seknox/trasa/server/models"
)

func (s GroupStore) Create(group *models.Group) error {
	_, err := s.DB.Exec(`INSERT into groups (id, org_id, name,type,  status, created_at,updated_at)
						values($1, $2, $3, $4, $5, $6, $7);`,
		&group.GroupID, &group.OrgID, &group.GroupName, &group.GroupType, &group.Status, &group.CreatedAt, &group.UpdatedAt)

	return err
}

func (s GroupStore) Get(groupID, orgID string) (models.Group, error) {
	var group models.Group
	err := s.DB.QueryRow(`SELECT id, org_id, name, type, status FROM groups WHERE id=$1 AND org_id=$2`,
		groupID, orgID).
		Scan(&group.GroupID, &group.OrgID, &group.GroupName, &group.GroupType, &group.Status)

	return group, err
}

func (s GroupStore) GetAllUserGroups(orgID string) ([]models.Group, error) {

	groups := make([]models.Group, 0)

	rows, err := s.DB.Query(`SELECT id,name,status,created_at,updated_at FROM groups WHERE type=$1 AND org_id=$2`,
		"usergroup", orgID)
	if err != nil {
		return groups, err
	}
	for rows.Next() {
		var temp models.Group
		err = rows.Scan(&temp.GroupID, &temp.GroupName, &temp.Status, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return groups, err
		}
		//TODO @bhrg3se
		s.DB.QueryRow(`SELECT count(user_id) FROM user_group_maps WHERE group_id=$1`, temp.GroupID).Scan(&temp.MemberCount)
		groups = append(groups, temp)
	}

	return groups, nil
}

func (s GroupStore) GetAllServiceGroups(orgID string) ([]models.Group, error) {

	groups := make([]models.Group, 0)

	rows, err := s.DB.Query(`SELECT id,name,status,created_at,updated_at FROM groups WHERE type=$1 AND org_id=$2`,
		"servicegroup", orgID)

	if err != nil {
		return groups, err
	}

	for rows.Next() {
		var temp models.Group
		err = rows.Scan(&temp.GroupID, &temp.GroupName, &temp.Status, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return groups, err
		}

		//TODO @bhrg3se
		s.DB.QueryRow(`SELECT count(service_id) FROM service_group_maps WHERE group_id=$1`, temp.GroupID).Scan(&temp.MemberCount)
		groups = append(groups, temp)
	}

	return groups, nil
}

func (s GroupStore) Update(group *models.Group) error {
	_, err := s.DB.Exec(`UPDATE groups set name=$1, updated_at=$2 where id=$3 AND org_id=$4;`,
		&group.GroupName, &group.UpdatedAt, &group.GroupID, &group.OrgID)

	return err
}

// Delete  deletes a row from group table and delete every element from GroupMap table
func (s GroupStore) Delete(groupID, orgID string) (name string, err error) {

	err = s.DB.QueryRow(`DELETE FROM groups WHERE id=$1 AND org_id=$2 RETURNING name;`,
		groupID, orgID).
		Scan(&name)

	return name, err
}

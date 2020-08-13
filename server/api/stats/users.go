package stats

func (s statStore) GetTotalAdmins(orgID string) (int, error) {
	var total int
	err := s.DB.QueryRow(`SELECT count(*)from users where user_role=$1 AND org_id=$2`, "orgAdmin", orgID).Scan(&total)

	return total, err
}

func (s statStore) GetTotalDisabledUsers(orgID string) (int, error) {
	var total int
	err := s.DB.QueryRow(`SELECT count(*)from users where status=$1 AND org_id=$2`, false, orgID).Scan(&total)

	return total, err
}

func (s statStore) GetAggregatedIdpUsers(entityType, entityID, orgID string) (users totalUsers, err error) {
	users.Users = 0
	rows, err := s.DB.Query(`select count(*) as c,idp_name from users WHERE org_id=$1 GROUP BY idp_name ORDER BY c ASC`, orgID)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var idpUsers idpUsers
		err = rows.Scan(&idpUsers.Value, &idpUsers.Name)
		if err != nil {
			return users, err
		}
		users.Users += idpUsers.Value
		users.Idps = append(users.Idps, idpUsers)

	}
	return users, err
}

package stats

func (s StatStore) GetPoliciesStats(orgID string) (stat policyStat, err error) {
	//TODO check timezone
	rows, err := s.DB.Query(`select count(*) as c , expiry::timestamp < now() as is_expired from policies where org_id=$1 GROUP BY   expiry::timestamp < now();`, orgID)
	if err != nil {
		return stat, err
	}
	for rows.Next() {
		var c int
		var isExpired bool
		err = rows.Scan(&c, &isExpired)
		if err != nil {
			return
		}
		if isExpired {
			stat.Expired = c
		}
		stat.Total = stat.Total + c

	}

	return stat, err
}

type policyStat struct {
	Total   int `json:"total"`
	Expired int `json:"expired"`
}

func (s StatStore) GetTotalGroups(orgID, groupType string) (int, error) {
	var total int
	err := s.DB.QueryRow(`SELECT count(*)from groups where type=$1 AND org_id=$2`, groupType, orgID).Scan(&total)

	return total, err
}

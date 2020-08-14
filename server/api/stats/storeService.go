package stats

import (
	"database/sql"

	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/utils"
)

func (s statStore) GetTotalServices(orgID string) (count int64, err error) {
	err = s.DB.QueryRow(`select COUNT(*) from services where org_id=$1;`, orgID).Scan(&count)
	return
}

func (s statStore) GetTotalManagedUsers(entityType, entityID, orgID string) (count int, err error) {
	count = 0
	var managedAcc string
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(`managed_accounts`)
	sb.From("services")

	//TODO does user id makes sense here??
	if entityType == "service" {
		sb.Where(sb.Equal("id", entityID))
	} else if entityType == "user" {
		sb.Where(sb.Equal("user_id", entityID))
	}
	sb.Where(sb.Equal(`org_id`, orgID))
	sqlStr, args := sb.Build()
	//change ? into $
	sqlStr = utils.SqlReplacer(sqlStr)
	err = s.DB.QueryRow(sqlStr, args...).Scan(&managedAcc)

	return count, err
}

func (s statStore) GetPoliciesOfService(serviceID, orgID string) (count int, err error) {
	count = 0

	//TODO check this query
	err = s.DB.QueryRow(`SELECT count(DISTINCT policy_id) as count from
        (
            select policy_id FROM user_accessmaps where service_id=$2 AND  org_id=$1
            UNION ALL
            (SELECT policy_id FROM usergroup_accessmaps
             	JOIN service_group_maps ON usergroup_accessmaps.servicegroup_id=service_group_maps.group_id
            	JOIN services ON service_group_maps.service_id = services.id
             	where services.id=$2 AND map_type='servicegroup' AND  services.org_id=$1
             	)
             	UNION ALL (
             	SELECT policy_id FROM usergroup_accessmaps 
             	where servicegroup_id=$2 AND map_type='service' AND org_id=$1
             	)
        ) as pol ;`, orgID, serviceID).Scan(&count)

	return count, err
}

func (s statStore) GetTotalPrivilegesOfService(serviceID, orgID string) (count int, err error) {
	count = 0
	err = s.DB.QueryRow(`SELECT count(DISTINCT privilege) as count from
        (
            select privilege FROM user_accessmaps where service_id=$2 AND org_id=$1
            UNION ALL
            (SELECT privilege FROM usergroup_accessmaps
             	JOIN services ON usergroup_accessmaps.servicegroup_id=services.id 
             	where services.id=$2 AND map_type='servicegroup' AND services.org_id=$1
             	)
             	UNION ALL (
             	SELECT usergroup_accessmaps.privilege FROM usergroup_accessmaps 
             	where servicegroup_id=$2 AND map_type='service' AND org_id=$1
             	)
        ) privs  ;`, orgID, serviceID).Scan(&count)

	return count, err
}

func (s statStore) GetTotalGroupsServiceIsAssignedTo(serviceID, orgID string) (count int, err error) {
	count = 0
	err = s.DB.QueryRow(`SELECT count(DISTINCT group_id) as count  FROM service_group_maps 
             	where service_id=$2  AND service_group_maps.org_id=$1 ;`, orgID, serviceID).Scan(&count)

	return count, err
}

//func (s statStore) GetTotalGroupsServiceIsAssignedTo(serviceID, orgID string) (count int, err error) {
//	count = 0
//	err = s.DB.QueryRow(`SELECT count(DISTINCT appgroup_id) as count
// FROM appgroup_usergroup_mapv1
//             	JOIN servicesv1 ON appgroup_usergroup_mapv1.appgroup_id=servicesv1.service_id
//             	where service_id=$2 AND appgroup_type='appgroup' AND servicesv1.org_id=$1 ;`, orgID, serviceID).Scan(&count)
//
//	return count, err
//}

func (s statStore) GetTotalUsersAssignedToService(serviceID, orgID string) (count int, err error) {
	count = 0
	err = s.DB.QueryRow(`
SELECT count (DISTINCT user_id)
FROM(
	SELECT user_id from user_accessmaps where service_id=$1 AND org_id=$2
UNION 
	(
	SELECT user_id from usergroup_accessmaps 
	JOIN user_group_maps ON usergroup_accessmaps.usergroup_id=user_group_maps.group_id where usergroup_accessmaps.servicegroup_id=$1 AND usergroup_accessmaps.org_id=$2 AND map_type='service'
	)
UNION 
	(
	SELECT user_id from usergroup_accessmaps 
	JOIN user_group_maps ON usergroup_accessmaps.usergroup_id=user_group_maps.group_id 
	JOIN service_group_maps ON usergroup_accessmaps.servicegroup_id=service_group_maps.group_id where service_group_maps.service_id=$1 AND  service_group_maps.org_id=$2 AND map_type='servicegroup'
	)
) as usserassigned
 `, serviceID, orgID).Scan(&count)

	return count, err
}

func (s statStore) GetAggregatedServices(orgID string) (apps allServices, err error) {
	apps.ServicesByType = []nameValue{}
	apps.TotalServices = 0
	rows, err := s.DB.Query(`select count(*) as c,type from services WHERE org_id=$1 GROUP BY type ORDER BY c ASC`, orgID)
	if err != nil {
		return apps, err
	}
	for rows.Next() {
		var service nameValue
		err = rows.Scan(&service.Value, &service.Name)
		if err != nil {
			return apps, err
		}
		apps.TotalServices += service.Value
		apps.ServicesByType = append(apps.ServicesByType, service)
	}

	err = s.DB.QueryRow(`select count(*) as appgroup_count from groups where type='servicegroup' AND org_id=$1`, orgID).Scan(&apps.TotalGroups)
	if err != nil && (!errors.Is(err, sql.ErrNoRows)) {
		return apps, err
	}

	err = s.DB.QueryRow(`select count(*) as session_disabled from services where session_record=false AND org_id=$1;`, orgID).Scan(&apps.SessionRecordingDisabledCount)
	if err != nil && (!errors.Is(err, sql.ErrNoRows)) {
		return apps, err
	}

	//TODO
	//sett, err := s.GetGlobalServiceSettings(orgID)
	//if err != nil {
	//	logrus.Error(err)
	//}
	//apps.DynamicService = sett.Status

	return apps, err
}

func (s statStore) GetAggregatedIDPServices(idpName, orgID string) ([]nameValue, error) {
	appsByType := []nameValue{}

	rows, err := s.DB.Query(`select count(*) as c,type from services WHERE org_id=$1 AND external_provider_name=$2 GROUP BY type ORDER BY c ASC`, orgID, idpName)
	if err != nil {
		return appsByType, err
	}
	for rows.Next() {
		var service nameValue
		err = rows.Scan(&service.Value, &service.Name)
		if err != nil {
			return appsByType, err
		}
		appsByType = append(appsByType, service)
	}

	return appsByType, err
}

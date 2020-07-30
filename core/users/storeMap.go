package users

import (
	"database/sql"

	"github.com/seknox/trasa/models"
)

func (s UserStore) GetAccessMapDetails(userID, orgID string) ([]models.AccessMapDetail, error) {
	var accessMapDetail models.AccessMapDetail
	var mapDetails = make([]models.AccessMapDetail, 0)
	rows, err := s.DB.Query(`
SELECT DISTINCT services.org_id,services.name,services.type,services.hostname, services.id,gappusersv1.privilege, user_id, gappusersv1.policy_name ,gappusersv1.policy_id,record_session,file_transfer, services.created_at
		 
		FROM (
			--user assigned to service 	
		SELECT service_id,user_id,p.name as policy_name,privilege,p.id as policy_id,p.record_session as record_session,p.file_transfer as file_transfer  from 
			user_accessmaps
			JOIN policies p on user_accessmaps.policy_id=p.id  
			
			--usergroup assigned to app
		UNION SELECT servicegroup_id as service_id,user_id,p.name as policy_name,privilege,p.id as policy_id,p.record_session as record_session,p.file_transfer as file_transfer FROM
				(
				usergroup_accessmaps ag_ug 
				join policies p on ag_ug.policy_id=p.id 
				join user_group_maps ug on ug.group_id=ag_ug.usergroup_id 
				 
				) WHERE ag_ug.map_type='service'
				
			--usergroup assigned to appgroup				 
		UNION 	SELECT ag.service_id,user_id,p.name as policy_name,ag_ug.privilege,p.id as policy_id,p.record_session as record_session,p.file_transfer as file_transfer FROM
				(
				usergroup_accessmaps ag_ug 
				join policies p on ag_ug.policy_id=p.id 
				join user_group_maps ug on ug.group_id=ag_ug.usergroup_id 
				join service_group_maps ag on ag.group_id=ag_ug.servicegroup_id 
				) WHERE ag_ug.map_type='servicegroup'

		
		) as gappusersv1 
		 JOIN services ON  gappusersv1.service_id=services.id   --TODO why is this extra join necessary? 
		 
		WHERE gappusersv1.user_id=$1 AND services.org_id=$2;
`, userID, orgID)
	if err != nil {
		return mapDetails, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&accessMapDetail.OrgID,
			&accessMapDetail.ServiceName,
			&accessMapDetail.ServiceType,
			&accessMapDetail.Hostname,
			&accessMapDetail.ServiceID,
			&accessMapDetail.Privilege,
			&accessMapDetail.UserID,
			&accessMapDetail.Policy.PolicyName,
			&accessMapDetail.Policy.PolicyID,
			&accessMapDetail.Policy.RecordSession,
			&accessMapDetail.Policy.FileTransfer,
			&accessMapDetail.UserAddedAt)
		if err != nil {
			return mapDetails, err
		}

		mapDetails = append(mapDetails, accessMapDetail)
	}

	return mapDetails, err
}

//GetServicesAssigned returns Apps assigned to a user with the policy
func (s UserStore) GetAssignedServices(userID, orgID string) (services []models.Service, err error) {

	services = make([]models.Service, 0)
	var rows *sql.Rows

	rows, err = s.DB.Query(`

SELECT service_id FROM
(
--users directly assigned to apps (appusers) 
(SELECT service_id,user_id,org_id FROM user_accessmaps )

UNION 
--usergroup assigned to app
(SELECT servicegroup_id as service_id,user_id,ag_ug.org_id FROM usergroup_accessmaps ag_ug  
		join user_group_maps ug on ug.group_id=ag_ug.usergroup_id 		
	WHERE map_type='service'
		)
UNION 
--usergroup assigned to appgroup
(SELECT ag.service_id,user_id,ag_ug.org_id FROM usergroup_accessmaps ag_ug  
		join user_group_maps ug on ug.group_id=ag_ug.usergroup_id
		join service_group_maps ag on ag.group_id=ag_ug.servicegroup_id
	WHERE map_type='servicegroup'
) 
		 )WHERE user_id=$1 AND org_id=$2;`, userID, orgID)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var s models.Service
		err := rows.Scan(&s.ID)
		if err != nil {
			return services, err
		}
		services = append(services, s)
	}

	return

}

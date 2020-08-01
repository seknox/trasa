package policies

import (
	"encoding/json"

	"github.com/seknox/trasa/server/models"
)

func (s PolicyStore) GetPolicy(policyID, orgID string) (models.Policy, error) {
	var policy models.Policy
	var tempDayTimeStr string
	err := s.DB.QueryRow(`SELECT id, name, org_id, day_time,risk_threshold, tfa_enabled,file_transfer,record_session,ip_source,device_policy,expiry,allowed_countries , created_at, updated_at from policies WHERE id=$1 AND org_id=$2`, policyID, orgID).
		Scan(&policy.PolicyID, &policy.PolicyName, &policy.OrgID, &tempDayTimeStr, &policy.RiskThreshold, &policy.TfaRequired, &policy.FileTransfer, &policy.RecordSession, &policy.IPSource, &policy.DevicePolicy, &policy.Expiry, &policy.AllowedCountries, &policy.CreatedAt, &policy.UpdatedAt)
	if err != nil {
		return policy, err
	}

	err = json.Unmarshal([]byte(tempDayTimeStr), &policy.DayAndTime)

	return policy, err
}

// CreatePolicy in database
func (s PolicyStore) CreatePolicy(policy models.Policy) error {

	dayTime, err := json.Marshal(policy.DayAndTime)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(`INSERT into policies (id, name, org_id, day_time,risk_threshold, tfa_enabled,file_transfer,record_session,ip_source,device_policy,expiry,allowed_countries, created_at, updated_at)
	values($1, $2, $3, $4, $5, $6, $7, $8,$9,$10,$11,$12, $13,$14);`,
		policy.PolicyID, policy.PolicyName, policy.OrgID, string(dayTime), policy.RiskThreshold, policy.TfaRequired, policy.FileTransfer, policy.RecordSession, policy.IPSource, policy.DevicePolicy, policy.Expiry, policy.AllowedCountries, policy.CreatedAt, policy.UpdatedAt)

	return err
}

func (s PolicyStore) UpdatePolicy(policy models.Policy) error {
	day_time, err := json.Marshal(policy.DayAndTime)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(`Update  policies set id=$1, name=$2, org_id=$3, day_time=$4,risk_threshold=$5, tfa_enabled=$6,file_transfer=$7,ip_source=$8,device_policy=$9,expiry=$10,allowed_countries=$11, created_at=$12, updated_at=$13,record_session=$14
	Where org_id=$3 AND id=$1;`, policy.PolicyID, policy.PolicyName, policy.OrgID, string(day_time), policy.RiskThreshold, policy.TfaRequired, policy.FileTransfer, policy.IPSource, policy.DevicePolicy, policy.Expiry, policy.AllowedCountries, policy.CreatedAt, policy.UpdatedAt, policy.RecordSession)

	return err
}

func (s PolicyStore) GetAllPolicies(orgID string) ([]models.Policy, error) {
	policies := []models.Policy{}
	rows, err := s.DB.Query(`
SELECT coalesce(count,0) used_by,
       policies.id, name, org_id, day_time,risk_threshold, tfa_enabled,file_transfer,record_session,ip_source,device_policy,expiry,allowed_countries, created_at, updated_at
from policies
         LEFT JOIN (
    SELECT count(policy_id) as count,policy_id from
        (
            select policy_id FROM user_accessmaps 

            UNION ALL
            SELECT policy_id FROM usergroup_accessmaps
        ) as p

    GROUP BY policy_id

) AS policy_count ON policies.id=policy_count.policy_id WHERE org_id=$1
`, orgID)
	if err != nil {
		return policies, err
	}

	for rows.Next() {
		var policy models.Policy
		var tempDayTimeStr string
		err = rows.Scan(&policy.UsedBy, &policy.PolicyID, &policy.PolicyName, &policy.OrgID, &tempDayTimeStr, &policy.RiskThreshold, &policy.TfaRequired, &policy.FileTransfer, &policy.RecordSession, &policy.IPSource, &policy.DevicePolicy, &policy.Expiry, &policy.AllowedCountries, &policy.CreatedAt, &policy.UpdatedAt)
		if err != nil {
			return policies, nil
		}
		err = json.Unmarshal([]byte(tempDayTimeStr), &policy.DayAndTime)
		if err != nil {
			return policies, nil
		}
		policies = append(policies, policy)
	}

	return policies, err
}

func (s PolicyStore) DeletePolicy(policyID, orgID string) error {
	_, err := s.DB.Exec(`DELETE FROM policies WHERE id = $1 AND org_id=$2 RETURNING *`, policyID, orgID)

	return err
}

func (s PolicyStore) GetAccessPolicy(userID, serviceID, orgID string) (policy *models.Policy, privilege string, adhoc bool, err error) {
	policy, privilege, adhoc, err = s.getUserAccessPolicy(userID, serviceID, orgID)
	if err != nil {
		policy, privilege, adhoc, err = s.getUserGroupAccessPolicy(userID, serviceID, orgID)
		if err != nil {
			policy, privilege, adhoc, err = s.getServiceUserGroupAccessPolicy(userID, serviceID, orgID)
			if err != nil {
				return
			}
		}
	}
	return
}

func (s PolicyStore) getUserAccessPolicy(userID, serviceID, orgID string) (policy *models.Policy, privilege string, adhoc bool, err error) {

	policy = &models.Policy{}
	var day_time string

	err = s.DB.QueryRow(`
	SELECT policies.day_time,
			adhoc,
			policies.tfa_enabled,
			map.privilege,
			policies.record_session,
			policies.file_transfer,
			policies.ip_source,
			policies.risk_threshold,
			policies.expiry,
			policies.allowed_countries ,
			policies.device_policy                                   

FROM user_accessmaps map
JOIN policies ON map.policy_id=policies.id AND map.org_id=policies.org_id                             
JOIN services ON map.service_id=services.id AND map.org_id=services.org_id
WHERE map.user_id= $1 AND map.service_id= $2 AND map.org_id=$3;`, userID, serviceID, orgID).
		Scan(&day_time, &adhoc, &policy.TfaRequired, &privilege, &policy.RecordSession, &policy.FileTransfer, &policy.IPSource, &policy.RiskThreshold, &policy.Expiry, &policy.AllowedCountries, &policy.DevicePolicy)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(day_time), &policy.DayAndTime)

	return
}

//get permission of usergroup assigned to app
func (s PolicyStore) getUserGroupAccessPolicy(userID, serviceID, orgID string) (policy *models.Policy, privilege string, adhoc bool, err error) {
	policy = &models.Policy{}

	var day_time string

	//In this case appgroup_id of appgroup_usergroup_mapv1 table is used as service_id
	//because usergroup is assigned to single Service
	err = s.DB.QueryRow(`
			SELECT gappusersv1.day_time,
				   adhoc,
				   gappusersv1.tfa_enabled,
			       gappusersv1.privilege,
			       gappusersv1.record_session,
			       gappusersv1.file_transfer,
			       gappusersv1.ip_source,
			       gappusersv1.risk_threshold,
			       gappusersv1.expiry,
			       gappusersv1.allowed_countries,
			       gappusersv1.device_policy
			FROM (usergroup_accessmaps ag_ug 
					join policies p on ag_ug.policy_id=p.id
					join user_group_maps ug on ug.group_id=ag_ug.usergroup_id  
					) as gappusersv1  
						 JOIN services ON gappusersv1.servicegroup_id=services.id AND gappusersv1.org_id=services.org_id 
					WHERE gappusersv1.user_id= $1 AND gappusersv1.servicegroup_id= $2 AND gappusersv1.org_id=$3  AND gappusersv1.map_type='service';`, userID, serviceID, orgID).
		Scan(&day_time, &adhoc, &policy.TfaRequired, &privilege, &policy.RecordSession, &policy.FileTransfer, &policy.IPSource, &policy.RiskThreshold, &policy.Expiry, &policy.AllowedCountries, &policy.DevicePolicy)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(day_time), &policy.DayAndTime)

	return
}

//get permission of usergroup assigned to appgroup
func (s PolicyStore) getServiceUserGroupAccessPolicy(userID, serviceID, orgID string) (policy *models.Policy, privilege string, adhoc bool, err error) {
	policy = &models.Policy{}

	var day_time string

	err = s.DB.QueryRow(`SELECT gappusersv1.day_time,adhoc,gappusersv1.tfa_enabled,gappusersv1.privilege,gappusersv1.record_session,gappusersv1.file_transfer,gappusersv1.ip_source,gappusersv1.risk_threshold,gappusersv1.expiry,gappusersv1.allowed_countries,gappusersv1.device_policy 
FROM (usergroup_accessmaps ag_ug 
		join policies p on ag_ug.policy_id=p.id 
		join user_group_maps ug on ug.group_id=ag_ug.usergroup_id 
		join service_group_maps ag on ag.group_id=ag_ug.servicegroup_id 
		) as gappusersv1  
			 JOIN services ON gappusersv1.service_id=services.id AND gappusersv1.org_id=services.org_id
		WHERE gappusersv1.user_id= $1 AND gappusersv1.service_id= $2 AND gappusersv1.org_id=$3 AND gappusersv1.map_type='servicegroup';`, userID, serviceID, orgID).
		Scan(&day_time, &adhoc, &policy.TfaRequired, &privilege, &policy.RecordSession, &policy.FileTransfer, &policy.IPSource, &policy.RiskThreshold, &policy.Expiry, &policy.AllowedCountries, &policy.DevicePolicy)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(day_time), &policy.DayAndTime)

	return
}

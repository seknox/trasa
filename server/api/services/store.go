package services

import (
	"database/sql"
	"strings"

	"github.com/seknox/trasa/server/models"
	"github.com/sirupsen/logrus"
)

func (s serviceStore) GetFromID(serviceID string) (service *models.Service, err error) {
	service = &models.Service{}
	err = s.DB.QueryRow(`
SELECT id,org_id, name, secret_key, passthru,adhoc, created_at,hostname, type, managed_accounts,remoteapp_name,native_log,rdp_protocol,proxy_config
 FROM services
 
  WHERE services.id = $1`, serviceID).
		Scan(&service.ID, &service.OrgID, &service.Name, &service.SecretKey, &service.Passthru, &service.Adhoc, &service.CreatedAt, &service.Hostname, &service.Type, &service.ManagedAccounts, &service.RemoteAppName, &service.NativeLog, &service.RdpProtocol, &service.ProxyConfig)
	if err != nil {
		return service, err
	}

	return service, nil
}

func (s serviceStore) GetFromHostname(hostname, ServiceType, remoteAppName, orgID string) (service *models.Service, err error) {
	logrus.Debug(hostname, ServiceType, remoteAppName, orgID)
	service = &models.Service{}
	err = s.DB.QueryRow(`
SELECT org_id,id, name, secret_key, passthru,adhoc, created_at,hostname, type, managed_accounts,remoteapp_name,native_log,rdp_protocol, proxy_config
 FROM services

  WHERE hostname = $1 AND type = $2 AND remoteapp_name=$3 AND org_id = $4`,
		hostname, ServiceType, remoteAppName, orgID).
		Scan(&service.OrgID, &service.ID, &service.Name, &service.SecretKey, &service.Passthru, &service.Adhoc, &service.CreatedAt, &service.Hostname, &service.Type, &service.ManagedAccounts, &service.RemoteAppName, &service.NativeLog, &service.RdpProtocol, &service.ProxyConfig)

	if err != nil {
		return service, err
	}

	return service, nil
}

func (s serviceStore) Create(app *models.Service) (err error) {

	_, err = s.DB.Exec(`INSERT into services (id, org_id, name, secret_key,type, passthru, hostname, created_at, updated_at, deleted_at, managed_accounts,adhoc,remoteapp_name,rdp_protocol,native_log, proxy_config, external_provider_name, external_id, external_security_group, distro_name,distro_version, ip_details )
						 values($1, $2, $3, $4, $5, $6, $7, $8 ,$9, $10,$11, $12, $13,$14,$15,$16, $17, $18, $19, $20, $21, $22);`,
		app.ID, app.OrgID, app.Name, app.SecretKey, app.Type, app.Passthru, app.Hostname, app.CreatedAt, app.UpdatedAt, app.DeletedAt, app.ManagedAccounts, false, "", "nla", false, app.ProxyConfig, app.ExternalProviderName, app.ExternalID, app.ExternalSecurityGroup, app.DistroName, app.DistroVersion, app.IPDetails)

	return err
}

func (s serviceStore) Update(service *models.Service) error {
	_, err := s.DB.Exec(`UPDATE services SET name = $1, type = $2, passthru = $3, updated_at = $4, hostname=$7, adhoc=$8,remoteapp_name=$9,rdp_protocol=$10,native_log=$11 WHERE id = $5 AND org_id = $6;`,
		&service.Name, &service.Type, &service.Passthru, &service.UpdatedAt, &service.ID, &service.OrgID, &service.Hostname, &service.Adhoc, &service.RemoteAppName, &service.RdpProtocol, &service.NativeLog)

	if err != nil {
		return err
	}
	return err
}

func (s serviceStore) updateHttpProxy(serviceID, orgID string, time int64, proxyConfig models.ReverseProxy) error {
	_, err := s.DB.Exec(`UPDATE services SET proxy_config = $1, updated_at=$2 WHERE id = $3 AND org_id = $4;`,
		proxyConfig, time, serviceID, orgID)

	return err
}

//Delete deletes a service and returns its name
func (s serviceStore) Delete(serviceID, orgID string) (serviceName string, err error) {
	err = s.DB.QueryRow(`DELETE FROM services WHERE id = $1 AND org_id=$2 RETURNING name`, serviceID, orgID).Scan(&serviceName)
	if err != nil {
		return "", err
	}
	return serviceName, nil
}

func (s serviceStore) AddManagedAccounts(serviceID, orgID string, username string) error {
	// we first check if similer service already exists in database. This is searched in context of organization
	var managedAccounts string
	err := s.DB.QueryRow("SELECT managed_accounts FROM services WHERE id = $1 AND org_id=$2", serviceID, orgID).Scan(&managedAccounts)
	if err != nil {
		return err
	}

	val := strings.Split(managedAccounts, ",")

	for _, v := range val {
		if v == username {
			return nil
		}
	}

	val = append(val, username)

	valstring := strings.Join(val, ",")

	_, err = s.DB.Exec(`UPDATE services SET managed_accounts = $1  WHERE id = $2 AND org_id=$3 `, valstring, serviceID, orgID)

	return err
}

func (s serviceStore) RemoveManagedAccounts(serviceID, orgID string, username string) error {
	// we first check if similer service already exists in database. This is searched in context of organization
	var managedAccounts string
	err := s.DB.QueryRow("SELECT managed_accounts FROM services WHERE id = $1 AND org_id=$2 ",
		serviceID, orgID).Scan(&managedAccounts)
	if err != nil {
		return err
	}

	val := strings.Split(managedAccounts, ",")

	for i := 0; i < len(val); i++ {
		if val[i] == username {
			val = sliceManagedUsers(i, val)
		}
	}

	valstring := strings.Join(val, ",")

	_, err = s.DB.Exec(`UPDATE services SET managed_accounts = $1  WHERE id = $2 AND org_id=$3 ;`,
		valstring, serviceID, orgID)

	return err
}

func sliceManagedUsers(index int, array []string) []string {
	array[index] = array[len(array)-1]
	array[len(array)-1] = ""
	array = array[:len(array)-1]

	return array
}

//GetAllByType returns all services of org
// Service_id, Service_name, adhoc,hostname
func (s serviceStore) GetAllByType(serviceType, orgID string) (services []models.Service, err error) {
	services = make([]models.Service, 0)

	var rows *sql.Rows
	rows, err = s.DB.Query("SELECT id, name, adhoc,hostname, proxy_config FROM services WHERE org_id= $1 AND type=$2", orgID, serviceType)
	if err != nil {
		return services, err
	}
	defer rows.Close()
	for rows.Next() {
		var service models.Service
		err = rows.Scan(&service.ID, &service.Name, &service.Adhoc, &service.Hostname, &service.ProxyConfig)
		if err != nil {
			logrus.Error(err)
		}
		services = append(services, service)
	}
	return
}

// get service Details based on service id. Returns service id and service name.
func (s serviceStore) GetServiceNameFromID(serviceID string, orgID string) (appName string, err error) {
	err = s.DB.QueryRow("SELECT services.name FROM services WHERE id= $1 AND org_id=$2", serviceID, orgID).Scan(&appName)
	return
}

func (s serviceStore) GetTotalServiceUsers(serviceID string, orgID string) (int64, error) {
	var count int64
	row := s.DB.QueryRow(`
select COUNT(*) from (SELECT  DISTINCT * FROM 
( 
		SELECT user_id,service_id,org_id from 
			user_accessmaps		  
			--usergroup assigned to app
		UNION SELECT user_id,servicegroup_id as service_id,ag_ug.org_id FROM
				(
				usergroup_accessmaps ag_ug  
				join user_group_maps ug on ug.group_id=ag_ug.usergroup_id 
				 
				) WHERE ag_ug.map_type='service'
				
			--usergroup assigned to appgroup				 
		UNION 	SELECT user_id,ag.service_id,ag_ug.org_id FROM
				(
				usergroup_accessmaps ag_ug  
				join user_group_maps ug on ug.group_id=ag_ug.usergroup_id 
				join service_group_maps ag on ag.group_id=ag_ug.servicegroup_id 
				) WHERE ag_ug.map_type='servicegroup'

 ) AS cc
where service_id=$1 AND org_id=$2) AS c;`, serviceID, orgID)

	err := row.Scan(&count)

	return count, err
}

func (s serviceStore) GetAllHostnames(orgID, appType string) ([]string, error) {
	var hosts []string = make([]string, 0)

	rows, err := s.DB.Query(`SELECT hostname FROM services where org_id=$1 AND type=$2;`, orgID, appType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var host string
		err := rows.Scan(&host)
		if err != nil {
			return hosts, err
		}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

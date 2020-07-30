package stats

func (s StatStore) GetAggregatedDeviceUsers(entityType, entityID, deviceType, orgID string) (total int, devices []deviceByType, err error) {
	//TODO handle diffenent entities
	devices = []deviceByType{}
	total = 0
	rows, err := s.DB.Query(`select count(*),
       COALESCE(device_hygiene->'deviceOS'->>'osName','Unknown') as os_name,
       COALESCE(device_hygiene->'deviceOS'->>'osVersion','') as os_version
from devices WHERE org_id=$1 AND type=$2
GROUP BY   device_hygiene->'deviceOS'->>'osName',device_hygiene->'deviceOS'->>'osVersion'
`, orgID, deviceType)

	if err != nil {
		return -1, devices, err
	}

	defer rows.Close()

	for rows.Next() {
		var userDevice deviceByType
		err = rows.Scan(&userDevice.Value, &userDevice.Name, &userDevice.Version)
		if err != nil {
			return -1, devices, err
		}
		//total += userDevice.Value
		devices = append(devices, userDevice)
	}

	//TODO wrap error
	err = s.DB.QueryRow(`SELECT count(id) from devices where type=$1 AND org_id=$2`, deviceType, orgID).Scan(&total)

	return total, devices, err
}

func (s StatStore) GetAggregatedMobileDevices(entityType, entityID, orgID string) (total int, devices []deviceByType, err error) {
	//TODO handle diffenent entities
	devices = []deviceByType{}
	total = 0
	rows, err := s.DB.Query(`select count(*),
       COALESCE(device_hygiene->'deviceInfo'->>'brand','Unknown') as os_name,
       COALESCE(device_hygiene->'deviceInfo'->>'deviceModel','') as os_version
from devices WHERE org_id=$1 AND type=$2
GROUP BY   device_hygiene->'deviceInfo'->>'brand',device_hygiene->'deviceInfo'->>'deviceModel'
`, orgID, "mobile")

	if err != nil {
		return -1, devices, err
	}
	defer rows.Close()

	for rows.Next() {
		var userDevice deviceByType
		err = rows.Scan(&userDevice.Value, &userDevice.Name, &userDevice.Version)
		if err != nil {
			return -1, devices, err
		}
		//		total += userDevice.Value
		devices = append(devices, userDevice)
	}

	err = s.DB.QueryRow(`SELECT count(id) from devices where type='mobile' AND org_id=$1`, orgID).Scan(&total)

	return total, devices, err
}

func (s StatStore) GetAggregatedBrowsers(entityType, entityID, orgID string) (total int, devices []deviceByType, err error) {

	rows, err := s.DB.Query(`select
    COALESCE(name,'') as browser_name,
    COALESCE(version,'') as browser_version,
    count(*)

from browsers where org_id=$1
GROUP BY name,version ;`, orgID)
	if err != nil {
		return -1, devices, err
	}

	for rows.Next() {
		var userDevice deviceByType
		err = rows.Scan(&userDevice.Name, &userDevice.Version, &userDevice.Value)
		if err != nil {
			return -1, devices, err
		}

		//total += userDevice.Value
		devices = append(devices, userDevice)
	}

	err = s.DB.QueryRow(`SELECT count(id) from browsers where org_id=$1`, orgID).Scan(&total)

	return total, devices, err
}

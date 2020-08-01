package users

import (
	"github.com/seknox/trasa/server/models"
)

////GetDeviceDetails returns device details without tokens and keys
//func (s UserStore) GetAllDevicesByType(userID, orgID, deviceType string) ([]models.UserDevice, error) {
//	var device models.UserDevice
//	var userDevices = make([]models.UserDevice, 0)
//	rows, err := s.DB.Query("SELECT device_id, device_type,  device_finger,device_hygiene, added_at  FROM userdevicesv1 WHERE deleted != true AND user_id = $1 AND org_id=$2 AND device_type = $3",
//		userID, orgID, deviceType)
//	if err != nil {
//		return userDevices, err
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		err := rows.Scan(&device.DeviceID, &device.DeviceType, &device.DeviceFinger, &device.DeviceHygiene, &device.AddedAt)
//		if err != nil {
//			return userDevices, err
//		}
//		userDevices = append(userDevices, device)
//	}
//
//	return userDevices, nil
//}

func (s UserStore) GetAllDevices(userID, orgID string) ([]models.UserDevice, error) {
	var device models.UserDevice
	var userDevices = make([]models.UserDevice, 0)
	rows, err := s.DB.Query("SELECT id, type,trusted,  device_hygiene, added_at  FROM devices WHERE deleted != true AND user_id = $1 AND org_id=$2",
		userID, orgID)
	if err != nil {
		return userDevices, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&device.DeviceID, &device.DeviceType, &device.Trusted, &device.DeviceHygiene, &device.AddedAt)
		if err != nil {
			return userDevices, err
		}
		userDevices = append(userDevices, device)
	}

	return userDevices, nil
}

// deregister user specefic devices
func (s UserStore) DeregisterUserDevices(userID, orgID string) error {
	_, err := s.DB.Exec(`Update  devices SET deleted=true where user_id = $1 AND org_id=$2`, userID, orgID)

	return err
}

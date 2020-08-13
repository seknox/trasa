package devices

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"github.com/seknox/trasa/server/models"
)

// get user devices
func (s deviceStore) GetFromID(deviceID string) (*models.UserDevice, error) {
	var device models.UserDevice
	err := s.DB.QueryRow(
		"SELECT user_id, org_id, id, type,trusted, fcm_token, public_key, device_hygiene, added_at FROM devices WHERE deleted != true AND id = $1 AND type != $2",
		deviceID, "htoken").
		Scan(&device.UserID, &device.OrgID, &device.DeviceID, &device.DeviceType, &device.Trusted, &device.FcmToken, &device.PublicKey, &device.DeviceHygiene, &device.AddedAt)

	return &device, err
}

func (s deviceStore) Trust(trusted bool, deviceID, orgID string) error {
	_, err := s.DB.Exec(`UPDATE devices set trusted=$1 WHERE id=$2 AND org_id=$3`, trusted, deviceID, orgID)
	return err
}

func (s deviceStore) Register(device models.UserDevice) error {

	//If device is already registered with same machineID, just update device hygiene
	//var deviceID string

	//If machineID is not blank, update deviceHygiene based on that machineID and set deleted=false
	// if device.DeviceHygiene.DeviceInfo.MachineID != "" {
	// 	err := s.DB.QueryRow(`SELECT id from devices WHERE machine_id=$1 AND org_id=$2`, device.MachineID, device.OrgID).Scan(&deviceID)
	// 	if err == nil {
	// 		device.DeviceID = deviceID
	// 		err = s.ReRegisterDevice(device)
	// 		return err
	// 	}
	// }

	_, err := s.DB.Exec(`INSERT into devices (user_id, org_id, id,machine_id, type,trusted, fcm_token, device_hygiene, public_key,totpsec, added_at, deleted)
						values($1, $2, $3, $4, $5, $6, $7, $8,$9,$10,$11,$12);`,
		device.UserID, device.OrgID, device.DeviceID, device.MachineID, device.DeviceType, false, device.FcmToken, device.DeviceHygiene, device.PublicKey, device.TotpSec, device.AddedAt, false)

	return err

}

func (s deviceStore) ReRegisterDevice(device models.UserDevice) error {
	_, err := s.DB.Exec(`Update  devices SET 
										user_id=$1,
										org_id=$2, 
										fcm_token=$3,
										 device_hygiene=$4,
										  public_key=$5,
										  totpsec=$6,
										   added_at=$7,
											deleted=false
											WHERE id=$8 AND machine_id=$9 AND type=$10
    `, device.UserID, device.OrgID, device.FcmToken, device.DeviceHygiene, device.PublicKey, device.TotpSec, device.AddedAt,
		device.DeviceID, device.MachineID, device.DeviceType)
	return err
}

//UpdateWorkstationHygiene updates device hygiene based on deviceID and orgID
func (s deviceStore) UpdateWorkstationHygiene(deviceHyg models.DeviceHygiene, deviceID, orgID string) error {

	_, err := s.DB.Exec(`UPDATE devices SET device_hygiene=$1,deleted = false 
									WHERE id=$2 AND org_id=$3`,
		deviceHyg, deviceID, orgID)
	if err != nil {
		return err
	}

	return nil

}

//UpdateDeviceHygiene updates device hygiene based on machineID
func (s deviceStore) UpdateDeviceHygiene(deviceHyg models.DeviceHygiene, orgID string) (deviceID string, err error) {

	//logger.Debug(deviceHyg)
	deviceHyg.LastCheckedTime = time.Now().Unix()
	err = s.DB.QueryRow(`SELECT id from devices WHERE machine_id=$1  AND org_id=$2`,
		deviceHyg.DeviceInfo.MachineID, orgID).
		Scan(&deviceID)
	if err != nil {
		return "", err
	}

	_, err = s.DB.Exec(`UPDATE devices SET device_hygiene=$1,deleted = false 
									WHERE id=$2 AND org_id=$3`,
		deviceHyg, deviceID, orgID)
	if err != nil {
		return deviceID, err
	}

	return deviceID, nil

}

// BrowserStoreExtensionDetails stores details of other extensions installed in user's browser.
func (s deviceStore) BrowserStoreExtensionDetails(brsr models.BrowserExtensions, orgID, userID, deviceID string) error {

	_, err := s.DB.Exec(`INSERT into browser_ext (user_id, org_id, browser_id,ext_id, name, description, version, maydisable,enabled,install_type,type,perms, host_perms, isvuln, vuln_reason, last_checked)
	values($1, $2, $3, $4, $5, $6, $7, $8,$9,$10,$11,$12,$13, $14, $15, $16);`, userID, orgID, deviceID, brsr.ExtensionID, brsr.Name, brsr.Description, brsr.Version, brsr.MayDisable, brsr.Enabled, brsr.InstallType, brsr.Type, pq.Array(brsr.Permissions), pq.Array(brsr.HostPermissions), brsr.IsVulnerable, brsr.VulnReason, time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}

func (s deviceStore) Deregister(deviceID, orgID string) error {
	_, err := s.DB.Exec(`DELETE from devices where id = $1 AND org_id=$2`, deviceID, orgID)
	return err
}

// RegisterBrowser stores browser detail referencing deviceID of workstation
func (s deviceStore) RegisterBrowser(brsr models.DeviceBrowser) error {

	exts, err := json.Marshal(brsr.Extensions)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(`INSERT into browsers (id, org_id, device_id,user_agent, name, version, build, extensions)
						values($1, $2, $3, $4, $5, $6, $7, $8);`,
		brsr.ID, brsr.OrgID, brsr.DeviceID, brsr.UserAgent, brsr.Name, brsr.Version, brsr.Build, string(exts))

	return err

}

//UpdateBrowserHygiene updates device hygiene based on machineID
//UpdateWorkstationHygiene updates device hygiene based on deviceID and orgID
func (s deviceStore) UpdateBrowserHygiene(brsr models.DeviceBrowser, brsrID, orgID string) error {

	exts, err := json.Marshal(brsr.Extensions)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(`UPDATE browsers SET user_agent = $1, name = $2, version = $3, build = $4, extensions = $5
									WHERE id=$6 AND org_id=$7`,
		brsr.UserAgent, brsr.Name, brsr.Version, brsr.Build, string(exts), brsrID, orgID)
	if err != nil {
		return err
	}

	return nil

}

// CheckIfExtIsRegistered validates if extID is in database.
func (s deviceStore) CheckIfExtIsRegistered(extID string) (string, error) {
	var orgID string
	err := s.DB.QueryRow(
		"SELECT devices.org_id from devices JOIN browsers b ON devices.id= b.device_id WHERE b.id=$1",
		extID).
		Scan(&orgID)

	if err != nil {
		return "", err
	}

	return orgID, err
}

// GetDeviceAndOrgIDFromExtID from extID
func (s deviceStore) GetDeviceAndOrgIDFromExtID(extID string) (orgID, deviceID, userID string, err error) {

	err = s.DB.QueryRow(
		"SELECT browsers.org_id, browsers.device_id, devices.user_id from browsers JOIN devices ON browsers.device_id=devices.id WHERE browsers.id = $1",
		extID).
		Scan(&orgID, &deviceID, &userID)

	if err != nil {
		return
	}

	return
}

func (s deviceStore) GetDeviceIDFromExtID(machineID string) (deviceID string, err error) {
	err = s.DB.QueryRow(`SELECT id from devices where machine_id=$1`, machineID).Scan(&deviceID)
	return
}

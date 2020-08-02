package devices

import (
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
)

func InitStore(state *global.State) {
	Store = DeviceStore{state}
}

var Store DeviceAdapter

type DeviceStore struct {
	*global.State
}

type DeviceAdapter interface {
	GetFromID(deviceID string) (*models.UserDevice, error)
	Register(device models.UserDevice) error
	Deregister(deviceID, orgID string) error
	Trust(trusted bool, deviceID, orgID string) error
	ReRegisterDevice(device models.UserDevice) error

	UpdateDeviceHygiene(deviceHyg models.DeviceHygiene, orgID string) (deviceID string, err error)
	UpdateWorkstationHygiene(deviceHyg models.DeviceHygiene, deviceID, orgID string) error

	RegisterBrowser(brsr models.DeviceBrowser) error
	UpdateBrowserHygiene(brsr models.DeviceBrowser, brsrID, orgID string) error

	BrowserStoreExtensionDetails(brsr models.BrowserExtensions, orgID, userID, deviceID string) error

	GetDeviceAndOrgIDFromExtID(extID string) (orgID, deviceID, userID string, err error)
	CheckIfExtIsRegistered(extID string) (string, error)
	GetDeviceIDFromExtID(machineID string) (string, error)
}

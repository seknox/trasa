package users

import (
	"net/http"

	"github.com/seknox/trasa/server/api/devices"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

type devicesByType struct {
	Mobile      []models.UserDevice `json:"mobile"`
	Workstation []models.UserDevice `json:"workstation"`
	Browser     []models.UserDevice `json:"browser"`
	HToken      []models.UserDevice `json:"hToken"`
}

// GetUserDevicesByType returns all user devices under device types.
func GetUserDevicesByType(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	userID := chi.URLParam(r, "userID")

	//used in /my/devices
	if userID == "" {
		userID = userContext.User.ID
	}

	var resp = devicesByType{
		Mobile:      make([]models.UserDevice, 0),
		Workstation: make([]models.UserDevice, 0),
		Browser:     make([]models.UserDevice, 0),
		HToken:      make([]models.UserDevice, 0),
	}

	alldevices, err := Store.GetAllDevices(userID, userContext.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to fetch user devices.", "SingleUserDevices")
		return
	}

	for _, device := range alldevices {
		switch device.DeviceType {
		case "mobile":
			resp.Mobile = append(resp.Mobile, device)
		case "workstation":
			resp.Workstation = append(resp.Workstation, device)
		case "browser":
			resp.Browser = append(resp.Browser, device)
		case "htoken":
			resp.HToken = append(resp.HToken, device)
		default:
			logrus.Errorf("invalid device type %s: %v", device.DeviceType, device)
		}

	}

	utils.TrasaResponse(w, 200, "success", "devices fetched.", "SingleUserDevices", resp)
}

//RemoveUserDevice removes user device
func RemoveUserDevice(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	deviceID := chi.URLParam(r, "deviceID")

	err := devices.Store.Deregister(deviceID, userContext.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to delete user device.", "failed to remove user device", nil, nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "user device removed", "user device removed", nil, nil)

}

//TrustUserDevice marks certain user device as trusted
func TrustUserDevice(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	var dev models.UserDevice
	err := utils.ParseAndValidateRequest(r, &dev)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "failed to trust user device")
		return
	}

	err = devices.Store.Trust(dev.Trusted, dev.DeviceID, userContext.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to update user device.", "failed to trust user device", nil, nil)
		return
	}

	intent := "user device marked untrusted"
	if dev.Trusted {
		intent = "user device marked trusted"
	}

	utils.TrasaResponse(w, 200, "success", "user device updated", intent)

}

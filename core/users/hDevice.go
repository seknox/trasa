package users

import (
	"net/http"

	"github.com/seknox/trasa/core/devices"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/models"
	"github.com/seknox/trasa/utils"
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

func RemoveUserDevice(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	deviceID := chi.URLParam(r, "deviceID")

	err := devices.Store.Deregister(deviceID, userContext.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to delete user device.", "User device not removed", nil, nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "user device removed", "User device removed", nil, nil)

}

func TrustUserDevice(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	var dev models.UserDevice
	err := utils.ParseAndValidateRequest(r, &dev)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "User device not updated")
		return
	}

	err = devices.Store.Trust(dev.Trusted, dev.DeviceID, userContext.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to update user device.", "User device not update", nil, nil)
		return
	}

	intent := "user device untrusted"
	if dev.Trusted {
		intent = "user device untrusted"
	}

	utils.TrasaResponse(w, 200, "success", "user device updated", intent, nil, nil)

}

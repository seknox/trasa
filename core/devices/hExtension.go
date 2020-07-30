package devices

import (
	"encoding/json"
	"net/http"

	"github.com/seknox/trasa/models"
	"github.com/seknox/trasa/utils"
	logger "github.com/sirupsen/logrus"
)

type brsrDetailReq struct {
	// ExtToken is extension Token assigned by TRASA to extension
	ExtToken string                     `json:"extToken"`
	Details  []models.BrowserExtensions `json:"details"`
}

// GetBrsrDetails collects extension and browser details from trasa browser extension
func GetBrsrDetails(w http.ResponseWriter, r *http.Request) {

	var req brsrDetailReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.TrasaResponse(w, 200, "failed", "invalid request format", "GetBrsrDetails", nil, nil)
		return
	}
	// verify and get orgID from extToken

	deviceDetailFromDB, err := Store.GetFromID(req.ExtToken)
	if err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid DeviceID", "GetBrsrDetails", nil, nil)
		return
	}

	for _, v := range req.Details {
		check := IsKnownExts(v.ExtensionID)
		if check != true {
			err := Store.BrowserStoreExtensionDetails(v, deviceDetailFromDB.OrgID, deviceDetailFromDB.UserID, deviceDetailFromDB.DeviceID)
			if err != nil {
				// if we get error here, it means extensiondetails could not be store, alert admins here. TODO
				logger.Trace(err)
			}
		}

	}

	utils.TrasaResponse(w, 200, "success", "Extension Rergistered", "GetBrsrDetails", nil, nil)
	return
}

var knownExts = []string{"default-theme@mozilla.org", "firefox-compact-light@mozilla.org", "firefox-compact-dark@mozilla.org"}

// IsKnownExts checks default extensions found in browsers
func IsKnownExts(extID string) bool {
	for _, v := range knownExts {
		if v == extID {
			return true
		}
	}
	return false
}

package idps

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/crypt"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	logger "github.com/sirupsen/logrus"
)

// GetSyncDetail retrieves key or token from database. should fetch and return key tag rather than key value.
func GetSyncDetail(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	vendorID := chi.URLParam(r, "vendorID")

	key, err := Store.GetCloudSyncState(uc.User.OrgID, vendorID)
	if err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to retrieve sync data.", "GetSyncDetail-GetCloudSyncState", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "data fetched", "GetSyncDetail", key)
}

// SyncNow initiates TRASA sync with cloudIAAS provider
func SyncNow(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	//vendorID := chi.URLParam(r, "vendorID")

	key, err := crypt.Store.GetKeyOrTokenWithKeyval(uc.User.OrgID, consts.KEY_DOAPI)
	if err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to fetch api key.", "SyncNow", nil, nil)
		return
	}

	drops, err := GetDoDroplets(key.KeyVal)
	if err != nil {
		logger.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to fetch resources with provided key.", "SyncNow", nil, nil)
		return
	}

	var failedservices []string
	for _, v := range drops {
		var srvc = new(models.Service)
		srvc.ID = utils.GetUUID()
		srvc.Name = strings.ToLower(v.Name)
		srvc.Type = "ssh"
		srvc.OrgID = uc.User.OrgID
		srvc.SecretKey = utils.GetRandomString(7)
		srvc.Adhoc = false
		srvc.Passthru = false
		srvc.CreatedAt = time.Now().Unix() //.In(nep).String() // .UTC().Format(time.RFC3339)
		srvc.UpdatedAt = time.Now().Unix() //.UTC().Format(time.RFC3339)
		srvc.ExternalProviderName = "digital-ocean"
		srvc.ManagedAccounts = ""
		srvc.ProxyConfig = models.ReverseProxy{}
		srvc.IPDetails = models.IPDetails{}
		srvc.ExternalSecurityGroup = "{}"

		srvc.Hostname, _ = v.PublicIPv4()

		err := services.Store.Create(srvc)
		if err != nil {
			failedservices = append(failedservices, srvc.Name)
			logger.Debug(fmt.Sprintf("Failed to create app: %s", srvc.Name))
		}
	}

	logger.Debug("failedServices: ", failedservices)

	utils.TrasaResponse(w, 200, "success", "resources synced", "SyncNow", key)
}

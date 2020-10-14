package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

func GetServiceDetail(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "serviceID")

	//TODO make different get method with orgID
	appDetailFromDB, err := Store.GetFromID(serviceID)
	if err != nil {
		logrus.Error(err)
		return
	}

	utils.TrasaResponse(w, 200, "success", "app details fetched", "GetAppDetailFromserviceIDV2", appDetailFromDB)

}

type AllServicesByType struct {
	SSH    []models.Service `json:"ssh"`
	RDP    []models.Service `json:"rdp"`
	HTTP   []models.Service `json:"http"`
	DB     []models.Service `json:"db"`
	Radius []models.Service `json:"radius"`
	Other  []models.Service `json:"other"`
}

func GetAllServices(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	var resp AllServicesByType
	var err error
	resp.SSH, err = Store.GetAllByType("ssh", uc.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not get ssh services", "GetAllServices", resp)
		return
	}
	resp.RDP, err = Store.GetAllByType("rdp", uc.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not get rdp services", "GetAllServices", resp)
		return
	}
	resp.DB, err = Store.GetAllByType("db", uc.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not get db services", "GetAllServices", resp)
		return
	}
	resp.HTTP, err = Store.GetAllByType("http", uc.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not get http services", "GetAllServices", resp)
		return
	}
	resp.Radius, err = Store.GetAllByType("radius", uc.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not get radius services", "GetAllServices", resp)
		return
	}

	resp.Other, err = Store.GetAllByType("other", uc.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not get other services", "GetAllServices", resp)
		return
	}

	utils.TrasaResponse(w, 200, "success", "", "GetAllServices", resp)

}

// CreateApp creates new App. service here represents login (can be 2fa) endpoints.
// service details can include service name, passthru options and costum 2FA policies.
func CreateService(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	var req models.Service

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "failed to create service", nil)
		return
	}

	newService := fillInitialFields(&req)
	newService.OrgID = uc.Org.ID

	err := Store.Create(newService)
	if err != nil {
		logrus.Error(err)
		reason := utils.GetConstraintErrorMessage(err)
		utils.TrasaResponse(w, 200, "failed", reason, "failed to create service")
		return
	}

	utils.TrasaResponse(w, 200, "success", "service created", fmt.Sprintf("service named %s created. Type:%s Hostname:%s", req.Name, req.Type, req.Hostname), req)
}

func fillInitialFields(req *models.Service) *models.Service {
	req.Name = utils.NormalizeString((req.Name))
	req.ID = utils.GetUUID()

	req.CreatedAt = time.Now().Unix()
	req.UpdatedAt = time.Now().Unix()
	req.SecretKey = utils.GetRandomString(17) //hex.EncodeToString(apptoken)
	req.Passthru = false
	req.ManagedAccounts = ""
	req.Hostname = utils.NormalizeString(req.Hostname)
	req.ExternalProviderName = ""
	req.ExternalSecurityGroup = "{}"
	req.IPDetails = models.IPDetails{}
	req.ProxyConfig = models.ReverseProxy{}

	return req
}

// UpdateService should handle Service detail updates.
func UpdateService(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	var req models.Service

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		utils.TrasaResponse(w, 200, "failed", "invalid request", "service not updated", nil, nil)
		return
	}

	req.OrgID = userContext.Org.ID
	req.Hostname = utils.NormalizeString(req.Hostname)
	req.Name = utils.NormalizeString(req.Name)
	req.RemoteAppName = utils.NormalizeString(req.RemoteAppName)
	req.UpdatedAt = time.Now().Unix() //.In(zone).Format(time.RFC3339)

	// create new app
	err := Store.Update(&req)
	if err != nil {
		reason := utils.GetConstraintErrorMessage(err)
		utils.TrasaResponse(w, 200, "success", reason, fmt.Sprintf("service %s not updated", req.Name))
		return
	}
	utils.TrasaResponse(w, 200, "success", "Service profile updated", fmt.Sprintf("service %s updated", req.Name))

}

type ReverseProxyReq struct {
	ServiceID string              `json:"serviceID"`
	Name      string              `json:"name"`
	Proxy     models.ReverseProxy `json:"proxy"`
}

// UpdateHTTPProxy updates proxyConfig of http service.
func UpdateHTTPProxy(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	var req ReverseProxyReq

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "service not updated", nil, nil)
		return
	}

	// create new service
	err := Store.updateHttpProxy(req.ServiceID, uc.User.OrgID, time.Now().Unix(), req.Proxy)
	if err != nil {
		logrus.Error(err)
		reason := utils.GetConstraintErrorMessage(err)
		utils.TrasaResponse(w, 200, "success", reason, fmt.Sprintf("service %s not updated", req.Name))
		return
	}
	utils.TrasaResponse(w, 200, "success", "reverse proxy configured", fmt.Sprintf("service %s updated", req.Name))

}

// DeleteService deletes Service from database.
func DeleteService(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	var newApp models.Service

	if err := json.NewDecoder(r.Body).Decode(&newApp); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "failed to delete service", nil, nil)
		return
	}

	appName, err := Store.Delete(newApp.ID, userContext.User.OrgID)
	if err != nil {
		logrus.Errorf("delete app: %v", err)
		utils.TrasaResponse(w, 200, "failed", "failed to delete service", "failed to delete service ")
		return
	}

	utils.TrasaResponse(w, 200, "success", "", fmt.Sprintf("service  %s deleted", appName))

}

// CheckAppConfigs accepts service config request and validates configuration settigs, basically as integrationID and integrationKey
// set on client Trasa Connectors (trasa-win, ssh).
func CheckAppConfigs(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("CheckAppConfigs request received")
	var remoteLogin models.ServiceLogin

	if err := utils.ParseAndValidateRequest(r, &remoteLogin); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "service config check failed", "")
		return
	}

	// Get service secret from service ID
	serviceDetailFromDB, err := Store.GetFromID(remoteLogin.ServiceID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "invalid service ID", "service config check failed")
		return
	}

	if remoteLogin.ServiceKey != serviceDetailFromDB.SecretKey {
		logrus.Debug("invalid secret")
		utils.TrasaResponse(w, http.StatusOK, "failed", "invalid service secret or ID", "service config check failed", nil)
		return
	}

	// if we are here, this means config values were successfully validated and we return success response
	utils.TrasaResponse(w, http.StatusOK, "success", "valid integration data", "service config check successful", nil)
	return

}

package stats

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/seknox/trasa/server/api/crypt"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/utils"
	//"io"
)

//GetAggregatedUsers returns user aggregations
//It returns stats like total users in IDP, admin users, disabled users etc
func GetAggregatedUsers(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	entityType := chi.URLParam(r, "entitytype")
	entityID := chi.URLParam(r, "entityid")

	switch entityType {
	case "org":
		usersAgg, err := Store.GetAggregatedIdpUsers(entityType, entityID, uc.User.OrgID)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "Could not get user count", "GetTotalGroups", nil)
			return
		}

		usersAgg.Groups, err = Store.GetTotalGroups(uc.User.OrgID, "usergroup")
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "Could not get user count", "GetTotalGroups", nil)
			return
		}

		usersAgg.Admins, err = Store.GetTotalAdmins(uc.User.OrgID)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "Could not get user count", "GetTotalAdmins", nil)
			return
		}

		usersAgg.DisabledUsers, err = Store.GetTotalDisabledUsers(uc.User.OrgID)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "Could not get user count", "GetTotalDisabled", nil)
			return
		}

		utils.TrasaResponse(w, 200, "success", "", "", usersAgg)
		return
	}
	logrus.Error("unsupported entity ", entityType)
	utils.TrasaResponse(w, 200, "failed", "Not supported yet", "", nil)
}

//GetAggregatedDevices returns device stats (total devices, mobile devices, total browsers)
func GetAggregatedDevices(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	entityType := chi.URLParam(r, "entitytype")
	entityID := chi.URLParam(r, "entityid")
	var res allUserDevices
	var err error
	switch entityType {
	case "org":
		res.TotalBrowsers, res.BrowserByType, err = Store.GetAggregatedBrowsers(entityType, uc.Org.ID, uc.User.OrgID)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "Could not get aggregated browsers", "", nil)
			return
		}
		res.TotalMobiles, res.MobileByType, err = Store.GetAggregatedMobileDevices(entityType, entityID, uc.User.OrgID)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "Could not get aggregated mobiles", "", nil)
			return
		}

		res.TotalWorkstations, res.WorkstationByType, err = Store.GetAggregatedDeviceUsers(entityType, entityID, "workstation", uc.User.OrgID)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 200, "failed", "Could not get  aggregated workstation", "", nil)
			return
		}
		res.TotalUserdeivce = res.TotalWorkstations + res.TotalMobiles + res.TotalBrowsers

		utils.TrasaResponse(w, 200, "success", "", "", res)
		return
	}
	logrus.Error("unsupported entity ", entityType)
	utils.TrasaResponse(w, 200, "failed", "Not supported yet", "", nil)
}

//GetAggregatedFailedReasons aggregates authentication according to failed reasons
func GetAggregatedFailedReasons(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	entityType := chi.URLParam(r, "entitytype")
	entityID := chi.URLParam(r, "entityid")
	timeFilter := chi.URLParam(r, "timeFilter")

	reasonsAgg, err := Store.GetAggregatedLoginFails(entityType, entityID, uc.User.OrgID, uc.Org.Timezone, timeFilter)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not get failed reasons aggr", "", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "", "", reasonsAgg)
	return
}

//GetAggregatedServices aggregates services according to type
func GetAggregatedServices(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	//entityType := chi.URLParam(r, "entitytype")
	//entityID := chi.URLParam(r, "entityid")

	appsAgg, err := Store.GetAggregatedServices(uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not get service count", "GetTotalGroups", nil)
		return
	}
	utils.TrasaResponse(w, 200, "success", "", "", appsAgg)
	return

}

//GetAggregatedIDPServices aggregates services according to IDP
func GetAggregatedIDPServices(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	idpName := chi.URLParam(r, "idpname")

	appsAgg, err := Store.GetAggregatedIDPServices(idpName, uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not get service count", "GetTotalGroups", nil)
		return
	}
	utils.TrasaResponse(w, 200, "success", "", "", appsAgg)
	return

}

//GetAggregatedLoginHours aggregates authentications according to hour
func GetAggregatedLoginHours(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	timeFilter := chi.URLParam(r, "timeFilter")
	statusFilter := chi.URLParam(r, "statusFilter")

	entityType := chi.URLParam(r, "entitytype")
	entityID := chi.URLParam(r, "entityid")

	loginHours, err := Store.GetAggregatedLoginHours(entityType, entityID, uc.Org.Timezone, uc.Org.ID, timeFilter, statusFilter)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to get login hours", "GetAggregatedLoginHours", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "success", "GetAggregatedLoginHours", loginHours)

}

//GetTotalManagedUsers returns total managed users (password stored in vault)
func GetTotalManagedUsers(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	managedUsers, err := Store.GetTotalManagedUsers("org", uc.Org.ID, uc.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to get managed users", "GetTotalManagedUsers", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "success", "GetTotalManagedUsers", managedUsers)

}

//GetIPAggs aggregates authentications according to client IP
func GetIPAggs(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	timeFilter := chi.URLParam(r, "timeFilter")
	statusFilter := chi.URLParam(r, "statusFilter")

	entityType := chi.URLParam(r, "entitytype")
	entityID := chi.URLParam(r, "entityid")

	aggsIPs, err := Store.GetAggregatedIPs(entityType, entityID, uc.Org.ID, uc.Org.Timezone, timeFilter, statusFilter)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to get ips", "GetIPAggs", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "", "GetIPAggs", aggsIPs)

}

//GetLoginsByType aggregates authentications by service type
func GetLoginsByType(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	timeFilter := chi.URLParam(r, "timeFilter")
	statusFilter := strings.ToLower(chi.URLParam(r, "statusFilter"))

	entityType := chi.URLParam(r, "entitytype")
	entityID := chi.URLParam(r, "entityid")

	loginsByType, err := Store.GetLoginsByType(entityType, entityID, uc.Org.ID, uc.Org.Timezone, timeFilter, statusFilter)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to get ips", "GetIPAggs", nil)
		return
	}

	remoteAppCount, err := Store.GetRemoteAppCount(entityType, entityID, uc.Org.ID, uc.Org.Timezone, timeFilter, statusFilter)

	utils.TrasaResponse(w, 200, "success", "", "GetIPAggs", loginsByType, remoteAppCount)

}

// GetMapPlotData returns city name, total population (total count), and cooardinates of login source.
func GetMapPlotData(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	timeFilter := chi.URLParam(r, "timeFilter")
	statusFilter := chi.URLParam(r, "statusFilter")

	entityType := chi.URLParam(r, "entitytype")
	entityID := chi.URLParam(r, "entityid")

	//fmt.Println("org inside gettotal: ", userContext.Org.ID)
	geoPlotData, err := Store.SortLoginByCity(entityType, entityID, uc.Org.ID, uc.Org.Timezone, timeFilter, statusFilter)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to retrieve geo data", "GetUserMapPlotData", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "", "GetIPAggs", geoPlotData)
}

//GetPoliciesStats returns policy stats
func GetPoliciesStats(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	//fmt.Println("org inside gettotal: ", userContext.Org.ID)
	policyStats, err := Store.GetPoliciesStats(uc.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to retrieve policy stat", "GetPolicyStat")
		return
	}

	utils.TrasaResponse(w, 200, "success", "", "GetIPAggs", policyStats)
}

//func GetCAStats(w http.ResponseWriter, r *http.Request) {
//	uc := r.Context().Value("user").(models.UserContext)
//
//	var resp struct {
//		SshHostCA   bool `json:"sshHostCA"`
//		SshUserCA   bool `json:"sshUserCA"`
//		SshSystemCA bool `json:"sshSystemCA"`
//		HttpCA      bool `json:"httpCA"`
//	}
//	//fmt.Println("org inside gettotal: ", userContext.Org.ID)
//	_, err := Store.GetCertDetail(uc.Org.ID, "user", consts.CERT_TYPE_SSH_CA)
//	if err == nil {
//		resp.SshUserCA = true
//	}
//	_, err = Store.GetCertDetail(uc.Org.ID, "host", consts.CERT_TYPE_SSH_CA)
//	if err == nil {
//		resp.SshHostCA = true
//	}
//	_, err = Store.GetCertDetail(uc.Org.ID, "system", consts.CERT_TYPE_SSH_CA)
//	if err == nil {
//		resp.SshSystemCA = true
//	}
//	_, err = Store.GetCertDetail(uc.Org.ID, "ca", consts.CERT_TYPE_HTTP_CA)
//	if err == nil {
//		resp.HttpCA = true
//	}
//
//	utils.TrasaResponse(w, 200, "success", "", "GetCAstats", nil, resp)
//}

func GetAppPermStats(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	serviceID := chi.URLParam(r, "serviceID")

	var resp struct {
		Users      int `json:"users"`
		Policies   int `json:"policy"`
		Groups     int `json:"groups"`
		Secrets    int `json:"secrets"`
		Privileges int `json:"privileges"`
	}

	var err error
	//fmt.Println("org inside gettotal: ", userContext.Org.ID)
	resp.Secrets, err = Store.GetTotalManagedUsers("service", serviceID, uc.Org.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.Error(err)
	}
	resp.Policies, err = Store.GetPoliciesOfService(serviceID, uc.Org.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.Error(err)
	}
	resp.Privileges, err = Store.GetTotalPrivilegesOfService(serviceID, uc.Org.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.Error(err)
	}
	resp.Groups, err = Store.GetTotalGroupsServiceIsAssignedTo(serviceID, uc.Org.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.Error(err)
	}

	resp.Users, err = Store.GetTotalUsersAssignedToService(serviceID, uc.Org.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.Error(err)
	}

	utils.TrasaResponse(w, 200, "success", "", "GetCAstats", resp)
}

func HexaEvents(w http.ResponseWriter, r *http.Request) {

	userContext := r.Context().Value("user").(models.UserContext)

	entityType := chi.URLParam(r, "entitytype")
	entityID := chi.URLParam(r, "entityid")
	status := chi.URLParam(r, "status")

	events, err := Store.GetTodayHexaLoginEvents(entityType, entityID, userContext.Org.ID, status, userContext.Org.Timezone)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to get logs", "HexaEvents: GetTodayHexaLoginEvents")
		return
	}

	utils.TrasaResponse(w, 200, "success", "HexaEvents", "", events)
}

// GetSuccessAndFailedEvents returns total number of logins, total failed logins and total successful logins aggregated stats
func GetSuccessAndFailedEvents(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	entityType := chi.URLParam(r, "entitytype")
	entityID := chi.URLParam(r, "entityid")
	timeFilter := chi.URLParam(r, "timeFilter")
	if entityType == "org" {
		entityID = userContext.Org.ID
	}

	events, err := Store.GetAllAuthEventsByEntityType(entityType, entityID, timeFilter, userContext.Org.Timezone)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to fetch fetched aggregated events", "GetLoginEvents: checkIfExist")
		return
	}

	utils.TrasaResponse(w, 200, "success", "fetched auth events", "GetSuccessAndFailedEvents", events)
}

//GetCAStats returns CA stats
func GetCAStats(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	var resp struct {
		SshHostCA   bool `json:"sshHostCA"`
		SshUserCA   bool `json:"sshUserCA"`
		SshSystemCA bool `json:"sshSystemCA"`
		HttpCA      bool `json:"httpCA"`
	}
	//fmt.Println("org inside gettotal: ", userContext.Org.ID)
	_, err := crypt.Store.GetCertDetail(uc.Org.ID, "user", consts.CERT_TYPE_SSH_CA)
	if err == nil {
		resp.SshUserCA = true
	}
	_, err = crypt.Store.GetCertDetail(uc.Org.ID, "host", consts.CERT_TYPE_SSH_CA)
	if err == nil {
		resp.SshHostCA = true
	}
	_, err = crypt.Store.GetCertDetail(uc.Org.ID, "system", consts.CERT_TYPE_SSH_CA)
	if err == nil {
		resp.SshSystemCA = true
	}
	_, err = crypt.Store.GetCertDetail(uc.Org.ID, "ca", consts.CERT_TYPE_HTTP_CA)
	if err == nil {
		resp.HttpCA = true
	}

	utils.TrasaResponse(w, 200, "success", "", "GetCAstats", resp)
}

// GetTotalLoginsByDate returns array of GetSuccessAndFailedEvents per day
func GetTotalLoginsByDate(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	entityType := chi.URLParam(r, "entitytype")
	entityID := chi.URLParam(r, "entityid")

	if entityType == "org" {
		entityID = userContext.User.OrgID
	}
	//fmt.Println("org inside gettotal: ", userContext.Org.ID)
	events, err := Store.GetTotalLoginsByDate(entityType, entityID, userContext.Org.ID, userContext.Org.Timezone)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to retrieve policy stat", "GetPolicyStat")
		return
	}
	utils.TrasaResponse(w, 200, "failed", "failed to retrieve policy stat", "GetPolicyStat", events)
}

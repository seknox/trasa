package main

import (
	"github.com/seknox/trasa/server/api/adhoc"
	"net/http"

	"github.com/seknox/trasa/server/accessproxy/rdpproxy"
	"github.com/seknox/trasa/server/accessproxy/sshproxy"
	"github.com/seknox/trasa/server/api/auth/serviceauth"
	"github.com/seknox/trasa/server/api/crypt"
	"github.com/seknox/trasa/server/api/providers/ca"
	"github.com/seknox/trasa/server/api/providers/sidp"
	"github.com/seknox/trasa/server/api/providers/uidp"

	"github.com/seknox/trasa/server/middlewares"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/accessmap"
	"github.com/seknox/trasa/server/api/auth"
	"github.com/seknox/trasa/server/api/auth/tfa"
	"github.com/seknox/trasa/server/api/devices"
	"github.com/seknox/trasa/server/api/groups"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/my"
	"github.com/seknox/trasa/server/api/notif"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/policies"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/api/stats"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/api/users"
	"github.com/sirupsen/logrus"
)

// CoreAPIRoutes holds api route declarations for trasa-server
func CoreAPIRoutes(r *chi.Mux) *chi.Mux {

	//logLevel := utils.NormalizeString(global.GetConfig().Logging.Level)
	//if logLevel == "trace" {
	//	r.Use(middlewares.Dumper{}.Handler)
	//}

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		logrus.Debug("NOT FOUND URL in core api: ", req.URL)
		FileServer(r, req.URL.Path)

	})

	r.Route("/auth", func(r chi.Router) {
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			logrus.Debugf("Reached not found in auth: %s", r.URL)

		})

		r.Post("/identity", auth.LoginHandler)
		r.Post("/external/saml/{orgid}/{vendorname}", auth.SAMLLoginHandler)
		r.Post("/tfa", auth.TfaHandler)

		r.Post("/accessproxy/http", serviceauth.AuthHTTPAccessProxy)

		r.Post("/crypto/kex", crypt.Kex)
		r.Post("/device/register", auth.RegisterUserDevice)
		r.Post("/device/ext/sync", auth.SyncExtension)

		r.Post("/agent/nix", serviceauth.AgentLogin)
		r.Post("/agent/win", serviceauth.AgentLogin)
		r.Post("/agent/db", serviceauth.DBLogin)
		r.Post("/agent/checkconfig", services.CheckAppConfigs)
	})

	r.Route("/idp", func(r chi.Router) {
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			logrus.Debugf("Reached not found in auth: %s", r.URL)
		})

		r.Post("/login", auth.LoginHandler)
		r.Post("/login/tfa", auth.TfaHandler)
		r.Post("/login/checktfa", auth.ConfirmTOTPAndSave)

	})

	r.Route("/api/woa", func(r chi.Router) {
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			logrus.Debugf("Reached not found in auth: %s", r.URL)
		})
		r.Get("/verify/{verifytoken}", my.VerifyAccount)
		r.Post("/setup/password", my.FirstTimePasswordSetup)

		r.Post("/proxy/http/getsession", serviceauth.GetHttpSession)
		r.Post("/proxy/http/rmsession", serviceauth.DestroyHttpSession)
		r.Post("/devices/brsrdetail", devices.GetBrsrDetails)

		r.Post("/forgotpass", my.ForgotPassword)
		r.Get("/providers/uidp/all", uidp.GetAllIdpsWoa)

	})
	r.Get("/api/v1/my/download_file/get/{fileName}/{sskey}", my.FileDownloadHandler)

	// ///////////////////////////////////////////////////
	// *** DEPCRECATED ***
	// Below api only kept for older version compatibility
	r.Post("/api/v1/remote/auth/win", serviceauth.AgentLogin)
	r.Post("/api/v1/remote/auth/nix", serviceauth.AgentLogin)

	r.Get("/api/v1/logs/livesessions", middlewares.SessionValidatorWS(logs.GetLiveSessions))
	r.Get("/api/v1/logs/vsessionlog", logs.GetVideoLog) //SessionValidator(GetSessionLog))

	authmiddleware := middlewares.AuthMiddleware{}
	inapptrailmiddleware := middlewares.InAppTrail{}

	gproxy := rdpproxy.NewProxy()

	r.Route("/accessproxy", func(r chi.Router) {
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			logrus.Debugf("Reached not found in auth: %s", r.URL)
		})
		//r.Use(authmiddleware.Handler)
		r.Route("/rdp", func(r chi.Router) {
			r.Get("/tunnel", middlewares.SessionValidatorWS(gproxy.ServeWS))
			//r.Get("/ping", oproxy.Ping)

		})

		r.Route("/ssh", func(r chi.Router) {
			r.Get("/connect", middlewares.SessionValidatorWS(sshproxy.ConnectNewSSH))
			r.Get("/join", middlewares.SessionValidatorWS(sshproxy.JoinSSHSession))
		})

	})

	r.Route("/api/v1", func(r chi.Router) {
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			logrus.Debugf("Reached not found in auth: %s", r.URL)
		})

		r.Use(authmiddleware.Handler)
		r.Use(inapptrailmiddleware.Handler)

		r.Post("/devicedetailpipe", devices.DeviceDetailPipe)
		r.Post("/passmydevicedetail", devices.PassMyDeviceDetail)

		r.Get("/org/detail", orgs.Get)
		r.Post("/org/update", orgs.Update)

		r.Get("/my", my.GetMyDetail)
		r.Post("/my/forgotpass", my.ForgotPassword)
		r.Post("/my/forgotpasstfa", auth.TfaHandler)
		r.Post("/my/changepass", my.ChangePassword)
		r.Get("/my/services", my.GetMyServicesDetail)
		r.Post("/my/services/adhoc/request", adhoc.AdhocReq)
		r.Get("/my/services/adhoc/getadmins", adhoc.GetAdmins)

		r.Post("/my/generatekey", my.GenerateKeyPair)
		//r.Post("/setup/password/{setpasswordtoken}", users.FirstTimePasswordSetup)

		r.Get("/my/account-details", my.MyAccountDetails)
		r.Get("/my/devices", users.GetUserDevicesByType)
		r.Get("/my/auth/log", my.GetMyEvents)
		r.Get("/my/auth/log/{page}/{size}/{dateFrom}/{dateTo}", my.GetMyEventsByPage)

		r.Post("/my/upload_file", my.FileUploadHandler)
		r.Get("/my/download_file/list", my.GetDownloadableFileList)
		r.Post("/my/download_file/token", my.GetFileDownloadToken)
		r.Post("/my/delete_file/{fileName}", my.FileDeleteHandler)

		//r.Post("/testdevicesync", mdlwr.SessionValidator(TestDeviceSync))
		//Enroll device needs both password and session
		//After device is enrolled
		r.Post("/my/enroldevice", auth.Enrol2FADevice)
		r.Post("/my/devices/delete/{deviceID}", users.RemoveUserDevice)
		r.Get("/my/authmeta/{appID}/{username}", my.GetAuthMeta)

		r.Get("/my/notifs", notif.GetPendingNotif)
		r.Post("/my/notif/resolve", notif.ResolveNotif)
		r.Delete("/my/logout", auth.LogoutHandler)
		//User Crud
		r.Get("/user/{userID}", users.GetUserDetails)
		r.Get("/user/all", users.GetAllUsers)
		r.Post("/user/create", users.CreateUser)
		r.Post("/user/update", users.UpdateUser)
		r.Post("/user/delete/{userID}", users.DeleteUser)

		r.Get("/user/devices/all/{userID}", users.GetUserDevicesByType)
		r.Post("/user/devices/delete/{deviceID}", users.RemoveUserDevice)
		r.Post("/user/devices/trust", users.TrustUserDevice)

		r.Get("/user/assignedgroups/{userID}", users.GetGroupsAssignedToUser)

		//Services
		r.Post("/services/create", services.CreateService)
		r.Get("/services/all", services.GetAllServices)
		r.Get("/services/{serviceID}", services.GetServiceDetail)
		r.Post("/services/update", services.UpdateService)
		r.Post("/services/httpproxy/update", services.UpdateHTTPProxy)
		r.Post("/services/delete", services.DeleteService)

		r.Post("/services/creds/store", services.StoreServiceCredentials)
		r.Post("/services/creds/view", services.ViewCreds)
		r.Post("/services/creds/delete", services.DeleteCreds)
		//TODO use vault
		r.Post("/services/sslcerts/update/{serviceID}", services.UpdateSSLCerts)
		r.Post("/services/hostcerts/update", services.UpdateHostCerts)
		r.Get("/services/hostcerts/download/{serviceID}", services.DownloadHostCerts)

		//Groups

		r.Get("/groups/user/{groupid}", groups.GetUserGroup)
		r.Get("/groups/service/{groupID}", groups.GetServiceGroup)
		r.Get("/groups/{groupType}", groups.GetAllGroups)

		r.Post("/groups/create", groups.CreateGroup)
		r.Post("/groups/update", groups.UpdateGroup)
		r.Post("/groups/delete/{groupID}", groups.DeleteGroup)

		r.Post("/groups/service/update", groups.UpdateServiceGroup)
		r.Post("/groups/user/update", groups.UpdateUsersGroup)

		//Access Maps

		r.Get("/accessmap/service/usergroup/{serviceID}", accessmap.GetUserGroupServiceGroupAccessMaps)
		r.Get("/accessmap/servicegroup/usergroup/{serviceGroupID}", accessmap.GetUserGroupServiceGroupAccessMaps)

		r.Post("/accessmap/servicegroup/usergroup/create", accessmap.CreateServiceGroupUserGroupMap)
		r.Post("/accessmap/servicegroup/usergroup/update", accessmap.UpdateServiceGroupUserGroup)
		r.Post("/accessmap/servicegroup/usergroup/delete", accessmap.DeleteServiceGroupUserGroupMap)
		r.Get("/accessmap/servicegroup/addedusergroups/{groupID}", accessmap.GetUserGroupsAssignedToServiceGroups)
		r.Get("/accessmap/servicegroup/usergroupstoadd", accessmap.UserGroupsToAdd)

		r.Get("/accessmap/service/user/{serviceID}", accessmap.GetUserAccessMaps)
		r.Post("/accessmap/service/user/create", accessmap.CreateServiceUserMap)
		r.Post("/accessmap/service/user/update", accessmap.UpdateServiceUserMap)
		r.Post("/accessmap/service/user/delete", accessmap.DeleteServiceUserMap)

		//Devices

		//////////////////// 	POLICY 		/////////////////////////
		r.Post("/groups/policy/create", policies.CreatePolicy)
		r.Post("/groups/policy/update", policies.UpdatePolicy)
		r.Get("/groups/policy/all", policies.GetPolicies)
		r.Get("/groups/policy/{policyID}", policies.GetPolicy)
		r.Post("/groups/policy/delete", policies.DeletePolicies)

		//AD Hoc
		r.Get("/iam/adhoc/requests/my", adhoc.AdhocReqAssignedToMe)
		r.Post("/iam/adhoc/respond", adhoc.GrantOrDenyAdhoc)
		r.Get("/iam/adhoc/requests/all", adhoc.GetAllAdhoqRequests)

		//Logs
		r.Get("/logs/auth/{entitytype}/{entityid}", logs.GetLoginEvents)
		r.Get("/logs/auth/{entitytype}/{entityid}/{page}/{size}", logs.GetLoginEventsByPage)
		r.Get("/logs/auth/{entitytype}/{entityid}/{page}/{size}/{dateFrom}/{dateTo}", logs.GetLoginEventsByPage)
		r.Get("/logs/inapptrails/org", logs.GetAllInAppTrails)
		r.Get("/logs/inapptrails/org/{page}/{size}/{dateFrom}/{dateTo}", logs.GetAllInAppTrails)
		r.Get("/logs/sessionlog", logs.GetRawLog)

		//Aggregrations and Stats
		r.Get("/stats/users/{entitytype}/{entityid}", stats.GetAggregatedUsers)
		r.Get("/stats/services/{entitytype}/{entityid}", stats.GetAggregatedServices)
		r.Get("/stats/serviceidp/{idpname}", stats.GetAggregatedIDPServices)
		r.Get("/stats/devices/{entitytype}/{entityid}", stats.GetAggregatedDevices)
		r.Get("/stats/policies", stats.GetPoliciesStats)

		r.Get("/stats/failedreasons/{entitytype}/{entityid}/{timeFilter}", stats.GetAggregatedFailedReasons)
		r.Get("/stats/loginhours/{entitytype}/{entityid}/{timeFilter}/{statusFilter}", stats.GetAggregatedLoginHours)
		r.Get("/stats/totalmanagedusers/{entitytype}/{entityid}", stats.GetTotalManagedUsers)
		r.Get("/stats/ips/{entitytype}/{entityid}/{timeFilter}/{statusFilter}", stats.GetIPAggs)
		r.Get("/stats/mapplot/{entitytype}/{entityid}/{timeFilter}/{statusFilter}", stats.GetMapPlotData)
		r.Get("/stats/todayauths/{entitytype}/{entityid}/{status}", stats.HexaEvents)
		r.Get("/stats/loginsbytype/{entitytype}/{entityid}/{timeFilter}/{statusFilter}", stats.GetLoginsByType)
		r.Get("/stats/total/{entitytype}/{entityid}/{timeFilter}", stats.GetSuccessAndFailedEvents)
		r.Get("/stats/ca", stats.GetCAStats)
		r.Get("/stats/appperms/{serviceID}", stats.GetServicePermStats)
		r.Get("/events/stats/{entitytype}/{entityid}/byday", stats.GetTotalLoginsByDate)

		r.Get("/system/status", system.SystemStatus)

		r.Post("/remote/auth/u2f", tfa.U2fHandler)

		//r.Post("/remote/auth/radius", services.RadiusLogin)

		//

		// Global System settings
		r.Get("/system/settings/all", system.GlobalSettings)
		r.Get("/system/security/rules", system.SecurityRules)
		r.Post("/system/security/rule/update", system.UpdateSecurityRule)
		r.Post("/system/settings/email/update", system.UpdateEmailSetting)
		r.Post("/system/settings/devicehygienecheck/update", system.UpdateDeviceHygieneSetting)
		r.Post("/system/settings/dynamicaccess/update", system.UpdateDynamicAccessSetting)
		r.Post("/system/settings/cloudproxy/access", system.StoreCloudProxyKey)

		r.Get("/system/welcome-note", system.WelcomeNote)

		// Identity Providers
		r.Post("/providers/uidp/create", uidp.CreateIdp)
		r.Post("/providers/uidp/update", uidp.UpdateIdp)
		r.Get("/providers/uidp/all", uidp.GetAllIdps)
		r.Post("/providers/uidp/generatescimtoken/{idpID}", uidp.GenerateSCIMAuthToken)
		r.Post("/providers/uidp/ldap/importusers", uidp.ImportLdapUsers)
		r.Post("/providers/uidp/activateordisable", uidp.ActivateOrDisableIdp)
		r.Get("/providers/uidp/users/all/{idpname}", uidp.GetAllUsersForIdp)
		r.Post("/providers/uidp/users/transfer", uidp.TransferUserToGivenIdp)

		r.Get("/providers/sidp/syncstatus/{vendorID}", sidp.GetSyncDetail)
		r.Post("/providers/sidp/syncnow/{vendorID}", sidp.SyncNow)

		//// TSxVault Operations
		r.Get("/providers/vault/tsxvault/key/{vendorID}", system.Getkey)
		r.Post("/providers/vault/tsxvault/store/key", system.StoreKey)
		r.Post("/providers/vault/tsxvault/init", system.TsxvaultInit)
		r.Delete("/providers/vault/tsxvault/reinit", system.ReInit)
		r.Get("/providers/vault/tsxvault/status", system.Status)
		r.Post("/providers/vault/tsxvault/decrypt", system.DecryptKey)

		r.Post("/providers/vault/credprov", system.UpdateCredProv)

		r.Post("/providers/ca/tsxca/init", ca.InitCA)
		r.Post("/providers/ca/tsxca/ssh/init/{type}", ca.InitSSHCA)
		r.Post("/providers/ca/tsxca/upload", ca.UploadCA)
		r.Get("/providers/ca/tsxca/http/detail", ca.GetHttpCADetail)
		r.Get("/providers/ca/tsxca/all", ca.GetAllCAs)
		r.Get("/providers/ca/tsxca/ssh/{type}", ca.DownloadSshCA)

	})

	// SCIM endpoint allows scim providers to manage users and user groups in TRASA.
	scimMiddleware := middlewares.SCIMAuth{}
	r.Route("/scim/v2", func(r chi.Router) {

		// 	User     /Users           GET (Section 3.4.1),   Retrieve, add,
		// 							  POST (Section 3.3),    modify Users.
		// 							  PUT (Section 3.5.1),
		// 							  PATCH (Section 3.5.2),
		// 							  DELETE (Section 3.6)
		// get single user detail from userID

		r.Use(scimMiddleware.Handler)
		r.Get("/Users/{userID}", uidp.SCIMGetSingleUser)
		r.Get("/Users", uidp.SCIMGetSingleUsersWithFilter)
		r.Post("/Users", uidp.SCIMCreateUser)
		r.Put("/Users/{userID}", uidp.SCIMPutSingleUser)
		r.Patch("/Users/{userID}", uidp.SCIMPatchSingleUser)
		r.Delete("/Users/{userID}", uidp.SCIMDeleteUser)
	})

	return r
}

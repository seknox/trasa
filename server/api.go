package main

import (
	"fmt"
	"net/http"

	"github.com/seknox/trasa/server/api/auth/serviceauth"

	"github.com/seknox/trasa/server/accessproxy/rdpproxy"
	"github.com/seknox/trasa/server/accessproxy/sshproxy"

	"github.com/seknox/trasa/server/api/crypt"

	"github.com/seknox/trasa/server/middlewares"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/accessmap"
	"github.com/seknox/trasa/server/api/auth"
	"github.com/seknox/trasa/server/api/auth/tfa"
	"github.com/seknox/trasa/server/api/devices"
	"github.com/seknox/trasa/server/api/groups"
	"github.com/seknox/trasa/server/api/idps"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/my"
	"github.com/seknox/trasa/server/api/notif"
	"github.com/seknox/trasa/server/api/policies"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/api/stats"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/api/users/passwordpolicy"
	"github.com/sirupsen/logrus"
)

// CoreAPIRoutes holds api route declarations for trasa-server
func CoreAPIRoutes(r *chi.Mux) *chi.Mux {

	r.Use(middlewares.Dumper{}.Handler)

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		logrus.Debug("NOT FOUND URL in core api: ", req.URL)
		FileServer(r, req.URL.Path)

	})

	r.Route("/auth", func(r chi.Router) {
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Reached not found in auth ")
			fmt.Println(r.URL)
		})

		r.Post("/identity", auth.LoginHandler)
		r.Post("/tfa", auth.TfaHandler)
		r.Delete("/logout", auth.LogoutHandler)

		r.Post("/enrol/ext", auth.EnrolBrowserExtension)
		r.Post("/accessproxy/http", serviceauth.AuthHTTPAccessProxy)

		r.Post("/crypto/kex", crypt.Kex)
		r.Post("/device/register", auth.RegisterUserDevice)
		r.Post("/device/ext/sync", auth.SyncExtension)

		r.Post("/agent/nix", serviceauth.AgentLogin)
		r.Post("/agent/win", serviceauth.AgentLogin)
		r.Post("/agent/checkconfig", services.CheckAppConfigs)
	})

	r.Route("/idp", func(r chi.Router) {
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Reached not found in idp ")
			fmt.Println(r.URL)
		})

		r.Post("/login", auth.LoginHandler)
		r.Post("/login/tfa", auth.TfaHandler)
		r.Delete("/logout", auth.LogoutHandler)

	})

	r.Route("/api/woa", func(r chi.Router) {
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Reached not found in File server ")
			fmt.Println(r.URL)
		})
		r.Get("/verify/{verifytoken}", my.VerifyAccount)
		r.Post("/setup/password", my.FirstTimePasswordSetup)

		r.Post("/proxy/http/getsession", serviceauth.GetHttpSession)
		r.Post("/proxy/http/rmsession", serviceauth.DestroyHttpSession)
		r.Post("/devices/brsrdetail", devices.GetBrsrDetails)

		r.Post("/forgotpass", my.ForgotPassword)

	})
	r.Get("/api/v1/my/download_file/get/{fileName}/{sskey}", my.FileDownloadHandler)

	// ///////////////////////////////////////////////////
	// *** DEPCRECATED ***
	// Below api only kept for older version compatibility
	r.Post("/api/v1/remote/auth/win", serviceauth.AgentLogin)
	r.Post("/api/v1/remote/auth/nix", serviceauth.AgentLogin)
	// r.Post("/api/v1/remote/auth/win", serviceauth.AgentLogin)
	// r.Post("/api/v1/remote/auth/checkconfig", services.CheckAppConfigs)
	// ////////////////////////////////////////////////////////

	r.Get("/api/v1/logs/livesessions", middlewares.SessionValidatorWS(logs.GetLiveSessions))
	r.Get("/api/v1/logs/vsessionlog", logs.GetVideoLog) //SessionValidator(GetSessionLog))

	authmiddleware := middlewares.AuthMiddleware{}
	inapptrailmiddleware := middlewares.InAppTrail{}
	gproxy := rdpproxy.NewProxy()

	r.Route("/trasagw", func(r chi.Router) {
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Reached not found in File server ")
			fmt.Println(r.URL)
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
			fmt.Println("Reached not found in File server ")
			fmt.Println(r.URL)
		})

		r.Use(authmiddleware.Handler)
		r.Use(inapptrailmiddleware.Handler)

		//public APIs
		//TODO @sshah verify these

		r.Post("/devicedetailpipe", devices.DeviceDetailPipe)
		r.Post("/passmydevicedetail", devices.PassMyDeviceDetail)

		//r.Post("/gateway/ext/sync", gateway.SyncExtension)

		r.Get("/my", my.GetMyDetail)
		r.Post("/my/forgotpass", my.ForgotPassword)
		r.Post("/my/forgotpasstfa", auth.TfaHandler)
		r.Post("/my/changepass", my.ChangePassword)
		r.Get("/my/services", my.GetMyServicesDetail)

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
		r.Get("/my/authmeta/{appID}/{username}", my.GetAuthMeta)

		//r.Get("/my/regreq", mdlwr.SessionValidator(Authorization(crypto.RegisterRequest)))
		//r.Post("/my/regresp", mdlwr.SessionValidator(Authorization(crypto.RegisterResponse)))
		//r.Post("/my/signreq", mdlwr.SessionValidator(Authorization(crypto.SignRequest)))
		//r.Post("/my/signresp", mdlwr.SessionValidator(Authorization(crypto.SignResponseHandler)))

		r.Get("/my/notifs", notif.GetPendingNotif)
		r.Post("/my/notif/resolve", notif.ResolveNotif)

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
		r.Get("/stats/appperms/{serviceID}", stats.GetAppPermStats)
		r.Get("/events/stats/{entitytype}/{entityid}/byday", stats.GetTotalLoginsByDate)

		r.Get("/system/status", system.SystemStatus)

		r.Post("/remote/auth/u2f", tfa.U2fHandler)

		//r.Post("/remote/auth/radius", services.RadiusLogin)

		//
		//// TSxVault Operations
		r.Get("/crypto/key/{vendorID}", system.Getkey)
		r.Post("/crypto/store/key", system.StoreKey)
		r.Post("/crypto/vault/init", system.TsxvaultInit)
		r.Delete("/crypto/vault/reinit", system.ReInit)
		r.Get("/crypto/vault/status", system.Status)
		r.Post("/crypto/vault/decrypt", system.DecryptKey)

		r.Post("/system/ca/init", crypt.InitCA)
		r.Post("/system/sshca/init/{type}", crypt.InitSSHCA)
		r.Post("/system/ca/upload", crypt.UploadCA)
		r.Get("/system/ca/detail", crypt.GetHttpCADetail)
		r.Get("/system/ca/all", crypt.GetAllCAs)
		r.Get("/system/ca/ssh/{type}", crypt.DownloadSshCA)

		// Global System settings
		r.Get("/system/settings/all", system.GlobalSettings)
		r.Get("/system/security/rules", system.SecurityRules)
		r.Post("/system/security/rule/update", system.UpdateSecurityRule)
		r.Post("/system/settings/passwordpolicy/update", system.UpdatePasswordPolicy)
		r.Get("/system/settings/passwordpolicy/enforce", passwordpolicy.EnforcePasswordPolicyNow)
		r.Post("/system/settings/email/update", system.UpdateEmailSetting)
		r.Post("/system/settings/sshcert/update", system.UpdateSSHCertSetting)
		r.Post("/system/settings/devicehygienecheck/update", system.UpdateDeviceHygieneSetting)
		r.Post("/system/settings/dynamicaccess/update", system.UpdateDynamicAccessSetting)

		// Identity Providers
		r.Post("/idp/external/create", idps.CreateIdp)
		r.Post("/idp/external/update", idps.UpdateIdp)
		r.Get("/idp/external/cloudiaas/syncstatus/{vendorID}", idps.GetSyncDetail)
		r.Post("/idp/external/cloudiaas/syncnow/{vendorID}", idps.SyncNow)
		r.Get("/idp/external/all", idps.GetAllIdps)
		r.Post("/idp/external/generatescimtoken/{idpID}", idps.GenerateSCIMAuthToken)

		r.Post("/idp/external/ldap/importusers", idps.ImportLdapUsers)
		r.Post("/idp/external/activateordisable", idps.ActivateOrDisableIdp)
		r.Get("/user/idp/{idpname}", idps.GetAllUsersForIdp)
		r.Post("/user/idp/transfer", idps.TransferUserToGivenIdp)

		//// we are using http session validator in this case
		//r.Post("/gateway/getpass", mdlwr.HttpSessionValidator(gateway.ValidateUserAndGetPass))

		//r.Get("/gateway/internalhosts", mdlwr.SessionValidator(Authorization(gateway.GetInternalHosts)))
		//r.Post("/gateway/updateinternalhosts", mdlwr.SessionValidator(Authorization(mdlwr.InAppTrailWrapper(gateway.UpdateInternalHosts, consts.UPDATE_HTTP_PROXY, true))))

	})

	return r
}

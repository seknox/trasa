package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/seknox/trasa/global"
	"github.com/seknox/trasa/models"

	"github.com/sirupsen/logrus"

	"github.com/seknox/trasa/utils"
)

type roleSwitch struct {
	orgAdmin        bool
	securityAnalyst bool
	userManager     bool
	selfUser        bool
}

// permissionChecker checks user role and returns true or false based on permitted actions
func permissionChecker(assignedRole string, requestEndpoint string, method string) bool {
	// we switch valid user role here.
	var switchUserRole roleSwitch
	allow := false

	//var roles roleType
	// our hardcoded roles. maybe can be done better ?
	roleTypes := []string{"orgAdmin", "securityAnalyst", "userManager", "selfUser"}
	// get user role from request context
	userAsignedRole := assignedRole

	userRole, status := identifyUserRole(userAsignedRole, roleTypes)
	if status == false {
		logrus.Errorf("error in identifyuserRole: %s\n", userRole)
		return false // invalid role
	}

	switch {
	case strings.Compare("orgAdmin", userRole) == 0:
		switchUserRole.orgAdmin = true
	case strings.Compare("securityAnalyst", userRole) == 0:
		switchUserRole.securityAnalyst = true
	case strings.Compare("userManager", userRole) == 0:
		switchUserRole.userManager = true
	default:
		switchUserRole.selfUser = true
	}

	//userRole.userManager = true

	//requestEndpoint :=  requestEndpoint  //urlParser(r.URL)
	endpoints := []string{"services", "stats", "user", "my", "iam", "gateway", "crypto", "org", "system", "groups", "globalsettings", "idp", "logs", "accessmap", "devices", "trasagw"}
	if validateEndpoint(requestEndpoint, endpoints) == true {
		switch {
		case switchUserRole.orgAdmin:
			if method == "POST" || method == "PATCH" || method == "GET" || method == "DELETE" {
				allow = true
			}
		case switchUserRole.securityAnalyst:
			if method == "GET" {
				allow = true
			}
		case switchUserRole.userManager:
			if strings.Compare(requestEndpoint, "user") == 0 {
				if method == "POST" || method == "PATCH" || method == "GET" || method == "DELETE" {
					allow = true

				}
			}
		case switchUserRole.selfUser:
			if strings.Compare(requestEndpoint, "my") == 0 {
				if method == "PATCH" || method == "GET" || method == "POST" {
					allow = true

				}
			}

		}

		if allow == false {
			return false
		}

		return true
	}
	return false
}

// secureDirectObjectReference is anti IDOR (Insecure Direct object Reference) to prevent malicious actors from playing things
// they are not supposed to. we do this by getting valid organization context from user session. since roles and permissions
// will be verified with permissionChecker function, we only need to verify if the entity user is trying to access/update/create/delete
// resources that are not related to scoped organization.
func secureDirectObjectReference() {

}

// Authorization has one of the most important security function for TRASA.
// It handles roles and permission for user specific to organization and sandbox crud operations between inter Ogranization.
// 1 - we extract request method
// 2 - we extract api endpoint from url path
// 3 - we extract user role and permissions
// 4 - we verify user role (based on group) to incoming request endpoint and method.
// 5 - we verify if the requested crud operation belongs to same organization. orgID derived from session ID is relevant source of truth here.
func Authorization(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get context value from sessionValidator
		uc := r.Context().Value("user").(models.UserContext)
		// get error struct
		var msg models.ErrorStrings
		isSaas := global.GetConfig().Platform.Base != "private"

		reqEndpoint := urlParser(r.URL)
		onpremEndpoints := []string{"system", "crypto", "upload_file", "download_file"}

		if isSaas && uc.Org.ID != uc.User.OrgID {
			for _, ep := range onpremEndpoints {
				if reqEndpoint == ep {

					utils.TrasaResponse(w, 403, "failed", "Authorization failed", "Authorization", nil)
					logrus.Trace(" api endpoint not authorized")
					return
				}
			}
		}

		permissionStatus := permissionChecker(uc.User.UserRole, reqEndpoint, r.Method)
		if permissionStatus == false {
			logrus.Trace("request validation. trying to elevate privilege??")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			msg.Status = "failed"
			msg.Reason = "Permission denied"
			msg.Intent = "request validation. trying to elevate privilege??"
			//logrus.ErrorLogs(msg)

			errmsg, _ := json.Marshal(msg)

			w.Write(errmsg)
			return
		}

		ctx := context.WithValue(r.Context(), "user", uc)
		next(w, r.WithContext(ctx))

	})

}

// urlParser parses incoming request URL and extracts api endpoint, which is
// 3 element from slice returned from url.Path seperated by "/".
func urlParser(urlPath *url.URL) string {
	path := strings.Split(urlPath.Path, "/")
	return path[3]
}

// validateEndpoint checks extracted api endpoint from urlParser to check if it is correct.
// if it turns out not-found, this must be valid security incident (attacker trying to brute force URL)
// return Not-found error on false.
func validateEndpoint(requestEndpoint string, endpoints []string) bool {
	for _, v := range endpoints {

		if strings.Compare(v, requestEndpoint) == 0 {
			return true
		}
	}
	return false
}

// identifyUserRole returns role value and a boolean value that will be used to switch user role scope
func identifyUserRole(role string, roleType []string) (string, bool) {
	for _, v := range roleType {

		if strings.Compare(v, role) == 0 {
			return v, true
		}
	}
	return "", false

}

/*
/Service
/user
/event

4 - crud
3 - ru
2 - r
1 = self(ru)

	userRoles.admin.apps = 3
	userRoles.admin.events = 3
	userRoles.admin.user = 3

	userRoles.analyst.apps = 1
	userRoles.analyst.events = 1
	userRoles.analyst.user = 2

	userRoles.user.apps = 0
	userRoles.user.events = 0
	userRoles.user.user = 1

*/

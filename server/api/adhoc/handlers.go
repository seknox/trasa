package adhoc

import (
	"encoding/json"
	"fmt"
	"github.com/seknox/trasa/server/api/notif"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type nilPolicy struct {
	Unassigned string `json:"unassigned"`
}

func AdhocReq(w http.ResponseWriter, r *http.Request) {
	var req models.AdhocPermission

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "AdhocReq", nil, nil)
		return
	}

	userContext := r.Context().Value("user").(models.UserContext)

	nep, err := time.LoadLocation(userContext.Org.Timezone)
	if err != nil {
		logrus.Error(err)
		return
	}

	req.RequestID = utils.GetRandomString(5)
	req.RequesterID = userContext.User.ID
	req.RequestedOn = time.Now().Unix()
	req.OrgID = userContext.User.OrgID
	req.SessionID = make([]string, 0) //append(req.SessionID, "zero")
	req.IsExpired = false
	req.AuthorizedPolicy = models.Policy{}

	// we check if active adhoc request for this user already exists in database. If found, we check if authorized perioed has not expired.
	// we return active response of not expired else we expire the active request and let user create new one.
	checkActiveAdhocReq, err := Store.GetActiveReqOfUser(userContext.User.ID, req.ServiceID, userContext.User.OrgID)
	if err == nil {
		if checkActiveAdhocReq.AuthorizedPeriod == 0 {
			utils.TrasaResponse(w, 200, "active", "you already have active adhoc request", "AdhocReq", checkActiveAdhocReq.AuthorizedPeriod)
			return
		} else {
			err := Store.Expire(checkActiveAdhocReq.RequestID, userContext.User.OrgID)
			if err != nil {
				//TODO check if trasaresponse/return is needed
				logrus.Error(err)
			}
		}

	}

	req.RequesterID = userContext.User.ID
	// storeRequestToDB
	if err = Store.Create(req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error generating request", "AdhocReq")
		return
	}

	serviceName, err := services.Store.GetServiceNameFromID(req.ServiceID, userContext.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid service ID", "AdhocReq")
		return
	}

	var inAppNotification models.InAppNotification
	inAppNotification.NotificationID = utils.GetRandomString(5)
	inAppNotification.OrgID = userContext.User.OrgID
	inAppNotification.NotificationLabel = "access-request"
	inAppNotification.NotificationText = fmt.Sprintf("User %s has requested access to %s.", userContext.User.Email, serviceName)
	inAppNotification.CreatedOn = time.Now().Unix() //time.Now().In(nep).Unix()
	inAppNotification.EmitterID = req.RequestID
	inAppNotification.IsResolved = false
	inAppNotification.ResolvedOn = 0
	inAppNotification.UserID = req.RequesteeID

	// we call notification store to store notification for this event in user id who is to be notified.
	err = notif.Store.StoreNotif(inAppNotification)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Debug(req.RequesteeID)
	requestee, err := users.Store.GetFromID(req.RequesteeID, userContext.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid requestee", "")
		return
	}

	uname := strings.Split(requestee.Email, "@")
	dashboardPath := fmt.Sprintf(`https://%`, global.GetConfig().Trasa.ListenAddr)
	dashLink := fmt.Sprintf("%s/control/access-requests", dashboardPath)
	var tmplt models.EmailAdhoc
	tmplt.Requester = userContext.User.Email
	tmplt.Requestee = uname[0]
	tmplt.ReceiverEmail = requestee.Email
	tmplt.DashLink = dashLink
	tmplt.App = serviceName
	tmplt.Reason = req.RequestTxt
	tmplt.Status = ""
	tmplt.Subject = "Adhoc Access Request"
	tmplt.CC = []string{}
	tmplt.Time = time.Now().In(nep).Format(time.RFC3339)
	tmplt.Req = true
	err = notif.Store.SendEmail(userContext.User.OrgID, consts.EMAIL_ADHOC, tmplt)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "success", "Adhoc request was success but we could not send email to administrator", "AdhocReq")
		return
	}

	utils.TrasaResponse(w, 200, "success", "request successfull", "AdhocReq")

}

func AdhocReqAssignedToMe(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	getRequests, err := Store.GetReqsAssignedToUser(userContext.User.ID, userContext.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error generating request", "AdhocReqAssignedToMe")
		return
	}

	utils.TrasaResponse(w, 200, "success", "request successfull", "AdhocReqAssignedToMe", getRequests)
}

func GrantOrDenyAdhoc(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	//nep, _ := time.LoadLocation("Asia/Kathmandu")
	var req models.AdhocPermission

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request format", "Could not respond to adhoc request")
		return
	}

	req.AuthorizedPolicy = models.Policy{}

	req.RequesteeID = userContext.User.ID
	req.OrgID = userContext.User.OrgID

	if req.IsAuthorized {
		req.IsExpired = false
		req.AuthorizedOn = time.Now().Unix()

	} else {
		req.IsExpired = true
		req.AuthorizedOn = 0
		req.AuthorizedPeriod = 0
	}

	// storeRequestToDB
	if err := Store.GrantOrReject(req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error generating request", "Could not respond to adhoc request")
		return
	}

	var inAppNotification models.InAppNotification
	inAppNotification.EmitterID = req.RequestID
	inAppNotification.OrgID = userContext.User.OrgID
	inAppNotification.IsResolved = true
	inAppNotification.ResolvedOn = time.Now().Unix()

	// we call notification store to store notification for this event in user id who is to be notified.
	err := notif.Store.UpdateNotif(inAppNotification)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "error updating notif", "Could not respond to adhoc request")
		return
	}

	adhoc, err := Store.GetAdhocDetail(req.RequestID, userContext.Org.ID)

	if err != nil {
		logrus.Errorf("GetAdhocDetail: %v", err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "Could not respond to adhoc request")
		return
	}

	uname := strings.Split(userContext.User.Email, "@")
	dashboardPath := fmt.Sprintf(`https://%`, global.GetConfig().Trasa.ListenAddr)
	dashLink := fmt.Sprintf("%s/control/access-requests", dashboardPath)

	tm := time.Unix(req.AuthorizedPeriod/1000, 0)

	status := "rejected"

	nep, err := time.LoadLocation(userContext.Org.Timezone)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid timezone", "Could not respond to adhoc request")
		return
	}
	till := tm.In(nep).Format(time.RFC3339)
	if req.IsAuthorized {
		status = "granted"
		//
	} else {
		till = "not authorized"
	}

	var tmplt models.EmailAdhoc
	tmplt.Requester = adhoc.RequesterEmail
	tmplt.Requestee = uname[0]
	tmplt.ReceiverEmail = adhoc.RequesterEmail
	tmplt.DashLink = dashLink
	tmplt.App = adhoc.ServiceName
	tmplt.Reason = req.RequestTxt
	tmplt.Status = status
	tmplt.Subject = "Your Adhoc Access Request Status"
	tmplt.CC = []string{}
	tmplt.Time = till
	tmplt.Req = false
	err = notif.Store.SendEmail(userContext.User.OrgID, consts.EMAIL_ADHOC, tmplt)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "success", "Adhoc request was success but we could not send email", fmt.Sprintf(`Adhoc request %s`, status))
		return
	}

	utils.TrasaResponse(w, 200, "success", "request successfull", fmt.Sprintf(`Adhoc request %s`, status))

}

func GetAllAdhoqRequests(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	getAllAdhocRequest, err := Store.GetAll(userContext.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "request history fetched", "GetAllAdhoqRequests", getAllAdhocRequest)
		return
	}

	utils.TrasaResponse(w, 200, "success", "request history fetched", "GetAllAdhoqRequests", getAllAdhocRequest)
}

func GetAdmins(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	//role := request.URL.Query().Get("userRole")

	admins, err := Store.GetAdmins(userContext.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "success", "fetched suers", "GetUsetsBasedOnUserRoleSecure", admins)

	}

	utils.TrasaResponse(w, 200, "success", "fetched suers", "GetUsetsBasedOnUserRoleSecure", admins)
}

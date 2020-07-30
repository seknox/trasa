package logs

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/seknox/trasa/models"
	"github.com/seknox/trasa/utils"
	"github.com/sirupsen/logrus"
)

// GetLoginEvents returns every loging events (failed and passed) with detailed log values.
// we can either check for idor here or query elasticsearch with both orgID and serviceID.
// opting to check idor here based on orgID.
func GetLoginEvents(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	entityType := chi.URLParam(r, "entitytype")
	entityID := chi.URLParam(r, "entityid")

	//TODO is CheckIfEntityIsWithinOrg necessary?

	if entityType == "org" {
		entityID = userContext.User.OrgID
	}

	events, err := Store.GetLoginEvents(entityType, entityID, userContext.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "could not get logs", "get org logs by page")
		return

	}
	utils.TrasaResponse(w, http.StatusOK, "success", "invalid size or page", "get org logs by page", events)
}

func GetLoginEventsByPage(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	entityType := chi.URLParam(r, "entitytype")
	entityID := chi.URLParam(r, "entityid")

	if entityType == "org" {
		entityID = userContext.User.OrgID
	}

	//orgID := uc.User.OrgID
	page, err1 := strconv.ParseInt(chi.URLParam(r, "page"), 10, 32)
	size, err2 := strconv.ParseInt(chi.URLParam(r, "size"), 10, 32)

	loc, err := time.LoadLocation(userContext.Org.Timezone)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to load location", "GetLoginEvents: loadlocation", err)
		return

	}

	//Date format 2020-05-18

	dateFromTime, err3 := time.Parse("2006-01-02", chi.URLParam(r, "dateFrom"))
	dateToTime, err4 := time.Parse("2006-01-02", chi.URLParam(r, "dateTo"))

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		utils.TrasaResponse(w, http.StatusBadRequest, "failed", "invalid size or page", "get org logs by page", nil)
		return
	}

	dateFrom := dateFromTime.In(loc).UnixNano()
	dateTo := dateToTime.In(loc).UnixNano()

	events, err := Store.GetLoginEventsByPage(entityType, entityID, userContext.Org.ID, int(page), int(size), dateFrom, dateTo)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "could not get logs", "get org logs by page")
		return

	}
	utils.TrasaResponse(w, http.StatusOK, "success", "", "get org logs by page", events)

}

func GetAllInAppTrails(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	page, err1 := strconv.ParseInt(chi.URLParam(r, "page"), 10, 32)
	size, err2 := strconv.ParseInt(chi.URLParam(r, "size"), 10, 32)

	dateFrom, err3 := strconv.ParseInt(chi.URLParam(r, "dateFrom"), 10, 32)
	dateTo, err4 := strconv.ParseInt(chi.URLParam(r, "dateTo"), 10, 32)

	if err3 != nil || err4 != nil {
		dateFrom = -1
		dateTo = -1
	}

	//dateFrom := chi.URLParam(r, "dateFrom")
	//dateTo := chi.URLParam(r, "dateTo")

	if err1 != nil || err2 != nil {
		page = 0
		size = 100
	}

	events, err := Store.GetOrgInAppTrails(userContext.Org.ID, int(page), int(size), dateFrom, dateTo)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not get logs", "GetAllInAppTrails")
		return
	}
	utils.TrasaResponse(w, 200, "success", "", "GetAllInAppTrails", events)
}

func GetLiveSessions(params models.ConnectionParams, uc models.UserContext, ws *websocket.Conn) {

	if uc.User.UserRole != "orgAdmin" {
		ws.WriteMessage(websocket.CloseMessage, nil)
		ws.Close()
		return
	}

	Store.ServeLiveSessions(ws)

}

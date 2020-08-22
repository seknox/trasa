package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"testing"
)

func GetTotpCode(secret string) string {
	_, t, _ := utils.CalculateTotp(secret)
	return t
}

// AddTestUserContext is a middleware that adds  mock userContext
func AddTestUserContext(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userContext := models.UserContext{
			User: &models.User{
				ID:         "13c45cfb-72ca-4177-b968-03604cab6a27",
				OrgID:      "153f7582-5ae2-46ba-8c1c-79ef73fe296e",
				UserName:   "root",
				FirstName:  "Bhargab",
				MiddleName: "",
				LastName:   "Acharya",
				Email:      "bhargab@seknox.com",
				Groups:     nil,
				UserRole:   "orgAdmin",
				Status:     true,
				IdpName:    "trasa",
			},
			Org: models.Org{
				ID:             "153f7582-5ae2-46ba-8c1c-79ef73fe296e",
				OrgName:        "Trasa",
				Domain:         "trasa.io",
				PrimaryContact: "",
				Timezone:       "Asia/Kathmandu",
				PhoneNumber:    "",
				CreatedAt:      0,
				PlatformBase:   "",
				License:        models.License{},
			},
			DeviceID:  "",
			BrowserID: "",
		}
		ctx := context.WithValue(r.Context(), "user", userContext)
		next(w, r.WithContext(ctx))

	})

}

// AddTestUserContextWS is a middleware that adds  mock userContext to ws handlers
func AddTestUserContextWS(next func(params models.ConnectionParams, uc models.UserContext, ws *websocket.Conn)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		conn, err := Mockupgrader.Upgrade(w, r, nil)
		if err != nil {
			logrus.Error(err)
			return
		}
		//defer conn.Close()

		//TODO use different generic model for session validation
		var params models.ConnectionParams
		err = conn.ReadJSON(&params)
		if err != nil {
			logrus.Error(err)
			conn.WriteMessage(1, []byte(err.Error()))
			conn.Close()
			return
		}

		userContext := models.UserContext{
			User: &models.User{
				ID:         "13c45cfb-72ca-4177-b968-03604cab6a27",
				OrgID:      "153f7582-5ae2-46ba-8c1c-79ef73fe296e",
				UserName:   "root",
				FirstName:  "Bhargab",
				MiddleName: "",
				LastName:   "Acharya",
				Email:      "bhargab@seknox.com",
				Groups:     nil,
				UserRole:   "orgAdmin",
				Status:     true,
				IdpName:    "trasa",
			},
			Org: models.Org{
				ID:             "153f7582-5ae2-46ba-8c1c-79ef73fe296e",
				OrgName:        "Trasa",
				Domain:         "trasa.io",
				PrimaryContact: "",
				Timezone:       "Asia/Kathmandu",
				PhoneNumber:    "",
				CreatedAt:      0,
				PlatformBase:   "",
				License:        models.License{},
			},
			DeviceID:  "",
			BrowserID: "",
		}

		params.UserIP = utils.GetIp(r)

		next(params, userContext, conn)

	})

}

func GetReqWithBody(t *testing.T, body interface{}) *http.Request {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	return req
}

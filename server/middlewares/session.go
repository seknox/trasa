package middlewares

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/nacl/secretbox"
)

// SessionValidator is a middleware that checks for csrf tokens and session cookies
func SessionValidator(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//get session cookie and csrf tokens
		sessionToken := r.Header.Get("X-SESSION")
		csrfToken := r.Header.Get("X-CSRF")

		userContext, err := getUserContext(sessionToken, csrfToken)
		if err != nil {
			logrus.Error(err)
			utils.TrasaResponse(w, 403, "failed", "failed to verify token", "SessionValidator", nil, nil)
			return
		}
		ctx := context.WithValue(r.Context(), "user", userContext)
		next(w, r.WithContext(ctx))

	})

}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		//TODO
		return true
	},
	Subprotocols: []string{"trasa", "guacamole", "livesessions", "xterm"},
}

// SessionValidator is a middleware that checks for csrf tokens and session cookies
func SessionValidatorWS(next func(params models.ConnectionParams, uc models.UserContext, ws *websocket.Conn)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
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

		//uc := r.Context().Value("user").(models.UserContext)
		uc, err := getUserContext(params.SESSION, params.CSRF)
		if err != nil {
			logrus.Debug(err)
			conn.WriteMessage(1, []byte("Invalid session. Try logging in again."))
			conn.WriteMessage(websocket.CloseMessage, nil)
			conn.Close()
			return
		}
		logrus.Trace(uc.DeviceID)
		params.UserIP = utils.GetIp(r)

		next(params, uc, conn)

	})

}

func getUserContext(sessionToken, csrfToken string) (models.UserContext, error) {
	var userContext models.UserContext

	if sessionToken == "" || sessionToken == "null" {
		return userContext, errors.New("no session token")
	}

	if csrfToken == "" || csrfToken == "null" {
		return userContext, errors.New("no csrf token")
	}

	// if we are here, it means cookies are present. we look for auth token with redis store
	userID, orgID, deviceID, browserID, authToken, err := redis.Store.GetSession(sessionToken)
	if err != nil {
		return userContext, errors.Errorf("get session redis: %v", err)
	}
	//	_ = dbSession
	// first lets check if provided cookie value does exist in redis
	if authToken == "" {
		return userContext, errors.New("invalid session token")
	}
	// lastly check for valid csrf tokens

	var decretkey [32]byte
	decodedkey, _ := base64.StdEncoding.DecodeString(authToken)
	copy(decretkey[:], decodedkey)
	var decryptNonce [24]byte
	decodedCsrf, _ := base64.StdEncoding.DecodeString(csrfToken)
	copy(decryptNonce[:], decodedCsrf[:24])
	_, ok := secretbox.Open(nil, []byte(decodedCsrf[24:]), &decryptNonce, &decretkey)
	if !ok {
		return userContext, errors.Errorf("error in decryption")
	}

	user, err := users.Store.GetFromID(userID, orgID)
	if err != nil {
		return userContext, errors.Errorf(`get user  from db: %v`, err)
	}

	org, err := orgs.Store.Get(user.OrgID)
	if err != nil {
		return userContext, errors.Errorf("get org from db: %v", err)

	}
	//logrus.Trace(deviceID)

	userContext.User = user
	userContext.Org = org
	userContext.DeviceID = deviceID
	userContext.BrowserID = browserID
	userContext.Org.PlatformBase = global.GetConfig().Platform.Base

	return userContext, nil

}

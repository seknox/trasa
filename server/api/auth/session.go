package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/nacl/secretbox"
)

//LogoutHandler handles logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	sessionToken, err := r.Cookie("X-SESSION")
	if sessionToken.Value == "" || err != nil {
		return
	}

	err = Store.Logout(sessionToken.Value)
	if err != nil {
		logrus.Error(err)
		return
	}

	xSESSION := http.Cookie{
		Name:     "X-SESSION",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
	}

	http.SetCookie(w, &xSESSION)

	utils.TrasaResponse(w, http.StatusOK, "success", "", "")
}

type userAuthSessionResp struct {
	User      models.User `json:"user"`
	CSRFToken string      `json:"CSRFToken"`
}

// type tokens struct {
// 	Session string `json:"session"`
// 	Csrf    string `json:"csrf"`
// }

func sessionResponse(uc models.UserContext) (sessionToken string, response userAuthSessionResp, err error) {

	var csrfToken string
	sessionToken, csrfToken, err = SetSession(uc)
	if err != nil {
		return
	}

	response.User = *uc.User
	response.CSRFToken = csrfToken

	return sessionToken, response, nil
}

// SetSession sets, encrypts and serializes session cookies and csrf tokens
func SetSession(uc models.UserContext) (string, string, error) {

	sessionKey := utils.GetRandomBytes(17)
	authKey := utils.GetRandomBytes(17)

	// Insert user session values in database
	orgusr := fmt.Sprintf("%v:%v", uc.Org.ID, uc.User.ID)

	encodedSession := hex.EncodeToString(sessionKey)
	encodedAuth := base64.StdEncoding.EncodeToString(authKey)

	uc.Org.PlatformBase = global.GetConfig().Platform.Base
	err := redis.Store.SetSessionWithUserContext(encodedSession, time.Second*900, encodedAuth, uc)
	if err != nil {
		return "", "", err
	}

	// create csrf token
	var encryptionKey [32]byte
	// we use sessionKey as secret key for encryption
	copy(encryptionKey[:], authKey)

	// generate random nonce
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return "", "", err
	}

	// we encrypt keyVal with sessionKey and set cipher text as csrf token.
	// To verify csrf token, we simply decrypt csrf token with session key
	// If decryption fails, it suggest csrf token has been tampered and is invalid.
	csrfToken := secretbox.Seal(nonce[:], []byte(orgusr), &nonce, &encryptionKey)

	//	csrfcookie := http.Cookie{Name: "x-csrf", Value: base64.StdEncoding.EncodeToString(csrfToken), Path: "/"}

	return encodedSession, base64.StdEncoding.EncodeToString(csrfToken), nil

}

func deleteSession(key string) bool {
	//fmt.Println("enter deletesession in sessions.go")
	err := redis.Store.Delete(key)
	if err != nil {
		return false
	}
	return true
}

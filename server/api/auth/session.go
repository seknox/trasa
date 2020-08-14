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
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/nacl/secretbox"
)

//LogoutHandler handles logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionToken := r.Header.Get("X-SESSION")
	if sessionToken == "" || sessionToken == "null" {
		return
	}

	err := Store.Logout(sessionToken)
	if err != nil {
		logrus.Error(err)
		return
	}
	utils.TrasaResponse(w, http.StatusOK, "success", "", "")
}

type userAuthSessionResp struct {
	User   models.User `json:"user"`
	Tokens tokens      `json:"tokens"`
}

type tokens struct {
	Session string `json:"session"`
	Csrf    string `json:"csrf"`
}

func sessionResponse(userDetails *models.User, deviceID, browserID string) (response userAuthSessionResp, err error) {

	var sessionToken, csrfToken string
	sessionToken, csrfToken, err = SetSession(userDetails.ID, userDetails.OrgID, deviceID, browserID)
	if err != nil {
		return
	}

	var tokenval tokens
	tokenval.Session = sessionToken
	tokenval.Csrf = csrfToken

	response.User = *userDetails
	response.Tokens = tokenval

	return response, nil
}

// SetSession sets, encrypts and serializes session cookies and csrf tokens
func SetSession(userID, orgID, deviceID, browserID string) (string, string, error) {

	sessionKey := utils.GetRandomBytes(17)
	authKey := utils.GetRandomBytes(17)

	// Insert user session values in database
	orgusr := fmt.Sprintf("%v:%v", orgID, userID)
	//fmt.Printf("orgusr value is %s\n", orgusr)
	//nep, _ := time.LoadLocation(timezone)

	//timerVal := time.Now().In(nep) //.Format(time.RFC3339)
	encodedSession := hex.EncodeToString(sessionKey)
	encodedAuth := base64.StdEncoding.EncodeToString(authKey)

	err := redis.Store.Set(encodedSession, time.Second*900,
		"userID", userID,
		"orgID", orgID,
		"deviceID", deviceID,
		"browserID", browserID,
		"auth", encodedAuth)
	if err != nil {
		return "", "", err
	}

	// generate cookies
	//sessoinCookie := http.Cookie{Name: "auth", Value: encodedSession, Path: "/"}
	//userCookie := http.Cookie{Name: "user", Value: keyVal}

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

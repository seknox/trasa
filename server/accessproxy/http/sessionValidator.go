package http

import (
	"encoding/base64"
	"fmt"

	"net/http"
	"strings"

	"github.com/seknox/trasa/server/api/redis"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/nacl/secretbox"
)

// tokenValidator validates extoken and sessionID.
// If tokens are not valid, tokenValidator generates and assign new one.
// Once validated, It is responsibility of tokenValidator to strip off headers
// before passing request to upstream handler.
func tokenValidator(r *http.Request, userName string, isSSO bool) error {

	//get session cookie and csrf tokens
	sessionToken := r.Header.Get("TRASA-X-SESSION")
	if sessionToken == "" || sessionToken == "null" {
		return fmt.Errorf("No TRASA-X-SESSION Header")
	}

	csrfToken := r.Header.Get("TRASA-X-CSRF")
	if csrfToken == "" || sessionToken == "null" {
		return fmt.Errorf("No TRASA-X-CSRF Header")
	}

	logger.Tracef("HTTP-Session-Req: %s , %s", sessionToken, csrfToken)

	// 2) validate the token. sesionID is stored in redis with key as sessionID:hostname
	// this key holds extoken value, and secret nacl seal that contains userID and orgID.
	// this userID and orgID is then verified with extoken values which also contains same.
	// If userID and orgID doesnot match, false is return. else tokenValidator returns true.
	orguser, sessionData, sRecord, err := redis.Store.GetHTTPGatewaySession(sessionToken) // (sessionToken, "user", "auth", "sessionRecord")
	logger.Tracef("HTTP-Session-Redis: %s , %s , %s , %v", orguser, sessionData, sRecord, err)
	if err != nil || orguser == "" {
		return fmt.Errorf("Invalid Session")
	}

	session := strings.Split(sessionData, ":")

	secretKey := session[0]

	var decretkey [32]byte
	decodedkey, _ := base64.StdEncoding.DecodeString(secretKey)
	copy(decretkey[:], decodedkey)
	var decryptNonce [24]byte
	decodedCsrf, _ := base64.StdEncoding.DecodeString(csrfToken)
	copy(decryptNonce[:], decodedCsrf[:24])
	decrypted, ok := secretbox.Open(nil, []byte(decodedCsrf[24:]), &decryptNonce, &decretkey)
	if !ok {
		return fmt.Errorf("Invalid CSRF Token")

	}

	_ = decrypted

	passwordManAndLogger(r, sessionToken, csrfToken, userName, isSSO, sRecord)

	// If session validation is Okay, we remove these headers before sending to upstream server.
	for header := range r.Header {
		if header == "TRASA-X-SESSION" {
			r.Header.Del(header)
		}

		if header == "TRASA-X-CSRF" {
			r.Header.Del(header)
		}

		if header == "TRASA-X-USER" {
			r.Header.Del(header)
		}

	}

	return nil

}

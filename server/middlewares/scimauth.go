package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/providers/uidp"
	"github.com/seknox/trasa/server/api/providers/vault"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// SCIMTokenValidator validates incoming authorization token for scim requests
// refer to handler GenerateSCIMAuthToken to check how the token is generated.
func SCIMTokenValidator(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || authHeader == "null" {
			logrus.Error("no scim authorization header")
			utils.TrasaResponse(w, 403, "failed", "no Authorization header", "SCIMTokenValidator", nil, nil)
			return
		}

		authorization := strings.Split(authHeader, " ")
		if len(authorization) < 2 {
			logrus.Error("no authorization token")
			utils.TrasaResponse(w, 403, "failed", "no Authorization token", "SCIMTokenValidator", nil, nil)
			return
		}

		// devode base64 encoded token
		decodedkey, err := utils.DecodeBase64(authorization[1])
		if err != nil {
			logrus.Error("failed to decode token")
			utils.TrasaResponse(w, 403, "failed", "invalid key", "SCIMTokenValidator", nil, nil)
			return
		}

		orgPass := strings.Split(string(decodedkey), ":")

		// fetch and check from database
		keyVal, err := vault.Store.GetKeyOrTokenWithKeyval(orgPass[0], consts.SCIMKEY)
		if err != nil {
			logrus.Error("failed to retreive token")
			utils.TrasaResponse(w, 403, "failed", "invalid key", "SCIMTokenValidator", nil, nil)
			return

		}

		err = bcrypt.CompareHashAndPassword([]byte(keyVal.KeyVal), []byte(string(decodedkey)))

		if err != nil {
			logrus.Error("invalid token")
			utils.TrasaResponse(w, 403, "failed", "invalid key", "SCIMTokenValidator", nil, nil)
			return

		}

		// if we are here means token is valid.
		// fetch org detail from orgID and store it in scimContext

		orgdetailfromDB, err := orgs.Store.Get(orgPass[0])
		if err != nil {
			logrus.Error("invalid orgtoken")
			utils.TrasaResponse(w, 403, "failed", "invalid key", "SCIMTokenValidator", nil, nil)
			return
		}

		idpDetail, err := uidp.Store.GetByID(orgdetailfromDB.ID, keyVal.KeyID)

		var scmProv models.ScimContext
		scmProv.OrgID = orgdetailfromDB.ID
		scmProv.Orgname = orgdetailfromDB.OrgName
		scmProv.IdpID = idpDetail.IdpID
		scmProv.IdpName = idpDetail.IdpName
		scmProv.TimeZone = orgdetailfromDB.Timezone

		ctx := context.WithValue(r.Context(), "scimprov", scmProv)
		next(w, r.WithContext(ctx))

	})

}

package redis

import (
	"fmt"
	"time"

	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/utils"
)

// SetVerifyToken takes userID, orgID and sets unique token in redis which can be used to verify certain tasks
// such as password setup token.
func SetVerifyToken(orgID, userID string) (string, error) {
	// this will allow us to verify user account as well and let the user setup password as soon as the token is validated.
	verifyToken := utils.GetRandomString(12)

	// config fie should provide full url scheme

	err := Store.Set(
		verifyToken,
		consts.TOKEN_EXPIRY_SIGNUP,
		"orguser", fmt.Sprintf("%s:%s", orgID, userID),
		"intent", string(consts.VERIFY_TOKEN_CHANGEPASS),
		"createdAt", time.Now().String(),
	)

	if err != nil {

		return "", err
	}

	return verifyToken, nil
}

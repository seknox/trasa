package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/api/notif"
	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
)

// forgotPassTfaResp initiates forget password sequence after tfa validation
func forgotPassTfaResp(userDetails models.User) (err error) {
	normalizedEmail := strings.ToLower(userDetails.Email)

	dashboardPath := global.GetConfig().Trasa.Dashboard

	// while previously we created temporary password for users and send it to them via email,
	// we now will generate a short lived token which will be presented in a link.
	// this will allow us to verify user account as well and let the user setup password as soon as the token is validated.
	verifyToken := utils.GetRandomString(12)

	// config fie should provide full url scheme
	verifyURL := fmt.Sprintf("%s/woa/verify#token=%s", dashboardPath, verifyToken)

	orgUser := fmt.Sprintf("%s:%s", userDetails.OrgID, userDetails.ID)

	err = redis.Store.Set(verifyToken, consts.TOKEN_EXPIRY_CHANGEPASS,
		"orguser", orgUser,
		"intent", string(consts.VERIFY_TOKEN_CHANGEPASS),
		"createdAt", time.Now().String())

	if err != nil {
		return errors.Errorf("failed to store verify token: %v", err)
	}

	var tmplt models.EmailUserCrud
	tmplt.ReceiverEmail = normalizedEmail
	tmplt.Username = userDetails.FirstName
	tmplt.VerifyUrl = verifyURL
	tmplt.NewM = false

	err = notif.Store.SendEmail(userDetails.OrgID, consts.EMAIL_USER_CRUD, tmplt)
	if err != nil {
		return errors.Errorf("could not send password email: %v", err)
	}

	return nil
}

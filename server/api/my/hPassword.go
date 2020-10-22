package my

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/auth"
	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"github.com/trustelem/zxcvbn"
	"golang.org/x/crypto/bcrypt"
)

//GetMyDetail returns current user details
func GetMyDetail(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	utils.TrasaResponse(w, 200, "success", "my", "my", userContext)
}

type forgotPassReq struct {
	Email string `json:"email"`
}

//ForgotPassword starts forgot password process
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req forgotPassReq
	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "email doesn't exist", "ForgotPassword")
		return
	}

	user, err := auth.Store.GetLoginDetails(req.Email, "")
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "email doesn't exist", "Reset Password")
		return
	}

	orgUser := fmt.Sprintf("%s:%s", user.OrgID, user.ID)

	// if we are here, it means username password validation succeeded. we will generate a unique token attached to this user.
	// and send this token in response. this token will be used to validate tfa request and retrieve userID for this user.
	// this userID will be used to retrieve user detail and send in response to successful authentication.
	token := utils.GetRandomString(15)
	//TODO @sshah is the intent corrent here?
	err = redis.Store.Set(token, time.Second*400, "orgUser", orgUser, "intent", "AUTH_REQ_FORGOT_PASS")
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "second step verification required", "")
		return
	}
	// At this point, we shold now only return TfaRequired response.
	utils.TrasaResponse(w, http.StatusOK, "success", "second step verification required", consts.AUTH_RESP_TFA_REQUIRED, token)

}

type setPassword struct {
	Token     string `json:"token"`
	Password  string `json:"password"`
	CPassword string `json:"cpassword"`
}

// FirstTimePasswordSetup used in case of forget password and after user account is created for first time.
// Any change password process in active user session should use ChangePassword instead
func FirstTimePasswordSetup(w http.ResponseWriter, r *http.Request) {
	var req setPassword

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to validate input", "ChangePassword")
		return
	}

	if strings.Compare(req.Password, req.CPassword) != 0 {
		utils.TrasaResponse(w, 200, "failed", "password mismatch", "ChangePassword")
		return
	}

	res, err := redis.Store.MGet(req.Token, "orguser", "intent")
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid token", "ChangePassword")
		return
	}

	// verifyIntent
	if string(consts.VERIFY_TOKEN_CHANGEPASS) != res[1] {
		utils.TrasaResponse(w, 200, "failed", "invalid token", "ChangePassword")
		return
	}

	orgUser := strings.Split(res[0], ":")

	err = updatePassword(orgUser[1], orgUser[0], req)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", err.Error(), "ChangePassword")
		return
	}

	utils.TrasaResponse(w, 200, "success", "password updated", "Change password", nil)

}

// ChangePassword is same as FirstTimePasswordSetup except for this handler is passed through session validation and user details are
// taked from user session Context
func ChangePassword(w http.ResponseWriter, r *http.Request) {

	uc := r.Context().Value("user").(models.UserContext)

	var req setPassword

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "json parse error", "ChangePassword")
		return
	}

	if strings.Compare(req.Password, req.CPassword) != 0 {
		utils.TrasaResponse(w, 200, "failed", "password mismatch", "ChangePassword")
		return
	}

	err := updatePassword(uc.User.ID, uc.User.OrgID, req)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", err.Error(), "ChangePassword")
		return
	}

	utils.TrasaResponse(w, 200, "success", "password updated", "Change password", nil)
}

func updatePassword(userID, orgID string, changePassReq setPassword) error {

	err := redis.Store.VerifyIntent(changePassReq.Token, consts.VERIFY_TOKEN_CHANGEPASS)
	if err != nil {
		return fmt.Errorf("token is invalid")
	}

	passwordStreanght := zxcvbn.PasswordStrength(changePassReq.Password, nil)

	// first, we retrieve password policy from global settings
	gsetting, err := system.Store.GetGlobalSetting(orgID, consts.GLOBAL_PASSWORD_CONFIG)
	if err != nil {
		return fmt.Errorf("unable to get global setting")
	}

	var policy models.PasswordPolicy
	err = json.Unmarshal([]byte(gsetting.SettingValue), &policy)
	if err != nil {
		return fmt.Errorf("invalid global setting found")
	}

	if policy.EnforceStrongPass == true {
		if passwordStreanght.Score < policy.ZxcvbnScore {
			return fmt.Errorf("Low entropy. Choose strong password")
		}
	}

	if len(changePassReq.Password) < policy.MinimumChars {
		return fmt.Errorf("does not match maximum character requirement")
	}

	hashedpass, err := bcrypt.GenerateFromPassword([]byte(changePassReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Failed to Generate secure hash. Contact Administrator")
	}

	///////////////////// Begin PassState checker  ///////////////////
	//  new password should not match existing one.
	// Should not be same as last 4 pass.?
	//////////////////// End of PassState checker   ////////////////

	// set user password
	err = users.Store.UpdatePassword(userID, string(hashedpass))
	if err != nil {
		return fmt.Errorf("Failed to store password. Contact your administrator")
	}

	// once password is updated, we store this hashed in passwordState database as lastpass.
	err = users.Store.UpdatePasswordState(userID, orgID, string(hashedpass), time.Now().Unix())
	if err != nil {
		logrus.Error(err)
		//TODO @sshah return here???
	}

	// validate token with redis. this action should also retrieve user email address and user id so that we can match user password
	// with respective account.

	err = redis.Store.Delete(changePassReq.Token)
	if err != nil {
		// TODO @bharg3se if it fails here means trasa failed to remove change password token from redis. What should we do here? alert admin?
		return fmt.Errorf("password updated but failed to process further. Contact Administrator")
	}

	// if password is updates, we check if user has pending password change policy.
	// If yes, we delete policy row from database.
	err = users.Store.DeleteActivePolicy(userID, orgID, consts.ChangePassword)
	if err != nil {
		logrus.Error(err)
		return fmt.Errorf("password is updated but system failed to update policy")
	}

	return nil
}

// VerifyAccount does two things. validate verify account link. If the link is valid and is not expired,
// respond with success value with intent of setup password. Dashboard ui should then present set password dialogue and send request to
// FirstTimePasswordSetup handler.
// Format for validation link: http://trasa.seknox.com/verify/{link}
func VerifyAccount(w http.ResponseWriter, r *http.Request) {

	// Get token value from request
	verifyToken := chi.URLParam(r, "verifytoken")

	//TODO
	// validate token with redis
	err := redis.Store.VerifyIntent(verifyToken, consts.VERIFY_TOKEN_CHANGEPASS)
	if err != nil {
		logrus.Debug(err)
		utils.TrasaResponse(w, 200, "failed", "token not valid", "setpassword", verifyToken)
		return
	}

	// once the token is verified, dont just delete it yet. send http redirect response with this token as data value.
	// same token will be used to setup users password for first time login.

	// initiate response. we dont redirect client here. rather send 200 success request. client shold parse this response value
	// and dynamically create client page for password setting including verify token.
	utils.TrasaResponse(w, 200, "success", "token validated", "setpassword", verifyToken)

}

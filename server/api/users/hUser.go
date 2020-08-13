package users

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/notif"
	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

type userDetails struct {
	User           models.User              `json:"user"`
	UserAccessMaps []models.AccessMapDetail `json:"userAccessMaps"`
	UserDevices    []models.UserDevice      `json:"userDevices"`
	UserGroups     []models.Group           `json:"userGroups"`
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("Request Received")
	uc := r.Context().Value("user").(models.UserContext)
	userID := chi.URLParam(r, "userID")

	var resp userDetails

	user, err := Store.GetFromID(userID, uc.User.OrgID)
	if err != nil {
		logrus.Errorf("get user: %v", err)
		utils.TrasaResponse(w, 200, "failed", "failed to get user details", "GetSingleUserDetail", resp)
		return
	}
	resp.User = *user

	accessmaps, err := Store.GetAccessMapDetails(userID, uc.Org.ID)
	if err != nil {
		logrus.Errorf("get access map: %v", err)
		utils.TrasaResponse(w, 200, "failed", "failed to get access map", "GetSingleUserDetail", resp)
		return
	}

	resp.UserAccessMaps = accessmaps

	udevices, err := Store.GetAllDevices(userID, uc.Org.ID)
	if err != nil {
		logrus.Errorf("get user devices: %v", err)
		utils.TrasaResponse(w, 200, "failed", "failed to get devices", "GetSingleUserDetail", resp)
		return
	}
	resp.UserDevices = udevices

	ugroups, err := Store.GetGroups(userID, uc.Org.ID)
	if err != nil {
		logrus.Errorf("get user groups: %v", err)
		utils.TrasaResponse(w, 200, "failed", "failed to get user groups", "GetSingleUserDetail", resp)
		return
	}
	resp.UserGroups = ugroups

	utils.TrasaResponse(w, 200, "success", "success", "GetSingleUserDetail", resp)
}

// GetAllUsers returns json array of user list.
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	val, err := Store.GetAll(uc.Org.ID)

	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "users not fetched", "GetUsersAll")
		return
	}
	utils.TrasaResponse(w, 200, "success", "users fetched", "GetUsersAll", val)
}

type createUserReq struct {
	User           models.UserWithPass `json:"user"`
	PasswordMethod string              `json:"passMethod"`
}
type createUserResp struct {
	User             models.User `json:"user"`
	ConfirmationLink string      `json:"confirmLink"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	// declare request struct
	var request createUserReq

	// parse json value into struct
	if err := utils.ParseAndValidateRequest(r, &request); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "CreateUser", nil, nil)
		return
	}

	// get userContext
	uc := r.Context().Value("user").(models.UserContext)
	request.User.ID = utils.GetUUID()
	request.User.CreatedAt = time.Now().Unix()
	prepareUserStruct(&request, uc)
	// If request has password method selfPassSetup, we generate activation link which will reqiure user to setup their own password.
	// If password method is autoGenPass, we generate password and return it in response.

	// generate password hash
	// hashedpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	panic(err)
	// }

	// while previously we created temporary password for users and send it to them via email,
	// we now will generate a short lived token which will be presented in a link.
	// this will allow us to verify user account as well and let the user setup password as soon as the token is validated.
	verifyToken := utils.GetRandomString(12)

	// config fie should provide full url scheme

	err := redis.Store.Set(
		verifyToken,
		consts.TOKEN_EXPIRY_SIGNUP,
		"orguser", fmt.Sprintf("%s:%s", uc.Org.ID, request.User.ID),
		"intent", string(consts.VERIFY_TOKEN_CHANGEPASS),
		"createdAt", time.Now().String(),
	)

	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not set token", fmt.Sprintf(`user "%s" not created`, request.User.Email))
		return
	}

	tmplt := newUserEmailTemplate(&request.User, verifyToken)

	errEmail := notif.Store.SendEmail(uc.Org.ID, consts.EMAIL_USER_CRUD, tmplt)
	if errEmail != nil {
		logrus.Error(errEmail)
	}

	err = Store.Create(&request.User) //createUser(&user)
	if err != nil {
		reason := utils.GetConstraintErrorMessage(err)
		utils.TrasaResponse(w, 200, "failed", reason, fmt.Sprintf(`user "%s" not created`, request.User.Email))
		return
	}

	var resp createUserResp
	resp.User = models.CopyUserWithoutPass(request.User)
	resp.ConfirmationLink = tmplt.VerifyUrl

	msg := `User created`
	if errEmail != nil {
		msg = msg + " but email not sent"
	}

	utils.TrasaResponse(w, 200, "success", msg, fmt.Sprintf(`user "%s" created`, request.User.Email), resp)

	constName := consts.CREATE_USER
	if request.User.UserRole == "orgAdmin" {
		constName = consts.CREATE_ADMIN_USER
	}

	go notif.CheckAndFireSecurityRule(uc.Org.ID, constName, request.User.Email)

}

// UpdateUser updates TRASA user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var request createUserReq

	if err := utils.ParseAndValidateRequest(r, &request); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "UpdateUser", nil, nil)
		return
	}
	uc := r.Context().Value("user").(models.UserContext)

	prepareUserStruct(&request, uc)

	// generate password hash
	// if user.Password != "" {
	// 	hashedpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	password := string(hashedpass)
	// 	//fmt.Println("hashed pass", password)
	// 	err = Store.UpdatePassword(user.ID, password)
	// 	if err != nil {
	// 		utils.TrasaResponse(w, 200, "failed", "failed updating user password.", "user  not updated")
	// 		return
	// 	}

	// }

	err := Store.Update(models.CopyUserWithoutPass(request.User)) //createUser(&user)
	if err != nil {
		reason := utils.GetConstraintErrorMessage(err)
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", reason, "user not updated")
		return
	}
	utils.TrasaResponse(w, 200, "success", "", fmt.Sprintf(`user %s updated`, request.User.Email))

	//TODO @sshah
	// check and fire security rule. Here only interesting if user is admin or is granted admin privilege
	userDetail, err := Store.GetFromID(uc.User.ID, uc.User.OrgID)
	if userDetail.UserRole == "selfUser" && request.User.UserRole == "orgAdmin" {
		go notif.CheckAndFireSecurityRule(uc.Org.ID, consts.GRANT_ADMIN_PRIVILEGE, request.User.Email)
		return
	}
	if request.User.UserRole == "orgAdmin" {
		go notif.CheckAndFireSecurityRule(uc.Org.ID, consts.ADMIN_PROFILE_EDITED, request.User.Email)

	}

}

// DeleteUser should be atomic transaction - a single user delete call should delete users detail from
// every database tables.
//TODO make delete user atomic
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	userID := chi.URLParam(r, "userID")
	uc := r.Context().Value("user").(models.UserContext)
	email, userRole, err := Store.Delete(userID, uc.Org.ID) //createUser(&user)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed deleting user", "delete user", nil)
		return
	}

	err = Store.DeleteAllUserAccessMaps(userID, uc.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed deleting user.", "delete user", nil)
		return
	}

	err = Store.DeregisterUserDevices(userID, uc.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed deleting user.", "delete user", nil)
		return
	}

	utils.TrasaResponse(w, http.StatusOK, "success", "successfully deleted user", fmt.Sprintf(`user "%s" deleted`, email), nil)

	//TODO @sshah
	_ = userRole
	if userRole == "orgAdmin" {
		go notif.CheckAndFireSecurityRule(uc.Org.ID, consts.DELETE_ADMIN_USER, email)

	} else {
		go notif.CheckAndFireSecurityRule(uc.Org.ID, consts.DELETE_USER, email)
	}

}

func newUserEmailTemplate(user *models.UserWithPass, token string) models.EmailUserCrud {
	dashboardPath := global.GetConfig().Trasa.Dashboard
	verifyURL := fmt.Sprintf("%s/woa/verify#token=%s", dashboardPath, token)
	var tmplt models.EmailUserCrud
	tmplt.ReceiverEmail = user.Email
	tmplt.Username = user.FirstName
	tmplt.VerifyUrl = verifyURL
	tmplt.NewM = true
	return tmplt
}

func prepareUserStruct(request *createUserReq, uc models.UserContext) {

	request.User.Email = strings.ToLower(request.User.Email)
	request.User.UserName = strings.ToLower(request.User.UserName)
	request.User.UpdatedAt = time.Now().Unix()
	request.User.OrgID = uc.Org.ID

	if request.User.UserRole == "" {
		request.User.UserRole = "selfUser"
	}

	request.User.IdpName = "trasa"
	request.User.ExternalID = request.User.ID

	return
}

func GetGroupsAssignedToUser(w http.ResponseWriter, r *http.Request) {

	userContext := r.Context().Value("user").(models.UserContext)
	userID := chi.URLParam(r, "userID")
	groups, err := Store.GetGroups(userID, userContext.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "could not get groups", "GetGroupsAssignedToUser", nil)
		return

	}
	utils.TrasaResponse(w, http.StatusOK, "success", "", "GetGroupsAssignedToUser", groups)

}

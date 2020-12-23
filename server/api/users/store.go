package users

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/models"
)

//GetFromID returns user details from user ID
func (s userStore) GetFromID(userID, orgID string) (user *models.User, err error) {
	user = &models.User{}
	err = s.DB.QueryRow(`
			SELECT org_id, id,username, first_name,middle_name, last_name, email, user_role,status, created_at, updated_at, idp_name,
			       ARRAY(select name from groups JOIN user_group_maps ugm on groups.id = ugm.group_id WHERE user_id=$1)
			FROM users 
			WHERE id = $1 AND org_id=$2`,
		userID, orgID).
		Scan(&user.OrgID, &user.ID, &user.UserName, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.UserRole, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.IdpName, pq.Array(&user.Groups))

	return
}

//GetFromWithLimit returns user details from with supplied limit
func (s userStore) GetFromWithLimit(orgID string, limit int) (user *models.User, err error) {
	user = &models.User{}
	err = s.DB.QueryRow("SELECT org_id, id,username, first_name,middle_name, last_name, email, user_role,status, created_at, updated_at, idp_name FROM users WHERE org_id=$1 LIMIT $2", orgID, limit).
		Scan(&user.OrgID, &user.ID, &user.UserName, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.UserRole, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.IdpName)

	return
}

// GetFromTRASAID returns user details from user trasaID (username or email address)
func (s userStore) GetFromTRASAID(trasaID, orgID string) (*models.User, error) {

	isTrasaIDEmail := strings.Contains(trasaID, "@")

	//TODO use domain

	sqlStr := ``

	if isTrasaIDEmail {
		sqlStr = `SELECT users.org_id, users.id, first_name, email, idp_name, user_role, status 
				FROM users
				JOIN org ON users.org_id=org.id
				WHERE users.email=$1 AND users.org_id=$2`
	} else {
		sqlStr = `SELECT users.org_id, users.id, first_name, email, idp_name, user_role, status
				FROM users
				JOIN org ON users.org_id=org.id
				WHERE users.username=$1 AND users.org_id=$2`
	}

	var user models.User
	err := s.DB.QueryRow(sqlStr, trasaID, orgID).Scan(&user.OrgID, &user.ID, &user.FirstName, &user.Email, &user.IdpName, &user.UserRole, &user.Status)

	return &user, err
}

//GetAll returns all users of an organization
func (s userStore) GetAll(orgID string) ([]models.User, error) {
	var users = make([]models.User, 0)

	rows, err := s.DB.Query(`SELECT users.org_id, users.id,  username, first_name, middle_name,
								   last_name, email, user_role,
								   users.status AND COALESCE(idp.is_enabled,true) as status,
								   created_at, updated_at, users.idp_name
							FROM users
							LEFT JOIN idp  on users.idp_name = idp.name WHERE users.org_id = $1`, orgID)

	if err != nil {
		return users, err
	}

	defer rows.Close()
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.OrgID, &user.ID,
			&user.UserName, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.UserRole, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.IdpName)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, err

}

//GetByLimit returns all users of an organization based on limit. Only count is supported for now.
func (s userStore) GetByLimit(orgID string) ([]models.User, error) {
	var users = make([]models.User, 0)

	rows, err := s.DB.Query(`SELECT users.org_id, users.id,  username, first_name, middle_name,
								   last_name, email, user_role,
								   users.status AND COALESCE(idp.is_enabled,true) as status,
								   created_at, updated_at, users.idp_name
							FROM users
							LEFT JOIN idp  on users.idp_name = idp.name WHERE users.org_id = $1`, orgID)

	if err != nil {
		return users, err
	}

	defer rows.Close()
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.OrgID, &user.ID,
			&user.UserName, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.UserRole, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.IdpName)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, err

}

//GetAdminEmails returns email of all admins of an organization
func (s userStore) GetAdminEmails(orgID string) ([]string, error) {
	var users = make([]string, 0)

	rows, err := s.DB.Query(`SELECT users.email
							FROM users
							 WHERE users.user_role = $1 AND users.org_id = $2`, "orgAdmin", orgID)

	if err != nil {
		return users, err
	}

	defer rows.Close()
	for rows.Next() {
		var user string
		err := rows.Scan(&user)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, err

}

//Delete deletes a user and returns its email and user role
func (s userStore) Delete(userID, orgID string) (email string, userRole string, err error) {

	err = s.DB.QueryRow(`DELETE FROM users WHERE id = $1 AND org_id=$2 RETURNING email, user_role`,
		userID, orgID).
		Scan(&email, &userRole)

	return email, userRole, err
}

//Update updates UserName, FirstName, MiddleName, LastName, Email, UserRole, UpdatedAt and  Status of given user based in user ID
func (s userStore) Update(user models.User) error {
	_, err := s.DB.Exec(`UPDATE users SET username = $1, first_name = $2, middle_name = $3, last_name = $4, email = $5, user_role = $6, updated_at  = $7,status=$8 WHERE id = $9;`,
		user.UserName, user.FirstName, user.MiddleName, user.LastName, user.Email, user.UserRole, user.UpdatedAt, user.Status, user.ID)

	return err
}

// UpdateStatus change active or disabled status of user.
func (s userStore) UpdateStatus(state bool, userID, orgID string) error {
	_, err := s.DB.Exec(`UPDATE users SET status = $1 WHERE id = $2 AND org_id = $3;`,
		state, userID, orgID)

	if err != nil {
		return err
	}
	return nil
}

//UpdatePassword updates password of given user. It expects password to be already hashed
func (s userStore) UpdatePassword(userID, password string) error {
	_, err := s.DB.Exec(`UPDATE users SET password = $1 WHERE id = $2;`, password, userID)

	return err
}

// UpdatePasswordState maintains state of user password history
func (s userStore) UpdatePasswordState(userID, orgID, oldpassword string, time int64) error {
	op := []string{oldpassword}
	var passState models.PasswordState
	// first check if row already exists for userid orgid. If not then Insert OR Update
	err := s.DB.QueryRow(`SELECT user_id, org_id,  last_passwords, last_updated FROM password_state WHERE user_id = $1 and org_id = $2`,
		userID, orgID).
		Scan(&passState.UserID, &passState.OrgID, pq.Array(&passState.LastPasswords), &passState.LastUpdated)
	if errors.Is(err, sql.ErrNoRows) {
		// this means row does not exists we should insert row
		_, err := s.DB.Exec(`INSERT into password_state (user_id, org_id,  last_passwords, last_updated)
		 values($1, $2, $3, $4);`,
			userID, orgID, pq.Array(op), time)

		if err != nil {
			return fmt.Errorf("failed to  insert password state: %v", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to  get password state: %v", err)
	}

	passState.LastPasswords = append(passState.LastPasswords, oldpassword)

	// if we reach here means row exists and we should update it.
	_, err = s.DB.Exec(`UPDATE password_state SET last_passwords =$1, last_updated = $2 WHERE user_id = $3 AND org_id = $4;`,
		pq.Array(passState.LastPasswords), time, userID, orgID)

	if err != nil {
		return fmt.Errorf("failed to  update password state: %v", err)
	}
	return nil
}

//DeleteActivePolicy deletes an active policy
func (s userStore) DeleteActivePolicy(userID, orgID, enforceType string) error {
	_, err := s.DB.Exec(`DELETE FROM policy_enforcer WHERE user_id = $1 AND org_id=$2 AND type=$3 RETURNING *`,
		userID, orgID, enforceType)

	return err
}

// DeleteAppUserbasedOnUserID will remove a user from every apps
func (s userStore) DeleteAllUserAccessMaps(userID, orgID string) error {
	_, err := s.DB.Exec(`DELETE FROM user_accessmaps WHERE user_id = $1 AND org_id=$2`, userID, orgID)
	return err
}

//Create inserts  a new user into user table
func (s userStore) Create(user *models.UserWithPass) error {
	_, err := s.DB.Exec(`INSERT into users (org_id, id,  username, first_name, middle_name, last_name, email, password, user_role,status, created_at, updated_at, idp_name, external_id)
						 values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);`,
		user.OrgID, user.ID, user.UserName, user.FirstName,
		user.MiddleName, user.LastName, user.Email, user.Password, user.UserRole, user.Status, user.CreatedAt,
		user.UpdatedAt, user.IdpName, user.ExternalID)
	return err
}

//GetPasswordState returns password state of a user
func (s userStore) GetPasswordState(userID, orgID string) (models.PasswordState, error) {
	var passState models.PasswordState
	err := s.DB.QueryRow(`SELECT user_id, org_id,  last_passwords, last_updated FROM password_state WHERE user_id = $1 and org_id = $2`,
		userID, orgID).
		Scan(&passState.UserID, &passState.OrgID, pq.Array(&passState.LastPasswords), &passState.LastUpdated)
	return passState, err
}

//GetDevicesByType returns array of devices of given type
func (s userStore) GetDevicesByType(userID, deviceType, orgID string) ([]models.UserDevice, error) {
	var device models.UserDevice
	var userDevices = make([]models.UserDevice, 0)
	rows, err := s.DB.Query("SELECT id, type, fcm_token, public_key, device_hygiene, added_at  FROM devices WHERE deleted != true AND user_id = $1 AND org_id=$2 AND type = $3",
		userID, orgID, deviceType)
	if err != nil {
		return userDevices, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&device.DeviceID, &device.DeviceType, &device.FcmToken, &device.PublicKey, &device.DeviceFinger, &device.DeviceHygiene, &device.AddedAt)
		if err != nil {
			return userDevices, err
		}
		userDevices = append(userDevices, device)
	}

	return userDevices, nil
}

//GetTOTPDevices returns all totp devices(mobile or htoken) of a user. Onlu totpsec and fcm_token fields are retrieved
func (s userStore) GetTOTPDevices(userID, orgID string) ([]models.UserDevice, error) {
	devices := []models.UserDevice{}

	//var userHash *string
	rows, err := s.DB.Query("SELECT totpsec,id,fcm_token FROM devices WHERE (type='mobile' OR type='htoken') AND deleted != true AND user_id = $1 AND org_id=$2",
		userID, orgID)
	if err != nil {
		return devices, err
	}
	for rows.Next() {
		var dev models.UserDevice
		err = rows.Scan(&dev.TotpSec, &dev.DeviceID, &dev.FcmToken)
		if err != nil {
			continue
		}
		devices = append(devices, dev)
	}

	return devices, nil
}

func (s userStore) GetEnforcedPolicy(userID, orgID, enforceType string) (policy models.PolicyEnforcer, err error) {
	err = s.DB.QueryRow("SELECT id, user_id, org_id,pending, type FROM policy_enforcer WHERE user_id = $1 AND org_id = $2 AND type=$3;",
		userID, orgID, enforceType).
		Scan(&policy.EnforceID, &policy.UserID, &policy.OrgID, &policy.Pending, &policy.EnforceType)
	return
}

func (s userStore) GetGroups(userID, orgID string) ([]models.Group, error) {
	groups := make([]models.Group, 0)
	rows, err := s.DB.Query(`SELECT g.name,g.id,ug.created_at from groups g
							JOIN user_group_maps ug on g.id = ug.group_id
							WHERE ug.user_id=$1 AND ug.org_id=$2`, userID, orgID)
	if err != nil {
		return groups, err
	}

	for rows.Next() {
		var temp models.Group
		err = rows.Scan(&temp.GroupName, &temp.GroupID, &temp.CreatedAt)
		if err != nil {
			return groups, err
		}
		groups = append(groups, temp)
	}
	return groups, nil
}

func (s userStore) UpdatePublicKey(userID string, publicKey string) error {
	_, err := s.DB.Exec("Update users SET public_key=$1  WHERE id = $2 ;", publicKey, userID)

	return err
}

func (s userStore) EnforcePolicy(policy models.PolicyEnforcer) error {
	_, err := s.DB.Exec(`INSERT into policy_enforcer (id, user_id, org_id, type, pending, assigned_by, assigned_on, resolved_on) values($1, $2, $3, $4, $5, $6, $7, $8);`,
		policy.EnforceID, policy.UserID, policy.OrgID, policy.EnforceType, policy.Pending, policy.AssignedBy, policy.AssignedOn, policy.ResolvedOn)
	return err
}

func (s userStore) GetAllByIdp(orgID, idpName string) ([]models.User, error) {
	var usr []models.User = make([]models.User, 0)
	var user models.User

	rows, err := s.DB.Query(`SELECT id, first_name,
		last_name, email, idp_name FROM users WHERE org_id = $1 AND idp_name=$2`, orgID, idpName)
	if err != nil {
		return usr, errors.Errorf("query idp users: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.ID,
			&user.FirstName, &user.LastName, &user.Email, &user.IdpName)
		if err != nil {
			logrus.Errorf("scan idp users: %v", err)
		}
		usr = append(usr, user)
	}

	return usr, nil

}

// TransferUser transfers user(based on email) to another provided idp
func (s userStore) TransferUser(orgID, email, idpName string) error {
	_, err := s.DB.Exec(`UPDATE users SET idp_name = $1 WHERE org_id=$2 AND email = $3;`,
		idpName, orgID, email)

	if err != nil {
		return err
	}
	return nil
}

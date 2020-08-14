package auth

import (
	"context"
	"strings"

	"github.com/seknox/trasa/server/models"
)

func (s authStore) GetLoginDetails(trasaID, orgDomain string) (*models.UserWithPass, error) {
	//var users []models.UserWithPass = make([]models.UserWithPass, 0)
	isTrasaIDEmail := strings.Contains(trasaID, "@")

	//TODO use domain

	sqlStr := ``

	if isTrasaIDEmail {
		sqlStr = `SELECT users.org_id, users.id, first_name, email, password, idp_name, user_role, status ,org.org_name
				FROM users
				JOIN org ON users.org_id=org.id
				WHERE users.email=$1
`
	} else {
		sqlStr = `SELECT users.org_id, users.id, first_name, email, password, idp_name, user_role, status ,org.org_name
				FROM users
				JOIN org ON users.org_id=org.id
				WHERE users.username=$1`
	}

	var user models.UserWithPass
	err := s.DB.QueryRow(sqlStr, trasaID).Scan(&user.OrgID, &user.ID, &user.FirstName, &user.Email, &user.Password, &user.IdpName, &user.UserRole, &user.Status, &user.OrgName)

	return &user, err
}

func (s authStore) Logout(sessionID string) error {
	client := s.RedisClient
	ctx := context.Background()
	return client.Del(ctx, sessionID).Err()
}

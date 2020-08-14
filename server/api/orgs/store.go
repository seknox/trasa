package orgs

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/models"
	"github.com/sirupsen/logrus"
)

func (s orgStore) RemoveAllManagedAccounts(orgID string) error {

	_, err := s.DB.Exec(`UPDATE services SET managed_accounts = $1  WHERE org_id=$2`, "", orgID)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (s orgStore) CheckOrgExists() (orgID string, err error) {
	err = s.DB.QueryRow(`select id from org`).Scan(&orgID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return orgID, err

}

func (s orgStore) CreateOrg(org *models.Org) error {

	_, err := s.DB.Exec(`INSERT into org (id, org_name, domain, primary_contact,timezone, created_at,phone_number,license)
						 values($1, $2, $3, $4, $5,$6,$7,$8);`, org.ID, org.OrgName, org.Domain, org.PrimaryContact, org.Timezone, org.CreatedAt, org.PhoneNumber, org.License)

	return err
}

func (s orgStore) Get(orgID string) (models.Org, error) {
	var org models.Org
	err := s.DB.QueryRow("SELECT id, org_name, domain, primary_contact, timezone, phone_number,created_at FROM org WHERE id = $1", orgID).Scan(&org.ID, &org.OrgName, &org.Domain,
		&org.PrimaryContact, &org.Timezone, &org.PhoneNumber, &org.CreatedAt)

	return org, err
}

func (s orgStore) GetIDP(orgID, idpName string) (models.IdentityProvider, error) {
	var idp models.IdentityProvider
	err := s.DB.QueryRow("SELECT id, org_id, name,type, meta, is_enabled, redirect_url, audience_uri, endpoint, created_by , last_updated FROM idp WHERE org_id = $1 AND name=$2",
		orgID, idpName).
		Scan(&idp.IdpID, &idp.OrgID, &idp.IdpName, &idp.IdpType, &idp.IdpMeta, &idp.IsEnabled, &idp.RedirectURL, &idp.AudienceURI, &idp.Endpoint, &idp.CreatedBy, &idp.LastUpdated)

	return idp, err
}

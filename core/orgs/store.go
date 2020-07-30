package orgs

import (
	"github.com/seknox/trasa/models"
	"github.com/sirupsen/logrus"
)

func (s OrgStore) RemoveAllManagedAccounts(orgID string) error {

	_, err := s.DB.Exec(`UPDATE services SET managed_accounts = $1  WHERE org_id=$2`, "", orgID)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (s OrgStore) Get(orgID string) (models.Org, error) {
	var org models.Org
	err := s.DB.QueryRow("SELECT id, org_name, domain, primary_contact, timezone, phone_number,created_at FROM org WHERE id = $1", orgID).Scan(&org.ID, &org.OrgName, &org.Domain,
		&org.PrimaryContact, &org.Timezone, &org.PhoneNumber, &org.CreatedAt)

	return org, err
}

func (s OrgStore) GetIDP(orgID, idpName string) (models.IdentityProvider, error) {
	var idp models.IdentityProvider
	err := s.DB.QueryRow("SELECT id, org_id, name,type, meta, is_enabled, redirect_url, audience_uri, endpoint, created_by , last_updated FROM idp WHERE org_id = $1 AND name=$2",
		orgID, idpName).
		Scan(&idp.IdpID, &idp.OrgID, &idp.IdpName, &idp.IdpType, &idp.IdpMeta, &idp.IsEnabled, &idp.RedirectURL, &idp.AudienceURI, &idp.Endpoint, &idp.CreatedBy, &idp.LastUpdated)

	return idp, err
}

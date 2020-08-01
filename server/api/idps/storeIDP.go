package idps

import (
	"github.com/seknox/trasa/server/models"
	logger "github.com/sirupsen/logrus"
)

// GetAllIdps retreives all idps configured for organization
func (s IDPStore) GetAllIdps(orgID string) ([]models.IdentityProvider, error) {
	var idps []models.IdentityProvider = make([]models.IdentityProvider, 0)
	var idp models.IdentityProvider
	rows, err := s.DB.Query("SELECT id, org_id, name,type, meta, is_enabled, redirect_url, audience_uri, client_id, endpoint, created_by , last_updated, scim_endpoint FROM idp WHERE org_id = $1", orgID)

	if err != nil {
		return idps, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&idp.IdpID, &idp.OrgID, &idp.IdpName, &idp.IdpType, &idp.IdpMeta, &idp.IsEnabled, &idp.RedirectURL, &idp.AudienceURI, &idp.ClientID, &idp.Endpoint, &idp.CreatedBy, &idp.LastUpdated, &idp.SCIMEndpoint)
		if err != nil {
			logger.Errorf("scan error in IDPStore.GetAllIdps: %v", err)
		}
		idps = append(idps, idp)
	}

	return idps, err
}

// GetByID retreives IDP detail based on ID
func (s IDPStore) GetByID(orgID, idpID string) (models.IdentityProvider, error) {
	var idp models.IdentityProvider
	err := s.DB.QueryRow("SELECT id, org_id, name,type, meta, is_enabled, redirect_url, audience_uri, client_id, endpoint, created_by , last_updated FROM idp WHERE org_id = $1 AND id=$2",
		orgID, idpID).
		Scan(&idp.IdpID, &idp.OrgID, &idp.IdpName, &idp.IdpType, &idp.IdpMeta, &idp.IsEnabled, &idp.RedirectURL, &idp.AudienceURI, &idp.ClientID, &idp.Endpoint, &idp.CreatedBy, &idp.LastUpdated)
	return idp, err
}

// CreateIDP created new Identity Provider
func (s IDPStore) CreateIDP(idp *models.IdentityProvider) error {

	_, err := s.DB.Exec(`INSERT into idp (id, org_id, name,type, meta, is_enabled, redirect_url, audience_uri,client_id, endpoint, created_by , integration_type,scim_endpoint, last_updated )
						 values($1, $2, $3, $4, $5,$6,$7,$8, $9, $10, $11, $12, $13, $14);`, idp.IdpID, idp.OrgID, idp.IdpName, idp.IdpType, idp.IdpMeta, idp.IsEnabled, idp.RedirectURL, idp.AudienceURI, idp.ClientID, idp.Endpoint, idp.CreatedBy, idp.IntegrationType, idp.SCIMEndpoint, idp.LastUpdated)

	return err
}

func (s IDPStore) UpdateIDP(idp *models.IdentityProvider) error {

	_, err := s.DB.Exec(`UPDATE idp SET meta = $1, is_enabled = $2, endpoint = $3, created_by = $4 , last_updated = $5  WHERE org_id=$6 AND id=$7`, idp.IdpMeta, idp.IsEnabled, idp.Endpoint, idp.CreatedBy, idp.LastUpdated, idp.OrgID, idp.IdpID)

	return err
}

func (s IDPStore) UpdateLDAPIDP(idp *models.IdentityProvider) error {

	_, err := s.DB.Exec(`UPDATE idp SET meta = $1, is_enabled = $2, endpoint = $3, created_by = $4 , last_updated = $5, audience_uri=$6, client_id=$7  WHERE org_id=$8 AND id=$9`, idp.IdpMeta, idp.IsEnabled, idp.Endpoint, idp.CreatedBy, idp.LastUpdated, idp.AudienceURI, idp.ClientID, idp.OrgID, idp.IdpID)

	return err
}

func (s IDPStore) activateOrDisableIdp(orgID, idpID string, updateTime int64, updateVal bool) error {

	_, err := s.DB.Exec(`UPDATE idp SET is_enabled = $1, last_updated = $2  WHERE org_id=$3 AND id=$4`, updateVal, updateTime, orgID, idpID)

	return err
}

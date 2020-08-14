package idps

import "github.com/seknox/trasa/server/models"

// GetCloudSyncState fetches time on when was trasa last synced with cloudIaas provider.
func (s idpStore) GetCloudSyncState(orgID, cName string) (*models.CloudIaaSSync, error) {
	var sync models.CloudIaaSSync
	err := s.DB.QueryRow(`
		SELECT id, org_id, name, last_synced_by, last_synced_on FROM cloudiaas_sync WHERE org_id = $1 AND name = $2 ;`, orgID, cName).
		Scan(&sync.CloudIaasID, &sync.OrgID, &sync.CloudIaasName, &sync.LasgtSyncedBy, &sync.LastSyncedOn)

	return &sync, err
}

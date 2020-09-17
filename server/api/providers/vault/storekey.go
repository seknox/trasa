package vault

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/models"
)

func (s cryptStore) StoreKey(k models.KeysHolder) error {

	storedKey, err := s.GetKeyOrTokenWithTag(k.OrgID, k.KeyName)
	if errors.Is(err, sql.ErrNoRows) {
		_, err := s.DB.Exec(`INSERT into key_holder (id, org_id, name, value, tag, added_by, added_at)
		values($1, $2, $3, $4, $5, $6, $7);`, k.KeyID, k.OrgID, k.KeyName, k.KeyVal, k.KeyTag, k.AddedBy, k.AddedAt)
		if err != nil {
			return err
		}
	} else {
		_, err := s.DB.Exec(`UPDATE key_holder SET value = $1, tag = $2, added_by =$3, added_at = $4 WHERE id = $5 AND org_id =$6;`, k.KeyVal, k.KeyTag, k.AddedBy, k.AddedAt, storedKey.KeyID, storedKey.OrgID)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetKeyOrTokenWithTag returns key or token detail but without actual key value but rather tagged value.
func (s cryptStore) GetKeyOrTokenWithTag(orgID string, keyName string) (*models.KeysHolder, error) {
	var k models.KeysHolder
	err := s.DB.QueryRow(`
		SELECT id, org_id, name, tag, added_by, added_at FROM key_holder WHERE org_id = $1 AND name = $2 ;`, orgID, keyName).Scan(&k.KeyID, &k.OrgID, &k.KeyName, &k.KeyTag, &k.AddedBy, &k.AddedAt)

	return &k, err
}

// GetKeyOrTokenWithTag returns key or token detail but without actual key value but rather tagged value.
func (s cryptStore) GetKeyOrTokenWithKeyval(orgID, keyName string) (*models.KeysHolder, error) {
	var k models.KeysHolder
	err := s.DB.QueryRow(`
		SELECT id, org_id, name, value, added_by, added_at FROM key_holder WHERE org_id = $1 AND name = $2 ;`,
		orgID, keyName).
		Scan(&k.KeyID, &k.OrgID, &k.KeyName, &k.KeyVal, &k.AddedBy, &k.AddedAt)
	if err != nil {
		return &k, err
	}

	return &k, nil
}

// GetKeyOrTokenWithTag returns key or token detail but without actual key value but rather tagged value.
func (s cryptStore) GetKeyOrTokenWithKeyvalAndID(orgID, keyName, keyID string) (*models.KeysHolder, error) {
	var k models.KeysHolder
	err := s.DB.QueryRow(`
		SELECT id, org_id, name, value, added_by, added_at FROM key_holder WHERE org_id = $1 AND name = $2 AND id=$3 ;`,
		orgID, keyName, keyID).
		Scan(&k.KeyID, &k.OrgID, &k.KeyName, &k.KeyVal, &k.AddedBy, &k.AddedAt)

	return &k, err
}

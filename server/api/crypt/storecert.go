package crypt

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////				Certificate Authority						/////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// StoreCert inserts certificate detail in cert_holderv1. If cert for service_id or type already exists,
// StoreCert should update the value.
func (s CryptStore) StoreCert(ch models.CertHolder) error {

	storedCert, err := s.GetCertDetail(ch.OrgID, ch.EntityID, ch.CertType)
	if errors.Is(err, sql.ErrNoRows) {
		_, err := s.DB.Exec(`INSERT into cert_holder (id, org_id, entity_id, cert, key, csr, cert_type, created_at, last_updated)
		values($1, $2, $3, $4, $5, $6, $7, $8, $9);`, ch.CertID, ch.OrgID, ch.EntityID, ch.Cert, ch.Key, ch.Csr, ch.CertType, ch.CreatedAt, ch.LastUpdated)
		if err != nil {
			return err
		}
	} else {
		_, err := s.DB.Exec(`UPDATE cert_holder SET cert = $1, key = $2, csr =$3, cert_type = $4, last_updated = $5 WHERE id = $6 AND org_id =$7;`,
			ch.Cert, ch.Key, ch.Csr, storedCert.CertType, ch.LastUpdated, storedCert.CertID, storedCert.OrgID)
		if err != nil {
			return err
		}
	}
	return nil
}

// delete CA
func (s CryptStore) DelCA(userID, orgID string) error {
	_, err := s.DB.Exec(`DELETE from devices where user_id = $1 AND org_id=$2`, userID, orgID)

	return err
}

// get cert
func (s CryptStore) GetCertDetail(orgID, entityID, certType string) (*models.CertHolder, error) {
	var ch models.CertHolder
	err := s.DB.QueryRow(`
		SELECT id, org_id, entity_id, cert, cert_type, created_at, last_updated
			FROM cert_holder WHERE org_id = $1 AND entity_id = $2 AND cert_type = $3
	`, orgID, entityID, certType).Scan(&ch.CertID, &ch.OrgID, &ch.EntityID, &ch.Cert, &ch.CertType, &ch.CreatedAt, &ch.LastUpdated)

	return &ch, err
}

// get cert
func (s CryptStore) GetAllCAs(orgID string) ([]models.CertHolder, error) {
	var cas []models.CertHolder
	rows, err := s.DB.Query(`
		SELECT id, org_id, entity_id, cert, cert_type, created_at, last_updated
			FROM cert_holder WHERE org_id = $1 AND (cert_type=$2 OR cert_type=$3 )
	`, orgID, consts.CERT_TYPE_SSH_CA, consts.CERT_TYPE_HTTP_CA)
	if err != nil {
		return cas, err
	}

	for rows.Next() {
		var ch models.CertHolder
		err := rows.Scan(&ch.CertID, &ch.OrgID, &ch.EntityID, &ch.Cert, &ch.CertType, &ch.CreatedAt, &ch.LastUpdated)
		if err != nil {
			return cas, err
		}
		cas = append(cas, ch)
	}

	return cas, nil
}

// get user devices
//Deprecated
func (s CryptStore) GetCAkey(orgID string) ([]byte, error) {
	var key []byte
	err := s.DB.QueryRow(`
		SELECT key FROM cert_holder WHERE org_id = $1 AND cert_type = 'ca'
	`, orgID).Scan(&key)

	return key, err
}

func (s CryptStore) GetCertHolder(certType, entityID, orgID string) (models.CertHolder, error) {
	var certHolder models.CertHolder
	err := s.DB.QueryRow(`SELECT id, org_id, entity_id, cert, "key", csr, cert_type, created_at, last_updated FROM cert_holder WHERE entity_id=$1 AND org_id=$2 AND cert_type=$3`,
		entityID, orgID, certType).
		Scan(&certHolder.CertID, &certHolder.OrgID, &certHolder.EntityID, &certHolder.Cert, &certHolder.Key, &certHolder.Csr, &certHolder.CertType, &certHolder.CreatedAt, &certHolder.LastUpdated)
	if err != nil {
		return certHolder, err
	}
	return certHolder, nil
}

package services

import (
	"context"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/seknox/trasa/server/consts"
)

//TODO how to store tls/ssl certs
func (s ServiceStore) UpdateSSLCerts(caCert, caKey, clientCert, clientKey, serviceID, orgID string) error {

	certID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	tx, err := s.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM cert_holder WHERE entity_id=$1 AND org_id=$2 `, serviceID, orgID)
	if err != nil {
		tx.Rollback()
		return err
	}

	//CA certs
	_, err = tx.Exec(`
INSERT INTO cert_holder (id, org_id, entity_id, cert, key, csr, cert_type, created_at, last_updated) 
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, certID.String(), orgID, serviceID, caCert, caKey, "", "ca", time.Now().Unix(), time.Now().Unix())
	if err != nil {
		tx.Rollback()
		return err
	}

	certID, err = uuid.NewV4()
	if err != nil {
		tx.Rollback()
		return err
	}

	//client certs
	_, err = tx.Exec(`
INSERT INTO cert_holder (id, org_id, entity_id, cert, key, csr, cert_type, created_at, last_updated) 
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, certID.String(), orgID, serviceID, clientCert, clientKey, "", "client certs", time.Now().Unix(), time.Now().Unix())
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()

}

func (s ServiceStore) UpdateHostCert(hostCert, serviceID, orgID string) error {
	certID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	tx, err := s.DB.BeginTx(context.Background(), nil)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM cert_holder where cert_type=$3 AND entity_id=$1 AND org_id=$2`,
		serviceID, orgID, consts.CERT_TYPE_SSH_HOST_KEY)
	if err != nil {
		return errors.Errorf("could not delete old cert: %v", err)
	}
	_, err = tx.Exec(`
			INSERT INTO cert_holder (id, org_id, entity_id, cert, key, csr, cert_type, created_at, last_updated)
 			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		certID.String(), orgID, serviceID, hostCert, "", "", consts.CERT_TYPE_SSH_HOST_KEY, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		tx.Rollback()
		return errors.Errorf("could not insert cert: %v", err)
	}

	return tx.Commit()

}

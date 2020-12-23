package tsxvault

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/models"

	"github.com/seknox/trasa/server/utils"
	logger "github.com/sirupsen/logrus"
)

// First we check secret storage setting. If credStorage specefies tsxvalut, we store it in our database
// else we use api token to store it in secret storage provider.
// Configuration for vault should be stored in trasa_featuresv1 table.
func (s vaultStore) StoreSecret(key models.ServiceSecretVault) error {
	// feature, err := Connect.CRDBGetOrgFeatureStatus(s.OrgID, "vault")
	// if err != nil {
	// 	logger.Error(err)
	// 	return err
	// }
	// var vaultConfig utils.CredProvProps
	// json.Unmarshal([]byte(feature.Config), &vaultConfig)
	// if vaultConfig.CredStorage == "tsxvault" {
	// 	// store it in tsxvault
	// 	err := Connect.TsxvStoreSecret(s)
	// 	if err != nil {
	// 		logger.Error(err)
	// 		return err
	// 	}

	// } else {
	// 	//store is in specified provider.
	// }

	// return fmt.Errorf("%v", 1)

	if s.TsxvKey.State == false {
		return fmt.Errorf("encryption key is not retrieved yet.")
	}
	//	fmt.Println("ENCRYPTION KEY: ", hex.EncodeToString(TsxvKey.Key[:]))

	ct, err := utils.AESEncrypt(s.TsxvKey.Key[:], key.Secret)
	if err != nil {
		return err
	}
	key.Secret = ct

	err = s.TsxvStoreSecret(key)
	if err != nil {
		return err
	}

	return nil

}

// GetSecret is single method for retreiving secret either from tsxvault or 3rd party storage provider.
func (s vaultStore) GetSecret(orgID, serviceID, secretType, secretid string) (string, error) {

	if s.TsxvKey.State == false {
		return "", fmt.Errorf("encryption key is not retrieved yet")
	}

	secret, err := s.TsxvGetSecret(orgID, serviceID, secretType, secretid)
	if err != nil {
		return "", errors.Errorf("tsxv get secret: %v", err)
	}

	pt, err := utils.AESDecrypt(s.TsxvKey.Key[:], secret)
	if err != nil {
		return "", errors.Errorf("decrypt aes: %v", err)
	}
	return string(pt), nil
}

// TsxvStoreSecret stores secret in Service_keyvaultv1 table.
// It should be receiving already encrypted secret from the caller.
func (s vaultStore) TsxvStoreSecret(secret models.ServiceSecretVault) error {

	_, err := s.TsxvGetSecret(secret.OrgID, secret.ServiceID, secret.SecretType, secret.SecretID)
	if err != nil {
		_, err := s.DB.Exec(`INSERT into service_keyvault (id, org_id, service_id, secret_type, secret_id, secret, added_at,last_updated, updated_by)
		values($1, $2, $3, $4, $5, $6, $7, $8, $9);`, secret.KeyID, secret.OrgID, secret.ServiceID, secret.SecretType, secret.SecretID, secret.Secret, secret.AddedAt, secret.LastUpdated, secret.UpdatedBy)
		if err != nil {
			return errors.Errorf("could not insert: %v", err)
		}
	} else {
		_, err := s.DB.Exec(`UPDATE service_keyvault SET secret = $1, last_updated = $2, updated_by =$3 WHERE org_id = $4 AND service_id =$5 AND secret_type=$6;`,
			secret.Secret, secret.LastUpdated, secret.UpdatedBy, secret.OrgID, secret.ServiceID, secret.SecretType)
		if err != nil {
			return errors.Errorf("could not update: %v", err)
		}
	}
	return nil
}

// GetSecret detail retrieves secretdetails  from tsxVault
// "ServiceSecretVault": `create table IF NOT EXISTS service_keyvaultv1(
// 	key_id VARCHAR PRIMARY KEY NOT NULL,
// 	org_id VARCHAR REFERENCES orgv1(org_id) ON DELETE CASCADE ,
// 	service_id VARCHAR REFERENCES servicesv1(service_id) ON DELETE CASCADE ,
// 	secretid VARCHAR,
// 	secret BYTEA,
// 	added_at INT,
// 	last_updated INT
// )`,
func (s vaultStore) TsxvGetSecretDetail(orgID, serviceID, appType, secretID string) (*models.ServiceSecretVault, error) {
	var secret models.ServiceSecretVault
	err := s.DB.QueryRow(`
		SELECT id, org_id, service_id, secret_type, secret_id, secret, added_at,last_updated, updated_by FROM service_keyvault WHERE org_id = $1 AND service_id = $2 AND secret_id=$3 AND secret_type=$4`,
		orgID, serviceID, secretID, appType).
		Scan(&secret.KeyID, &secret.OrgID, &secret.ServiceID, &secret.SecretType, &secret.SecretID, &secret.Secret, &secret.AddedAt, &secret.LastUpdated, secret.UpdatedBy)

	return &secret, err
}

// GetSecret is only returns secret value from tsxVault
func (s vaultStore) TsxvGetSecret(orgID, serviceID, appType, secretid string) ([]byte, error) {
	var secret []byte
	err := s.DB.QueryRow(`SELECT secret FROM service_keyvault WHERE org_id = $1 AND service_id = $2 AND secret_id=$3 AND secret_type=$4`,
		orgID, serviceID, secretid, appType).
		Scan(&secret)

	return secret, err
}

// GetSecret is only returns secret value from tsxVault
func (s vaultStore) TsxvDeleteSecret(orgID, serviceID, secretType, secretid string) error {
	_, err := s.DB.Exec(`DELETE FROM service_keyvault WHERE org_id = $1 AND service_id = $2 AND secret_id=$3 AND secret_type=$4`,
		orgID, serviceID, secretid, secretType)

	return err
}

// GetSecret is only returns secret value from tsxVault
func (s vaultStore) TsxvDeleteAllSecret(orgID string) error {
	_, err := s.DB.Exec(`DELETE FROM service_keyvault WHERE org_id = $1`, orgID)

	return err
}

// StoreEncKeyHash stores sha512 hash of encryption key in database.
// This helps to verify whether user supplied encryption key is correct.
// Every retreival of key shold be appended in audit log.
// "KeyLog": `create table IF NOT EXISTS keylogv1(
// 	key_id VARCHAR PRIMARY KEY NOT NULL,
// 	org_id VARCHAR REFERENCES orgv1(org_id) ON DELETE CASCADE ,
// 	key_hash VARCHAR NOT NULL,
// 	generated_at VARCHAR,
//  status BOOL,
// 	log JSONB
// )`,
func (s vaultStore) TsxvStoreEncKeyHash(secret models.EncryptionKeyLog) error {

	_, err := s.DB.Exec(`INSERT into keylog (id, org_id, hash, generated_at, status, last_updated)
		values($1, $2, $3, $4, $5, $6);`,
		secret.KeyID, secret.OrgID, secret.KeyHash, secret.GeneratedAt, secret.Status, secret.LastUpdated)

	return err
}

func (s vaultStore) TsxvGetEncKeyHash(orgID, keyHash string) (models.EncryptionKeyLog, error) {

	var kl models.EncryptionKeyLog
	err := s.DB.QueryRow(`SELECT id, org_id, hash, generated_at, status, last_updated FROM keylog WHERE org_id = $1 AND hash =$2`,
		orgID, keyHash).
		Scan(&kl.KeyID, &kl.OrgID, &kl.KeyHash, &kl.GeneratedAt, &kl.Status, &kl.LastUpdated)

	return kl, err

}

func (s vaultStore) TsxvUpdateEncKeyHashLog(orgID, keyHash string, time int64, status bool) error {

	_, err := s.DB.Exec(`UPDATE keylog SET status = $1, last_updated = $2 WHERE org_id = $3 AND hash =$4;`,
		status, time, orgID, keyHash)

	return err
}

func (s vaultStore) TsxvdeactivateAllKeys(orgID string, time int64) error {
	_, err := s.DB.Exec(`UPDATE keylog SET status = $1, last_updated = $2 WHERE org_id = $3`,
		false, time, orgID)

	return err
}

// TsxVaultTestEncrypter stores test data in database. this should be encrypted by
// newly created encryption token and should be decrypted in compared when
// token key retreival process is finished to verify if key works.
func (s vaultStore) TsxvTestEncrypter(key []byte) error {

	secretValue := "testsecret"
	var secret models.ServiceSecretVault
	secret.SecretID = "testid"
	secret.Secret = []byte(secretValue)
	secret.ServiceID = "testapp"
	secret.OrgID = "testorg"
	secret.UpdatedBy = "trasatest"
	secret.AddedAt = time.Now().Unix()
	secret.LastUpdated = time.Now().Unix()

	//logger.Debug(fmt.Sprintf("the plain text is: %s", string(secret.Secret)))
	// encrypt the secret.
	ct, err := utils.AESEncrypt(key, secret.Secret)
	if err != nil {
		return fmt.Errorf("failed to pass encryption test")
	}

	//logger.Debug(fmt.Sprintf("the cipher text is: %s", string(ct)))
	secret.Secret = ct

	// store it in database
	err = s.TsxvStoreSecret(secret)
	if err != nil {
		return fmt.Errorf("failed to store encrypted value in vault: %v", err)
	}

	return nil
}

func (s vaultStore) TsxvTestDecrypter(key []byte) error {
	// retrieve from database
	ct, err := s.TsxvGetSecret("testorg", "testapp", "testssh", "testid")
	if err != nil {
		return fmt.Errorf("failed to retrieve sipher text: %v", err)
	}

	// retrieve plaintext.
	pt, err := utils.AESDecrypt(key, ct)
	if err != nil {
		return fmt.Errorf("decryption failed failed: %v", err)
	}

	// compare pt and ct
	if string(pt) == "testsecret" {
		return fmt.Errorf("ct and pt does not match. ")
	}

	return nil
}

// TsxVaultTester performd encryption, secret storage in db, retrieving secret and perform decryption
// to test every thing is working in tsxVault.
// This test should be performed after shamir key reducer returns successful key retreival.
func (s vaultStore) TsxVaultTester() error {

	if s.TsxvKey.State == false {
		return fmt.Errorf("Encryption key is not retrieved yet.")
	}

	secretValue := "testsecret"
	var secret models.ServiceSecretVault
	secret.SecretID = "testid"
	secret.Secret = []byte(secretValue)
	secret.ServiceID = "testapp"
	secret.OrgID = "testorg"
	secret.UpdatedBy = "trasatest"
	secret.AddedAt = time.Now().Unix()
	secret.LastUpdated = time.Now().Unix()

	// encrypt the secret.
	ct, err := utils.AESEncrypt(s.TsxvKey.Key[:], secret.Secret)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("failed to pass encryption test")
	}

	secret.Secret = ct

	// store it in database
	err = s.TsxvStoreSecret(secret)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("failed to store encrypted value in vault")
	}

	// retrieve from database
	ct, err = s.TsxvGetSecret(secret.OrgID, secret.ServiceID, secret.SecretType, secret.SecretID)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("failed to retrieve sipher text")
	}

	// retrieve plaintext.
	pt, err := utils.AESDecrypt(s.TsxvKey.Key[:], ct)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("decryption failed failed. ")
	}

	// compare pt and ct
	if string(pt) == secretValue {
		return fmt.Errorf("ct and pt does not match. ")
	}

	return nil
}

// AesEncrypt receives message to be encrypted from caller, retrieves key from global key store and performs encryption.
// This function should be used instead of directly using utils.AESEncrypt
func (s vaultStore) AesEncrypt(message []byte) ([]byte, error) {
	// get decrypted email password or key.
	if s.TsxvKey.State == false {
		return nil, fmt.Errorf("encryption key is not retrieved yet")
	}

	ct, err := utils.AESEncrypt(s.TsxvKey.Key[:], message)
	if err != nil {
		return nil, err
	}

	return ct, nil
}

// AesDecrypt receives message to be decrypted from caller, retrieves key from global key store and performs encryption.
// This function should be used instead of directly using utils.AESDecrypt
func (s vaultStore) AesDecrypt(message []byte) ([]byte, error) {
	// get decrypted email password or key.
	if s.TsxvKey.State == false {
		return nil, fmt.Errorf("encryption key is not retrieved yet")
	}

	pt, err := utils.AESDecrypt(s.TsxvKey.Key[:], message)
	if err != nil {
		return nil, err
	}

	return pt, nil
}

// GetKeyOrTokenWithTag returns key or token detail but without actual key value but rather tagged value.
func (s vaultStore) GetKeyOrTokenWithTag(orgID string, keyName string) (*models.KeysHolder, error) {
	var k models.KeysHolder
	err := s.DB.QueryRow(`
		SELECT id, org_id, name, tag, added_by, added_at FROM key_holder WHERE org_id = $1 AND name = $2 ;`, orgID, keyName).Scan(&k.KeyID, &k.OrgID, &k.KeyName, &k.KeyTag, &k.AddedBy, &k.AddedAt)

	return &k, err
}

// StoreKeyOrTokens inserts certificate detail in cert_holderv1. If cert for service_id or type already exists,
// StoreKeyOrTokens should update the value.
// In case of Idp SCIM key, keyID is idpID from idpv1
func (s vaultStore) StoreKeyOrTokens(k models.KeysHolder) error {

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

// GenAndStoreKey generates encryption keys, store it in database.
func (s vaultStore) GenAndStoreKey(orgID string) (*[32]byte, error) {
	encryptionKey, err := utils.AESGenKey()
	if err != nil {
		return encryptionKey, err
	}

	s.TsxvKey.Key = encryptionKey
	s.TsxvKey.State = true
	hash := sha512.New()

	//fmt.Println("Original Key: ", hex.EncodeToString(encryptionKey[:]))
	hashed := hash.Sum(encryptionKey[:])

	// store key has in database
	var kl models.EncryptionKeyLog
	kl.OrgID = orgID
	kl.KeyID = utils.GetUUID()
	kl.KeyHash = hex.EncodeToString(hashed)
	kl.GeneratedAt = time.Now().Unix()
	kl.LastUpdated = time.Now().Unix()
	kl.Status = true
	err = s.TsxvStoreEncKeyHash(kl)
	if err != nil {
		return nil, err
	}

	return encryptionKey, nil
}

// GetTsxVaultKey returns retreival status of encryption key
func (s vaultStore) GetTsxVaultKey() (*[32]byte, bool) {
	return s.TsxvKey.Key, s.TsxvKey.State
}

// SetTsxVaultKey returns retreival status of encryption key
func (s vaultStore) SetTsxVaultKey(key *[32]byte, status bool, credprov models.CredProvProps) {

	s.TsxvKey.Key = key
	s.TsxvKey.State = status
	s.TsxvKey.CredProv = credprov
}

// SetTsxCPxyKey assigns retreived cloud prxy api key in global state
func (s vaultStore) SetTsxCPxyKey(key string) {

	s.TsxCPxyKey = key

}

// GetTsxCPxyKey retreives cloud prxy api key from global state
func (s vaultStore) GetTsxCPxyKey() string {

	return s.TsxCPxyKey

}

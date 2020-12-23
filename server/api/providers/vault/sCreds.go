package vault

import (
	"github.com/seknox/trasa/server/api/providers/vault/tsxvault"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/sirupsen/logrus"
)

// StoreCred saves user credential on specified secret provider.
func (s cryptStore) StoreCred(key models.ServiceSecretVault) error {

	logrus.Trace("Writing to CredProv: ", s.TsxvKey.CredProv.ProviderName)
	if s.TsxvKey.CredProv.ProviderName == consts.CREDPROV_HCVAULT {
		err := Store.HCVStoreCred(key)
		if err != nil {
		 	return err
		}
		return nil
	}

	err := tsxvault.Store.StoreSecret(key)
	if err != nil {
		return err
	}

	return nil

}

// ReadCred returns saved decrypted credential
func (s cryptStore) ReadCred(orgID, serviceID, secretType, secretID string) (string, error) {
	logrus.Trace("Reading from CredProv: ", s.TsxvKey.CredProv.ProviderName)
	if s.TsxvKey.CredProv.ProviderName == consts.CREDPROV_HCVAULT {
		cred, err := Store.HCVReadCred(orgID, serviceID, secretID)
		if err != nil {
		 	return "", err
		}
		return cred, nil
	}

	cred, err := tsxvault.Store.GetSecret(orgID, serviceID, secretType, secretID)
	if err != nil {
		return "", err
	}

	return cred, nil

}

// RemoveCred removes particular credential from secret storage
func (s cryptStore) RemoveCred(orgID, serviceID, secretType, secretID string) error {
	logrus.Trace("Removing from CredProv: ", s.TsxvKey.CredProv.ProviderName)
	if s.TsxvKey.CredProv.ProviderName == consts.CREDPROV_HCVAULT {
		err := Store.HCVRemoveCred(orgID, serviceID, secretID)
		if err != nil {
		 	return  err
		}
		return nil
	}

	err := tsxvault.Store.TsxvDeleteSecret(orgID, serviceID, secretType, secretID)
	if err != nil {
		return err
	}

	return nil

}

// DeleteCred only used if we use hashicorp vault
func (s cryptStore) DeleteCreds(orgID, serviceID string) error {
	logrus.Trace("Deleting from CredProv: ", s.TsxvKey.CredProv.ProviderName)
	if s.TsxvKey.CredProv.ProviderName == consts.CREDPROV_HCVAULT {
		err := Store.HCVDeleteForService(orgID, serviceID)
		if err != nil {
		 	return  err
		}
		return nil
	}

return nil

}
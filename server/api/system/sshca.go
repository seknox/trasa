package system

import (
	"database/sql"
	"errors"
	"github.com/seknox/trasa/server/api/providers/ca"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"time"
)

func initCA(caType string) {
	_, err := ca.Store.GetCertHolder(consts.CERT_TYPE_SSH_CA, caType, global.GetConfig().Trasa.OrgId)
	if err == nil {
		logrus.Debugf("ssh %s CA already initialised", caType)
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		logrus.Error(err)
		return
	}

	privateKey, err := utils.GeneratePrivateKey(4096)
	if err != nil {
		logrus.Error(err)
		return
	}
	pubKey, err := utils.ConvertPublicKeyToSSHFormat(&privateKey.PublicKey)
	if err != nil {
		logrus.Error(err)
		return
	}

	privateKeyBytes := utils.EncodePrivateKeyToPEM(privateKey)

	caCert := models.CertHolder{
		CertID:      utils.GetUUID(),
		OrgID:       global.GetConfig().Trasa.OrgId,
		EntityID:    caType,
		Cert:        pubKey,
		Key:         privateKeyBytes,
		Csr:         nil,
		CertType:    consts.CERT_TYPE_SSH_CA,
		CreatedAt:   time.Now().Unix(),
		CertMeta:    "",
		LastUpdated: time.Now().Unix(),
	}
	err = ca.Store.StoreCert(caCert)
	if err != nil {
		logrus.Error(err)
		return
	}
}

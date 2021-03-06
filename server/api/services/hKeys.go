package services

import (
	"archive/zip"
	"crypto/rand"
	"encoding/binary"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/api/providers/ca"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type UpdateSSLCertsReq struct {
	SslKey    string `json:"sslKey"`
	SslCert   string `json:"sslCert"`
	CaCert    string `json:"caCert"`
	Username  string `json:"username"`
	ServiceID string `json:"serviceID"`
}

func UpdateSSLCerts(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	var req UpdateSSLCertsReq
	err := utils.ParseAndValidateRequest(r, &req)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "json parse error", "failed to update SSL certs", nil)
		return
	}

	req.SslKey = utils.NormalizeString(req.SslKey)
	req.SslCert = utils.NormalizeString(req.SslCert)
	req.CaCert = utils.NormalizeString(req.CaCert)

	//TODO Use vault to store key
	err = Store.UpdateSSLCerts(req.CaCert, "", req.SslCert, "", req.ServiceID, userContext.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not update certs", "failed to update SSL certs", nil)
		return
	}
	//var s models.ServiceSecretVault
	//s.KeyID = utils.GetRandomString(7)
	//s.ServiceID = req.ServiceID
	//s.SecretType = "db"
	//s.OrgID = userContext.Org.ID
	//
	//s.Secret = []byte(req.SslKey)
	//s.SecretID = ""
	//s.AddedAt = time.Now().Unix()
	//s.LastUpdated = time.Now().Unix()
	//
	//err:=vault.Store.StoreKey(s)

	utils.TrasaResponse(w, 200, "success", "", "SSL certs updated", nil)
	return

}

type UpdateHostCertsReq struct {
	CertVal   string `json:"certVal"`
	ServiceID string `json:"serviceID" validate:"required"`
}

func UpdateHostCerts(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	var req UpdateHostCertsReq

	err := utils.ParseAndValidateRequest(r, &req)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", err.Error(), "failed to update host certs", nil)
		return
	}

	_, _, _, _, err = ssh.ParseAuthorizedKey([]byte(req.CertVal))
	if err != nil && req.CertVal != "" {
		logrus.Debug(err)
		utils.TrasaResponse(w, 200, "failed", "Invalid format. Make sure it is in ssh known hosts format", "failed to update host certs", nil)
		return
	}

	err = Store.UpdateHostCert(req.CertVal, req.ServiceID, userContext.Org.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Could not update certs", "failed to update host certs", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "successfully updated host key", "host certs updated")
	return

}

func DownloadHostCerts(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)

	serviceID := chi.URLParam(r, "serviceID")

	bitSize := 4096
	privateKey, err := utils.GeneratePrivateKey(bitSize)
	if err != nil {
		logrus.Errorf(`could not generate private key: %v`, err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "could not generate private key", "failed to download SSH host cert", nil)
		return
	}

	sshHostCA, err := ca.Store.GetCertHolder(consts.CERT_TYPE_SSH_CA, "host", userContext.Org.ID)
	if err != nil {
		logrus.Debugf(`could not get CA key: %v`, err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "could not get CA key", "failed to download SSH host cert", nil)
		return
	}
	caKeyStr := sshHostCA.Key

	publicKeySSH, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		logrus.Errorf(`could not generate public key: %v`, err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "could not generate public key", "failed to download SSH host cert", nil)
		return
	}

	publicKeyBytes := ssh.MarshalAuthorizedKey(publicKeySSH)

	caKey, err := ssh.ParsePrivateKey(caKeyStr)
	if err != nil {
		logrus.Errorf(`Could not parse CA private key: %v`, err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "Could not parse CA private key", "failed to download SSH host cert", nil)
		return
	}

	buf := make([]byte, 8)
	_, err = rand.Read(buf)
	if err != nil {
		logrus.Errorf("failed to read random bytes: %v", err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "failed to read random bytes", "failed to download SSH host cert", nil)

		return
	}
	serial := binary.LittleEndian.Uint64(buf)

	extensions := make(map[string]string)
	extensions = map[string]string{}

	appDetail, err := Store.GetFromID(serviceID)
	if err != nil {
		logrus.Errorf("failed to get host principals: %v", err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "failed to get host principals", "failed to download SSH host cert", nil)
		return
	}

	principals := []string{appDetail.Hostname}

	cert := ssh.Certificate{
		Key:             publicKeySSH,
		Serial:          serial,
		CertType:        ssh.HostCert,
		KeyId:           serviceID,
		ValidPrincipals: principals,
		ValidAfter:      uint64(time.Now().Unix()),
		ValidBefore:     uint64(time.Now().Add(time.Hour * 24 * 30).Unix()),
		Permissions: ssh.Permissions{
			Extensions: extensions,
		},
	}

	err = cert.SignCert(rand.Reader, caKey)
	if err != nil {
		logrus.Errorf(`could not sign public key: %v`, err)
		utils.TrasaResponse(w, http.StatusOK, "failed", "could not sign public key", "failed to download SSH host cert", nil)
		return
	}

	privateKeyBytes := utils.EncodePrivateKeyToPEM(privateKey)
	certBytes := ssh.MarshalAuthorizedKey(&cert)
	if len(certBytes) == 0 {
		logrus.Errorf("failed to marshal signed certificate, empty result")
		utils.TrasaResponse(w, http.StatusOK, "failed", "failed to marshal signed certificate, empty result", "failed to download SSH host cert", nil)
		return
	}

	// Create a buffer to write our archive to.
	//buffer := new(bytes.Buffer)

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=server-certs.zip")

	// Create a new zip archive.
	zipWriter := zip.NewWriter(w)

	// Add some files to the archive.
	var files = []struct {
		Name string
		Body []byte
	}{
		{"id_rsa", privateKeyBytes},
		{"id_rsa.pub", publicKeyBytes},
		{"id_rsa-cert.pub", certBytes},
	}
	for _, file := range files {

		zipFile, err := zipWriter.Create(file.Name)
		if err != nil {
			logrus.Errorf("create host cert zip: %v", err)
		}
		_, err = zipFile.Write(file.Body)
		if err != nil {
			logrus.Errorf("write host cert zip: %v", err)
		}

	}

	//fmt.Println(zipWriter.Flush())

	// Make sure to check the error on Close.
	err = zipWriter.Close()
	if err != nil {
		logrus.Error(err)
	}

	//buffer.WriteTo(w)
	//http.ServeContent(w,r,"id_rsa.zip",time.Now(),bytes.NewReader(buffer.Bytes()))
	//w.Write(buffer)
	//	io.Copy(w,buf)
	//http.ServeContent(w,r,"id_rsa.zip",time.Now(),*buf)
	//utils.TrasaResponse(w, http.StatusOK, "success", "keypair generated and saved", "GenerateKeyPair", nil, buf.String())
	return

}

package ca

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cloudflare/cfssl/certinfo"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	"github.com/go-chi/chi"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

//TODO complete init ssh CA

//InitSSHCA creates SSH CA of given type
func InitSSHCA(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	caType := chi.URLParam(r, "type")

	privateKey, err := utils.GeneratePrivateKey(4096)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", " CA not generated", nil, nil)
		return
	}

	pubKey, err := utils.GeneratePublicKey(&privateKey.PublicKey)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", " CA not generated", nil, nil)
		return
	}

	privateKeyBytes := utils.EncodePrivateKeyToPEM(privateKey)

	ca := models.CertHolder{
		CertID:      utils.GetUUID(),
		OrgID:       uc.Org.ID,
		EntityID:    caType,
		Cert:        pubKey,
		Key:         privateKeyBytes,
		Csr:         nil,
		CertType:    consts.CERT_TYPE_SSH_CA,
		CreatedAt:   time.Now().Unix(),
		CertMeta:    "",
		LastUpdated: time.Now().Unix(),
	}
	err = Store.StoreCert(ca)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", " CA not generated", nil, nil)
		return
	}
	utils.TrasaResponse(w, 200, "success", "CA successfully generated", " CA generated", nil, nil)

}

//InitCA initialises HTTP CA of given type
func InitCA(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	logrus.Trace("request received")
	req := new(csr.CertificateRequest)
	req.KeyRequest = csr.NewKeyRequest()

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "failed to save password")
		return
	}

	//	req.CN = "TestCA"
	//	req.Names = []csr.Name{{C: "NP", ST: "3", L: "BAGMATI", O: "TestCorp", OU: "ITDEP"}}
	req.KeyRequest.A = "rsa"

	//TODO make it 4096??
	req.KeyRequest.S = 2048

	cert, csr, key, err := initca.New(req)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "InitCA-initca.New")
		return
	}

	var ca models.CertHolder
	ca.CertID = utils.GetRandomString(10)
	ca.EntityID = "HTTP_CA"
	ca.OrgID = uc.User.OrgID
	ca.Cert = cert
	ca.Csr = csr
	ca.Key = key
	ca.CertType = "HTTP_CA"
	ca.CreatedAt = time.Now().Unix()
	ca.LastUpdated = time.Now().Unix()

	err = Store.StoreCert(ca)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "InitCA-initca.New")
		return
	}

	ca.Csr = []byte("")
	ca.Key = []byte("")
	ca.Cert = []byte("")
	utils.TrasaResponse(w, 200, "success", "CA created", "InitCA-initca.New", ca)
}

//GetHttpCADetail returns HTTP CA details
func GetHttpCADetail(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	logrus.Trace("request received")
	cert, err := Store.GetCertDetail(uc.User.OrgID, "ca", consts.CERT_TYPE_HTTP_CA)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to fetch data", "GetCADetail-GetCertDetail", nil, nil)
		return
	}

	certDetail, err := certinfo.ParseCertificatePEM(cert.Cert)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to fetch cert data", "GetCADetail-certinfo.ParseCertificatePEM(")
		return
	}

	var certResp certHolderResponse
	certResp.CreatedAt = cert.CreatedAt
	certResp.LastUpdated = cert.LastUpdated
	certResp.Cert = certDetail
	certResp.CertID = cert.CertID
	certResp.CertType = cert.CertType

	utils.TrasaResponse(w, 200, "success", "CA created", "GetCADetail", certResp)

}

//GetAllCAs returns all CAs of an organization
func GetAllCAs(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	logrus.Trace("request received")
	cas, err := Store.GetAllCAs(uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to fetch data", "GetAllCA-GetCertDetail", nil, nil)
		return
	}
	var certList []certHolderResponse
	for _, cert := range cas {
		var certResp certHolderResponse
		if cert.CertType == consts.CERT_TYPE_HTTP_CA {
			certDetail, err := certinfo.ParseCertificatePEM(cert.Cert)
			if err != nil {
				logrus.Error(err)
				//TODO @sshah why is this commented?

				//utils.TrasaResponse(w, 200, "failed", "failed to fetch cert data", "GetCADetail-certinfo.ParseCertificatePEM(", nil, nil)
				//			return
			}

			certResp.Cert = certDetail

		} else {

		}

		certResp.CreatedAt = cert.CreatedAt
		certResp.LastUpdated = cert.LastUpdated
		certResp.CertID = cert.CertID
		certResp.CertType = cert.CertType
		certResp.EntityID = cert.EntityID

		certList = append(certList, certResp)

	}

	utils.TrasaResponse(w, 200, "success", "CA created", "GetCADetail", certList)

}

func DownloadSshCA(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	logrus.Trace("request received")
	caType := chi.URLParam(r, "type")

	cert, err := Store.GetCertDetail(userContext.User.OrgID, caType, consts.CERT_TYPE_SSH_CA)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to fetch data", "GetCADetail-GetCertDetail")
		return
	}

	w.Header().Set("Content-Type", "application/x-pem-file")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-ca-cert.pem", caType))

	w.Write(cert.Cert)

}

type certHolderResponse struct {
	CertID   string                `json:"certID"`
	OrgID    string                `json:"orgID"`
	EntityID string                `json:"entityID"`
	Cert     *certinfo.Certificate `json:"cert"`
	// CertificateType should be constant representing CA, intermediate CA or Service(for http?) cert others
	CertType    string `json:"certType"`
	CreatedAt   int64  `json:"createdAt"`
	LastUpdated int64  `json:"lastUpdated"`
}

//UploadCA uploads new HTTP CA from user
func UploadCA(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	var req struct {
		CertVal string `json:"certVal"`
		KeyVal  string `json:"keyVal"`
		CsrVal  string `json:"csrVal"`
	}

	err := utils.ParseAndValidateRequest(r, &req)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "Json parse error", "Upload CA")
		return
	}

	var ca models.CertHolder
	ca.CertID = utils.GetRandomString(10)
	ca.EntityID = "ca"
	ca.OrgID = uc.User.OrgID
	ca.Cert = []byte(req.CertVal)
	ca.Csr = []byte(req.CsrVal)
	ca.Key = []byte(req.KeyVal)
	ca.CertType = consts.CERT_TYPE_HTTP_CA
	ca.CreatedAt = time.Now().Unix()
	ca.LastUpdated = time.Now().Unix()

	err = Store.StoreCert(ca)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "UploadCA-UploadCA.New")
		return
	}

	ca.Csr = []byte("")
	ca.Key = []byte("")
	ca.Cert = []byte("")
	utils.TrasaResponse(w, 200, "success", "CA created", "IniUploadCAtCA-UploadCA.New", ca)
}

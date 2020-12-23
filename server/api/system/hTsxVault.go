package system

import (
	"crypto/sha512"
	"fmt"
	"net/http"
	"strings"
	"time"

	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	"github.com/go-chi/chi"
	hashicorpVault "github.com/hashicorp/vault/api"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/providers/vault"
	"github.com/seknox/trasa/server/api/providers/vault/tsxvault"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

// StoreKey stores keys in database.
// Keys should be encrypted and tag value must be generated.
// Before storing key, check if the key is valid and working.
func StoreKey(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	var req models.KeysHolderReq

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request.", "StoreKey", "")
		return
	}

	var store models.KeysHolder
	store.OrgID = uc.User.OrgID
	store.KeyID = utils.GetRandomString(5)
	store.KeyTag = fmt.Sprintf("%s-xxxx-xxxx...", req.KeyVal[0:4])
	store.AddedBy = uc.User.ID
	store.AddedAt = time.Now().Unix()
	store.KeyVal = []byte(req.KeyVal)
	store.KeyName = req.KeyName

	_, err := EncryptAndStoreKeyOrToken(store)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", err.Error(), "StoreKey", nil)
		return
	}
	utils.TrasaResponse(w, 200, "success", "keys stored", "StoreKey", req.KeyTag)
}

// Getkey retrieves key or token from database. should fetch and return key tag rather than key value.
func Getkey(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	vendorID := chi.URLParam(r, "vendorID")

	key, err := tsxvault.Store.GetKeyOrTokenWithTag(uc.User.OrgID, vendorID)
	if err != nil {
		logrus.Error(err, vendorID)
		utils.TrasaResponse(w, 200, "failed", "failed to get token.", "Getkey-GetKeyOrTokenWithTag", nil)
		return
	}

	var resp models.KeysHolderReq
	resp.KeyName = key.KeyName
	resp.KeyID = key.KeyID
	resp.KeyVal = key.KeyTag
	resp.KeyTag = key.KeyTag
	resp.AddedAt = key.AddedAt
	resp.AddedBy = key.AddedBy

	utils.TrasaResponse(w, 200, "success", "key fetched", "Getkey", resp)
}

//EncryptAndStoreKeyOrToken is helper function which encrypts key or token and store it in database.
func EncryptAndStoreKeyOrToken(req models.KeysHolder) ([]byte, error) {

	ct, err := tsxvault.Store.AesEncrypt([]byte(req.KeyVal))
	req.KeyVal = ct
	if err != nil {
		//logrus.Error(err)
		return nil, err
	}

	err = tsxvault.Store.StoreKeyOrTokens(req)
	if err != nil {
		logrus.Error(err)
		return nil, fmt.Errorf("failed to store token")
	}
	return req.KeyVal, nil
}

type VaultInit struct {
	SecretShares    int `json:"secretShares"`
	SecretThreshold int `json:"secretThreshold"`
}
type VaultInitResp struct {
	UnsealKeys   []string `json:"unsealKeys"`
	DecryptKeys  []string `json:"decryptKeys"`
	EncRootToken string   `json:"encRootToken"`
	Tsxvault     bool     `json:"tsxvault"`
}

// TsxvaultInit initializes TRASA built in secure storage. master key for encryption is
// Shamir'ed into 5 keys with minimum 3 keys threshold and responded back to administrator.
func TsxvaultInit(w http.ResponseWriter, r *http.Request) {

	uc := r.Context().Value("user").(models.UserContext)

	var req VaultInit

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "Vault not initialised", nil)
		return
	}

	// Initing tsxvault process: (1) Create encryption key. (2) Create master key and encrypt encryption key with it.
	// (3) shard master key and give it to user. (4) store encrypted encryption key in database.
	encKeyShards, err := InitTsxvault(uc.User.OrgID, uc.User.ID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", err.Error(), "failed to rinitialize", nil)
		return
	}

	var resp VaultInitResp
	resp.DecryptKeys = encKeyShards
	resp.Tsxvault = true

	utils.TrasaResponse(w, 200, "success", "sucess", "Vault initialised", resp)

}

// InitTsxvault generates aws encryption key. shards it with sharder and returns sharded key. It also updates global setting.
func InitTsxvault(orgID, userID string) ([]string, error) {

	encKey, err := tsxvault.Store.GenAndStoreKey(orgID)
	if err != nil {
		return nil, err
	}

	shardedKeys := utils.ShamirSharder(encKey[:], 5, 3)

	var store models.GlobalSettings
	store.SettingID = utils.GetUUID()
	store.OrgID = orgID
	store.Status = true
	store.SettingType = consts.GLOBAL_TSXVAULT

	var vaultFeature models.CredProvProps
	vaultFeature.ProviderName = consts.CREDPROV_TSXVAULT
	vaultFeature.ProviderAddr = ""
	vaultFeature.ProviderAccessToken = ""
	jsonV, err := json.Marshal(vaultFeature)
	if err != nil {
		return nil, err
	}

	store.SettingValue = string(jsonV)
	store.UpdatedBy = userID
	store.UpdatedOn = time.Now().Unix()

	err = Store.UpdateGlobalSetting(store)
	if err != nil {

		return nil, err
	}

	return shardedKeys, nil
}

// ReInit purpose is to delete exisiting vault configs and instances from database.
// Clients should immediately send another request to vault init when this handler returns success response.
func ReInit(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	// (1) remove all managed users for this organization.
	err := orgs.Store.RemoveAllManagedAccounts(uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to remove manged users but vault storage is removed.", "Vault not reinitialised", nil)
		return
	}

	err = tsxvault.Store.TsxvdeactivateAllKeys(uc.User.OrgID, time.Now().Unix())
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to remove manged users but vault storage is removed.", "Vault not reinitialised", nil)
		return
	}

	// delete all rows from Service_keyvaultv1
	err = tsxvault.Store.TsxvDeleteAllSecret(uc.User.OrgID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to remove manged users but vault storage is removed.", "Vault not reinitialised", nil)
		return
	}

	utils.TrasaResponse(w, 200, "success", "submit-init", "Vault reinitialised", nil)

}

type VaultStatus struct {
	InitStatus  models.GlobalSettings              `json:"initStatus"`
	SealStatus  *hashicorpVault.SealStatusResponse `json:"sealStatus"`
	TokenStatus hashicorpVault.SealStatusResponse  `json:"tokenStatus"`
	// TsxVault is TRASA's built in tsxvault. if false, caller should assume hashicorp vault is used instead.
	Tsxvault bool   `json:"tsxvault"`
	Setting  string `json:"setting"`
}

// Status returns vault's current status.
func Status(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)

	vaultInitStatus, err := Store.GetGlobalSetting(uc.Org.ID, consts.GLOBAL_TSXVAULT)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "successfully retrieved org details", "GetOrgDetail", nil)
		return
	}

	var resp VaultStatus
	resp.InitStatus = vaultInitStatus

	resp.Tsxvault = global.GetConfig().Vault.Tsxvault

	resp.Setting = vaultInitStatus.SettingValue

	_, status := tsxvault.Store.GetTsxVaultKey()
	if status == false {
		resp.TokenStatus = hashicorpVault.SealStatusResponse{Sealed: true}
		utils.TrasaResponse(w, 200, "failed", "This key is not retrieved", "Vault not decrypted", resp)
		return
	}

	resp.TokenStatus = hashicorpVault.SealStatusResponse{Sealed: false}

	utils.TrasaResponse(w, 200, "success", "successfully retrieved org details", "GetOrgDetail", resp)

}

// HoldDecryptShard hods state of encryption key retreival during shamir deduce function
var HoldDecryptShard [][]byte

type unseal struct {
	Key string `json:"key"`
}

// DecryptKey retrieves token from vaultDecrypt function and
// store it in vaultEncryption Token. This is only available option for tsxtsxvault.
// TODO @sshahcodes compose this handler to smaller functions
func DecryptKey(w http.ResponseWriter, r *http.Request) {
	uc := r.Context().Value("user").(models.UserContext)
	var req unseal
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "Vault not decrypted", nil)
		return
	}

	trimmed := strings.TrimSpace(req.Key)
	key, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "decode error", "Vault not decrypted", nil)
		return
	}

	var resp hashicorpVault.SealStatusResponse

	HoldDecryptShard = append(HoldDecryptShard, key)

	if len(HoldDecryptShard) < 3 {

		resp.Sealed = true
		resp.Progress = len(HoldDecryptShard)
		utils.TrasaResponse(w, 200, "success", "next-key", "Vault decrypted", resp)
		return
	}
	// try to deduce the key. if deducing is success, also test if the key works. else return failed.

	deducedVal, err := utils.ShamirDeducer(HoldDecryptShard)
	if err != nil {
		logrus.Error("ShamirDeducer ", err)
		HoldDecryptShard = nil
		utils.TrasaResponse(w, 200, "failed", "unable to deduce key. Retry again from 1st key", "Vault decrypted", resp)
		return
	}

	hash := sha512.New()

	var buf []byte
	buf = append(buf, deducedVal[:]...)
	hashed := hash.Sum([]byte(buf))

	// check if token is valid.
	// this is verified by fetching value from TsxvGetEncKeyHash. if this fails or returned key is expired,
	// we return failed response. other wise will store it in TsxVaultKey.
	getKey, err := tsxvault.Store.TsxvGetEncKeyHash(uc.User.OrgID, hex.EncodeToString(hashed))
	if err != nil {
		logrus.Error(err)
		HoldDecryptShard = nil
		utils.TrasaResponse(w, 200, "failed", "This key is not registered", "Vault not decrypted", nil)
		return
	}
	if getKey.Status == false {
		logrus.Error("failed status")
		HoldDecryptShard = nil
		utils.TrasaResponse(w, 200, "failed", "This key is expired", "Vault not decrypted", nil)
		return
	}

	if hex.EncodeToString(hashed) != getKey.KeyHash {
		logrus.Error("hash mismatch")
		HoldDecryptShard = nil
		utils.TrasaResponse(w, 200, "failed", "This key is not valid", "Vault not decrypted", nil)
		return
	}

	// reaching here means we can set encryption key and state
	nkey := new([32]byte)
	copy(nkey[:], deducedVal)

	// Get global vault settings
	vaultsetting, err := Store.GetGlobalSetting(uc.Org.ID, consts.GLOBAL_TSXVAULT)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "unable retrieved org details", "DecryptKey", nil)
		return
	}



	// store cred prov setting in global tsxvkey struct
	var cred models.CredProvProps
	err = json.Unmarshal([]byte(vaultsetting.SettingValue), &cred)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "unable to unmarshal setting values", "DecryptKey", nil)
		return
	}

		// get access token from keyholder if credprov is hashicorp vault
		if cred.ProviderName == consts.CREDPROV_HCVAULT {
			ct, err := vault.Store.GetKeyOrTokenWithKeyval(uc.User.OrgID, string(consts.CREDPROV_HCVAULT_TOKEN))
			if err != nil {
				logrus.Error(err)
			}
		
			pt, err := utils.AESDecrypt(nkey[:], ct.KeyVal)
			if err != nil {
				logrus.Error(err)
			}

			cred.ProviderAccessToken = string(pt)
		}
	

	tsxvault.Store.SetTsxVaultKey(nkey, true, cred)

	HoldDecryptShard = HoldDecryptShard[:0]

	resp.Sealed = false
	resp.Progress = len(HoldDecryptShard)

	// retreive trasaCPxy api key here

	// get key ct from database.
	apikey, err := vault.Store.GetKeyOrTokenWithKeyval(uc.User.OrgID, consts.GLOBAL_CLOUDPROXY_APIKEY)
	if err != nil {
		logrus.Error(err)
	}

	pt, err := utils.AESDecrypt(nkey[:], apikey.KeyVal)
	if err != nil {
		logrus.Error(err)
	}

	tsxvault.Store.SetTsxCPxyKey(string(pt))

	utils.TrasaResponse(w, 200, "success", "token retrieved", "Vault decrypted", resp)

}

// UpdateCredProv changes vault credential provider setting. (e.g. where to store service credentials, tsxVault or external secret provider)
func UpdateCredProv(w http.ResponseWriter, r *http.Request) {

	uc := r.Context().Value("user").(models.UserContext)

	var req models.CredProvProps

	if err := utils.ParseAndValidateRequest(r, &req); err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "invalid request", "secret storage not updated", nil)
		return
	}

	// update crdprov config in global variable
	tsxvault.Store.UpdateTsxVaultKeyCredProvConfig(req)

	// check and store access token
	start := ""
	if len(req.ProviderAccessToken) > 4 {
		start = req.ProviderAccessToken[0:4]
	}

	var key models.KeysHolder
	key.OrgID = uc.User.OrgID
	key.KeyID = utils.GetRandomString(5)
	key.KeyTag = fmt.Sprintf("%sxxxx-xxxx...", start)
	key.AddedBy = uc.User.ID
	key.AddedAt = time.Now().Unix()
	key.KeyName = string(consts.CREDPROV_HCVAULT_TOKEN)
	key.KeyVal = []byte(req.ProviderAccessToken)
	_, err := EncryptAndStoreKeyOrToken(key)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "vault key is not retreived yet.", "could not store access token")
		return
	}

	req.ProviderAccessToken = key.KeyTag
	
	// update global settting
	var store models.GlobalSettings
	store.OrgID = uc.Org.ID
	store.Status = true
	store.SettingType = consts.GLOBAL_TSXVAULT


	jsonV, err := json.Marshal(req)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "failed to marshal vault features", "secret storage not updated", nil)
		return
	}

	store.SettingValue = string(jsonV)
	store.UpdatedBy = uc.User.ID
	store.UpdatedOn = time.Now().Unix()

	err = Store.UpdateGlobalSetting(store)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", err.Error(), "secret storage not updated", nil)
		return
	}




	utils.TrasaResponse(w, 200, "success", "sucess", "Vault initialised", nil)

}

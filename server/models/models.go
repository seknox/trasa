package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/seknox/trasa/server/consts"

	firebase "firebase.google.com/go"
	"github.com/minio/minio-go"
	"github.com/oschwald/geoip2-golang"
	"github.com/seknox/trasa/server/global"
	"github.com/tstranex/u2f"

	"github.com/go-redis/redis/v8"
)

//Gstate is a global state struct which contains database connections, configurations etc
type Gstate struct {
	db             *sql.DB
	geoip          *geoip2.Reader
	firebaseClient *firebase.App
	minioClient    *minio.Client
	config         global.Config
	redisClient    *redis.Client
}

type GeoLocation struct {
	IsoCountryCode string    `json:"isoCountryCode"`
	City           string    `json:"city"`
	Country        string    `json:"country"`
	TimeZone       string    `json:"timeZone"`
	Location       []float64 `json:"location"`
}

func (a GeoLocation) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *GeoLocation) Scan(value interface{}) error {
	if value == nil {
		*a = GeoLocation{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}

	return nil
}

type TrasaResponseStruct struct {
	Status string      `json:"status"`
	Error  error       `json:"error,omitempty"`
	Reason string      `json:"reason,omitempty"`
	Intent string      `json:"intent,omitempty"`
	Data   interface{} `json:"data"`
}

type TrasaResponseStructWIthDataString struct {
	Status string `json:"status"`
	Error  error  `json:"error,omitempty"`
	Reason string `json:"reason,omitempty"`
	Intent string `json:"intent,omitempty"`
	Data   string `json:"data"`
}

type ScimContext struct {
	OrgID    string `json:"orgID"`
	Orgname  string `json:"orgName"`
	IdpID    string `json:"idpID"`
	IdpName  string `json:"idpName"`
	TimeZone string `json:"timeZone"`
}

type InitSignup struct {
	OrgID          string `json:"orgID"`
	UserID         string `json:"userID"`
	OrgName        string `json:"orgName"`
	PrimaryContact string `json:"primaryContact" valid:"email"`
	UserName       string `json:"userName" valid:"alphanum"`
	FirstName      string `json:"firstName" valid:"alpha"`
	MiddleName     string `json:"middleName" valid:"alpha"`
	LastName       string `json:"lastName" valid:"alpha"`
	Email          string `json:"email" valid:"email"`
	Password       string `json:"password"`
	UserRole       string `json:"userRole"`
	Company        string `json:"companyName"`
	JobTitle       string `json:"jobTitle"`
	PhoneNumber    string `json:"phoneNumber"`
	Country        string `json:"country"`
	Timezone       string `json:"timezone"`
	Reference      string `json:"reference"`
	LicenseType    string `json:"licenseType"`
	CreatedAt      string
	UpdatedAt      string
	DeletedAt      string
}

// AppLogin mocks request structure which ssh logins and rdp logins generates
type ServiceLogin struct {
	ServiceID       string           `json:"serviceID"`
	DynamicService  bool             `json:"dynamicService"`
	ServiceKey      string           `json:"serviceKey"`
	User            string           `json:"user"`
	Password        string           `json:"password"`
	PublicKey       []byte           `json:"publicKey"`
	TfaMethod       string           `json:"tfaMethod"`
	TotpCode        string           `json:"totpCode"`
	UserIP          string           `json:"userIP"`
	UserWorkstation string           `json:"workstation"`
	TrasaID         string           `json:"trasaID"`
	SessionID       string           `json:"sessionID"`
	IsSharedSession bool             `json:"isSharedSession"`
	AppType         string           `json:"appType"`
	RdpProtocol     string           `json:"rdpProto"`
	OrgID           string           `json:"orgID"`
	Hostname        string           `json:"hostname"`
	Skip2FA         bool             `json:"skip2FA"`
	SignResponse    u2f.SignResponse `json:"signResponse"`
	DeviceHygiene   DeviceHygiene    `json:"deviceHygiene"`
}

type IPDetails struct {
	IpAddress      string `json:"IPAddr"`
	NetMask        string `json:"netMask"`
	DefaultGateway string `json:"defaultGateway"`
}

func (a IPDetails) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *IPDetails) Scan(value interface{}) error {
	if value == nil {
		*a = IPDetails{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}

	return nil
}

// func (a App) Value() (driver.Value, error) {
// 	return json.Marshal(a)
// }

// // Make the Attrs struct implement the sql.Scanner interface. This method
// // simply decodes a JSON-encoded value into the struct fields.
// func (a *App) Scan(value interface{}) error {
// 	b, ok := value.([]byte)
// 	if !ok {
// 		return errors.New("type assertion to []byte failed")
// 	}

// 	err := json.Unmarshal(b, &a)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// AuthRequest struct. Authentication Request:
type AuthRequest struct {
	RequestID    string
	Scopes       []string
	ClientID     string
	RedirectURI  string
	ResponseType []string
	state        string
}

// Authorization struct defines authorization event. server reference this event to generate access tokens
type Authorization struct {
	AuthorizationID string
	serviceID       string
	UserID          string
	Scopes          []string
	Nonce           string
	CreatedAt       string
}

// AccessToken defines api access token structures.
type AccessToken struct {
	GrantID     string
	serviceID   string
	UserID      string
	AccessToken string
	IDToken     string
	Scopes      []string
	CreatedAt   string
	TTLValue    string
}

type ServiceUserMap struct {
	MapID     string `json:"mapID"`
	ServiceID string `json:"serviceID"`
	OrgID     string `json:"orgID"`
	UserID    string `json:"userID"`
	PolicyID  string `json:"policyID"`
	Privilege string `json:"username"`
	AddedAt   int64  `json:"addedAt"`
}

type ServiceGroupUserGroupMap struct {
	MapID          string `json:"mapID"`
	ServiceGroupID string `json:"serviceGroupID"`
	MapType        string `json:"mapType"`
	UserGroupID    string `json:"userGroupID"`
	Privilege      string `json:"privilege"`
	OrgID          string `json:"orgID"`
	PolicyID       string `json:"policyID"`
	CreatedAt      int64  `json:"createdAt"`
}

type InAppNotification struct {
	NotificationID    string `json:"notificationID"`
	UserID            string `json:"userID"`
	EmitterID         string `json:"emitterID"`
	OrgID             string `json:"orgID"`
	NotificationLabel string `json:"notificationLabel"`
	NotificationText  string `json:"notificationText"`
	CreatedOn         int64  `json:"createdOn"`
	IsResolved        bool   `json:"isResolved"`
	ResolvedOn        int64  `json:"resolvedOn"`
}

type VaultCredStorageEvent struct {
	OrgID              string `json:"orgID"`
	AccessedBy         string `json:"accessedBy"`
	FetchedForApp      string `json:"fetchedForApp"`
	FetchedForUsername string `json:"fetchedForUsername"`
	AccessedOn         int64  `json:"accessdOn"`
}

type VaultAccessLogs struct {
	OrgID              string `json:"orgID"`
	AccessedBy         string `json:"accessedBy"`
	FetchedForApp      string `json:"fetchedForApp"`
	FetchedForUsername string `json:"fetchedForUsername"`
	AccessedOn         int64  `json:"accessdOn"`
}

type TRASAFeaturesStatus struct {
	OrgID     string `json:"orgID"`
	Feature   string `json:"feature"`
	InitBy    string `json:"initBy"`
	Status    bool   `json:"status"`
	Remarks   string `json:"remarks"`
	InitOn    int64  `json:"initOn"`
	UpdatedOn int64  `json:"updatedOn"`
	Config    string `json:"config"`
}

// VaultFeature stores information regarding where is the secret stored (or to be stored)
// For example VaultFeature.CredStorage value can be tsxvault or aws secret storage.
// If tsxvault is set, we store user credentials in our built in vault. if aws is set,
// we push secrets to aws secret storage. What happens if user wants to migrate secret from
// tsxvault to aws secret storage? --- migration code required...
// 3rd party api keys which is used by TRASA will always be stored in key_holderv1.
// Only one secret storage provider is supported at given time.
type VaultFeature struct {
	// CredStorage is for storing user credentials(uname:pass) or (uname:privatekey)
	CredStorage string `json:"credStorage"`
	// CertStorage determines where ca certificates and private eys are stored. it can be stored in cert_holder
	// or external ca storage.
	CertStorage string `json:"certStorage"`
}

type BackupPlan struct {
	OrgID              string   `json:"orgID"`
	BackupPlanID       string   `json:"backupPlanID"`
	BackupPlanName     string   `json:"backupPlanName"`
	BackupType         string   `json:"backupType"`
	ScheduleTime       int64    `json:"scheduleTime"`
	Interval           string   `json:"interval"`
	BackupServiceNames []string `json:"backupServiceNames"`
	CreatedAt          int64    `json:"createdAt"`
	UpdatedAt          int64    `json:"updatedAt"`
}

type Backup struct {
	OrgID      string `json:"orgID"`
	BackupID   string `json:"backupID"`
	BackupName string `json:"backupName"`
	BackupType string `json:"backupType"`
	CreatedAt  int64  `json:"createdAt"`
}

type MyService struct {
	AccessMapDetail
	Adhoc        bool                `json:"adhoc"`
	Usernames    []string            `json:"usernames"`
	IsAuthorised bool                `json:"isAuthorised"`
	Reason       consts.FailedReason `json:"reason"`
}

// PolicyEnforcer type is generic policy enforcement model which can be used to assign and track user's
// for specefic enforced action that assigned user must perform. eg. change password, change username etc...
type PolicyEnforcer struct {
	// EnforceID is unique id for the event
	EnforceID string `json:"enforceID"`
	// userID represents user who is effected by this policy
	UserID string `json:"userID"`
	OrgID  string `json:"orgID"`
	// EnforceType refers to unique constant for this event type. eg change password? username?
	EnforceType string `json:"enforceType"`
	// Status of the event. True means pending. False means resolved.
	Pending bool `json:"status"`
	// AssignedBy can be either system assigned or assigned by administrator.
	// In case of system assigned, use constat else the value must be userID of administrator.
	AssignedBy string `json:"assignedBy"`
	AssignedOn int64  `json:"assignedOn"`
	ResolvedOn int64  `json:"resolvedOn"`
}

// ComplianceViolation is triggered based on violation of compliance requirements.
type ComplianceViolation struct {
	// ViolationID is unique id for the event
	ViolationID string `json:"violationID"`
	OrgID       string `json:"orgID"`

	// EntityType can be either user or Service or any entity type
	EntityType string `json:"entityType"`
	// EntityID is unique ID of entityType in scope
	EntityID       string `json:"entityID"`
	ComplianceType string `json:"ComplianceType"`
	ComplReqID     string `json:"complReqID"`
	CompleReqDesc  string `json:"compleReqDesc"`
	// ViolationType should be based on constant value of violation
	ViolationType string `json:"violationType"`
	ReportedOn    int64  `json:"reportedOn"`
	ResolvedOn    int64  `json:"resolvedOn"`
}

// PasswordPolicy represents global policy for passwords that are used to log into TRASA dashboard.
// This is stored as settingValue in GlobalSettings for settingType as passwordPolicy
type PasswordPolicy struct {
	Expiry            string `json:"expiry"`
	MinimumChars      int    `json:"minimumChars"`
	EnforceStrongPass bool   `json:"enforceStrongPass"`
	ZxcvbnScore       int    `json:"zxcvbnScore"`
}

type AlertPolicy struct {
	PolicyID     string `json:"policyID"`
	PolicyName   string `json:"policyName"`
	OrgID        string `json:"orgID"`
	NotifCase    string `json:"notifCase"`
	NotifChannel string `json:"notifChannel"`
	NotifyTo     string `json:"notifyTo"`
	CreatedBy    string `json:"createdBy"`
	CreatedAt    int64  `json:"createdAt"`
	UpdatedAt    int64  `json:"updatedAt"`
}

type SecurityRule struct {
	RuleID       string `json:"ruleID"`
	OrgID        string `json:"orgID"`
	Name         string `json:"name"`
	ConstName    string `json:"constName"`
	Description  string `json:"description"`
	Scope        string `json:"scope"`
	Condition    string `json:"condition"`
	Status       bool   `json:"status"`
	Source       string `json:"source"`
	Action       string `json:"action"`
	CreatedBy    string `json:"createdBy"`
	CreatedAt    int64  `json:"createdAt"`
	LastModified int64  `json:"lastModified"`
}

type SecurityRuleViolationAction struct {
	ActionName       string   `json:"actionName"`
	ActionType       string   `json:"actionType"`
	AffectedGroups   []string `json:"affectedGroups"`
	AffectedEntities []string `json:"affectedEntities"`
}

type Entity struct {
	EntityType string
	EntityName string
	EntityDesc string
	EntityID   string
}

type Intent struct {
	IntentType    string   `json:"intentType"`
	MainEntity    Entity   `json:"mainEntity"`
	OtherEntities []Entity `json:"otherEntities"`
	DescString    string   `json:"descString"`
}

// PasswordState holds status for user passwords
type PasswordState struct {
	UserID        string   `json:"userID"`
	OrgID         string   `json:"orgID"`
	LastPasswords []string `json:"lastPasswords"`
	LastUpdated   int64    `json:"lastUpdated"`
}

type ErrorStrings struct {
	OrgId  string `json:"orgId"`
	UserId string `json:"userId"`
	Status string `json:"status"`
	Error  error  `json:"error,omitempty"`
	Reason string `json:"reason,omitempty"`
	Intent string `json:"intent,omitempty"`
}

type ResponseStruct struct {
	Status string `json:"status"`
	Error  error  `json:"error,omitempty"`
	Reason string `json:"reason,omitempty"`
	Intent string `json:"intent,omitempty"`
}

// KeyStore stores encryption keys
type KeyStore struct {
	KeyID       string `json:"keyID"`
	OrgID       string `json:"orgID"`
	Key         string `json:"key"`
	CreatedAt   int64  `json:"createdAt"`
	LastUpdated int64  `json:"lastUpdated"`
}

type ServiceSecretVault struct {
	KeyID       string `json:"keyID"`
	OrgID       string `json:"orgID"`
	ServiceID   string `json:"serviceID"`
	SecretType  string `json:"secretType"`
	UpdatedBy   string `json:"updatedBy"`
	SecretID    string `json:"secretID"`
	Secret      []byte `json:"secret"`
	AddedAt     int64  `json:"addedAt"`
	LastUpdated int64  `json:"lastUpdated"`
}

// AccountSecrets holds secrets for users which needs to be stored in trasaVault
type AccountSecrets struct {
	Username    string `json:"userName"`
	Secret      string `json:"secret"`
	AddedAt     int64  `json:"addedAt"`
	LastUpdated int64  `json:"lastUpdated"`
}

// CertHolder holds certificate data.
type CertHolder struct {
	CertID   string `json:"certID"`
	OrgID    string `json:"orgID"`
	EntityID string `json:"entityID"`
	Cert     []byte `json:"cert"`
	Key      []byte `json:"key"`
	Csr      []byte `json:"csr"`
	// CertificateType should be constant representing CA, intermediate CA or Service(for http?) cert others
	CertType  string `json:"certType"`
	CreatedAt int64  `json:"createdAt"`
	// CertMeta holds metadata for generating or signing other certs.
	// This metadata is only valid as default parameters and can be override by Service specefic metadata.
	// For example default generated client cert expiry time might be 24 hours but specefic Service can allow
	// access only for 1 hour or 1 time access as 1 minute validity.
	CertMeta    string `json:"certMeta"`
	LastUpdated int64  `json:"lastUpdated"`
}

// KeysHolder stores access key supplied by administrators for managing external resources.
// E.g. api keys, tokens etc.
type KeysHolder struct {
	KeyID   string `json:"keyID"`
	OrgID   string `json:"orgID"`
	KeyTag  string `json:"keyTag"`
	KeyName string `json:"keyName"`
	//KeyVal  string `json:"keyVal"`
	KeyVal  []byte `json:"keyVal"`
	AddedBy string `json:"addedBy"`
	AddedAt int64  `json:"addedAt"`
}

type KeysHolderReq struct {
	KeyID   string `json:"keyID"`
	OrgID   string `json:"orgID"`
	KeyTag  string `json:"keyTag"`
	KeyName string `json:"keyName"`
	KeyVal  string `json:"keyVal"`
	AddedBy string `json:"addedBy"`
	AddedAt int64  `json:"addedAt"`
}

// CloudIaaSSync tracks synchronization with cloud service provider.
type CloudIaaSSync struct {
	CloudIaasID   string `json:"cloudIaasID"`
	OrgID         string `json:"orgID"`
	CloudIaasName string `json:"cloudIaasName"`
	LasgtSyncedBy string `json:"LasgtSyncedBy"`
	LastSyncedOn  int64  `json:"keyTag"`
}

type EncryptionKeyLog struct {
	KeyID       string `json:"keyID"`
	OrgID       string `json:"orgID"`
	KeyHash     string `json:"keyHash"`
	GeneratedAt int64  `json:"generatedAt"`
	Status      bool   `json:"status"`
	LastUpdated int64  `json:"lastUpdated"`
}

type AccessMapDetail struct {
	MapID       string `json:"mapID"`
	ServiceID   string `json:"serviceID"`
	ServiceName string `json:"serviceName"`
	ServiceType string `json:"serviceType"`
	Hostname    string `json:"hostname"`
	OrgID       string `json:"orgID"`
	UserID      string `json:"userID"`
	Email       string `json:"email"`
	Policy      Policy `json:"policy"`
	Privilege   string `json:"privilege"`
	UserAddedAt int64  `json:"userAddedAt"`
}

type GlobalEmailSetting struct {
	EmailSettingID    string                 `json:"emailSettingID"`
	IntegrationType   string                 `json:"integrationType"`
	IntegrationConfig EmailIntegrationConfig `json:"emailIntegrationConfig"`
	IsEnabled         bool                   `json:"isEnabled"`
	UpdatedAt         int64                  `json:"updatedAt"`
}

type EmailIntegrationConfig struct {
	IntegrationType string `json:"integrationType"`
	// AuthEmailAddr and AuthEmailPass is email:pass  that will be used for smtp authentication.
	// Incase of api integration, this holds api key and api keyvalue respectively .
	AuthKey       string `json:"authKey"`
	AuthPass      string `json:"authPass"`
	ServerAddress string `json:"serverAddress"`
	ServerPort    string `json:"serverPort"`
	SenderAddress string `json:"senderAddress"`
}

type GlobalTrasaSshAuth struct {
	MandatoryCertAuth bool `json:"mandatoryCertAuth"`
}

type EmailAdhoc struct {
	Requester     string   `json:"requester"`
	Requestee     string   `json:"requestee"`
	ReceiverEmail string   `json:"receiverEmail"`
	CC            []string `json:"cc"`
	DashLink      string   `json:"dashLink"`
	App           string   `json:"app"`
	Reason        string   `json:"reason"`
	Status        string   `json:"status"`
	Time          string   `json:"time"`
	Subject       string   `json:"subject"`
	Req           bool     `json:"req"`
}

type EmailDynamicAccess struct {
	User          string   `json:"user"`
	AppType       string   `json:"appType"`
	Hostname      string   `json:"hostname"`
	ReceiverEmail string   `json:"receiverEmail"`
	TimeInt       int64    `json:"timeInt"`
	CC            []string `json:"cc"`
}

type EmailSecurityAlert struct {
	ReceiverEmail     string   `json:"receiverEmail"`
	SecurityRuleTitle string   `json:"securityRuleTitle"`
	SecurityRuleText  string   `json:"securitRuleText"`
	EntityName        string   `json:"entityName"`
	CC                []string `json:"cc"`
}

type EmailUserCrud struct {
	ReceiverEmail string   `json:"receiverEmail"`
	Username      string   `json:"username"`
	VerifyUrl     string   `json:"verifyUrl"`
	NewM          bool     `json:"newM"`
	CC            []string `json:"cc"`
}

//////////////////////////////////////////////
///////// FCM server config
//////////////////////////////////////////////
type FcmConfig struct {
	ServerAddr string `json:"serverAddr"` // should be fqdn with path
	FcmApiKey  string `json:"fcmApiKey"`
}

type AppUserPermission struct {
	Policy       Policy
	Username     string
	AdHocEnabled bool
}

// Current Database version
type DBVersion struct {
	DBVersion string `json:"dbVersion"`
	CreatedOn int64  `json:"createdOn"`
}

type AdhocPermission struct {
	RequestID        string   `json:"reqID"`
	RequesterID      string   `json:"requesterID"`
	OrgID            string   `json:"orgID"`
	ServiceID        string   `json:"serviceID"`
	RequesteeID      string   `json:"requesteeID"`
	RequestTxt       string   `json:"requestTxt"`
	RequestedOn      int64    `json:"reqTime"`
	IsAuthorized     bool     `json:"isAuthorized"`
	AuthorizedOn     int64    `json:"authorizedOn"`
	AuthorizedPeriod int64    `json:"authorizedPeriod"`
	AuthorizedPolicy Policy   `json:"authorizedPolicy"`
	IsExpired        bool     `json:"isExpired"`
	SessionID        []string `json:"sessionID"`
}

type AdhocDetails struct {
	AdhocPermission
	ServiceName    string `json:"serviceName"`
	ServiceType    string `json:"serviceType"`
	RequesterEmail string `json:"requesterEmail"`
	RequesteeEmail string `json:"requesteeEmail"`
}

type MyServiceDetails struct {
	MyService
	Adhoc          bool                `json:"adhoc"`
	ServiceType    string              `json:"serviceType"`
	Hostname       string              `json:"hostname"`
	IsAdmin        bool                `json:"isAdmin"`
	Usernames      []string            `json:"usernames"`
	AuthorizedTill int64               `json:"authorizedTill"`
	AuthorizedOn   int64               `json:"authorizedOn"`
	RequestedOn    int64               `json:"requestedOn"`
	IsAuthorised   bool                `json:"isAuthorised"`
	Reason         consts.FailedReason `json:"reason"`
}

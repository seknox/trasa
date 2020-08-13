package models

//NewEmptyUserStruct returns a empty User struct
func NewEmptyUserStruct() User {
	return User{
		ID:         "",
		OrgID:      "",
		UserName:   "",
		FirstName:  "",
		MiddleName: "",
		LastName:   "",
		Email:      "",
		Groups:     nil,
		UserRole:   "",
		Status:     false,
		IdpName:    "",
		ExternalID: "",
		CreatedAt:  0,
		UpdatedAt:  0,
	}
}

//User Model stores behaviours related to single user
type User struct {
	ID         string   `json:"ID" `
	OrgID      string   `json:"orgId"`
	UserName   string   `json:"userName" validate:"alphanum"`
	FirstName  string   `json:"firstName" validate:"alpha"`
	MiddleName string   `json:"middleName" validate:"omitempty,alpha"`
	LastName   string   `json:"lastName" validate:"alpha"`
	Email      string   `json:"email" validate:"email"`
	Groups     []string `json:"groups"`
	UserRole   string   `json:"userRole"  valid:"alpha"`
	Status     bool     `json:"status"`
	// IdpName is name of identity provider for user. can be 'trasa' or 'okta' etc..
	IdpName string `json:"idpName"`
	// ExternalID is ID of service that exists outside of trasa. (eg, okta, onelogin)
	ExternalID string `json:"externalID"`
	CreatedAt  int64
	UpdatedAt  int64
}

//UserWithPass is a user struct with password.
type UserWithPass struct {
	User
	OrgName  string `json:"orgName"` //needed for org select
	Password string `json:"password"`
}

//CopyUserWithoutPass converts UserWithPass struct to User
func CopyUserWithoutPass(user UserWithPass) User {
	return User{
		ID:         user.ID,
		OrgID:      user.OrgID,
		UserName:   user.UserName,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		Email:      user.Email,
		Groups:     user.Groups,
		UserRole:   user.UserRole,
		Status:     user.Status,
		IdpName:    user.IdpName,
		ExternalID: user.ExternalID,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

// IdentityProvider holds details for OpenID connect Identity Provider
// CONSTRAINT unique_appproxy UNIQUE(org_id,service_id)
type IdentityProvider struct {
	IdpID   string `json:"idpID"`
	OrgID   string `json:"orgID"`
	IdpName string `json:"idpName"`
	// IdpType can be saml2 or openID or ldap provider
	IdpType string `json:"idpType"`
	// IDP meta can be saml2 xml metadata for saml or base for ldap
	IdpMeta   string `json:"idpMeta"`
	IsEnabled bool   `json:"isEnabled"`
	// Client ID and secret can be openid(oauth) credentials or ldap service account credentials
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	//AudienceURI for saml2 or user search base for ldap
	AudienceURI string `json:"audienceURI"`
	// RedirectURL is where idp would return code or callback
	RedirectURL string `json:"redirectURL"`
	// Endpoint can be openid endpoint or saml embed link
	Endpoint        string `json:"endpoint"`
	IntegrationType string `json:"string"`
	SCIMEndpoint    string `json:"scimEndpoint"`
	ApiKey          string `json:"apiKey"`
	// CreatedBy holds administrator user id
	CreatedBy   string `json:"createdBy"`
	LastUpdated int64  `json:"lastUpdated"`
}

type SAML struct {
	IdpName     string `json:"idpName"`
	IdpMeta     string `json:"idpMeta"`
	EmbedLink   string `json:"embedLink"`
	CallbackURL string `json:"callbackURL"`
}

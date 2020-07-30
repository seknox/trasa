package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

// Org stores behavious related to every Tenants
type Org struct {
	ID             string `json:"ID"`
	OrgName        string `json:"orgName"`
	Domain         string `json:"domain"`
	PrimaryContact string `json:"primaryContact"`
	Timezone       string `json:"timezone"`
	PhoneNumber    string `json:"phoneNumber"`
	CreatedAt      int64
	PlatformBase   string  `json:"platformBase"`
	License        License `json:"license"`
}

type UserContext struct {
	User      *User
	Org       Org
	DeviceID  string
	BrowserID string
}

type TrasaFeatures struct {
	Vault          bool `json:"vault"`
	DynamicService bool `json:"dynamicAuthApp"`
	OrgSignupCount int  `json:"orgSignupCount"`
	AllowRDP       bool `json:"allowRDP"`
}

type License struct {
	Features      TrasaFeatures `json:"features"`
	Expires       int64         `json:"expires"`
	MachineID     string        `json:"machineID"`
	AdminLimit    int           `json:"adminLimit"`
	ProxyAppLimit int           `json:"proxyAppLimit"`
	NodeLimit     int           `json:"nodeLimit"`
	UserLimit     int           `json:"userLimit"`
}

func (a License) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *License) Scan(value interface{}) error {
	if value == nil {
		*a = License{}
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

package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DayAndTimePolicy struct {
	Days     []string `json:"days"`
	FromTime string   `json:"fromTime"`
	ToTime   string   `json:"toTime"`
}

func (d DayAndTimePolicy) Value() (driver.Value, error) {
	data, err := json.Marshal(d)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	logrus.Debug(string(data))
	return string(data), nil
}
func (d DayAndTimePolicy) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &d)
}

// type DevicePolicy struct {
// 	Workstation  string `json:"workstation"`
// 	MobileDevice string `json:"mobileDevice"`
// }

func (d DevicePolicy) Value() (driver.Value, error) {
	devicePerm, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	logrus.Debug(string(devicePerm))
	return string(devicePerm), nil
}
func (d *DevicePolicy) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	//dArr := make([]DevicePolicy, 0)
	return json.Unmarshal(b, &d)
}

type Policy struct {
	PolicyID         string             `json:"policyID" db:"policy_id"`
	OrgID            string             `json:"orgID" db:"org_id"`
	PolicyName       string             `json:"policyName" db:"policy_name"`
	DayAndTime       []DayAndTimePolicy `json:"dayAndTime" db:"day_time"`
	TfaRequired      bool               `json:"tfaRequired" db:"tfa_enabled"`
	RecordSession    bool               `json:"recordSession" db:"record_session"`
	FileTransfer     bool               `json:"fileTransfer" db:"file_transfer"`
	IPSource         string             `json:"ipSource" db:"ip_source"`
	AllowedCountries string             `json:"allowed_countries" db:"allowed_countries"`
	DevicePolicy     DevicePolicy       `json:"devicePolicy" db:"device_policy"`
	RiskThreshold    float32            `json:"riskThreshold" db:"risk_threshold"`
	CreatedAt        int64              `json:"createdAt" db:"created_at"`
	UpdatedAt        int64              `json:"updatedAt" db:"updated_at"`
	Expiry           string             `json:"expiry" db:"expiry"`
	IsExpired        bool               `json:"isExpired"`
	UsedBy           int                `json:"usedBy"`
}

func (d Policy) Value() (driver.Value, error) {
	devicePerm, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return string(devicePerm), nil
}
func (d *Policy) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	//dArr := make([]DevicePolicy, 0)
	return json.Unmarshal(b, &d)
}

package accesscontrol

import (
	"testing"
	"time"

	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
)

var (
	fullDayTime = []models.DayAndTimePolicy{{
		Days:     []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
		FromTime: "01:00",
		ToTime:   "23:59",
	}}
)

func TestCheckTrasaUAC(t *testing.T) {

	type args struct {
		timezone string
		clientip string
		policy   *models.Policy
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 consts.FailedReason
	}{
		{
			name: "full policy",
			args: args{"Asia/Kathmandu", "1.1.1.1", &models.Policy{
				PolicyName: "full",
				DayAndTime: []models.DayAndTimePolicy{{
					Days:     []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
					FromTime: "01:00",
					ToTime:   "23:59",
				}},
				IPSource:         "0.0.0.0/0",
				AllowedCountries: "",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2090-01-01",
			}},
			want:  true,
			want1: "",
		},
		{
			name: "today only",
			args: args{"Asia/Kathmandu", "1.1.1.1", &models.Policy{
				PolicyName: "full",
				DayAndTime: []models.DayAndTimePolicy{{
					Days:     []string{time.Now().Weekday().String()},
					FromTime: "01:00",
					ToTime:   "23:59",
				}},
				IPSource:         "0.0.0.0/0",
				AllowedCountries: "",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2090-01-01",
			}},
			want:  true,
			want1: "",
		},
		{
			name: "no weekdays",
			args: args{"Asia/Kathmandu", "1.1.1.1", &models.Policy{
				PolicyName: "full",
				DayAndTime: []models.DayAndTimePolicy{{
					Days:     []string{},
					FromTime: "01:00",
					ToTime:   "23:59",
				}},
				IPSource:         "0.0.0.0/0",
				AllowedCountries: "",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2090-01-01",
			}},
			want:  false,
			want1: consts.REASON_TIME_POLICY_FAILED,
		},

		{name: "time policy should fail",
			args: args{"Asia/Kathmandu", "1.1.1.1", &models.Policy{
				PolicyName: "full",
				DayAndTime: []models.DayAndTimePolicy{{
					Days:     []string{},
					FromTime: "01:00",
					ToTime:   "01:00",
				}},
				IPSource:         "0.0.0.0/0",
				AllowedCountries: "",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2090-01-01",
			}},
			want:  false,
			want1: consts.REASON_TIME_POLICY_FAILED,
		},

		{name: "expired",
			args: args{"Asia/Kathmandu", "1.1.1.1", &models.Policy{
				DayAndTime:       fullDayTime,
				IPSource:         "0.0.0.0/0",
				AllowedCountries: "",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2020-01-01",
			}},
			want:  false,
			want1: consts.REASON_POLICY_EXPIRED,
		},
		{name: "empty IP",
			args: args{"Asia/Kathmandu", "", &models.Policy{
				DayAndTime:       fullDayTime,
				IPSource:         "0.0.0.0/0",
				AllowedCountries: "",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2090-01-01",
			}},
			want:  false,
			want1: consts.REASON_IP_POLICY_FAILED,
		},

		{name: " /32 IP policy should pass",
			args: args{"Asia/Kathmandu", "1.1.1.1", &models.Policy{
				DayAndTime:       fullDayTime,
				IPSource:         "1.1.1.1/32",
				AllowedCountries: "",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2090-01-01",
			}},
			want:  true,
			want1: "",
		},
		{name: " /24 IP policy should pass",
			args: args{"Asia/Kathmandu", "192.168.0.100", &models.Policy{
				DayAndTime:       fullDayTime,
				IPSource:         "192.168.0.0/24",
				AllowedCountries: "",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2090-01-01",
			}},
			want:  true,
			want1: "",
		},

		{name: " /16 IP policy should pass",
			args: args{"Asia/Kathmandu", "192.168.0.100", &models.Policy{
				DayAndTime:       fullDayTime,
				IPSource:         "192.168.0.0/16",
				AllowedCountries: "",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2090-01-01",
			}},
			want:  true,
			want1: "",
		},
		{name: " /8 IP policy should pass",
			args: args{"Asia/Kathmandu", "192.168.0.100", &models.Policy{
				DayAndTime:       fullDayTime,
				IPSource:         "192.0.0.0/8",
				AllowedCountries: "",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2090-01-01",
			}},
			want:  true,
			want1: "",
		},

		{name: " /0 IP policy should pass",
			args: args{"Asia/Kathmandu", "192.168.0.100", &models.Policy{
				DayAndTime:       fullDayTime,
				IPSource:         "192.0.0.0/0",
				AllowedCountries: "",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2090-01-01",
			}},
			want:  true,
			want1: "",
		},

		{name: "IP policy exact should pass",
			args: args{"Asia/Kathmandu", "192.168.0.100", &models.Policy{
				DayAndTime:       fullDayTime,
				IPSource:         "192.168.0.100",
				AllowedCountries: "",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2090-01-01",
			}},
			want:  true,
			want1: "",
		},

		{name: "invalid country",
			args: args{"Asia/Kathmandu", "1.1.1.1", &models.Policy{
				DayAndTime:       fullDayTime,
				IPSource:         "0.0.0.0/0",
				AllowedCountries: "abcee3",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2090-01-01",
			}},
			want:  false,
			want1: consts.REASON_COUNTRY_POLICY_FAILED,
		},

		{name: "country policy should fail",
			args: args{"Asia/Kathmandu", "1.1.1.1", &models.Policy{
				DayAndTime:       fullDayTime,
				IPSource:         "0.0.0.0/0",
				AllowedCountries: "NP",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2090-01-01",
			}},
			want:  false,
			want1: consts.REASON_COUNTRY_POLICY_FAILED,
		},

		{name: "country policy for private IP",
			args: args{"Asia/Kathmandu", "127.0.0.1", &models.Policy{
				DayAndTime:       fullDayTime,
				IPSource:         "0.0.0.0/0",
				AllowedCountries: "NP",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2090-01-01",
			}},
			want:  true,
			want1: "",
		},

		{name: "country policy should pass",
			args: args{"Asia/Kathmandu", "54.151.215.187", &models.Policy{
				DayAndTime:       fullDayTime,
				IPSource:         "0.0.0.0/0",
				AllowedCountries: "NP,SG",
				DevicePolicy:     models.DevicePolicy{},
				Expiry:           "2090-01-01",
			}},
			want:  true,
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CheckTrasaUAC(tt.args.timezone, tt.args.clientip, tt.args.policy)
			if got != tt.want {
				t.Errorf("CheckTrasaUAC() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CheckTrasaUAC() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

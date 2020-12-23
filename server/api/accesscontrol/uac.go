package accesscontrol

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/api/devices"
	"github.com/sirupsen/logrus"

	"net"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/now"
	"github.com/phuslu/iploc"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/models"
	"github.com/seknox/trasa/server/utils"
)

/*
cases:
- allow login only on monday
- allow login only between 5pm-6pm
- allow login between 10am-6pm on weekdays but allow login only between 5pm-6pm on saturday

*/

// TrasaUAC or User access control is handled when assigning user to application.
// User access control handles users timing for access, users workstation for access
// and user IP for access. Access can be assigned for limited time per day, limited day per week
// or limited access with expiry date.
// type TrasaUAC struct {
// 	Days        []string
// 	time        []int
// 	IP          string
// 	Workstation string
// 	ExpiryTimer string
// }
func CheckTrasaUAC(timezone, clientip string, policy *models.Policy) (bool, consts.FailedReason) {

	allow := false
	reason := consts.REASON_TIME_POLICY_FAILED
	for _, dayTimePolicy := range policy.DayAndTime {
		allow, reason = dayTimeExpiryChecker(timezone, dayTimePolicy)
		logrus.Tracef("day time check:%v , failed reason:%v", allow, reason)
	}

	// we check expiry
	current := time.Now()
	delta, err := time.Parse("2006-01-02", policy.Expiry)
	if err != nil {
		logrus.Error(err)
		return false, consts.REASON_UNKNOWN
	}
	if current.Before(delta) != true {
		logrus.Tracef("policy expired")
		return false, consts.REASON_POLICY_EXPIRED
	}

	//check ip policy
	allowedIps := strings.Split(policy.IPSource, ",")
	//logrus.Debug(allowedIps, clientip)

	chk, err := utils.NewChecker(allowedIps)
	if err != nil {
		logrus.Error(err)
	}
	err = chk.IsAuthorized(clientip)

	if err != nil {
		return false, consts.REASON_IP_POLICY_FAILED
	}

	//logrus.Debug(policy.AllowedCountries)
	//Country check
	allowedCountries := strings.Split(policy.AllowedCountries, ",")
	//logrus.Debug(allowedCountries)
	//logrus.Debug(geoip.Country(net.ParseIP(clientip))," country")

	//private ip and invalid ip has "ZZ" country
	//We alllow invalid or private ips
	if !utils.ArrayContainsString(allowedCountries, "ZZ") {
		allowedCountries = append(allowedCountries, "ZZ")
	}

	//len(allowedCountries)==1 means blank allowed countries
	if !utils.ArrayContainsString(allowedCountries, string(iploc.Country(net.ParseIP(clientip)))) && policy.AllowedCountries != "" {
		return false, consts.REASON_COUNTRY_POLICY_FAILED
	}
	return allow, reason
}

func dayTimeExpiryChecker(timezone string, perm models.DayAndTimePolicy) (bool, consts.FailedReason) {

	nep, err := time.LoadLocation(timezone)
	if err != nil {
		logrus.Errorf("load location: %v", err)
		nep, err = time.LoadLocation("UTC")
		if err != nil {
			logrus.Errorf("load location: %v", err)
			return false, consts.REASON_UNKNOWN
		}

	}
	local := time.Now().In(nep)

	//current := time.Now()
	hour := local.Hour()

	minute := local.Minute()
	// we check day

	daycheck := dayChecker(perm.Days)
	if daycheck != true {
		logrus.Trace("daycheck failed") // should return here
		return false, consts.REASON_TIME_POLICY_FAILED
	}

	// we check time

	fromstr := strings.Split(perm.FromTime, ":")
	tostr := strings.Split(perm.ToTime, ":")

	from := make([]int, len(fromstr))
	for i := range from {
		from[i], _ = strconv.Atoi(fromstr[i])
	}

	to := make([]int, len(tostr))
	for i := range to {
		to[i], _ = strconv.Atoi(tostr[i])
	}

	check := timeChecker(hour, minute, from, to)
	if check != true {
		logrus.Trace("failed at FromTO")
		return false, consts.REASON_TIME_POLICY_FAILED
	}

	return true, ""
}

func dayChecker(days []string) bool {
	date := now.BeginningOfHour()
	day := date.Weekday()
	for _, v := range days {
		if strings.Compare(v, day.String()) == 0 {
			return true
		}
	}
	return false
}

func timeChecker(hour, minute int, from, to []int) bool {

	if hour > from[0] && hour < to[0] {
		return true
	}

	if (hour >= from[0] && minute >= from[1]) && (hour <= to[0] && minute <= to[1]) {
		return true
	}
	return false
}

// CheckPolicy validates policy for user access
func CheckPolicy(params *models.ConnectionParams, policy *models.Policy, adHocSwitch bool) (bool, consts.FailedReason) {

	ok := false
	reason := consts.REASON_UNKNOWN

	//we check users policy
	ok, reason = CheckTrasaUAC(params.Timezone, params.UserIP, policy)
	if ok == true {
		return true, "user authorised by uac check"
	}

	return ok, reason
}

//CheckDevicePolicy checks if device hygiene of user is according to device policy
func CheckDevicePolicy(policy models.DevicePolicy, accessDeviceID, tfaDeviceID, orgID string) (consts.FailedReason, bool, error) {

	//Skipping device policy for now
	//return "", true, nil

	sett, err := system.Store.GetGlobalSetting(orgID, consts.GLOBAL_DEVICE_HYGIENE_CHECK)
	if err != nil {
		logrus.Error(err)
		return consts.REASON_UNKNOWN, false, errors.Errorf("could not get global settings: %v", err)
	}

	if !sett.Status {
		return "", true, nil
	}

	accessDevice, err := devices.Store.GetFromID(accessDeviceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return consts.REASON_DEVICE_NOT_ENROLLED, false, errors.Errorf("could not get access device detail of %s: %v", accessDeviceID, err)
		} else {
			return "Something is wrong", false, errors.Errorf("could not get access device detail: %v", err)
		}

	}

	accessDeviceHygiene := accessDevice.DeviceHygiene

	tfaDevice, err := devices.Store.GetFromID(tfaDeviceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return (consts.REASON_DEVICE_NOT_ENROLLED), false, err
		} else {
			return consts.REASON_UNKNOWN, false, errors.Errorf("could not get access device detail: %v", err)
		}
	}
	tfaDeviceHygiene := tfaDevice.DeviceHygiene

	//
	//pol,_:=json.Marshal(policy)
	//logrus.Trace("policy ",string(pol))
	//
	//pol,_=json.Marshal(accessDeviceHygiene)
	//logrus.Trace("acc ",string(pol))
	//
	//pol,_=json.Marshal(tfaDeviceHygiene)
	//logrus.Trace("tfa ",string(pol))
	//

	if policy.BlockUntrustedDevices {
		if !accessDevice.Trusted {
			logrus.Trace("!Trusted")
			return "Device policy failed: device not trusted by admin", false, nil
		}
		if !tfaDevice.Trusted {
			logrus.Trace("!Trusted")
			return "Device policy failed: device not trusted by admin", false, nil
		}
	}

	if policy.BlockRemoteLoginEnabled {
		if accessDeviceHygiene.LoginSecurity.RemoteLoginEnabled {
			logrus.Trace("RemoteLoginEnabled")
			return "Device policy failed: remote login enabled", false, nil
		}
	}

	if policy.BlockEncryptionNotSet {
		if !accessDeviceHygiene.EndpointSecurity.DeviceEncryptionEnabled { //|| !tfaDeviceHygiene.EndpointSecurity.DeviceEncryptionEnabled {
			logrus.Trace("!DeviceEncryptionEnabled")
			return "Device policy failed: device encryption not enabled", false, nil
		}
	}

	if policy.BlockAntivirusDisabled {
		if !accessDeviceHygiene.EndpointSecurity.EpsConfigured {
			logrus.Trace("!EpsConfigured")
			return "Device policy failed: endpoint security/antivirus not enabled", false, nil
		}
	}

	if policy.BlockFirewallDisabled {
		logrus.Trace("BlockFirewallDisabled")
		if !accessDeviceHygiene.EndpointSecurity.FirewallEnabled {
			logrus.Trace("!FirewallEnabled")
			return "Device policy failed: firewall not enabled", false, nil
		}
	}

	if policy.BlockCriticalAutoUpdateDisabled {
		logrus.Trace("BlockCriticalAutoUpdateDisabled")
		if !accessDeviceHygiene.DeviceOS.AutoUpdate {
			logrus.Trace("!AutoUpdate")
			return "Device policy failed: autoUpdate disabled", false, nil
		}
	}

	if policy.BlockIdleScreenLockDisabled {
		if !accessDeviceHygiene.LoginSecurity.IdleDeviceScreenLock {
			return "Device policy failed: screen lock disabled", false, nil
		}
	}

	if policy.BlockAutologinEnabled {
		if accessDeviceHygiene.LoginSecurity.AutologinEnabled || tfaDeviceHygiene.LoginSecurity.AutologinEnabled {
			logrus.Trace("!AutologinEnabled")
			return "Device policy failed: autologin enabled", false, nil
		}
	}

	if policy.BlockEmulated {
		if tfaDeviceHygiene.DeviceOS.IsEmulator {
			logrus.Trace("IsEmulator")
			return "Device policy failed: tfa device is a emulator", false, nil
		}
	}

	if policy.BlockJailBroken {
		if tfaDeviceHygiene.DeviceOS.JailBroken {
			logrus.Trace("JailBroken")
			return "Device policy failed: tfa device is jail broken", false, nil
		}
	}

	return "", true, nil

	if policy.BlockDebuggingEnabled {
		if accessDeviceHygiene.DeviceOS.DebugModeEnabled || tfaDeviceHygiene.DeviceOS.DebugModeEnabled {
			logrus.Trace("DebugModeEnabled")
			return "Device policy failed: debug mode enabled", false, nil
		}
	}

	if policy.BlockOpenWifiConn {
		if accessDeviceHygiene.NetworkInfo.OpenWifiConn || tfaDeviceHygiene.NetworkInfo.OpenWifiConn {
			logrus.Trace("OpenWifiConn")
			return "Device policy failed: connected to open WiFi", false, nil
		}
	}

	if policy.BlockTfaNotConfigured {
		if !accessDeviceHygiene.LoginSecurity.TfaConfigured {
			logrus.Trace("!TfaConfigured")
			return "Device policy failed: second factor authentication not configured in device", false, nil
		}

	}

	return "", true, nil
}

func checkVersion(policy string, ver string) (bool, error) {
	splittedPolicy := strings.Split(policy, ".")
	splittedVersion := strings.Split(ver, ".")

	for i, v := range splittedPolicy {
		polInt, err := strconv.Atoi(v)
		if err != nil {
			return false, errors.New("invalid policy: not a number")
		}
		if len(splittedVersion) < i+1 {
			return false, errors.Errorf("length mismatch: %d %d", len(splittedVersion), i+1)
		}

		verInt, err := strconv.Atoi(splittedVersion[i])
		if err != nil {
			return false, errors.New("invalid version: not a number")
		}
		if polInt > verInt {
			return false, nil
		} else if polInt < verInt {
			return true, nil
		}

	}

	return true, nil
}

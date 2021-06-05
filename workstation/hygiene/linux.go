package hygiene

/**
 * Copyright (C) 2020 Seknox Pte Ltd.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/shirou/gopsutil/net"
)

type DeviceLinux struct {
}

//
//ubuntu P 01/08/2020 0 99999 7 -1 can be used to retreive user password state.

// IsAutoLoginEnabled checks if auto logon is enabled
//  TODO check in centos and other linux distros
func (d DeviceLinux) IsAutoLoginEnabled() (bool, error) {
	n, _, _, _ := d.GetOSNameVersion()

	if normalizeString(n) == "ubuntu" {
		out, err := exec.Command("/bin/sh", "-c", `sed -n '/^\s*AutomaticLoginEnable\s*=/s///p' /etc/gdm3/custom.conf`).Output()
		if err != nil {
			return false, err
		}

		if normalizeString(string(out)) == "true" {
			return true, nil
		}
		return false, fmt.Errorf("Not enabled")

	} else if normalizeString(n) == "centos" {
		out, err := exec.Command("/bin/sh", "-c", `sed -n '/^\s*AutomaticLoginEnable\s*=/s///p' /etc/gdm/custom.conf`).Output()
		if err != nil {
			return false, err
		}

		if normalizeString(string(out)) == "true" {
			return true, nil
		}
		return false, fmt.Errorf("Not enabled")
	}

	return false, fmt.Errorf("Not enabled")

}

func (d DeviceLinux) IsFireWallSet() (bool, error) {
	cmd := exec.Command("lsmod")
	out, err := cmd.CombinedOutput()

	if err != nil {
		return false, fmt.Errorf(`%v : %s`, err, string(out))
	}

	matched, err := regexp.Match(`ip6*t_REJECT\s+\d+\s+(\d)`, out)
	if err != nil {
		return false, errors.New("Invalid output")
	}

	return matched, nil
}

func (d DeviceLinux) IsDeviceEncrypted() (bool, error) {
	checkCryptSetupStatusCmd := exec.Command("lsblk", "-f")
	out, err := checkCryptSetupStatusCmd.CombinedOutput()
	if err != nil && len(out) == 0 {
		return false, fmt.Errorf(`%v : %s`, err, string(out))
	}
	if strings.Contains(string(out), "crypt") {
		return true, nil
	}

	return false, nil
	//return false,errors.New("Not implemented")
}

func (d DeviceLinux) GetInstalledPackages() ([]string, error) {
	pkgListCmd := exec.Command("apt", "list", "--installed")
	out, err := pkgListCmd.CombinedOutput()
	if err != nil {
		return []string{}, fmt.Errorf(`%v : %s`, err, string(out))
	}
	pkgList := strings.Split(string(out), "\n")

	return pkgList, nil
}

// func (d DeviceLinux) GetDeviceName() (string, error) {
// 	osNameCmd := exec.Command(`hostname`)
// 	out, err := osNameCmd.CombinedOutput()
// 	if err != nil {
// 		return "", fmt.Errorf(`%v : %s`, err, string(out))
// 	}

// 	return strings.TrimSpace(string(out)), nil

// }

func (d DeviceLinux) GetOSNameVersion() (osName, osVersion, kernelVersion string, err error) {
	logrus.Info("___V@___")
	c := exec.Command(`lsb_release`, `-ir`)
	out, err := c.CombinedOutput()
	if err != nil {
		err = fmt.Errorf(`%v : %s`, err, string(out))
		return
	}

	lines := strings.Split(string(NormalizeNewlines(out)), "\n")
	for _, line := range lines {
		fmt.Println("->", line)
		if strings.Contains(line, "Distributor ID") {
			splitted := strings.Split(line, "Distributor ID:")
			if len(splitted) > 1 {
				osName = strings.TrimSpace(splitted[1])
			}
		}
		if strings.Contains(line, "Release") {
			splitted := strings.Split(line, "Release:")
			if len(splitted) > 1 {
				osVersion = strings.TrimSpace(splitted[1])
			}
		}

	}

	c = exec.Command("uname", "-r")
	out, err = c.CombinedOutput()
	if err != nil {
		err = fmt.Errorf(`%v : %s`, err, string(out))
		return
	}

	kernelVersion = string(out)

	return

}

func (DeviceLinux) GetPasswordLastUpdated() (string, error) {
	passChCmd := exec.Command(`/bin/sh`, `-c`, `chage -l $USER | grep Last`)
	out, err := passChCmd.CombinedOutput()
	if err != nil {
		return "", errors.New("Could not get last password change")
	}
	splitted := strings.Split(string(out), ":")
	if len(splitted) < 2 {
		return "", errors.New("Could not get last password change")
	}
	return strings.TrimSpace(splitted[1]), nil

}

func (d DeviceLinux) GetCriticalAutoUpdateStatus() (bool, error) {
	updateStatCmd := exec.Command(`/bin/sh`, `-c`, `cat /etc/apt/apt.conf.d/20auto-upgrades`)

	out, err := updateStatCmd.CombinedOutput()
	if err != nil {
		return false, errors.New("Could not get update status")
	}
	if strings.Contains(string(out), `APT::Periodic::Update-Package-Lists "1";`) && strings.Contains(string(out), `APT::Periodic::Unattended-Upgrade "1";`) {
		return true, nil
	} else if strings.Contains(string(out), "0") {
		return false, nil
	} else {
		return false, errors.New("Could not get update status")
	}

}

func (d DeviceLinux) GetPendingUpdates() ([]string, error) {
	updatesCmd := exec.Command(`/bin/sh`, `-c`, `apt list --upgradable | grep "\-security"`)
	out, err := updatesCmd.CombinedOutput()
	if err != nil {
		return []string{}, errors.New("Could not get pending updates")
	}

	updates := []string{}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "security") {
			updates = append(updates, line)
		}
	}

	return updates, nil

}

func (d DeviceLinux) IsRemoteConnectionEnabled() (bool, error) {
	c := exec.Command(`ss`, `-ltn`)
	out, err := c.CombinedOutput()
	if err != nil {
		return true, errors.New("Could not get remote connection status")
	}

	sshEnabled := false
	lines := strings.Split(string(out), "\n")

	for _, line := range lines {
		if strings.Contains(line, ":22") {
			sshEnabled = true
		}
	}

	return sshEnabled, nil
}

func (d DeviceLinux) ScreenLockEnabled() (bool, error) {
	//gsettings list-recursively org.gnome.desktop.screensaver
	c := exec.Command(`gsettings`, `list-recursively`, `org.gnome.desktop.screensaver`)
	out, err := c.CombinedOutput()
	if err != nil {
		return true, errors.New("Could not get screen lock status")
	}

	screenLockStatus := false
	lines := strings.Split(string(out), "\n")

	for _, line := range lines {
		if strings.Contains(line, "lock-enabled") && strings.Contains(line, "true") {
			screenLockStatus = true
		}
	}

	return screenLockStatus, nil
}

func (d DeviceLinux) GetNetwork() (string, string, string, error) {
	cmd := exec.Command("bash", "-c", "route | grep '^default' | grep -o '[^ ]*$'")
	defaultInterface, err := cmd.Output()
	if err != nil {
		return "", "", "", fmt.Errorf("failed finding default route %v", err)
	}

	// we retreive ip address and mac address for this route.
	var netVal net.InterfaceStat

	nets, err := net.Interfaces()
	if err != nil {
		return "", "", string(defaultInterface), err
	}
	for _, v := range nets {
		//fmt.Println(v.Name)
		//fmt.Println(v.String())
		if strings.TrimSpace(v.Name) == strings.TrimSpace(string(defaultInterface)) {
			netVal.Addrs = v.Addrs
			netVal.HardwareAddr = v.HardwareAddr
		}
	}

	addr := ""
	for _, n := range netVal.Addrs {
		addr = addr + "," + n.Addr
	}

	addr = strings.TrimLeft(addr, ",")

	return addr, netVal.HardwareAddr, string(defaultInterface), nil
}

func (d DeviceLinux) GetLatestSecurityPatch() (string, error) {
	return "", errors.New("Not supported yet")
}

func (d DeviceLinux) IdleDeviceScreenLockTime() (string, error) {
	return "", errors.New("Not supported yet")
}

func (d DeviceLinux) EndpointSecurity() (string, string, bool, error) {
	return "", "", true, errors.New("Not supported yet")
}

// normalizeString trims spaces and convert into lowercase
func normalizeString(s string) string {
	s = strings.TrimSpace(s)

	s = strings.ToLower(s)
	return s
}

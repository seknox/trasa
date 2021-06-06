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
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/net"

	logger "github.com/sirupsen/logrus"
)

type DeviceMac struct {
}

func (d DeviceMac) IsAutoLoginEnabled() (bool, error) {
	dsclCmd := exec.Command("dscl", ".", "list", "/Users")
	grepCmd := exec.Command("grep", "-v", "^_")

	inPipe, outPipe := io.Pipe()

	var buff bytes.Buffer
	grepCmd.Stdin = inPipe
	dsclCmd.Stdout = outPipe

	grepCmd.Stdout = &buff
	grepCmd.Stderr = os.Stderr
	dsclCmd.Start()
	grepCmd.Start()

	//grepCmd.Start()
	//grepPipe,err:=grepCmd.StdinPipe()
	//	var buff =make([]byte,100)
	//	go io.CopyBuffer(grepPipe,dsclPipe,buff)

	out := buff.String()

	//logger.Debug(out)
	//if err != nil {
	//	logger.Debug(string(out))
	//	return false,err
	//}
	//logger.Debug(out)

	outputStr := string(out)
	users := strings.Split(outputStr, "\n")
	for _, user := range users {
		dsclCmd = exec.Command("dscl", ".", "-authonly", user, "")
		err := dsclCmd.Run()
		if err == nil {
			logger.Debug("Password not set for user: ", user)
			return true, nil
		}

	}
	return false, nil

}

func (d DeviceMac) IsFireWallSet() (bool, error) {
	cmd := exec.Command("defaults", "read", "/Library/Preferences/com.apple.alf", "globalstate")
	out, err := cmd.Output()

	if err != nil {
		return false, err
	}
	status := strings.TrimSpace(string(out))
	if status == "1" {
		return true, nil
	} else if status == "0" {
		return false, nil
	}

	return false, errors.New("Could not get firewall status")

}

func (d DeviceMac) IsDeviceEncrypted() (bool, error) {
	cmd := exec.Command("fdesetup", "status")
	out, err := cmd.Output()

	if err != nil {
		return false, err
	}
	status := strings.TrimSpace(string(out))
	if strings.Contains(status, "On") {
		return true, nil
	} else if strings.Contains(status, "Off") {
		return false, nil
	}

	return false, errors.New("Could not get encryption status")
}

func (d DeviceMac) GetInstalledPackages() ([]string, error) {

	pkgListCmd := exec.Command("pkgutil", "--pkgs")
	out, err := pkgListCmd.Output()
	if err != nil {
		return []string{}, err
	}
	pkgList := strings.Split(string(out), "\n")

	return pkgList, nil
}

// func (d DeviceMac) GetDeviceName() (string, error) {
// 	osNameCmd := exec.Command(`hostname`)
// 	out, err := osNameCmd.CombinedOutput()
// 	if err != nil {
// 		return "", fmt.Errorf(`%v : %s`, err, string(out))
// 	}

// 	return strings.TrimSpace(string(out)), nil

// }

func (d DeviceMac) GetOSNameVersion() (osName, osVersion, kernelVersion string, err error) {
	osNameCmd := exec.Command(`sw_vers`)
	out, err := osNameCmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf(`%v : %s`, err, string(out))
		return
	}
	//osName := ""
	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		err = errors.New("Could not get OS name")
		return
	}
	productName := strings.Split(string(lines[0]), "ProductName:")
	version := strings.Split(string(lines[1]), "ProductVersion:")

	//fmt.Println(productName)
	//fmt.Println(version)

	if (len(productName) < 1) || len(version) < 1 {
		err = errors.New("Could not get OS name")
		return
	}
	//osName = fmt.Sprintf(`%s : %s`, strings.TrimSpace(productName[1]), strings.TrimSpace(version[1]))
	osName = strings.TrimSpace(productName[1])
	osVersion = strings.TrimSpace(version[1])

	c := exec.Command("uname", "-r")
	out, err = c.CombinedOutput()
	if err != nil {
		return
	}

	kernelVersion = strings.TrimSpace(string(out))

	return
}
func (d DeviceMac) GetPasswordLastUpdated() (string, error) {
	passCmd := exec.Command(`/bin/sh`, `-c`, `dscl . read /Users/$USER |  grep -A1 passwordLastSetTime | grep real`)
	out, err := passCmd.CombinedOutput()
	if err != nil {
		return "", errors.New("Could not get last password change")
	}

	splitted := strings.Split(string(out), "real")
	if len(splitted) > 2 {
		timeStr := strings.ReplaceAll(splitted[1], "</", "")
		timeStr = strings.ReplaceAll(timeStr, ">", "")
		//fmt.Println(timeStr)
		timeStrSpl := strings.Split(timeStr, ".")
		if len(timeStrSpl) != 2 {
			return "", errors.New("Could not get last password change")
		}
		m, s := 0, 0
		s, _ = strconv.Atoi(timeStrSpl[0])
		m, _ = strconv.Atoi(timeStrSpl[0])

		lastChange := time.Unix(int64(s), int64(m))
		return lastChange.Format(time.UnixDate), nil
	}

	return "", errors.New("Could not get last password change")
}

func (d DeviceMac) GetCriticalAutoUpdateStatus() (bool, error) {
	updateStatCmd := exec.Command(`/bin/sh`, `-c`, `defaults read /Library/Preferences/com.apple.SoftwareUpdate.plist | grep CriticalUpdateInstall`)
	out, err := updateStatCmd.CombinedOutput()
	if err != nil {
		return false, errors.New("Could not get update status")
	}
	if strings.Contains(string(out), "1") {
		return true, nil
	} else if strings.Contains(string(out), "0") {
		return false, nil
	} else {
		return false, errors.New("Could not get update status")
	}

}

func (d DeviceMac) GetPendingUpdates() ([]string, error) {
	updatesCmd := exec.Command(`softwareupdate`, `-l`, "--no-scan")
	out, err := updatesCmd.CombinedOutput()
	if err != nil {
		return []string{}, errors.New("Could not get pending updates")
	}

	updates := []string{}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Recommended") {
			updates = append(updates, line)
		}
	}

	return updates, nil

}

func (d DeviceMac) IsRemoteConnectionEnabled() (bool, error) {
	out, err := Execute(`netstat`, `-anv`)
	if err != nil {
		return false, err
	}
	extracted := Extract(`\*.(22|23)\s`, out)

	if len(extracted) > 0 {
		return true, nil
	} else {
		return false, nil
	}

}

func (d DeviceMac) ScreenLockEnabled() (bool, error) {
	return true, errors.New("Not supported yet")
}

func (d DeviceMac) GetNetwork() (string, string, string, error) {
	cmd := exec.Command("zsh", "-c", "route | grep '^default' | grep -o '[^ ]*$'")
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

func (d DeviceMac) GetLatestSecurityPatch() (string, error) {
	return "", errors.New("Not supported yet")
}

func (d DeviceMac) IdleDeviceScreenLockTime() (string, error) {
	return "", errors.New("Not supported yet")
}

func (d DeviceMac) EndpointSecurity() (string, string, bool, error) {
	return "", "", true, errors.New("Not supported yet")
}

func GetProduct() (name, version string, err error) {
	c := exec.Command("sw_vers")
	var out []byte
	out, err = c.CombinedOutput()
	if err != nil {
		return
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		splitted := strings.Split(line, ":")
		if len(splitted) == 2 && strings.TrimSpace(splitted[0]) == "ProductName" {
			name = strings.TrimSpace(splitted[1])
		} else if len(splitted) == 2 && strings.TrimSpace(splitted[0]) == "ProductVersion" {
			version = strings.TrimSpace(splitted[1])
		}
	}
	return
}

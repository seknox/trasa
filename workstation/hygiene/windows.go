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
	"net"
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	logger "github.com/sirupsen/logrus"
	// "golang.org/x/sys/windows/registry"
)

type DeviceWindows struct {
}

func (d DeviceWindows) IsAutoLoginEnabled() (bool, error) {

	//// check adminAutoLogin from registry.
	//// we will only return if value is true. else continue for net user command response.
	// k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon`, registry.QUERY_VALUE)
	// if err != nil {
	// 	logger.Trace(err)
	// 	// return false, err
	// }
	// defer k.Close()

	// s, _, err := k.GetStringValue("AutoAdminLogon")
	// if err != nil {
	// 	logger.Trace(err)
	// 	// return false, err
	// }

	// // value 1 is auto login enabled. 0 for disabled
	// if s == "1" {
	// 	return true, nil
	// }

	// If we reach here, it means either adminAutoLogin is not set or is set 0 value.
	// Check regular windows password not required via net user command

	userc, err := user.Current()
	if err != nil {
		return false, err
	}

	t := strings.Split(userc.Username, `\`)
	username := ""
	if len(t) == 2 {
		username = t[1]
	}

	c := exec.Command("net", "config", "workstation")
	out, err := c.CombinedOutput()
	if err != nil {
		logger.Errorf("net user workstation ", err)
		return false, err
	}

	lines := strings.Split(string(NormalizeNewlines(out)), "\n")

	domain := ""
	for _, line := range lines {
		// "No" means password is not required hence auto login is enabled.
		if strings.Contains(line, "Workstation domain") {
			splitted := strings.Split(line, "Workstation domain")
			if len(splitted) < 2 {
				return false, errors.New("could not get domain")
			}
			domain = strings.TrimSpace(splitted[1])
			//return true, nil
		}
	}

	c = exec.Command(
		"C:/WINDOWS/system32/cscript.exe",
		`C:\\Program Files\Fireser\checkpass.vbs`,
		fmt.Sprintf(`/user:%s`, username),
		fmt.Sprintf(`/workgroup:%s`, domain),
	)
	out, err = c.CombinedOutput()
	if err != nil {
		// return false in case of error
		logger.Error(err)
		return false, err
	}

	if strings.Contains(string(out), "YES") {
		return true, nil
	} else if strings.Contains(string(out), "NO") {
		return false, nil
	}

	return false, errors.New("could not get password state")

}

//needs admin privilege
// func (d DeviceWindows) IsDeviceEncrypted() (bool, error) {
// 	// Fully Decrypted

// 	c := exec.Command(`manage-bde`, `-status`)
// 	out, err := c.CombinedOutput()
// 	if err != nil {
// 		return false, err
// 	}
// 	str := string(NormalizeNewlines(out))
// 	if strings.Contains(str, "Fully Decrypted") {
// 		return false, nil
// 	} else {
// 		return true, nil
// 	}

// 	return false, errors.New("Not Supported Yet")
// }

func (d DeviceWindows) IsDeviceEncrypted() (bool, error) {
	// Fully Decrypted

	c := exec.Command(`C:\\Program Files\Fireser\bitlocker-status.exe`)
	out, err := c.Output()
	if err != nil {
		return false, err
	}
	str := string(NormalizeNewlines(out))
	if strings.Contains(str, "OFF") {
		return false, nil
	} else if strings.Contains(str, "ON") {
		return true, nil
	}

	return false, errors.New("Invalid output from bitlocker-status.exe")
}

func (d DeviceWindows) GetInstalledPackages() ([]string, error) {
	return nil, errors.New("Not Supported Yet")
}

//slower
func getOSNameVersionFromSysinfo() (string, string, string, error) {
	// fmt.Println("sysinfo")
	c := exec.Command(`systeminfo`)
	out, err := c.CombinedOutput()
	if err != nil {
		return "", "", "", err
	}
	osName := ""
	osVersion := ""
	lines := strings.Split(string(NormalizeNewlines(out)), "\n")
	for _, line := range lines {
		if strings.Contains(line, "OS Name") {
			temp := strings.Split(line, ":")
			if len(temp) == 2 {
				osName = strings.TrimSpace(temp[1])
			}
		} else if strings.Contains(line, "OS Version") && !strings.Contains(line, "BIOS Version") {
			temp := strings.Split(line, ":")
			if len(temp) == 2 {
				osVersion = strings.TrimSpace(temp[1])
			}
		}

	}

	return osName, osVersion, osVersion, nil
}

func getOSNameVersionFromPS() (string, string, string, error) {
	c := exec.Command(`powershell`, `-Command`, `(Get-WMIObject win32_operatingsystem).name`)
	out, err := c.CombinedOutput()
	if err != nil {
		return getOSNameVersionFromSysinfo()
	}
	spll := strings.Split(string(out), "|")
	if len(spll) < 1 {
		return getOSNameVersionFromSysinfo()
	}

	osName := spll[0]

	c = exec.Command(`powershell`, `-Command`, `(Get-WMIObject win32_operatingsystem).version`)
	out, err = c.CombinedOutput()
	if err != nil {
		return getOSNameVersionFromSysinfo()
	}
	osVersion := strings.TrimSpace(string(out))

	return osName, osVersion, osVersion, nil
}

func (d DeviceWindows) GetOSNameVersion() (string, string, string, error) {
	c := exec.Command("cmd.exe", `ver`)
	out, err := c.CombinedOutput()
	if err != nil {
		return getOSNameVersionFromPS()
	}

	//    Microsoft Windows [Version 10.0.15063]
	spll := strings.Split(string(out), "[")
	if len(spll) < 2 {
		return getOSNameVersionFromPS()
	}

	osName := strings.TrimSpace(spll[0])

	//    [Version 10.0.15063]
	splitted := strings.Split(spll[1], " ")
	if len(splitted) < 2 {
		return getOSNameVersionFromPS()
	}

	//   10.0.15063]
	splitted2 := strings.Split(splitted[1], "]")
	if len(splitted2) < 2 {
		return getOSNameVersionFromPS()
	}

	osVersion := strings.TrimSpace(splitted2[0])

	return osName, osVersion, osVersion, nil
}

func (d DeviceWindows) IsRemoteConnectionEnabled() (bool, error) {
	var isRDP = true
	var isRPC = true
	c1 := exec.Command(`reg`, `query`, `HKLM\System\CurrentControlSet\Control\Terminal Server`)
	out, err := c1.CombinedOutput()
	if err != nil {
		return true, err
	}
	lines := strings.Split(string(NormalizeNewlines(out)), "\n")
	for _, line := range lines {
		if strings.Contains(line, "AllowRemoteRPC") {
			if strings.Contains(line, "0x0") {
				isRPC = false
			}
		}
		if strings.Contains(line, "fDenyTSConnections") {
			if strings.Contains(line, "0x1") {
				isRDP = false
			}
		}

	}

	return isRDP || isRPC, nil

}

func (d DeviceWindows) EndpointSecurity() (string, string, bool, error) {

	//
	//Sample output
	// Antivirusenabled              : True
	// AMServiceEnabled              : True
	// AntispywareEnabled            : True
	// BehaviorMonitorEnabled        : True
	// IoavProtectionEnabled         : True
	// NISEnabled                    : True
	// OnAccessProtectionEnabled     : True
	// RealTimeProtectionEnabled     : True
	// AntivirusSignatureLastUpdated : 6/7/2020 6:01:31 AM

	c := exec.Command("powershell", "-command", `Get-MpComputerStatus | Select-Object -Property Antivirusenabled,AMServiceEnabled,AntispywareEnabled,BehaviorMonitorEnabled,IoavProtectionEnabled,PSComputerName,NISEnabled,OnAccessProtectionEnabled,RealTimeProtectionEnabled,AntivirusSignatureLastUpdated,NISSignatureLastUpdated,NISSignatureVersion,NISSignatureAge,NISEngineVersion,ComputerState,ComputerID,AntivirusSignatureVersion,AntispywareSignatureVersion,AntispywareSignatureLastUpdated,AntispywareSignatureAge,AMServiceVersion,AMProductVersion,AMEngineVersion | ConvertTo-Json`)
	out, err := c.CombinedOutput()
	if err != nil {
		return "", "", false, err
	}
	epsMeta := string(out)
	Antivirusenabled := false
	AMServiceEnabled := false
	AntispywareEnabled := false
	BehaviorMonitorEnabled := false
	IoavProtectionEnabled := false
	NISEnabled := false
	OnAccessProtectionEnabled := false
	RealTimeProtectionEnabled := false
	//AntivirusSignatureLastUpdated:="false"

	lines := strings.Split(string(NormalizeNewlines(out)), "\n")
	for _, line := range lines {
		if ContainsIgnoreCase(line, "true") {
			if ContainsIgnoreCase(line, "Antivirusenabled") {
				Antivirusenabled = true
			} else if ContainsIgnoreCase(line, "AMServiceEnabled") {
				AMServiceEnabled = true
			} else if ContainsIgnoreCase(line, "AntispywareEnabled") {
				AntispywareEnabled = true
			} else if ContainsIgnoreCase(line, "BehaviorMonitorEnabled") {
				BehaviorMonitorEnabled = true
			} else if ContainsIgnoreCase(line, "IoavProtectionEnabled") {
				IoavProtectionEnabled = true
			} else if ContainsIgnoreCase(line, "NISEnabled") {
				NISEnabled = true
			} else if ContainsIgnoreCase(line, "OnAccessProtectionEnabled") {
				OnAccessProtectionEnabled = true
			} else if ContainsIgnoreCase(line, "RealTimeProtectionEnabled") {
				RealTimeProtectionEnabled = true
			}
		}

	}

	epsEnabled := Antivirusenabled && AMServiceEnabled && AntispywareEnabled && BehaviorMonitorEnabled && IoavProtectionEnabled && NISEnabled && OnAccessProtectionEnabled && RealTimeProtectionEnabled

	c = exec.Command(`WMIC`, `/Node:localhost`, `/Namespace:\\root\SecurityCenter2`, `Path`, `AntiVirusProduct`, `Get`, `displayName`, "/Format:List")
	out, err = c.CombinedOutput()
	if err != nil {
		return "", epsMeta, false, err
	}

	lines = strings.Split(string(NormalizeNewlines(out)), "\n")

	avNames := ""
	for _, line := range lines {
		//	fmt.Println(line)
		if strings.Contains(line, "displayName") {
			spl := strings.Split(line, "=")
			if len(spl) == 2 {
				avNames = avNames + spl[1] + " "
			}
		}
	}

	return avNames, epsMeta, epsEnabled, nil
}

func (d DeviceWindows) ScreenLockEnabled() (bool, error) {
	lockTime, err := d.IdleDeviceScreenLockTime()
	logger.Debug(lockTime != "0")
	return lockTime != "0", err
}

func (d DeviceWindows) GetPasswordLastUpdated() (string, error) {

	userc, err := user.Current()
	if err != nil {
		return "", err
	}

	//In windows utils.Username is Computer-Name/username
	//hence it should be splitted
	t := strings.Split(userc.Username, `\`)
	username := ""
	if len(t) == 2 {
		username = t[1]
	}

	c := exec.Command("net", "user", username)
	out, err := c.CombinedOutput()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(NormalizeNewlines(out)), "\n")

	for _, line := range lines {
		if strings.Contains(line, "Password last set") {
			splittedLine := strings.Split(line, "Password last set")
			if len(splittedLine) == 2 {
				return strings.TrimSpace(splittedLine[1]), nil
			} else {
				return "", errors.New("password last updated not found")
			}
		}
	}
	return "", errors.New("password last updated not found")

}
func (d DeviceWindows) GetCriticalAutoUpdateStatus() (bool, error) {

	/*
		1. never check for updates
		2. check for update but let user choose to download
		3. auto download but lets user choose to install
		4. auto download, auto install
	*/

	c := exec.Command(`powershell`, `-Command`, `$AUSettings=(New-Object -com Microsoft.Update.AutoUpdate).Settings;echo $AUSettings.NotificationLevel`)
	out, err := c.CombinedOutput()
	if err != nil {
		return false, err
	}
	notificationLevel, err := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)
	if err != nil {
		return false, err
	}

	return notificationLevel == 4, nil

}
func (d DeviceWindows) GetPendingUpdates() ([]string, error) {
	c := exec.Command(`powershell`, `-Command`, `$UpdateSession = New-Object -ComObject Microsoft.Update.Session;$UpdateSearcher = $UpdateSession.CreateupdateSearcher();$Updates = @($UpdateSearcher.Search("IsHidden=0 and IsInstalled=0").Updates);$Updates | Select-Object Title`)

	updates := []string{}
	out, err := c.CombinedOutput()
	if err != nil {
		return updates, err
	}
	//fmt.Println(string(out))
	str := string(NormalizeNewlines(out))
	lines := strings.Split(str, "\n")

	for _, line := range lines {
		if !strings.Contains(line, "Title") && !strings.Contains(line, "----") && strings.TrimSpace(line) != "" {
			updates = append(updates, strings.TrimSpace(line))
		}
	}

	return updates, nil
}

func (d DeviceWindows) IsFireWallSet() (bool, error) {
	domainProfileCmd := exec.Command(`netsh`, `advfirewall`, `show`, `domain`)
	out, err := domainProfileCmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	str := string(NormalizeNewlines(out))
	lines := strings.Split(str, "\n")
	isfirewallOnDomain := false
	for _, line := range lines {
		if strings.Contains(line, "State") {
			if strings.Contains(line, "ON") {
				isfirewallOnDomain = true
			}
		}
	}
	publicProfileCmd := exec.Command(`netsh`, `advfirewall`, `show`, `public`)
	out, err = publicProfileCmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	str = string(NormalizeNewlines(out))
	lines = strings.Split(str, "\n")
	isfirewallOnPublic := false
	for _, line := range lines {
		if strings.Contains(line, "State") {
			if strings.Contains(line, "ON") {
				isfirewallOnPublic = true
			}
		}
	}
	privateProfileCmd := exec.Command(`netsh`, `advfirewall`, `show`, `private`)
	out, err = privateProfileCmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	str = string(NormalizeNewlines(out))
	lines = strings.Split(str, "\n")
	isfirewallOnPrivate := false
	for _, line := range lines {
		if strings.Contains(line, "State") {
			if strings.Contains(line, "ON") {
				isfirewallOnPrivate = true
			}
		}
	}
	return isfirewallOnDomain && isfirewallOnPrivate && isfirewallOnPublic, nil

}

func (d DeviceWindows) GetNetwork() (string, string, string, error) {
	//netsh interface show interface

	// sample output
	// IPAddress         : 192.168.100.91
	// InterfaceIndex    : 22
	// InterfaceAlias    : Wi-Fi
	// AddressFamily     : IPv4
	// Type              : Unicast
	// PrefixLength      : 24
	// PrefixOrigin      : Manual
	// SuffixOrigin      : Manual
	// AddressState      : Deprecated
	// ValidLifetime     : Infinite ([TimeSpan]::MaxValue)
	// PreferredLifetime : Infinite ([TimeSpan]::MaxValue)
	// SkipAsSource      : False
	// PolicyStore       : ActiveStore

	ipAddr := ""

	// hardwareAddr := ""
	defaultInterface := ""

	netshCmd := exec.Command("powershell", "-Command", "Find-NetRoute -RemoteIPAddress \"1.1.1.1\"")
	out, err := netshCmd.CombinedOutput()
	if err == nil {

		str := string(NormalizeNewlines(out))

		//	fmt.Println(str)
		lines := strings.Split(str, "\n")

		for _, line := range lines {
			splitted := strings.Split(line, ":")
			//	fmt.Println("splitted ", splitted)
			if len(splitted) == 2 && strings.Contains(splitted[0], "IPAddress") {
				ipAddr = strings.TrimSpace(splitted[1])
			}

			if len(splitted) == 2 && strings.Contains(splitted[0], "InterfaceAlias") {
				defaultInterface = strings.TrimSpace(splitted[1])
			}

		}
		return "", "", "", err
	}

	// fmt.Println(hardwareAddr)
	// fmt.Println(defaultInterface)
	// we retreive ip address and mac address for this route.

	// network, err := ghw.Network()

	// for _, nic := range net.NICs {
	// 	fmt.Printf(" %v\n", nic)

	// 	enabledCaps := make([]int, 0)
	// 	for x, cap := range nic.Capabilities {
	// 		if cap.IsEnabled {
	// 			enabledCaps = append(enabledCaps, x)
	// 		}
	// 	}
	// 	if len(enabledCaps) > 0 {
	// 		fmt.Printf("  enabled capabilities:\n")
	// 		for _, x := range enabledCaps {
	// 			fmt.Printf("   - %s\n", nic.Capabilities[x].Name)
	// 		}
	// 	}
	// }

	// var netVal net.InterfaceStat
	// nets, err := net.Interfaces()
	// if err != nil {
	// 	return "", "", "", err
	// }
	// for _, v := range nets {
	// 	//fmt.Println(v.Name)
	// 	fmt.Println(v.String())
	// 	if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	// 		netVal.Addrs = v.Addrs
	// 		netVal.HardwareAddr = v.HardwareAddr
	// 	}
	// }

	// addr := ""
	// for _, n := range netVal.Addrs {
	// 	addr = addr + "," + n.Addr
	// }

	// addr = strings.TrimLeft(addr, ",")

	// return addr, netVal.HardwareAddr, string(defaultInterface), nil

	interfaces, err := net.Interfaces()

	hardwareAddr := ""

	if defaultInterface != "" {
		for _, v := range interfaces {

			if v.Name == defaultInterface {
				hardwareAddr = v.HardwareAddr.String()
			}
		}

	} else {
		for _, v := range interfaces {

			if v.Flags&net.FlagUp != 0 && v.Flags&net.FlagLoopback == 0 {
				hardwareAddr = v.HardwareAddr.String()
				defaultInterface = v.Name
				ips, err := v.Addrs()
				if err == nil {
					for _, ip := range ips {
						ipAddr = ipAddr + ip.String() + "  "
					}
				}

			}

		}

	}

	return ipAddr, hardwareAddr, defaultInterface, nil
}

func (d DeviceWindows) GetLatestSecurityPatch() (string, error) {
	return "", errors.New("Not supported yet")
}

func (d DeviceWindows) IdleDeviceScreenLockTime() (string, error) {
	out, err := Execute("powercfg", "/getactivescheme")
	if err != nil {
		//fmt.Println(err)
		return "", err
	}

	powerSchemeGuids := Extract(`[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`, out)
	if len(powerSchemeGuids) != 1 {
		return "", errors.New("not found")
	}

	str := ""
	regex := ""
	out, err = Execute("powercfg", "/query", powerSchemeGuids[0], "7516b95f-f776-4464-8c53-06167f40cc99", "3c0bc021-c8a8-4e07-a973-6b14cbcb2b7e")
	if err != nil {
		out, err = Execute("powercfg", "/query", powerSchemeGuids[0], "7516b95f-f776-4464-8c53-06167f40cc99")
		if err != nil {
			return "", err
		}
		spll := strings.Split(string(out), "Power Setting GUID:")
		for _, powerSett := range spll {
			if strings.Contains(powerSett, "3c0bc021-c8a8-4e07-a973-6b14cbcb2b7e") {
				str = powerSett
			}
		}
		regex = `AC Power Setting Index: .*0x([0-9a-fA-F]{8,8})`

	} else {
		str = string(out)
		regex = `AC.*0x([0-9a-fA-F]{8,8})`
	}

	ac := Extract(regex, str)
	if len(ac) < 2 {
		// logger.Trace("not found: ")
		return "", errors.New("power config not found")
	}

	t, err := strconv.ParseInt(ac[1], 16, 64)
	if err != nil {
		logger.Trace(err)
		return "", err
	}
	//fmt.Println(t)
	return fmt.Sprintf("%d", t), nil
}

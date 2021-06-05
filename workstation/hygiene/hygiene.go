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
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/jaypipes/ghw"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true)

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: false,
		//PadLevelText:              false,
		FieldMap: nil,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return filepath.Base(frame.Function), fmt.Sprintf(`%s:%d`, filepath.Base(frame.File), frame.Line)
		},
	})
}

type Device interface {
	IsAutoLoginEnabled() (bool, error)
	IsFireWallSet() (bool, error)
	IsDeviceEncrypted() (bool, error)
	GetInstalledPackages() ([]string, error)
	//GetDeviceName() (string, error)
	GetOSNameVersion() (osName, osVersion, kernelVersion string, err error)
	GetPasswordLastUpdated() (string, error)
	GetCriticalAutoUpdateStatus() (bool, error)
	GetPendingUpdates() ([]string, error)
	IsRemoteConnectionEnabled() (bool, error)
	ScreenLockEnabled() (bool, error)
	GetNetwork() (string, string, string, error)
	GetLatestSecurityPatch() (string, error)
	IdleDeviceScreenLockTime() (string, error)
	EndpointSecurity() (string, string, bool, error)
}

func newDevice(deviceType string) Device {
	switch deviceType {
	case "darwin":
		return DeviceMac{}
	case "linux":
		return DeviceLinux{}
	case "windows":
		return DeviceWindows{}
	default:
		return nil
	}

}

// GetDeviceHygiene returns os relevant device hygiene
func GetDeviceHygiene(osType string) DeviceHygiene {

	device := newDevice(osType)

	machineID, err := machineid.ID()
	if err != nil {
		logrus.Warn(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		logrus.Warn(err)
	}

	start := time.Now()
	isAutoLoginEnabled, err := device.IsAutoLoginEnabled()
	if err != nil {
		logrus.Warn(err)
	}
	//logrus.Trace("autologin---------------------------------------------------------->", isAutoLoginEnabled, err)
	diff := time.Now().Sub(start)
	start = time.Now()
	logrus.Tracef("IsAutoLoginEnabled: %v", diff.Milliseconds())

	scrnLockTime, err := device.IdleDeviceScreenLockTime()
	if err != nil {
		logrus.Warn(err)
	}

	diff = time.Now().Sub(start)
	start = time.Now()
	logrus.Tracef("IdleDeviceScreenLockTime: %v", diff.Milliseconds())

	epsName, epsMeta, epsEnabled, err := device.EndpointSecurity()

	if err != nil {
		logrus.Warn(err)
	}

	diff = time.Now().Sub(start)
	start = time.Now()
	logrus.Tracef("EndpointSecurity: %v", diff.Milliseconds())

	isFwSet, err := device.IsFireWallSet()
	if err != nil {
		logrus.Warn(err)
	}

	diff = time.Now().Sub(start)
	start = time.Now()

	logrus.Tracef("IsFireWallSet: %v", diff.Milliseconds())

	isDeviceEncryptionSet, err := device.IsDeviceEncrypted()
	if err != nil {
		logrus.Warn(err)
	}

	diff = time.Now().Sub(start)
	start = time.Now()
	logrus.Tracef("IsDeviceEncrypted: %v", diff.Milliseconds())

	// deviceName, err := device.GetDeviceName()
	// if err != nil {
	// 	logrus.Warn(err)
	// }

	osName, osVersion, kernelVer, err := device.GetOSNameVersion()
	if err != nil {
		logrus.Warn(err)
	}

	diff = time.Now().Sub(start)
	start = time.Now()
	logrus.Tracef("GetOSNameVersion: %v", diff.Milliseconds())

	lastPassUpdated, err := device.GetPasswordLastUpdated()
	if err != nil {
		logrus.Warn(err)
	}

	diff = time.Now().Sub(start)
	start = time.Now()

	logrus.Tracef("GetPasswordLastUpdated: %v", diff.Milliseconds())

	autoUpdateStatus, err := device.GetCriticalAutoUpdateStatus()
	if err != nil {
		logrus.Warn(err)
	}

	diff = time.Now().Sub(start)
	start = time.Now()
	logrus.Tracef("GetCriticalAutoUpdateStatus: %v", diff.Milliseconds())

	//istalledApplications,err:=device.GetInstalledPackages()
	//if err != nil {
	//	logrus.Info(err)
	//}

	screenLock, err := device.ScreenLockEnabled()
	if err != nil {
		logrus.Warn(err)
	}

	diff = time.Now().Sub(start)
	start = time.Now()

	logrus.Tracef("ScreenLockEnabled: %v", diff.Milliseconds())

	remotelogin, err := device.IsRemoteConnectionEnabled()
	if err != nil {
		logrus.Warn(err)
	}

	diff = time.Now().Sub(start)
	start = time.Now()

	logrus.Tracef("IsRemoteConnectionEnabled: %v", diff.Milliseconds())

	product, err := ghw.Product()
	if err != nil {
		logrus.Warn(err)
		product = &ghw.ProductInfo{}
		if osType == "darwin" {
			product.Vendor = "Apple"
			product.Name, product.Version, err = GetProduct()
			if err != nil {
				logrus.Warn(err)
			}
		}

	}

	diff = time.Now().Sub(start)
	start = time.Now()

	logrus.Tracef("Product: %v", diff.Milliseconds())

	ipAddr, macAddr, intName, err := GetNetwork()
	if err != nil {
		logrus.Warn(err)
	}

	diff = time.Now().Sub(start)
	start = time.Now()
	logrus.Tracef("GetNetwork: %v", diff.Milliseconds())

	// logrus.Trace("GetPendingUpdates")
	//pendingUpdates, err := device.GetPendingUpdates()
	//
	//if err != nil {
	//	logrus.Warn(err)
	//}

	devHyg := DeviceHygiene{
		DeviceInfo: DeviceInfo{
			DeviceName:    product.Name,
			DeviceVersion: product.Version,
			MachineID:     machineID,
			Brand:         product.Vendor,
			Manufacturer:  product.Vendor,
			DeviceModel:   product.Family,
		},
		DeviceOS: DeviceOS{
			OSName:              osName,
			OSVersion:           osVersion,
			KernelType:          osType,
			KernelVersion:       kernelVer,
			ReadableVersion:     "",
			LatestSecurityPatch: "",
			AutoUpdate:          autoUpdateStatus,
			//PendingUpdates:      pendingUpdates,
		},
		LoginSecurity: LoginSecurity{
			AutologinEnabled:         isAutoLoginEnabled,
			LoginMethod:              "",
			IdleDeviceScreenLockTime: scrnLockTime,
			IdleDeviceScreenLock:     screenLock,
			PasswordLastUpdated:      lastPassUpdated,
			RemoteLoginEnabled:       remotelogin,
		},
		NetworkInfo: NetworkInfo{
			Hostname:        hostname,
			InterfaceName:   intName,
			IPAddress:       ipAddr,
			MacAddress:      macAddr,
			WirelessNetwork: false,
			NetworkName:     "",
			NetworkSecurity: "",
		},
		EndpointSecurity: EndpointSecurity{
			EpsConfigured:           epsEnabled,
			EpsVendorName:           epsName,
			EpsMeta:                 epsMeta,
			FirewallEnabled:         isFwSet,
			DeviceEncryptionEnabled: isDeviceEncryptionSet,
		},
	}
	//dev.DeviceFinger.SecurityStatus.IsFirewallSet=true
	return devHyg

}

func GetNetwork() (string, string, string, error) {

	intfs := ""
	ips := ""
	macs := ""

	ifaces, err := net.Interfaces()
	if err != nil {
		return ips, macs, intfs, err
	}
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", "", "", err
		}

		//check if interface is up and is not loopback
		if (i.Flags&net.FlagLoopback != 0) || (i.Flags&net.FlagUp == 0) {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.IsLoopback() || ip.To4() == nil {
				continue
			}
			ips = ips + ip.String() + ","
		}

		if len(addrs) > 0 {
			if i.HardwareAddr.String() != "" {
				macs = macs + i.HardwareAddr.String() + ","
			}
			intfs = intfs + i.Name + ","
		}
	}

	ips = strings.TrimRight(ips, ",")
	macs = strings.TrimRight(macs, ",")
	intfs = strings.TrimRight(intfs, ",")
	return ips, macs, intfs, nil
}

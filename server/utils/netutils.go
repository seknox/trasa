package utils

import (
	"github.com/sirupsen/logrus"
	"net"
	"strings"
)

//I sPrivateIP returns boolean value based on ip type
func IsPrivateIP(ip net.IP) bool {
	if ip == nil {
		logrus.Error("is private ip: ip is nil")
		return true
	}
	var privateIPBlocks []*net.IPNet
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)

	}
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

//GetIPFromAddr returns IP address as string from net.Addr type
func GetIPFromAddr(addr net.Addr) string {
	if addr == nil {
		return "0.0.0.0"
	}
	addrStr := addr.String()
	if addrStr == "" {
		return "0.0.0.0"
	}
	splitted := strings.Split(addr.String(), ":")
	if len(splitted) < 1 {
		return "0.0.0.0"
	} else {
		return splitted[0]
	}
}

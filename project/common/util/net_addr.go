package util

import (
	"net"
	"strings"
)

func IPToInt(ip string) uint32 {
	var ipInt uint32 = 0
	ips := strings.Split(ip, ":")
	netIP := net.ParseIP(ips[0])
	ipInt += uint32(netIP[15]) << 24
	ipInt += uint32(netIP[14]) << 16
	ipInt += uint32(netIP[13]) << 8
	ipInt += uint32(netIP[12])
	return ipInt
}

func IPToStr(ip uint32) string {
	var netIP [net.IPv4len]byte

	netIP[0] = byte(ip & 0xff)
	netIP[1] = byte((ip >> 8) & 0xff)
	netIP[2] = byte((ip >> 16) & 0xff)
	netIP[3] = byte((ip >> 24) & 0xff)
	return net.IP(netIP[:]).String()
}

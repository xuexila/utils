package tools

import "net"

// IsIP 判断输入是否是IP
func IsIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

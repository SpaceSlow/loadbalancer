package networks

import "net"

func ParseIP(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return "unknown"
	}
	return ip
}

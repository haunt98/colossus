package network

import (
	"fmt"
	"net"
)

// Random copy on the Internet
// God please forgives me
func GetIP() (string, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("failed to get interface addresses: %w", err)
	}

	for _, address := range addresses {
		ipNet, ok := address.(*net.IPNet)
		if ok && ipNet.IP.IsGlobalUnicast() && !ipNet.IP.IsLoopback() &&
			ipNet.IP.To4() != nil && ipNet.IP.To16() != nil {
			return ipNet.IP.String(), nil
		}
	}
	return "", nil
}

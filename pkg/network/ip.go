package network

import (
	"fmt"
	"net"
)

func GetIP() (string, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("failed to get interface addresses: %w", err)
	}

	for _, address := range addresses {
		ipNet, ok := address.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ip := ipNet.IP.String()
			return ip, nil
		}
	}
	return "", nil
}

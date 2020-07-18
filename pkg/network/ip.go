package network

import (
	"fmt"
	"net"

	"go.uber.org/zap"
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

func GetIPByDial(sugar *zap.SugaredLogger) (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("failed to dial connection: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			sugar.Error(err)
		}
	}()

	// localAddr := conn.LocalAddr().(*net.UDPAddr)
	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return "", fmt.Errorf("failed to bind udp address")
	}

	return addr.String(), nil
}

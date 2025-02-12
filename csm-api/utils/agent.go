package utils

import (
	"fmt"
	"net"
	"os"
)

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "UNKNOWN"
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return "UNKNOWN"
}

func GetAgent() string {
	host := getLocalIP()        // 현재 IP 주소
	osUser := os.Getenv("USER") // Unix 계열 OS
	if osUser == "" {
		osUser = os.Getenv("USERNAME") // Windows OS
	}
	module := "go http server" // 하드코딩된 모듈명

	// 최종 포맷된 데이터 생성
	result := fmt.Sprintf("HOST:%s/OS_USER:%s/MODULE:%s", host, osUser, module)

	return result
}

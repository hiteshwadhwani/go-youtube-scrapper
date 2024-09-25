package utils

import (
	"net"
	"net/http"
)

func GetUserIP(r *http.Request) net.IP {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	return net.ParseIP(IPAddress)
}

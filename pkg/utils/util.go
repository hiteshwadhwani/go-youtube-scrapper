package utils

import (
	"net"
	"net/http"
	"strconv"

	"github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/constants"
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

func GetPageOffsetAndLimit(req *http.Request) (int, int) {
	pageLimit := req.URL.Query().Get("limit")
	pageOffset := req.URL.Query().Get("offset")

	if pageLimit == "" {
		pageLimit = constants.DEFAULT_PAGE_LIMIT
	}

	if pageOffset == "" {
		pageOffset = constants.DEFAULT_PAGE_OFFSET
	}

	pageLimitInt, _ := strconv.Atoi(pageLimit)
	pageOffsetInt, _ := strconv.Atoi(pageOffset)

	return pageLimitInt, pageOffsetInt
}

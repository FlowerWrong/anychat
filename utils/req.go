package utils

import (
	"net"
	"net/http"
	"strings"
)

const remoteAddrHeader = "REMOTE_ADDR"

// Host tries its best to return the request host.
func Host(r *http.Request) string {
	if r.URL.IsAbs() {
		host := r.Host
		// Slice off any port information.
		if i := strings.Index(host, ":"); i != -1 {
			host = host[:i]
		}
		return host
	}
	return r.URL.Host
}


// FetchHeaders extracts specified headers from request
func FetchHeaders(r *http.Request, list []string) map[string]string {
	res := make(map[string]string)

	for _, header := range list {
		res[header] = r.Header.Get(header)
	}
	res[remoteAddrHeader], _, _ = net.SplitHostPort(r.RemoteAddr)
	return res
}

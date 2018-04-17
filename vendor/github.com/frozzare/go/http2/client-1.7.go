// +build go1.7

package http2

import (
	"net"
	"net/http"
	"time"
)

// DefaultTransport is a custom default implementation of Transport.
var DefaultTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).Dial,
	TLSHandshakeTimeout:   10 * time.Second,
	ResponseHeaderTimeout: 10 * time.Second,
	ExpectContinueTimeout: 10 * time.Second,
	IdleConnTimeout:       60 * time.Second,
	MaxIdleConns:          5,
	MaxIdleConnsPerHost:   5,
	DisableKeepAlives:     false,
}

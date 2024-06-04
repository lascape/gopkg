package httpx

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

var (
	hdrContentTypeKey   = http.CanonicalHeaderKey("Content-Type")
	hdrAuthorizationKey = http.CanonicalHeaderKey("Authorization")

	jsonContentType = "application/json"
	formContentType = "application/x-www-form-urlencoded"
)

// 全局http client ,使用全局client，端口复用
var defaultHttpClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			d := net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}
			conn, e = d.DialContext(ctx, network, addr)
			return
		},
		TLSHandshakeTimeout:   3 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
		MaxIdleConnsPerHost:   512,
		MaxConnsPerHost:       1000,
		IdleConnTimeout:       90 * time.Second,
	},
}

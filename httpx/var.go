package httpx

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/lascape/gopkg/envx"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var defaultHttp = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			d := net.Dialer{Timeout: 10 * time.Second, KeepAlive: 10 * time.Second}
			addr = tryDnsFormGlobalHost(addr)
			conn, e = d.DialContext(ctx, network, addr)
			return
		},
		TLSHandshakeTimeout:   3 * time.Second,  // 增加到 3 秒，避免网络稍差时频繁超时
		ResponseHeaderTimeout: 15 * time.Second, // 稍增大以应对慢响应
		ExpectContinueTimeout: 2 * time.Second,  // 增加到 2 秒，提高容错性
		MaxIdleConns:          512,              // 保持不变，适合多目标访问
		IdleConnTimeout:       60 * time.Second, // 降低到 60 秒，更适合短时间活跃的连接
		MaxIdleConnsPerHost:   100,              // 设置为 100，限制单主机连接池资源占用
		MaxConnsPerHost:       128,              // 启用并限制每主机最大连接数
		Proxy:                 proxyURL(),
	},
}

func proxyURL() func(*http.Request) (*url.URL, error) {
	urls := strings.Split(envx.PKG_HTTPX_ENGRESS, ",")
	var ips []*url.URL
	for _, s := range urls {
		parse, err := url.Parse(s)
		if err != nil {
			continue
		}
		if parse.Scheme != "http" && parse.Scheme != "https" {
			continue
		}

		ips = append(ips, parse)
	}
	if len(ips) == 0 {
		return nil
	}
	noproxy := strings.Split(envx.PKG_HTTPX_ENGRESS_IGNORE, ",")
	r := rand.NewSource(time.Now().Unix())
	return func(req *http.Request) (*url.URL, error) {
		if len(ips) == 0 {
			return nil, nil
		}
		for _, no := range noproxy {
			if no == req.URL.Host {
				return nil, nil
			}
		}
		ip := ips[r.Int63()%int64(len(ips))]
		fmt.Printf("req host:%s,proxy:%s\n", req.URL.Host, ip)
		return ip, nil
	}
}

var globalHost = new(sync.Map)

func AddGlobalHost(domain, ip string) {
	globalHost.Store(domain, ip)
}

func tryDnsFormGlobalHost(addr string) string {
	if !strings.Contains(addr, ":") {
		return addr
	}
	t := strings.Split(addr, ":")
	if len(t) != 2 {
		return addr
	}
	domain := t[0]
	ip, ok := globalHost.Load(domain)
	if !ok {
		return addr
	}
	_ip, ok := ip.(string)
	if !ok {
		return addr
	}
	addr = strings.Replace(addr, domain, _ip, 1)
	return addr
}

func SetTransport(t *http.Transport) {
	t.Proxy = proxyURL()
	defaultHttp.Transport = t
}

func SetCheckRedirect(f func(req *http.Request, via []*http.Request) error) {
	defaultHttp.CheckRedirect = f
}

// 去掉302自动跳转
//CheckRedirect: func(req *http.Request, via []*http.Request) error {
//	return http.ErrUseLastResponse
//},

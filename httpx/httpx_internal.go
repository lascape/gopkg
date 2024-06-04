package httpx

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
)

// 真正用于发送请求的函数
func (c *Client) request(ctx context.Context, method string) *Response {
	//记录请求方式
	c.method = method
	var (
		r = &Response{
			client:  c,
			errs:    new(strings.Builder),
			Request: &Request{},
		}
		err error
	)

	r.Request.RawRequest, err = http.NewRequestWithContext(ctx, method, c.uri, bytes.NewReader(c.body))
	if err != nil {
		c.errs.WriteString(err.Error() + ";")
		return r
	}

	//设置请求头
	r.Request.RawRequest.Header = c.header

	//设置Cookie
	for _, cookie := range c.cookies {
		r.Request.RawRequest.AddCookie(cookie)
	}

	for _, middleware := range c.beforeRequest {
		if err := middleware(c, r.Request); err != nil {
			r.errs.WriteString(err.Error() + ";")
			return r
		}
	}

	//发送请求
	r.RawResponse, err = c.client.Do(r.Request.RawRequest)
	if err != nil {
		r.errs.WriteString(err.Error() + ";")
		return r
	}

	for _, middleware := range c.afterResponse {
		if err := middleware(c, r); err != nil {
			r.errs.WriteString(err.Error() + ";")
			return r
		}
	}
	return r
}

// 构建CURL请求链接
func (c *Client) buildCurl() string {
	curl := fmt.Sprintf("curl  -X %s '%s'", c.method, c.uri)
	for k := range c.header {
		v := c.header.Get(k)
		if len(v) > 0 {
			curl += fmt.Sprintf(" -H '%s:%s'", k, v)
		}
	}

	if len(c.cookies) > 0 {
		curl += " -H 'Cookie: "
		for _, v := range c.cookies {
			curl += fmt.Sprintf("%s=%s;", v.Name, v.Value)
		}
		curl += "'"
	}

	if len(c.body) > 0 && len(c.body) < 2048 {
		curl += fmt.Sprintf(" -d '%s'", c.body)
	}

	return curl
}

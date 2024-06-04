package httpx

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	client        *http.Client
	header        http.Header
	cookies       []*http.Cookie
	body          []byte
	errs          *strings.Builder
	timeout       time.Duration
	uri           string
	method        string
	beforeRequest []RequestMiddleware
	afterResponse []ResponseMiddleware
}

type Option func(c *Client)

func MustClient(opts ...Option) interface{ New(uri string) *Client } {
	o := &Client{
		errs:   new(strings.Builder),
		client: defaultHttpClient,
		header: http.Header{},
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func WithTransport(transport *http.Transport) Option {
	return func(c *Client) {
		if transport == nil {
			return
		}
		c.client.Transport = transport
	}
}

func New(uri string) *Client {
	return &Client{
		errs:   new(strings.Builder),
		uri:    uri,
		client: defaultHttpClient,
		header: http.Header{},
	}
}

// SetDefaultTransport 设置defaultHttpClient的Transport
func SetDefaultTransport(transport *http.Transport) {
	defaultHttpClient.Transport = transport
}

func (c *Client) Get(ctx context.Context) *Response {
	return c.request(ctx, http.MethodGet)
}

func (c *Client) Post(ctx context.Context) *Response {
	return c.request(ctx, http.MethodPost)
}

func (c *Client) Delete(ctx context.Context) *Response {
	return c.request(ctx, http.MethodDelete)
}

func (c *Client) Head(ctx context.Context) *Response {
	return c.request(ctx, http.MethodHead)
}

func (c *Client) Put(ctx context.Context) *Response {
	return c.request(ctx, http.MethodPut)
}

func (c *Client) Options(ctx context.Context) *Response {
	return c.request(ctx, http.MethodOptions)
}

func (c *Client) New(uri string) *Client {
	c.uri = uri
	return c
}

// OnBeforeRequest 设置请求前对Request的设置
func (c *Client) OnBeforeRequest(m ...RequestMiddleware) *Client {
	c.beforeRequest = append(c.beforeRequest, m...)
	return c
}

// OnAfterResponse 设置接收到响应体之后执行的函数，
func (c *Client) OnAfterResponse(m ...ResponseMiddleware) *Client {
	c.afterResponse = append(c.afterResponse, m...)
	return c
}

// SetTimeout 设置请求超时断开的时间
func (c *Client) SetTimeout(t time.Duration) *Client {
	c.timeout = t
	return c
}

func (c *Client) SetContentTypeJSON() *Client {
	c.header.Set(hdrContentTypeKey, jsonContentType)
	return c
}

func (c *Client) SetContentTypeFORM() *Client {
	c.header.Set(hdrContentTypeKey, formContentType)
	return c
}

func (c *Client) SetHeader(k, v string) *Client {
	c.header.Set(k, v)
	return c
}

func (c *Client) SetBodyIo(reader io.Reader) *Client {
	b, err := io.ReadAll(reader)
	if err != nil {
		c.errs.WriteString(err.Error() + ";")
		return c
	}
	c.SetBody(b)
	return c
}

// SetBodyJson 默认设置Content-Type = "application/json"
// 建议在 SetBody* 之后执行，否则可能被默认替换
func (c *Client) SetBodyJson(data any) *Client {
	js, err := json.Marshal(data)
	if err != nil {
		c.errs.WriteString(err.Error() + ";")
		return c
	}
	c.SetBody(js)
	c.SetContentTypeJSON()
	return c
}

// SetBodyUrlValues 默认设置Content-Type = "application/x-www-form-urlencoded"
// 建议在 SetBody* 之后执行，否则可能被默认替换
func (c *Client) SetBodyUrlValues(params url.Values) *Client {
	c.SetBody([]byte(params.Encode()))
	c.SetContentTypeFORM()
	return c
}

func (c *Client) SetBody(b []byte) *Client {
	c.body = b
	return c
}

// SetQueryValues 从url.Values解析请求参数，并设置在URI的?之后
func (c *Client) SetQueryValues(params url.Values) *Client {
	if len(params) == 0 {
		return nil
	}
	c.SetQuery(params.Encode())
	return c
}

// SetQuery 设置URI的?之后的参数
func (c *Client) SetQuery(params string) *Client {
	if len(params) == 0 {
		return c
	}
	if strings.Index(c.uri, "?") >= 0 {
		c.uri = c.uri + "&" + params
	} else {
		c.uri = c.uri + "?" + params
	}
	return c
}

func (c *Client) SetHeaders(headers map[string]string) *Client {
	for h, v := range headers {
		c.SetHeader(h, v)
	}
	return c
}

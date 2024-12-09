package httpx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/lascape/gopkg/envx"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	req          *http.Request
	body         io.Reader
	errs         *strings.Builder
	timeout      time.Duration
	uri          string
	startTime    time.Time
	mustHttpCode int
}

func NewWithRequest(req *http.Request, uri string) *Client {
	return &Client{req: req, uri: uri, startTime: time.Now(), errs: &strings.Builder{}}
}

func New(uri string) *Client {
	return &Client{req: &http.Request{Header: make(http.Header)}, uri: uri, startTime: time.Now(), errs: &strings.Builder{}}
}

func (h *Client) MustCode(code int) *Client {
	h.mustHttpCode = code
	return h
}

func (h *Client) SetHeader(key, value string) *Client {
	h.req.Header.Set(key, value)
	return h
}

func (h *Client) SetHeaders(hs map[string]string) *Client {
	for key, value := range hs {
		h.SetHeader(key, value)
	}
	return h
}

func (h *Client) SetTimeout(t time.Duration) *Client {
	h.timeout = t
	return h
}

func (h *Client) SetContentTypeJson() *Client {
	h.SetHeader("Content-Type", "application/json")
	return h
}

func (h *Client) SetContentTypeFormUrlencoed() *Client {
	h.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	return h
}

func (h *Client) AddCookie(cs ...*http.Cookie) *Client {
	for _, c := range cs {
		h.req.AddCookie(c)
	}
	return h
}

func (h *Client) SetBodyIO(b io.Reader) *Client {
	h.body = b
	return h
}

func (h *Client) SetBodyBytes(b []byte) *Client {
	h.body = bytes.NewReader(b)
	return h
}

func (h *Client) SetBodyString(s string) *Client {
	h.body = strings.NewReader(s)
	return h
}

func (h *Client) SetBodyUrlValues(params url.Values) *Client {
	h.body = strings.NewReader(params.Encode())
	return h
}

func (h *Client) SetBodyJson(data interface{}) *Client {
	js, err := json.Marshal(data)
	if err != nil {
		h.errs.WriteString(err.Error())
		h.errs.WriteString(";")
		return h
	}
	h.body = bytes.NewReader(js)
	return h
}

func (h *Client) SetQueryString(params string) *Client {
	if params != "" {
		if strings.Index(h.uri, "?") >= 0 {
			h.uri = h.uri + "&" + params
		} else {
			h.uri = h.uri + "?" + params
		}
	}
	return h
}

func (h *Client) SetQuery(params url.Values) *Client {
	return h.SetQueryString(params.Encode())
}

func (h *Client) Error() error {
	if h.errs.Len() == 0 {
		return nil
	}
	return errors.New(h.errs.String())
}

func (h *Client) Get(ctx context.Context) *Response {
	return h.request(ctx, "GET")
}

func (h *Client) Post(ctx context.Context) *Response {
	return h.request(ctx, "POST")
}

func (h *Client) Head(ctx context.Context) *Response {
	return h.request(ctx, "HEAD")
}

func (h *Client) Delete(ctx context.Context) *Response {
	return h.request(ctx, "DELETE")
}

func (h *Client) Put(ctx context.Context) *Response {
	return h.request(ctx, "PUT")
}

func (h *Client) Options(ctx context.Context) *Response {
	return h.request(ctx, "OPTIONS")
}

func (h *Client) request(ctx context.Context, method string) *Response {
	if h.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.timeout)
		defer cancel()
	}

	var curl string
	if envx.PKG_HTTPX_BCURL {
		var data = ""
		if h.body != nil {
			t, _ := io.ReadAll(h.body)
			data = string(t)
			h.body = bytes.NewReader(t)
		}

		curl = buildCurl(h.uri, method, data, h.req.Header, h.req.Cookies())
	}

	resp := &Response{
		curl:      curl,
		startTime: h.startTime,
		errs:      new(strings.Builder),
	}

	req, err := http.NewRequestWithContext(ctx, method, h.uri, h.body)
	if err != nil {
		h.errs.WriteString(err.Error())
		h.errs.WriteString(";")
		resp.errs.WriteString(h.errs.String())
		return resp
	}

	req.Header = h.req.Header
	for _, cookie := range h.req.Cookies() {
		req.AddCookie(cookie)
	}

	res, err := defaultHttp.Do(req)
	if err != nil {
		// https://stackoverflow.com/questions/28046100/golang-http-concurrent-requests-post-eof
		if err == io.EOF {
			for i := 0; i < 2; i++ {
				res, err = defaultHttp.Do(req)
				if err != io.EOF {
					break
				}
			}
		}
		if err != nil {
			h.errs.WriteString(err.Error())
			h.errs.WriteString(";")
			resp.errs.WriteString(h.errs.String())
			return resp
		}
	}
	// 如果限定必须要返回的状态码
	// 对不上则认为是错误的
	if h.mustHttpCode > 0 {
		if res.StatusCode != h.mustHttpCode {
			h.errs.WriteString("bad http status ")
			h.errs.WriteString(fmt.Sprintf("%d", res.StatusCode))
			h.errs.WriteString(";")
		}
	}

	resp.body = h.readBody(res)
	resp.res = res
	resp.errs.WriteString(h.errs.String())
	return resp
}

func (h *Client) readBody(req *http.Response) []byte {
	b, e := io.ReadAll(req.Body)
	if e != nil {
		h.errs.WriteString(e.Error())
		h.errs.WriteString(";")
	}
	_, _ = io.Copy(io.Discard, req.Body)
	_ = req.Body.Close() //  must close
	req.Body = io.NopCloser(bytes.NewBuffer(b))
	return b
}

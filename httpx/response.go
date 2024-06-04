package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	client      *Client
	Request     *Request
	RawResponse *http.Response
	curl        string
	errs        *strings.Builder //[]error
	body        []byte
	startTime   time.Time
}

func (r *Response) Curl() string {
	if r.curl == "" {
		r.curl = r.client.buildCurl()
	}
	return r.curl
}

func (r *Response) GetBodyString() string {
	return string(r.GetBodyByte())
}

func (r *Response) GetBodyByte() []byte {
	if r.body != nil {
		return r.body
	}
	var err error
	if r.body, err = io.ReadAll(r.RawResponse.Body); err != nil {
		r.errs.WriteString(err.Error() + ";")
		return r.body
	}
	_ = r.RawResponse.Body.Close()
	r.RawResponse.Body = io.NopCloser(bytes.NewReader(r.body))
	return r.body
}

func (r *Response) Error() error {
	if r.errs.Len() <= 0 {
		return nil
	}
	return errors.New(r.errs.String())
}

func (r *Response) Response() *http.Response {
	if r.RawResponse == nil {
		r.RawResponse = new(http.Response)
	}
	return r.RawResponse
}

func (r *Response) StatusCode() int {
	if r.RawResponse == nil {
		return 0
	}
	return r.RawResponse.StatusCode
}

func (r *Response) Unmarshal(v interface{}) *Response {
	err := json.Unmarshal(r.GetBodyByte(), v)
	if err != nil {
		r.errs.WriteString(err.Error() + ";")
	}
	return r
}

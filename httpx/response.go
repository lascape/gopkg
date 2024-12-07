package httpx

import (
	"encoding/xml"
	"errors"
	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"net/http"
	"strings"
	"time"
)

func init() {
	extra.RegisterFuzzyDecoders()
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Response struct {
	res       *http.Response
	curl      string
	errs      *strings.Builder //[]error
	body      []byte
	startTime time.Time
}

func (r *Response) String() string {
	return r.curl
}
func (r *Response) Curl() string {
	return r.curl
}
func (r *Response) GetBodyString() string {
	if r.body == nil {
		return ""
	}
	return string(r.body)
}

func (r *Response) GetBodyByte() []byte {
	return r.body
}

func (r *Response) Error() error {
	if r.errs.Len() <= 0 {
		return nil
	}
	return errors.New(r.errs.String())
}

func (r *Response) Response() *http.Response {
	if r.res == nil {
		r.res = &http.Response{
			Status:           "",
			StatusCode:       0,
			Proto:            "",
			ProtoMajor:       0,
			ProtoMinor:       0,
			Header:           nil,
			Body:             nil,
			ContentLength:    0,
			TransferEncoding: nil,
			Close:            false,
			Uncompressed:     false,
			Trailer:          nil,
			Request:          nil,
			TLS:              nil,
		}
	}
	return r.res
}

func (r *Response) HttpCode() int {
	if r.res == nil {
		return 0
	}
	return r.res.StatusCode
}

func (r *Response) Unmarshal(v interface{}) *Response {
	if r.errs.Len() > 0 {
		return r
	}
	if len(r.body) < 2 { //"{}, []"
		r.errs.WriteString("Unmarshal fail, body not a json:" + string(r.body))
		return r
	}
	err := json.Unmarshal(r.body, v)
	if err != nil {
		r.errs.WriteString(err.Error())
		r.errs.WriteString(";")
		//r.errs = append(r.errs, err)
	}
	return r
}

func (r *Response) UnmarshalXml(v interface{}) *Response {
	if r.errs.Len() > 0 {
		return r
	}
	if len(r.body) < 2 { //"{}, []"
		r.errs.WriteString("Unmarshal fail, body not a json:" + string(r.body))
		return r
	}
	err := xml.Unmarshal(r.body, v)
	if err != nil {
		r.errs.WriteString(err.Error())
		r.errs.WriteString(";")
		//r.errs = append(r.errs, err)
	}
	return r
}

func (r *Response) UseTime() time.Duration {
	return time.Since(r.startTime)
}

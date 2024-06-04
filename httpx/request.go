package httpx

import "net/http"

type Request struct {
	RawRequest *http.Request
}

package httpx

import (
	"fmt"
	"net/http"
)

func buildCurl(uri string, method string, data string, header http.Header, cookies []*http.Cookie) string {
	c := fmt.Sprintf("curl  -X %s '%s'", method, uri)

	for k := range header {
		c += fmt.Sprintf(" -H '%s:%s'", k, header.Get(k))
	}

	if len(cookies) > 0 {
		c += " -H 'Cookie: "
		for _, v := range cookies {
			c += fmt.Sprintf("%s=%s;", v.Name, v.Value)
		}
		c += "'"
	}

	if data != "" {
		if len(data) > 1024*3 {
			data = data[:1024*3]
		}
		c += fmt.Sprintf(" -d '%s'", data)
	}
	return c
}

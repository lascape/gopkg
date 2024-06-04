package httpx

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
)

type baseServer struct {
	Age   int    `json:"age"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestMustClient(t *testing.T) {
	server := createBaseServer(t)
	defer server.Close()
	{
		resp := MustClient().New(server.URL+"/profile").SetQueryValues(url.Values{"name": []string{"jerry"}}).SetHeader(hdrAuthorizationKey, "Bear 123").Get(context.Background())
		assert.JSONEq(t, resp.GetBodyString(), "{\"name\":\"jerry\",\"email\":\"jerry001@gmail.cn\",\"age\":18}")
		assert.Equal(t, resp.StatusCode(), 200)
		assert.NoError(t, resp.Error())
		assert.Condition(t, func() (success bool) {
			v := regexp.MustCompile("curl {2}-X GET 'http://127.0.0.1:.*/profile\\?.*' -H 'Authorization:Bear 123'")
			return v.MatchString(resp.Curl())
		})
		assert.Condition(t, func() (success bool) {
			var s baseServer
			resp.Unmarshal(&s)
			return assert.Equal(t, s, baseServer{Name: "jerry", Age: 18, Email: "jerry001@gmail.cn"})
		})
	}
}

func TestClient_Get(t *testing.T) {
	server := createBaseServer(t)
	defer server.Close()

	{
		resp := New(server.URL+"/profile").SetQueryValues(url.Values{"name": []string{"jerry"}}).SetHeader(hdrAuthorizationKey, "Bear 123").Get(context.Background())
		assert.JSONEq(t, resp.GetBodyString(), "{\"name\":\"jerry\",\"email\":\"jerry001@gmail.cn\",\"age\":18}")
		assert.Equal(t, resp.StatusCode(), 200)
		assert.NoError(t, resp.Error())
		assert.Condition(t, func() (success bool) {
			v := regexp.MustCompile("curl {2}-X GET 'http://127.0.0.1:.*/profile\\?.*' -H 'Authorization:Bear 123'")
			return v.MatchString(resp.Curl())
		})
		assert.Condition(t, func() (success bool) {
			var s baseServer
			resp.Unmarshal(&s)
			return assert.Equal(t, s, baseServer{Name: "jerry", Age: 18, Email: "jerry001@gmail.cn"})
		})
	}
}

func TestClient_Post(t *testing.T) {
	server := createBaseServer(t)
	defer server.Close()

	{
		var s = baseServer{
			Age:   18,
			Name:  "jerry",
			Email: "jerry001@gmail.cn",
		}
		resp := New(server.URL + "/update").SetBodyJson(s).Post(context.Background())
		assert.JSONEq(t, resp.GetBodyString(), "{\"name\":\"jerry\",\"email\":\"jerry001@gmail.cn\",\"age\":20}")
		assert.Equal(t, resp.StatusCode(), 200)
		assert.NoError(t, resp.Error())
		assert.Condition(t, func() (success bool) {
			v := regexp.MustCompile("curl {2}-X POST 'http://127.0.0.1:.*/update' -H 'Content-Type:application/json' -d ")
			return v.MatchString(resp.Curl())
		})
		assert.Condition(t, func() (success bool) {
			var s baseServer
			resp.Unmarshal(&s)
			return assert.Equal(t, s, baseServer{Name: "jerry", Age: 20, Email: "jerry001@gmail.cn"})
		})
	}
}

func createBaseServer(t *testing.T) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Method: %v", r.Method)
		t.Logf("Path: %v", r.URL.Path)
		t.Logf("header: %v", r.Header)
		if r.Method == http.MethodGet {
			if r.URL.Path == "/profile" {
				if r.URL.Query().Get("name") == "jerry" {
					w.Header().Set(hdrContentTypeKey, jsonContentType)
					_, _ = w.Write([]byte(`{"age":18,"email":"jerry001@gmail.cn","name":"jerry"}`))
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		} else if r.Method == http.MethodPost {
			if r.URL.Path == "/update" {
				b, _ := io.ReadAll(r.Body)
				var base baseServer
				err := json.Unmarshal(b, &base)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				base.Age = 20
				marshal, err := json.Marshal(base)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				_, _ = w.Write(marshal)
			}
		}
	})
	return httptest.NewServer(handler)
}

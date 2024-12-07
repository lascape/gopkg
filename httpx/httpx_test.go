package httpx

import (
	"context"
	"testing"
)

func TestGet(t *testing.T) {
	resp := New("http://www.baidu.com").Get(context.Background())
	t.Log(resp.GetBodyString())
}

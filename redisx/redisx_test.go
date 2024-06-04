package redisx

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestMust(t *testing.T) {
	cmdable := Must(WithConfig(Config{
		Addrs: []string{"127.0.0.1:6379"},
	}))

	key := fmt.Sprintf("%d", time.Now().Unix())
	cmdable.Set(context.Background(), key, "aaabbb123", time.Second*10)
	result, err := cmdable.Get(context.Background(), key).Result()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}

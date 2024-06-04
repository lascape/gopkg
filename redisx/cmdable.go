package redisx

import (
	"context"
	redislock "github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cmdable interface {
	redis.Cmdable
	Obtain(ctx context.Context, key string, ttl time.Duration) (*redislock.Lock, error)
}

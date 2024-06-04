package redisx

import (
	"context"
	"fmt"
	"github.com/bsm/redislock"
	_ "github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	Source   string   `yaml:"source" json:"source"`
	Addrs    []string `yaml:"addrs" json:"addrs"`
	Password string   `yaml:"password" json:"password" cipher:"true"`
	Db       int      `yaml:"db" json:"db"`
}

type Redis struct {
	redis.Cmdable
	lock  *redislock.Client
	conf  Config
	hooks []redis.Hook
}

type Option func(*Redis)

func WithHooks(hks ...redis.Hook) Option {
	return func(r *Redis) {
		for _, hk := range hks {
			if hk == nil {
				continue
			}
			r.hooks = append(r.hooks, hk)
		}
	}
}

func WithConfig(c Config) Option {
	return func(r *Redis) {
		r.conf = c
	}
}

func Must(opts ...Option) Cmdable {
	cmdable, err := mustRedis(opts...)
	if err != nil {
		panic("mustRedis " + err.Error())
	}
	return cmdable
}

func mustRedis(opts ...Option) (Cmdable, error) {
	r := new(Redis)
	for _, opt := range opts {
		opt(r)
	}

	var err error

	//Standard only
	if r.conf.Source != "" || len(r.conf.Addrs) == 1 {
		var redisOptions = new(redis.Options)
		if r.conf.Source != "" {
			redisOptions, err = redis.ParseURL(r.conf.Source)
			if err != nil {
				return nil, err
			}
		}
		if len(r.conf.Addrs) > 0 {
			redisOptions.Addr = r.conf.Addrs[0]
		}
		if len(r.conf.Password) > 0 {
			redisOptions.Password = r.conf.Password
		}
		if r.conf.Db > 0 {
			redisOptions.DB = r.conf.Db
		}
		client := redis.NewClient(redisOptions)
		for _, hook := range r.hooks {
			client.AddHook(hook)
		}
		if _, err := client.Ping(context.Background()).Result(); err != nil {
			return nil, fmt.Errorf("init redisx err: %v", err)
		}
		r.Cmdable = client
	} else {
		var redisOptions = new(redis.ClusterOptions)
		if len(r.conf.Addrs) > 0 {
			redisOptions.Addrs = r.conf.Addrs
		}
		if len(r.conf.Password) > 0 {
			redisOptions.Password = r.conf.Password
		}
		client := redis.NewClusterClient(redisOptions)
		for _, hook := range r.hooks {
			client.AddHook(hook)
		}
		if _, err := client.Ping(context.Background()).Result(); err != nil {
			return nil, fmt.Errorf("init redisx err: %v", err)
		}
		r.Cmdable = client
	}

	return r, nil
}

func (r *Redis) Obtain(ctx context.Context, key string, ttl time.Duration) (*redislock.Lock, error) {
	if r.lock == nil {
		r.lock = redislock.New(r.Cmdable)
	}
	return r.lock.Obtain(ctx, key, ttl, nil)
}

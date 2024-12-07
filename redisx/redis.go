package redisx

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type envHook struct {
}

func (e envHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return nil, nil
}

func (e envHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return nil
}

func (e envHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (e envHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

type Config struct {
	Source   string `yaml:"source" json:"source"`
	Addr     string `yaml:"addr" json:"addr"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
}

type Options struct {
	conf  Config
	debug bool
}

type Option func(*Options)

func WithConfig(conf Config) Option {
	return func(o *Options) {
		o.conf = conf
	}
}

func Must(opts ...Option) *redis.Client {
	client, err := must(opts...)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return client
}

func must(opts ...Option) (*redis.Client, error) {
	var (
		o            Options
		redisOptions = &redis.Options{}
		err          error
	)
	for _, opt := range opts {
		opt(&o)
	}
	if o.conf.Source != "" {
		redisOptions, err = redis.ParseURL(o.conf.Source)
	} else {
		redisOptions.Addr = o.conf.Addr
		redisOptions.Password = o.conf.Password
	}
	redisOptions.DB = o.conf.DB

	client := redis.NewClient(redisOptions)
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()
	err = client.Ping(timeout).Err()
	return client, err
}

package server

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Option func(o *Options)

type Options struct {
	sigs      []os.Signal
	servers   []Server
	delayTime time.Duration
	conf      Config
}

type Config struct {
	ServiceName string `yaml:"service_name"`
	Env         string `yaml:"env"`
}

func Run(opts ...Option) {
	options := Options{
		sigs:      []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		delayTime: 10,
	}

	for _, o := range opts {
		o(&options)
	}
	options.run()
}

func WithServer(srv ...Server) Option {
	return func(o *Options) { o.servers = append(o.servers, srv...) }
}

func (o *Options) run() {
	defer func() {
		if e := recover(); e != nil {
			logrus.Panic(e)
		}
	}()
	ctx, cancelFunc := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	for _, srv := range o.servers {
		func(srv Server) {
			eg.Go(func() error {
				<-ctx.Done() // wait for stop signal
				_ctx, _cancelFunc := context.WithTimeout(context.Background(), o.delayTime*time.Second)
				defer _cancelFunc()
				return srv.Shutdown(_ctx)
			})
			go func() {
				defer func() {
					if e := recover(); e != nil {
						logrus.Panic(e)
					}
				}()
				srv.Start()
			}()
		}(srv)
	}
	sig := <-Signal()
	logrus.Infof("the received signal is %s", sig.String())
	cancelFunc()

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		logrus.Fatal(fmt.Sprintf("service failed to run: %+v", err))
	}
}

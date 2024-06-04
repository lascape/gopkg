package serverx

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Interface interface {
	Start()
	Shutdown(ctx context.Context) error
}

type Config struct {
	ServiceName string `json:"service_name" yaml:"service_name"`
	Env         string `json:"env" yaml:"env"`
	Deadline    int    `json:"deadline" yaml:"deadline"`
}

type Server struct {
	Deadline int `json:"deadline" yaml:"deadline"`
	servers  []Interface
}

type Option func(s *Server)

func WithServer(svrs ...Interface) Option {
	return func(s *Server) {
		for _, srv := range svrs {
			if srv == nil {
				continue
			}
			s.servers = append(s.servers, srv)
		}
	}
}

func WithDeadline(deadline int) Option {
	return func(s *Server) {
		s.Deadline = deadline
	}
}

func Must(opts ...Option) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Run() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	wgShutdown := &sync.WaitGroup{}
	for _, server := range s.servers {
		wgShutdown.Add(1)
		go s.Shutdown(server.Shutdown, ctx, wgShutdown)
		go s.Start(server.Start)
	}
	sig := <-Signal()
	logrus.WithField("signal", sig).Info("Server.Run shutting down")
	cancelFunc()
	wgShutdown.Wait()
}

func (s *Server) Start(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("error", err).Error("Server.Start panic")
		}
	}()
	fn()
}

func (s *Server) Shutdown(fn func(ctx context.Context) error, ctx context.Context, wg *sync.WaitGroup) {
	<-ctx.Done()
	logrus.Info("Server.Shutdown called")
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Duration(s.Deadline)*time.Second)
	defer cancelFunc()
	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("error", err).Error("Server.Shutdown panic")
		}
		wg.Done()
	}()
	if err := fn(timeout); err != nil {
		logrus.WithField("error", err).Error("Server.Shutdown error")
	}
}

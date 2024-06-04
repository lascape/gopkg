package ginserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/lascape/gopkg/internal/serverx/ginserver"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Config struct {
	Addr           string `json:"addr" yaml:"addr"`
	ReadTimeout    int    `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout   int    `json:"write_timeout" yaml:"write_timeout"`
	IdleTimeout    int    `json:"idle_timeout" yaml:"idle_timeout"`
	MaxHeaderBytes int    `json:"max_header_bytes" yaml:"max_header_bytes"`
	Mode           string `json:"mode" yaml:"mode"`
}

type Option func(server *Server)

type Server struct {
	conf   Config
	server *http.Server
	Engine *gin.Engine

	checkRelativePath string
	checkHandler      gin.HandlerFunc
	recoveryHandler   gin.HandlerFunc
	corsHandler       gin.HandlerFunc
	loggerHandler     gin.HandlerFunc
}

func (s *Server) Start() {
	logrus.Info("start gin server")
	if err := s.server.ListenAndServe(); err != nil {
		logrus.WithField("error", err).Error("start gin server err")
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func New(f func(server *Server), opts ...Option) *Server {
	server := newServer(opts...)

	f(server)

	return server
}

func newServer(opts ...Option) *Server {
	server := &Server{server: &http.Server{}, checkRelativePath: "/ping"}
	for _, opt := range opts {
		opt(server)
	}
	{ //initialize server configuration
		if server.conf.Addr == "" {
			panic("gin server addr is empty")
		}
		server.server.Addr = server.conf.Addr

		if server.conf.ReadTimeout > 0 {
			server.server.ReadTimeout = time.Duration(server.conf.ReadTimeout) * time.Millisecond
		}

		if server.conf.WriteTimeout > 0 {
			server.server.WriteTimeout = time.Duration(server.conf.WriteTimeout) * time.Millisecond
		}

		if server.conf.IdleTimeout > 0 {
			server.server.IdleTimeout = time.Duration(server.conf.IdleTimeout) * time.Millisecond
		}

		if server.conf.MaxHeaderBytes > 0 {
			server.server.MaxHeaderBytes = server.conf.MaxHeaderBytes
		}
	}

	if server.conf.Mode != "" {
		gin.SetMode(server.conf.Mode)
	}

	server.Engine = gin.Default()
	if server.recoveryHandler != nil {
		server.Engine.Use(server.recoveryHandler)
	} else {
		server.Engine.Use(ginserver.Recovery())
	}

	if server.loggerHandler != nil {
		server.Engine.Use(server.loggerHandler)
	} else {
		server.Engine.Use(ginserver.Logger())
	}

	if server.checkHandler != nil {
		server.Engine.GET(server.checkRelativePath, server.checkHandler)
	} else {
		server.Engine.GET(server.checkRelativePath, ginserver.HealthyCheck())
	}

	if server.corsHandler != nil {
		server.Engine.Use(server.corsHandler)
	}

	server.server.Handler = server.Engine
	return server
}

func WithConfig(config Config) Option {
	return func(server *Server) {
		server.conf = config
	}
}

func WithAddr(addr string) Option {
	return func(server *Server) {
		server.conf.Addr = addr
	}
}

func WithHealthyCheck(relativePath string, handlerFunc gin.HandlerFunc) Option {
	return func(server *Server) {
		if relativePath == "" {
			relativePath = "/ping"
		}
		server.checkRelativePath = relativePath
		server.checkHandler = handlerFunc
	}
}

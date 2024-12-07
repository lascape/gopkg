package httpServer

import (
	"github.com/lascape/gopkg/server/http/mid"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lascape/gopkg/response"
)

type Option func(c *Server)

func WithConfig(c Config) Option {
	return func(s *Server) {
		s.conf = c
	}
}

type Config struct {
	Name         string `yaml:"name"`
	Addr         string `yaml:"addr"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

func NewServer(f func(*Server), opts ...Option) *Server {
	server := newServer()
	for _, opt := range opts {
		opt(server)
	}
	server.server = &http.Server{
		Addr:           server.conf.Addr,
		ReadTimeout:    time.Duration(server.conf.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(server.conf.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.Engine = gin.New()
	server.Use(gin.Recovery())
	server.Use(mid.Cors())
	server.Use(mid.XRealIp)
	server.Use(mid.Logger())
	f(server)
	for _, register := range registers {
		register(server.Engine)
	}
	addCheckStatus(server.Engine)
	server.server.Handler = server.Engine
	return server
}

func addCheckStatus(engine *gin.Engine) *gin.Engine {
	engine.GET("/check_status", func(ctx *gin.Context) {
		response.Success(ctx, nil)
	})
	return engine
}

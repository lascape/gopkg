package httpServer

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	server *http.Server
	*gin.Engine
	conf Config
}

func newServer() *Server {
	return &Server{}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Use(middleware ...gin.HandlerFunc) {
	s.Engine.Use(middleware...)
}

func (s *Server) Start() {
	if err := s.server.ListenAndServe(); err != nil {
		logrus.Panic(err)
	}
}

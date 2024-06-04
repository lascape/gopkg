package serverx

import (
	"context"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/lascape/gopkg/serverx/ginserver"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestMust(t *testing.T) {
	Must(
		WithServer(nil),
		WithServer(&serverTimestamp{}),
		WithServer(&serverPanic{}),
		WithDeadline(3),
	).Run()
}

type serverTimestamp struct{}

func (*serverTimestamp) Start() {
	for {
		logrus.WithField("timestamp", time.Now().Unix()).Info("serverTimestamp.Start")
		time.Sleep(time.Second)
	}
}

func (*serverTimestamp) Shutdown(ctx context.Context) error {
	logrus.WithField("msg", "shutdown after 5 second").Info("serverTimestamp.Shutdown")
	now := time.Now()
	<-ctx.Done()
	logrus.
		WithField("msg", "shutdown will be done").
		WithField("duration", time.Now().Sub(now).Seconds()).
		Info("serverTimestamp.Shutdown")
	return nil
}

type serverPanic struct{}

func (*serverPanic) Start() {
	panic("serverPanic")
}

func (*serverPanic) Shutdown(ctx context.Context) error {
	return nil
}

func TestWithServer(t *testing.T) {
	server := ginserver.New(func(s *ginserver.Server) {
		s.Engine.GET("/", func(c *gin.Context) { c.JSON(200, "ok") })
	}, ginserver.WithAddr("0.0.0.0:8210"))
	Must(WithServer(server)).Run()
}

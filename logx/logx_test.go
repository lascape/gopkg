package logx

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestMust(t *testing.T) {
	Must(WithConfig(Config{Filename: "./logx.log"}))
	logrus.Info("hello world<a>123</a>")
	logrus.WithField("bb", "yy").Info("hello world<a>123</a>")
	logrus.Warn("hello world")
	logrus.Error("hello world")
}

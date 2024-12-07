package logx

import (
	"encoding/json"
	"github.com/lascape/gopkg/envx"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLog(t *testing.T) {
	Init("path")
	logrus.SetLevel(envx.PKG_LOGX_LEVEL)
	logrus.Info("123")
	logrus.Warn("123")
	logrus.Error("123")
	select {}
}

func TestLogstash(t *testing.T) {
	Init("api")

	type XX struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	x := XX{
		Name: "zgl",
		Age:  18,
	}
	marshal, _ := json.Marshal(x)
	logrus.Info(string(marshal))
}

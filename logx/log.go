package logx

import (
	"github.com/lascape/gopkg/envx"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var writer io.Writer = os.Stdout

func Init(srvName string) {
	logFile := &lumberjack.Logger{
		Filename:   "logs/app.log", // 日志文件名
		MaxSize:    50,             // 每个日志文件的最大尺寸，单位MB
		MaxBackups: 3,              // 保留的旧日志文件最大数量
		MaxAge:     37,             // 保留的旧日志文件最大天数
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(envx.PKG_LOGX_LEVEL)
	logrus.AddHook(&Caller{})
	logrus.AddHook(NewLogstash(srvName))
	logrus.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

func Writer() io.Writer {
	return writer
}

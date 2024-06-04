package logx

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	_ "gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"runtime"
	"strconv"
	"time"
)

type Config struct {
	Filename        string `yaml:"filename" json:"filename"`         //文件名
	MaxSize         int    `yaml:"max_size" json:"max_size"`         // 单个日志文件的最大尺寸（以MB为单位）
	MaxBackups      int    `yaml:"max_backups" json:"max_backups"`   // 保留的旧日志文件的最大数量
	MaxAge          int    `yaml:"max_age" json:"max_age"`           // 保留旧日志文件的最大天数
	Compress        bool   `yaml:"compress" json:"compress"`         // 是否压缩旧日志文件
	CloseStd        bool   `yaml:"close_std" json:"close_std"`       //关闭控制台打印
	CloseCaller     bool   `yaml:"close_caller" json:"close_caller"` //关闭堆栈打印
	TimestampFormat string `yaml:"timestamp_format" json:"timestamp_format"`
}

type Log struct {
	conf Config
}

type Option func(l *Log)

func WithConfig(config Config) Option {
	return func(l *Log) {
		l.conf = config
	}
}

func Must(opts ...Option) {
	l := &Log{}
	for _, opt := range opts {
		opt(l)
	}

	if !l.conf.CloseCaller {
		logrus.SetReportCaller(true)
	}

	if l.conf.TimestampFormat == "" {
		l.conf.TimestampFormat = time.RFC3339
	}
	//日志格式
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   l.conf.TimestampFormat,
		DisableHTMLEscape: true,
		DataKey:           "fields",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "func",
		},
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, " " + frame.File + ":" + strconv.Itoa(frame.Line) + " "
		},
	})

	//文件写入+分割
	var ios []io.Writer
	if l.conf.Filename != "" {
		ios = append(ios, &lumberjack.Logger{
			Filename:   l.conf.Filename,   // 日志文件路径
			MaxSize:    l.conf.MaxSize,    // 单个日志文件的最大尺寸（以MB为单位）
			MaxBackups: l.conf.MaxBackups, // 保留的旧日志文件的最大数量
			MaxAge:     l.conf.MaxAge,     // 保留旧日志文件的最大天数
			Compress:   l.conf.Compress,   // 是否压缩旧日志文件
		})
	}
	if !l.conf.CloseStd {
		ios = append(ios, os.Stdout)
	}
	logrus.SetOutput(io.MultiWriter(ios...))
}

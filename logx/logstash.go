package logx

import (
	"github.com/sirupsen/logrus"
)

type Logstash struct {
	AppName string
}

func (l Logstash) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (l Logstash) Fire(entry *logrus.Entry) error {
	entry.Data["app_name"] = l.AppName
	return nil
}

func NewLogstash(name string) *Logstash {
	return &Logstash{AppName: name}
}

package envx

import "github.com/sirupsen/logrus"

var (
	PKG_HTTPX_BCURL          = ValueByEnv("PKG_HTTPX_BCURL", "true").Bool()
	PKG_HTTPX_ENGRESS        = ValueByEnv("PKG_HTTPX_ENGRESS", "").String()
	PKG_HTTPX_ENGRESS_IGNORE = ValueByEnv("PKG_HTTPX_ENGRESS_IGNORE", "").String()
	PKG_LOGX_LEVEL, _        = logrus.ParseLevel(ValueByEnv("PKG_LOGX_LEVEL", "debug").String())
)

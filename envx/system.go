package envx

import "strconv"
import "github.com/sirupsen/logrus"

var (
	PKG_HTTPX_BCURL, _       = strconv.ParseBool(ValueByEnv("PKG_HTTPX_BCURL", "false"))
	PKG_HTTPX_ENGRESS        = ValueByEnv("PKG_HTTPX_ENGRESS", "")
	PKG_HTTPX_ENGRESS_IGNORE = ValueByEnv("PKG_HTTPX_ENGRESS_IGNORE", "")
	PKG_LOGX_LEVEL, _        = logrus.ParseLevel(ValueByEnv("PKG_LOGX_LEVEL", "debug"))
)

package ginserver

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"runtime"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				buf := make([]byte, 2048)
				buf = buf[:runtime.Stack(buf, true)]
				logrus.WithFields(logrus.Fields{
					"error": e,
					"stack": string(buf),
				}).Panic("gin server panic is recovery")
			}
		}()
		ctx.Next()
	}
}

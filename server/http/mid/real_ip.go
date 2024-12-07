package mid

import (
	"context"
	"github.com/gin-gonic/gin"
)

const (
	ctxWithRealIp = "ctx-with-source-ip"
)

func ValueRealIP(ctx context.Context) string {
	if v := ctx.Value(ctxWithRealIp); v != nil {
		return v.(string)
	}
	return ""
}

func XRealIp(ctx *gin.Context) {
	ip := ctx.GetHeader("X-Forwarded-For")
	if ip == "" { //如果IP为空，则默认执行以下操作
		ip = "127.0.0.1"
	}
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), ctxWithRealIp, ip))
	ctx.Set(ctxWithRealIp, ip)
	ctx.Next()
}

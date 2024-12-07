package mid

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Entry struct {
	Path       string        `json:"path"`
	Method     string        `json:"method"`
	Body       string        `json:"body"`
	Latency    time.Duration `json:"latency"`
	ClientIP   string        `json:"client_ip"`
	StatusCode int           `json:"status_code"`
}

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if strings.HasSuffix(path, "check_status") {
			ctx.Next()
			return
		}
		start := time.Now()
		raw := ctx.Request.URL.RawQuery
		entry := Entry{
			Path:     ctx.Request.URL.Path,
			Method:   ctx.Request.Method,
			ClientIP: ValueRealIP(ctx),
		}
		bodyBytes, _ := io.ReadAll(ctx.Request.Body)
		_ = ctx.Request.Body.Close() //  must close
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		entry.Body = string(bodyBytes)

		{ //注入traceId
			traceId := uuid.NewString()
			ctx.Set("trace_id", traceId)
			ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "trace_id", traceId))
		}

		ctx.Next()

		entry.Latency = time.Now().Sub(start)
		entry.StatusCode = ctx.Writer.Status()
		if raw != "" {
			entry.Path = entry.Path + "?" + raw
		}
		if !strings.HasPrefix(ctx.Request.URL.Path, "/swagger") {
			marshal, _ := json.Marshal(entry)
			logrus.WithContext(ctx.Request.Context()).Infof("Gin Request:%+v", string(marshal))
		}
	}
}

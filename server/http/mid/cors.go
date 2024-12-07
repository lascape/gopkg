package mid

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(r *gin.Context) {
		origin := r.Request.Header.Get("Origin")

		// 增强判断：只有域名匹配到后缀为以下的才允许跨域，file://这个属于特殊情况(匹配到前缀)
		r.Header("Access-Control-Allow-Origin", origin)
		r.Header("Access-Control-Allow-Credentials", "true")
		r.Header("Access-Control-Allow-SetHeaders", "COOKIE,DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,appid,accept-language")
		r.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,HEAD")

		if r.Request.Method == "OPTIONS" {
			r.AbortWithStatus(http.StatusNoContent)
		}

		r.Next()
	}
}

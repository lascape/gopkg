package ginserver

import "github.com/gin-gonic/gin"

func HealthyCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, "ok")
	}
}

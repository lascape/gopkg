package httpServer

import "github.com/gin-gonic/gin"

func AddRouter(register Register) {
	registers = append(registers, register)
}

var registers []Register

type Register func(engine *gin.Engine)

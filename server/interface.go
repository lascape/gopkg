package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

//Server 各个服务，安全退出的接口定义
type Server interface {
	Start()
	Shutdown(ctx context.Context) error
}

//Signal 全局退出信号接收
func Signal() <-chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return quit
}

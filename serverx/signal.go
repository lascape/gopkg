package serverx

import (
	"os"
	"os/signal"
	"syscall"
)

func Signal() <-chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	return quit
}

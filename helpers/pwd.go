package helpers

import "os"

func Pwd() string {
	wd, _ := os.Getwd()
	return wd
}

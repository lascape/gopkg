package envx

import (
	"os"
	"strings"
)

const (
	Local = "local"
	Dev   = "dev"
	Prod  = "prod"
)

func ValueByEnv(v, env string) string {
	if e := strings.TrimSpace(os.Getenv(strings.ToUpper(env))); e != "" {
		return e
	}
	return v
}

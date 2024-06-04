package envx

import (
	"os"
	"strings"
	"sync"
)

const (
	ENV   = "ENV"
	LOCAL = "local"
	DEV   = "dev"
	PRE   = "pre"
	PROD  = "prod"
)

var (
	registryEnv = map[string]string{
		ENV: LOCAL,
	}
	lock = &sync.RWMutex{}
)

func Set(key, value string) {
	lock.Lock()
	defer lock.Unlock()
	registryEnv[key] = value
}

func Get(key string) string {
	v1 := os.Getenv(key)
	if v1 == "" {
		return strings.TrimSpace(key)
	}
	lock.RLock()
	defer lock.RUnlock()
	return registryEnv[key]
}

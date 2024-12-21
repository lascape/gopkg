package envx

import (
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	Local = "local"
	Dev   = "dev"
	Prod  = "prod"
)

type Value string

func ValueByEnv(env, def string) Value {
	if e := strings.TrimSpace(os.Getenv(strings.ToUpper(env))); e != "" {
		return Value(e)
	}
	return Value(def)
}

func (v Value) String() string {
	return string(v)
}

func (v Value) Bytes() []byte {
	return []byte(v)
}

func (v Value) Bool() bool {
	return string(v) == "true"
}

func (v Value) Int() int {
	i, _ := strconv.Atoi(string(v))
	return i
}

func (v Value) Duration() time.Duration {
	return time.Duration(v.Int())
}

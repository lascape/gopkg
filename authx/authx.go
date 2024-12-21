package authx

import (
	"sync"
	"time"
)

const (
	AuthNameJwt = "jwt"
)

func init() {
	RegisterAuth(AuthNameJwt, &JWT{})
}

type manager struct {
	m    map[string]Auth
	lock *sync.Mutex
}

var defaultManager = &manager{
	m:    make(map[string]Auth),
	lock: new(sync.Mutex),
}

// Auth 是一个认证接口，可以支持 Session 或 JWT
type Auth interface {
	// GenerateToken 用于生成认证 Token
	GenerateToken(kv map[string]interface{}, expired time.Duration) (string, error)

	// ValidateToken 用于验证 Token，并返回用户 ID
	ValidateToken(token string) (map[string]interface{}, error)
}

func RegisterAuth(name string, auth Auth) {
	defaultManager.lock.Lock()
	defer defaultManager.lock.Unlock()
	if _, ok := defaultManager.m[name]; ok {
		panic("auth already registered")
	}
	defaultManager.m[name] = auth
}

func GetAuth(name string) Auth {
	defaultManager.lock.Lock()
	defer defaultManager.lock.Unlock()

	auth, ok := defaultManager.m[name]
	if !ok {
		return &defaultAuth{}
	}
	return auth
}

type defaultAuth struct{}

func (*defaultAuth) GenerateToken(_ map[string]interface{}, _ time.Duration) (string, error) {
	return "", nil
}

func (*defaultAuth) ValidateToken(_ string) (map[string]interface{}, error) {
	return make(map[string]interface{}), nil
}

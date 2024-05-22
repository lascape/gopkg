package cryptox

import (
	"sync"
)

const (
	// ErrCipherNotFound is returned when there is no cipher registered for a name.
	ErrCipherNotFound = cipherError("cipher not found for this name")

	// ErrCipherNameAlreadyRegistered is returned when an encoder is already registered for a format.
	ErrCipherNameAlreadyRegistered = cipherError("cipher already registered for this name")
)

type Cipher interface {
	Encrypt(key, text string) (string, error)
	Decrypt(key, text string) (string, error)
}

// CipherRegistry can choose an appropriate Cipher based on the provided format.
type CipherRegistry struct {
	ciphers map[string]Cipher

	mu sync.RWMutex
}

// NewCipherRegistry returns a new, initialized CipherRegistry.
func NewCipherRegistry() *CipherRegistry {
	return &CipherRegistry{
		ciphers: make(map[string]Cipher),
	}
}

// RegisterCipher registers a Cipher for a name.
// Registering a Cipher for an already existing name is not supported.
func (e *CipherRegistry) RegisterCipher(name string, enc Cipher) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.ciphers[name]; ok {
		return ErrCipherNameAlreadyRegistered
	}

	e.ciphers[name] = enc

	return nil
}

func (e *CipherRegistry) Encrypt(method, key, text string) (string, error) {
	e.mu.RLock()
	ciph, ok := e.ciphers[method]
	e.mu.RUnlock()
	if !ok {
		return "", ErrCipherNotFound
	}
	return ciph.Encrypt(key, text)
}

func (e *CipherRegistry) Decrypt(method, key, text string) (string, error) {
	e.mu.RLock()
	ciph, ok := e.ciphers[method]
	e.mu.RUnlock()
	if !ok {
		return "", ErrCipherNotFound
	}
	return ciph.Decrypt(key, text)
}

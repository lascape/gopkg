package cryptox

import "github.com/lascape/gopkg/cryptox/aes"

type Crypto struct {
	cipherRegistry *CipherRegistry
}

func NewCrypto() *Crypto {
	c := new(Crypto)
	c.cipherRegistry = NewCipherRegistry()
	_ = c.cipherRegistry.RegisterCipher("aes", new(aes.Cipher))
	return c
}

// Encrypt a text for a key.
func (e *Crypto) Encrypt(name, key, text string) (string, error) {
	e.cipherRegistry.mu.RLock()
	encoder, ok := e.cipherRegistry.ciphers[name]
	e.cipherRegistry.mu.RUnlock()

	if !ok {
		return "", ErrCipherNotFound
	}

	return encoder.Encrypt(key, text)
}

// Decrypt a text for a key
func (e *Crypto) Decrypt(name, key, text string) (string, error) {
	e.cipherRegistry.mu.RLock()
	encoder, ok := e.cipherRegistry.ciphers[name]
	e.cipherRegistry.mu.RUnlock()

	if !ok {
		return "", ErrCipherNotFound
	}
	return encoder.Decrypt(key, text)
}

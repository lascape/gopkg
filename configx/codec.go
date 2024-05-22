package configx

import "github.com/lascape/gopkg/cryptox"

type Codec interface {
	cryptox.Cipher
	Unmarshal(data []byte, v interface{}) error
}

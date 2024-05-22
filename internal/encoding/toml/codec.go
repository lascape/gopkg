package toml

import (
	"github.com/pelletier/go-toml/v2"
)

// Codec implements the encoding.Encoder and encoding.Decoder interfaces for TOML encoding.
type Codec struct{}

func (Codec) Encode(v any) ([]byte, error) {
	return toml.Marshal(v)
}

func (Codec) Decode(b []byte, v any) error {
	return toml.Unmarshal(b, v)
}

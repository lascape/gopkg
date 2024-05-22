package json

import (
	"encoding/json"
)

// Codec implements the encoding.Encoder and encoding.Decoder interfaces for JSON encoding.
type Codec struct{}

func (Codec) Encode(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func (Codec) Decode(b []byte, v any) error {
	return json.Unmarshal(b, v)
}

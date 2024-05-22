package yaml

import "gopkg.in/yaml.v3"

// Codec implements the encoding.Encoder and encoding.Decoder interfaces for YAML encoding.
type Codec struct{}

func (Codec) Encode(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (Codec) Decode(b []byte, v any) error {
	return yaml.Unmarshal(b, v)
}

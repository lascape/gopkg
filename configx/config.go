package configx

import (
	"bytes"
	"github.com/lascape/gopkg/cryptox"
	"github.com/lascape/gopkg/cryptox/aes"
	"github.com/lascape/gopkg/internal/encoding"
	"github.com/lascape/gopkg/internal/encoding/json"
	"github.com/lascape/gopkg/internal/encoding/toml"
	"github.com/lascape/gopkg/internal/encoding/yaml"
	"io"
	"reflect"
)

type ConfigLoader interface {
	BeforeLoad() error
	AfterLoad() error
}

type Config struct {
	v               any
	cipherKey       string
	cipherType      string
	cipherRegistry  *cryptox.CipherRegistry
	decoderType     string
	decoderRegistry *encoding.DecoderRegistry
}

type Option func(c *Config)

func WithCipher(name string, key string) Option {
	return func(c *Config) {
		c.cipherType = name
		c.cipherKey = key
	}
}

func WithDecoder(name string) Option {
	return func(c *Config) {
		c.decoderType = name
	}
}

func Must(v any, opts ...Option) *Config {
	if v == nil {
		panic("v must not be nil")
	}
	var c = &Config{
		v: v,
	}
	for _, opt := range opts {
		opt(c)
	}

	c.resetDecoderRegistry()
	c.resetCipherRegistry()

	return c
}

func (c *Config) ReadFile(reader io.Reader) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	if err := c.decoderRegistry.Decode(c.decoderType, buf.Bytes(), c.v); err != nil {
		return err
	}
	if err := c.decryptStructFields(c.v); err != nil {
		return err
	}
	return nil
}

func (c *Config) resetDecoderRegistry() {
	if c.decoderType == "" {
		panic("decoder type must be specified")
	}
	c.decoderRegistry = encoding.NewDecoderRegistry()
	{
		codec := yaml.Codec{}
		c.decoderRegistry.RegisterDecoder("yaml", codec)
		c.decoderRegistry.RegisterDecoder("yml", codec)
	}
	{
		codec := json.Codec{}
		c.decoderRegistry.RegisterDecoder("json", codec)
	}
	{
		codec := toml.Codec{}
		c.decoderRegistry.RegisterDecoder("toml", codec)
	}
}

func (c *Config) resetCipherRegistry() {
	if c.cipherType == "" {
		return
	}
	c.cipherRegistry = cryptox.NewCipherRegistry()
	{
		cipher := aes.Cipher{}
		c.cipherRegistry.RegisterCipher("aes", cipher)
	}
}

// decryptStructFields 递归遍历结构体并解密带有 tag 为 cipher=true 的字符串字段
func (c *Config) decryptStructFields(s interface{}) error {
	if c.cipherType == "" {
		return nil
	}
	val := reflect.ValueOf(s)
	return c.decodeString(nil, val)
}

func (c *Config) isRecursionField(k reflect.Kind) bool {
	return k == reflect.Slice ||
		k == reflect.Array ||
		k == reflect.Struct ||
		k == reflect.Map ||
		k == reflect.Interface ||
		k == reflect.Pointer
}

func (c *Config) decodeString(st *reflect.StructField, v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	k := v.Kind()

	switch k {
	case reflect.Map:
		m := v.MapRange()
		for m.Next() {
			elem := reflect.New(m.Value().Type()).Elem()
			elem.Set(m.Value())
			err := c.decodeString(nil, elem)
			if err != nil {
				return err
			}
			v.SetMapIndex(m.Key(), elem)
		}

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			k := f.Kind()
			tf := v.Type().Field(i)
			if c.isRecursionField(k) {
				err := c.decodeString(&tf, f)
				if err != nil {
					return err
				}
			}
			if k == reflect.String {
				if err := c.decodeCipher(tf, f); err != nil {
					return err
				}
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			iv := v.Index(i)
			k := iv.Kind()
			if c.isRecursionField(k) {
				err := c.decodeString(nil, iv)
				if err != nil {
					return err
				}
			}
		}
	case reflect.Interface:

	case reflect.Pointer:
		if v.IsZero() {
			return nil
		}

		ele := v.Elem()
		k := ele.Kind()

		if k == reflect.String {
			if err := c.decodeCipher(*st, ele); err != nil {
				return err
			}
		}

		if !c.isRecursionField(k) {
			return nil
		}

		var num = ele.NumField()
		for i := 0; i < num; i++ {
			val := ele.Field(i)
			k := val.Kind()
			tf := ele.Type().Field(i)
			if c.isRecursionField(k) {
				err := c.decodeString(&tf, val)
				if err != nil {
					return err
				}
			}
			if k == reflect.String {
				if err := c.decodeCipher(tf, val); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *Config) decodeCipher(st reflect.StructField, val reflect.Value) error {
	if tag := st.Tag.Get("cipher"); tag != "true" {
		return nil
	}
	if !val.CanSet() {
		return nil
	}
	plaintext, err := c.cipherRegistry.Decrypt(c.cipherType, c.cipherKey, val.String())
	if err != nil {
		return err
	}
	val.SetString(plaintext)
	return nil
}

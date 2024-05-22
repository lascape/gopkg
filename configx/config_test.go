package configx

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestMust(t *testing.T) {
	// 示例结构体
	type NestedStruct struct {
		InnerField string `cipher:"true"`
	}

	type MyStruct struct {
		Field1 string `cipher:"true"`
		Field2 string
		Nested NestedStruct
		Map    map[string]NestedStruct
		Slice  []NestedStruct
		Array  [2]NestedStruct
	}
	data := &MyStruct{
		Field1: "198963fd4f0dd0e3dec4ae2bbada079e10f9b3ce", // "hello" 加密后的形式（每个字符的ASCII码值加一）
		Field2: "world",
		Nested: NestedStruct{
			InnerField: "198963fd4f0dd0e3dec4ae2bbada079e10f9b3ce", // "they" 加密后的形式（每个字符的ASCII码值加一）
		},
		Map: map[string]NestedStruct{
			"key1": {InnerField: "198963fd4f0dd0e3dec4ae2bbada079e10f9b3ce"}, // "way" 加密后的形式（每个字符的ASCII码值加一）
		},
		Slice: []NestedStruct{
			{InnerField: "198963fd4f0dd0e3dec4ae2bbada079e10f9b3ce"}, // "apple" 加密后的形式（每个字符的ASCII码值加一）
		},
		Array: [2]NestedStruct{
			{InnerField: "198963fd4f0dd0e3dec4ae2bbada079e10f9b3ce"}, // "light" 加密后的形式（每个字符的ASCII码值加一）
			{InnerField: "198963fd4f0dd0e3dec4ae2bbada079e10f9b3ce"}, // "food" 加密后的形式（每个字符的ASCII码值加一）
		},
	}

	dataRealy := &MyStruct{
		Field1: "root",
		Field2: "world",
		Nested: NestedStruct{
			InnerField: "root",
		},
		Map: map[string]NestedStruct{
			"key1": {InnerField: "root"},
		},
		Slice: []NestedStruct{
			{InnerField: "root"},
		},
		Array: [2]NestedStruct{
			{InnerField: "root"},
			{InnerField: "root"},
		},
	}

	marshal, _ := json.Marshal(&data)
	reader := bytes.NewReader(marshal)
	config := Must(data, WithCipher("aes", "0d9cc89ab6a9ff1e"), WithDecoder("json"))
	err := config.ReadFile(reader)
	require.NoError(t, err)
	if !reflect.DeepEqual(data, dataRealy) {
		require.Error(t, errors.New("parse  failure"))
	}
}

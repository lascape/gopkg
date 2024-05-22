package cryptox

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewCryptoForAES(t *testing.T) {
	tests := []struct {
		Key  string
		Text string
	}{
		{
			Key:  "0d9cc89ab6a9ff1e",
			Text: "root",
		},
		{
			Key:  "0d9cc14ab6a9ab1e",
			Text: "123456",
		},
	}
	crypto := NewCrypto()

	for _, test := range tests {
		encrypt, err := crypto.Encrypt("aes", test.Key, test.Text)
		require.NoError(t, err)
		decrypt, err := crypto.Decrypt("aes", test.Key, encrypt)
		require.NoError(t, err)
		require.Equal(t, test.Text, decrypt)
	}
}

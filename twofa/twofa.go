package twofa

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

func Verify2faCode(secret, code string) bool {
	add := time.Now().Add(time.Second * 30 * 11)
	for i := 0; i < 21; i++ {
		if Generate(secret, add).Code == code {
			return true
		}
		add = add.Add(-time.Second * 30)
	}
	return false
}

func RandomSecret() string {
	uid := strings.Replace(uuid.New().String(), "-", "", -1)[0:10]
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(uid))
}

// Generate generate authenticator code
func Generate(secret string, now time.Time) Authenticator {
	secret = strings.Replace(secret, " ", "", -1)
	secret = strings.ToUpper(secret)

	t := now.Unix()

	count := uint64(t) / 30
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return Authenticator{Error: err}
	}

	codeInt := hotp(key, count, 6)
	code := codeAlignment(codeInt)

	return Authenticator{
		Secret: secret,
		Expire: int(30 - (t % 30)),
		Code:   code,
	}
}

func codeAlignment(code int) string {
	intFormat := fmt.Sprintf("%%0%dd", 6)
	codeStr := fmt.Sprintf(intFormat, code)
	return codeStr
}

func hotp(key []byte, counter uint64, digits int) int {
	h := hmac.New(sha1.New, key)
	binary.Write(h, binary.BigEndian, counter)
	sum := h.Sum(nil)
	v := binary.BigEndian.Uint32(sum[sum[len(sum)-1]&0x0F:]) & 0x7FFFFFFF
	d := uint32(1)
	for i := 0; i < digits && i < 8; i++ {
		d *= 10
	}
	return int(v % d)
}

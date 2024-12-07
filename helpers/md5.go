package helpers

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(src string) string {
	h := md5.New()
	if _, err := h.Write([]byte(src)); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

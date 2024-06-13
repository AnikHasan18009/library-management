package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func Sha1Hash(s string) string {

	h := sha1.New()
	h.Write([]byte(s))
	hashValue := h.Sum(nil)
	return hex.EncodeToString(hashValue)
}

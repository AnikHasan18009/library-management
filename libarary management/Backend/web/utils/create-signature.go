package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func CreateSignature(header, payload, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(header + "." + payload))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

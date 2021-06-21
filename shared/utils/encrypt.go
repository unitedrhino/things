package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

func HmacSha256(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func HmacSha1(data string, secret []byte) string {
	h := hmac.New(sha1.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

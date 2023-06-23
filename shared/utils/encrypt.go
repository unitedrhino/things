package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/i-Things/things/shared/errors"
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

func PKCS5Padding(src []byte, blockSize int) []byte {
	padLen := blockSize - len(src)%blockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padding...)
}
func AesCbcBase64(src, productSecret string) (string, error) {
	if src == "" || productSecret == "" {
		return "", errors.Default.AddMsg("加密参数错误")
	}
	// 截取 productSecret 前 16 位作为密钥
	key := []byte(productSecret)[:16]
	// 以长度 16 的字符 "0" 作为偏移量
	iv := bytes.Repeat([]byte("0"), 16)

	data := []byte(src)

	// 使用 AES-CBC 加密
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 对补全后的数据进行加密
	blockSize := block.BlockSize()
	data = PKCS5Padding(data, blockSize)
	cryptData := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cryptData, data)

	// 进行 base64 编码
	return base64.StdEncoding.EncodeToString(cryptData), nil
}

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
	"sort"
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

/*
	    计算签名:
		1. 对参数params序列按字典序升序排序。
		2. 将以上参数，按参数名称 = 参数值 & 参数名称 = 参数值拼接成字符串。
		3. 使用 HMAC-sha1 算法对上一步中获得的字符串进行计算，密钥为 secret。
		4. 将生成的结果使用 Base64 进行编码，即可获得最终的签名串放入 signature。
*/
func GetSignature(secret string, params map[string]string) string {
	if secret == "" || len(params) == 0 {
		return ""
	}
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	queryStr := ""
	for _, k := range keys {
		if queryStr != "" {
			queryStr += "&"
		}
		queryStr += k + "=" + params[k]
	}
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(queryStr))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padLen := blockSize - len(src)%blockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padding...)
}
func AesCbcBase64(str, productSecret string) (string, error) {
	// 截取 productSecret 前 16 位作为密钥
	key := []byte(productSecret)[:16]
	// 以长度 16 的字符 "0" 作为偏移量
	iv := bytes.Repeat([]byte("0"), 16)

	data := []byte(str)

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

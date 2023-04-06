package utils

import (
	"encoding/base64"
	"math/rand"
	"time"
)

var (
	digits    string = "0123456789"
	specials1 string = "=+=+//"
	specials2 string = "~=+%^*/()[]{}/!@#$?|"
	letter    string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz"
	//密码强度
	strength = map[int]string{
		0: digits,
		1: digits + letter,
		2: digits + letter + specials1,
		3: digits + specials2 + letter,
		4: digits + specials2 + letter + specials2}
)

/*
@in len 密码的长度
@in level 密码的强度级别  0:包含数字和字母  1:包含数字字母和特殊字符
*/
func Random(length int, level int) string {
	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, length)
	str, ok := strength[level]
	if ok != true {
		panic("Random not support level")
	}
	for i := 0; i < length; i++ {
		buf[i] = str[rand.Intn(len(str))]
	}
	//rand.Shuffle(len(buf), func(i, j int) {
	//	buf[i], buf[j] = buf[j], buf[i]
	//})
	return string(buf)
}

func GetRandomBase64(length int) string {
	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = byte(rand.Intn(255))
	}
	return base64.StdEncoding.EncodeToString(buf)
}

package utils

import (
	"math/rand"
	"time"
)

var (
	digits string= "0123456789"
	specials1 string = "=+=+//"
	specials2 string= "~=+%^*/()[]{}/!@#$?|"
	letter string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz"
	//密码强度
	strength = map[int]string{
		0:digits+letter,
		1:digits+letter+specials1,
		2:digits+specials2+letter,
		3:digits+specials2+letter+specials2}
)


/*
@in len 密码的长度
@in level 密码的强度级别  0:包含数字和字母  1:包含数字字母和特殊字符
*/
func GetPassword(length int,level int) string{
	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, length)
	str,ok := strength[level]
	if  ok != true {
		panic("GetPassword not support level")
	}
	for i := 0; i < length; i++ {
		buf[i] = str[rand.Intn(len(str))]
	}
	//rand.Shuffle(len(buf), func(i, j int) {
	//	buf[i], buf[j] = buf[j], buf[i]
	//})
	return string(buf)
}
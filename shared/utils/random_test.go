package utils

import (
	"fmt"
	"testing"
)

func TestGetPassword(t *testing.T) {
	pass := GetPassword(32,1)
	fmt.Println(pass)
}


func TestGetPwdBase64(t *testing.T) {
	pass := GetPwdBase64(20)
	fmt.Println(pass)
}
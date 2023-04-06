package utils

import (
	"fmt"
	"testing"
)

func TestRandom(t *testing.T) {
	pass := Random(32, 2)
	fmt.Println(pass)
}

func TestGetRandomBase64(t *testing.T) {
	pass := GetRandomBase64(20)
	fmt.Println(pass)
}

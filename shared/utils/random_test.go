package utils

import (
	"fmt"
	"testing"
)

func TestGetPassword(t *testing.T) {
	pass := GetPassword(32,1)
	fmt.Println(pass)
}
package utils

import (
	"fmt"
	"testing"
)

func TestDecimalToAny(t *testing.T) {
	num := int64(1699865466064867328)
	bit62 := DecimalToAny(num, 62)
	ans := AnyToDecimal("100+=0000000", 62)
	if ans != num {
		t.Error("ans not equal bit64")
	}
	fmt.Println(bit62, ans)
}
func TestToLen(t *testing.T) {
	src := FillZeroToLen("123", 11)
	ans := AnyToDecimal("00000000i23", 62)
	bit62 := DecimalToAny(ans, 62)
	fmt.Println(src, ans, bit62)
}

package utils

import (
	"github.com/spf13/cast"
	"strings"
)

var tenToAny map[int64]string = map[int64]string{0: "0", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9",
	10: "a", 11: "b", 12: "c", 13: "d", 14: "e", 15: "f", 16: "g", 17: "h", 18: "i", 19: "j",
	20: "k", 21: "l", 22: "m", 23: "n", 24: "o", 25: "p", 26: "q", 27: "r", 28: "s", 29: "t",
	30: "u", 31: "v", 32: "w", 33: "x", 34: "y", 35: "z", 36: "A", 37: "B", 38: "C", 39: "D",
	40: "E", 41: "F", 42: "G", 43: "H", 44: "I", 45: "J", 46: "K", 47: "L", 48: "M", 49: "N",
	50: "O", 51: "P", 52: "Q", 53: "R", 54: "S", 55: "T", 56: "U", 57: "V", 58: "W", 59: "X", 60: "Y", 61: "Z"}

//func main() {
//	fmt.Println(decimalToAny(9999, 76))
//	fmt.Println(anyToDecimal("1F[", 76))
//}

// 10进制转任意进制
func DecimalToAny(num, n int64) string {
	if n > 62 {
		panic("decimalToAny not suppot bit")
	}
	new_num_str := ""
	var remainder int64
	var remainder_string string
	for num != 0 {
		remainder = num % n
		if 62 > remainder && remainder > 9 {
			remainder_string = tenToAny[remainder]
		} else {
			remainder_string = cast.ToString(remainder)
		}
		new_num_str = remainder_string + new_num_str
		num = num / n
	}
	return new_num_str
}

func FillZeroToLen(src string, length int) string {
	for len(src) < length {
		src = "0" + src
	}
	return src
}

// map根据value找key
func findkey(in string) int64 {
	result := int64(-1)
	for k, v := range tenToAny {
		if in == v {
			result = k
		}
	}
	return result
}

// 任意进制转10进制
func AnyToDecimal(num string, n int) int64 {
	if n > 62 {
		panic("decimalToAny not suppot bit")
	}
	var new_num int64 = 0
	nNum := len(strings.Split(num, "")) - 1
	for _, value := range strings.Split(num, "") {
		tmp := findkey(value)
		if tmp != -1 {
			new_num = new_num + tmp*pow(int64(n), int64(nNum))
			nNum = nNum - 1
		} else {
			break
		}
	}
	return int64(new_num)
}

//a的n次方
//超出uint64的部分会丢失
func pow(a, n int64) int64 {
	result := int64(1)
	for i := n; i > 0; i >>= 1 {
		if i&1 != 0 {
			result *= a
		}
		a *= a
	}
	return result
}

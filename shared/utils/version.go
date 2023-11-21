package utils

import (
	"strings"
)

// 版本号对比：v1 > v2 ==> 1 或 v1 < v2 ==> -1 或 v1 == v2 ==> 0
func VersionCompare(v1 string, v2 string) int {
	sv1 := strings.Split(v1, ".")
	sv2 := strings.Split(v2, ".")
	s1Appended, s2Appended := apeendZreo(sv1, sv2)
	for i := 0; i < len(s1Appended); i++ {
		if s1Appended[i] > s2Appended[i] {
			return 1
		}
		if s1Appended[i] < s2Appended[i] {
			return -1
		}
	}
	// 退出循环表示版本号相同
	return 0
}

// 补全切片
func apeendZreo(s1 []string, s2 []string) ([]string, []string) {
	var count int
	if len(s1) > len(s2) {
		count = len(s1) - len(s2)
		for i := 0; i < count; i++ {
			s2 = append(s2, "0")
		}
	}
	if len(s1) < len(s2) {
		count = len(s2) - len(s1)
		for i := 0; i < count; i++ {
			s1 = append(s1, "0")
		}
	}

	return s1, s2
}

// func strToInt64(str string) int64 {
// 	res, err := strconv.Atoi(str)
// 	if err != nil {
// 		fmt.Println("Invalid Number string")
// 		return -1
// 	}
// 	return int64(res)
// }

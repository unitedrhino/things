package utils

import "strings"

// 从目标串src中查找第n个目标字符c所在位置下标
func IndexN(src string, c byte, n int) int {
	var s []byte
	s = []byte(src)

	for i := 0; i < len(s); i++ {
		if n == 0 {
			return i
		}
		if s[i] == c {
			n--
		}
	}
	return -1
}

// SplitCutset 按数组 cuset 里的分隔符，对 str 进行切割
func SplitCutset(str, cutset string) []string {
	words := strings.FieldsFunc(str, func(r rune) bool {
		return strings.ContainsRune(cutset, r)
	})
	result := make([]string, 0, len(words))
	for _, w := range words {
		wd := strings.TrimSpace(w)
		if wd != "" {
			result = append(result, wd)
		}
	}
	return result
}

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstLower 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func NewFillString(num int, val string, sep string) string {
	sli := NewFillSlice(num, val)
	return strings.Join(sli, sep)
}

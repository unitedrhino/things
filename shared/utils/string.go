package utils

//从目标串src中查找第n个目标字符c所在位置下标
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

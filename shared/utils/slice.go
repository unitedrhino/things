package utils

//SliceIn 判断 in 是否在 cmp 中
func SliceIn[T comparable](in T, cmp ...T) bool {
	for _, v := range cmp {
		if in == v {
			return true
		}
	}
	return false
}

//SliceIndex 获取 T[index] 的值，否则返回默认 defaul
func SliceIndex[T any](slice []T, index int, defaul T) T {
	if index >= 0 && index < len(slice) {
		return slice[index]
	} else {
		return defaul
	}
}

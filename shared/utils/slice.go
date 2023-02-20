package utils

func SliceIn[t comparable](in t, cmp ...t) bool {
	for _, v := range cmp {
		if in == v {
			return true
		}
	}
	return false
}

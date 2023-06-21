package utils

import (
	"github.com/gogf/gf/v2/container/gset"
)

// SliceIn 判断 in 是否在 cmp 中
func SliceIn[T comparable](in T, cmp ...T) bool {
	for _, v := range cmp {
		if in == v {
			return true
		}
	}
	return false
}

// SliceIndex 获取 T[index] 的值，否则返回默认 defaul
func SliceIndex[T any](slice []T, index int, defaul T) T {
	if index >= 0 && index < len(slice) {
		return slice[index]
	} else {
		return defaul
	}
}

func NewFillSlice[T any](num int, val T) []T {
	sli := make([]T, num)
	for i := range sli {
		sli[i] = val
	}
	return sli
}

// SliceLeftDiff 判断 childs 是否包含在 bases 里（忽略 childs 多出来的值）；
// 如 bases=[1,2,3]， childs=[2,3,4]，则会返回true；
func SliceLeftDiff[T comparable](bases, childs []T) []any {
	childSet := gset.NewFrom(childs)
	baseSet := gset.NewFrom(bases)
	return baseSet.Diff(childSet).Slice()
}

// SliceLeftContain 判断 childs 是否包含在 bases 里（忽略 childs 多出来的值）；
// 如 bases=[1,2,3]， childs=[2,3,4]，则会返回true；
func SliceLeftContain[T comparable](bases, childs []T) bool {
	if len(SliceLeftDiff(bases, childs)) != 0 {
		return false
	} else {
		return true
	}
}

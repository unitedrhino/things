package utils

import (
	"fmt"
	"github.com/spf13/cast"
	"golang.org/x/exp/constraints"
	"math"
	"strconv"
)

type CanAdd interface {
	constraints.Integer | constraints.Float
}

func Sum[addT CanAdd](datas ...addT) (sum addT) {
	for _, v := range datas {
		sum += v
	}
	return
}

// 保留n位小数
func Decimal[valueType constraints.Float](value valueType, n int) valueType {
	if math.IsNaN(float64(value)) {
		return 0
	}
	v, _ := strconv.ParseFloat(fmt.Sprintf("%."+cast.ToString(n)+"f", value), 64)
	return valueType(v)
}
func Max[t CanAdd](in []t) t {
	if len(in) == 0 {
		return 0
	}
	var max t
	for _, v := range in {
		if v > max {
			max = v
		}
	}
	return max
}

func Min[t CanAdd](in []t) t {
	if len(in) == 0 {
		return 0
	}
	var min t
	for _, v := range in {
		if v < min {
			min = v
		}
	}
	return min
}

package automation

import (
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/utils"
	"github.com/spf13/cast"
	"golang.org/x/exp/constraints"
)

type CmpType string

const (
	CmpTypeEq  CmpType = "eq"  //相等
	CmpTypeNot CmpType = "not" //不相等
	CmpTypeBtw CmpType = "btw" //在xx之间
	CmpTypeGt  CmpType = "gt"  //大于
	CmpTypeGte CmpType = "gte" //大于等于
	CmpTypeLt  CmpType = "lt"  //小于
	CmpTypeLte CmpType = "lte" //小于等于
	CmpTypeIn  CmpType = "in"  //在xx值之中,可以有n个参数
)

func (t CmpType) IsHit(dataType schema.DataType, data any, values []string) bool {
	switch dataType {
	case schema.DataTypeFloat:
		return TermCompareAll(t, cast.ToFloat64(data), utils.SliceTo(values, cast.ToFloat64))
	case schema.DataTypeInt, schema.DataTypeTimestamp:
		return TermCompareAll(t, cast.ToInt64(data), utils.SliceTo(values, cast.ToInt64))
	case schema.DataTypeString, schema.DataTypeBool:
		return TermCompareAll(t, cast.ToString(data), values)
	default:
		return TermCompareAll(t, cast.ToString(data), values)
	}
}

func TermCompareAll[dt constraints.Ordered](t CmpType, data dt, values []dt) bool {
	switch t {
	case CmpTypeEq:
		return data == values[0]
	case CmpTypeNot:
		return data != values[0]
	case CmpTypeBtw:
		if len(values) < 2 {
			return false
		}
		return data >= values[0] && data <= values[1]
	case CmpTypeGt:
		return data > values[0]
	case CmpTypeGte:
		return data >= values[0]
	case CmpTypeLt:
		return data < values[0]
	case CmpTypeLte:
		return data <= values[0]
	case CmpTypeIn:
		return utils.SliceIn(data, values...)
	}
	return false
}

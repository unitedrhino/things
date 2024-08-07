package scene

import (
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
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
	CmpTypeAll CmpType = "all" //全部触发
)

func (t CmpType) IsHit(dataType schema.DataType, data any, values []string) bool {
	switch dataType {
	case schema.DataTypeFloat:
		return TermCompareAll(t, cast.ToFloat64(data), utils.SliceTo(values, cast.ToFloat64))
	case schema.DataTypeInt, schema.DataTypeTimestamp, schema.DataTypeBool:
		return TermCompareAll(t, cast.ToInt64(data), utils.SliceTo(values, cast.ToInt64))
	case schema.DataTypeString:
		return TermCompareAll(t, cast.ToString(data), values)
	default:
		return TermCompareAll(t, cast.ToString(data), values)
	}
}

func (t CmpType) Validate(values []string) error {
	if !utils.SliceIn(t, CmpTypeEq, CmpTypeNot, CmpTypeBtw, CmpTypeGt, CmpTypeGte, CmpTypeLt, CmpTypeLte, CmpTypeIn, CmpTypeAll) {
		return errors.Parameter.AddMsg("动态条件类型 类型不支持:" + string(t))
	}
	if len(values) == 0 && t != CmpTypeAll {
		return errors.Parameter.AddMsg("动态条件类型 需要填写参数")
	}
	if utils.SliceIn(t, CmpTypeIn, CmpTypeBtw) && len(values) != 2 {
		return errors.Parameter.AddMsgf("动态条件类型:%v 需要填写2个参数:%v", string(t), values)
	}
	return nil
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

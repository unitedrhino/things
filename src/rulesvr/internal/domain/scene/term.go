// Package scene 触发条件
package scene

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
)

type Terms []*Term

type TermConditionType string

const (
	TermConditionTypeOr  TermConditionType = "or"
	TermConditionTypeAnd TermConditionType = "and"
)

// eq: 相等  not:不相等  btw:在xx之间  gt: 大于  gte:大于等于 lt:小于  lte:小于等于   in:在xx值之间
type TermType string

const (
	TermTypeEq  TermType = "eq"
	TermTypeNot TermType = "not"
	TermTypeBtw TermType = "btw"
	TermTypeGt  TermType = "gt"
	TermTypeGte TermType = "gte"
	TermTypeLt  TermType = "lt"
	TermTypeLte TermType = "lte"
	TermTypeIn  TermType = "in"
)

type Term struct {
	ColumnType    TermColumnType    `json:"columnType"`    //字段类型 property:属性 event:事件 sysTime:系统时间
	ColumnSchema  *ColumnSchema     `json:"columnSchema"`  //物模型类型
	ColumnTime    *TimeRange        `json:"columnTime"`    //时间类型 只支持后面几种特殊字符:*  - ,
	ConditionType TermConditionType `json:"conditionType"` //多个条件关联类型  or  and
	Terms         Terms             `json:"terms"`         //嵌套条件
}

func (t Terms) Validate() error {
	if t == nil {
		return nil
	}
	for _, v := range t {
		err := v.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Term) Validate() error {
	if t == nil {
		return nil
	}
	err := t.ColumnType.Validate()
	if err != nil {
		return err
	}
	switch t.ColumnType {
	case TermColumnTypeProperty, TermColumnTypeEvent:
		if err := t.ColumnSchema.Validate(); err != nil {
			return err
		}
	case TermColumnTypeSysTime:
		if err := t.ColumnTime.Validate(); err != nil {
			return err
		}
	}
	if !utils.SliceIn(t.ConditionType, TermConditionTypeOr, TermConditionTypeAnd) {
		return errors.Parameter.AddMsg("触发条件中的多个条件关联类型不支持的类型:" + string(t.ConditionType))
	}
	//if !utils.SliceIn(t.TermType, TermTypeEq, TermTypeNot, TermTypeBtw, TermTypeGt, TermTypeGte, TermTypeLt, TermTypeLte, TermTypeIn) {
	//	return errors.Parameter.AddMsg("触发条件中的动态条件类型不支持的类型:" + string(t.TermType))
	//}
	for i := range t.Terms {
		err := t.Terms[i].Validate()
		if err != nil {
			return err
		}
	}
	return nil
}
func (t TermType) Validate(values []string) error {
	if !utils.SliceIn(t, TermTypeEq, TermTypeNot, TermTypeBtw, TermTypeGt, TermTypeGte, TermTypeLt, TermTypeLte, TermTypeIn) {
		return errors.Parameter.AddMsg("动态条件类型 类型不支持:" + string(t))
	}
	if len(values) == 0 {
		return errors.Parameter.AddMsg("动态条件类型 需要填写参数")
	}
	if utils.SliceIn(t, TermTypeIn, TermTypeBtw) && len(values) != 2 {
		return errors.Parameter.AddMsgf("动态条件类型:%v 需要填写2个参数:%v", string(t), values)
	}
	return nil
}

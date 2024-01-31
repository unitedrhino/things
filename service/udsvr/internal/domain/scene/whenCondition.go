// Package scene 触发条件
package scene

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
)

type Conditions struct {
	Terms []Term       `json:"terms"`
	Type  TermCondType `json:"type"`
}

type TermCondType string

const (
	TermConditionTypeOr  TermCondType = "or"
	TermConditionTypeAnd TermCondType = "and"
)

type Term struct {
	ColumnType TermColumnType `json:"columnType"` //字段类型 property:属性 weather:天气

	Property *TermProperty `json:"property"` //属性类型
	//todo 天气状态处于xxx
	Weather any `json:"weather"`
}

func (t *Conditions) Validate() error {
	if t == nil {
		return nil
	}
	if len(t.Terms) > 1 && !utils.SliceIn(t.Type, TermConditionTypeOr, TermConditionTypeAnd) {
		return errors.Parameter.AddMsg("填写多个条件需要填写条件类型:" + string(t.Type))
	}
	for _, v := range t.Terms {
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
	case TermColumnTypeProperty:
		if err := t.Property.Validate(); err != nil {
			return err
		}
	}
	return nil
}
func (t CmpType) Validate(values []string) error {
	if !utils.SliceIn(t, CmpTypeEq, CmpTypeNot, CmpTypeBtw, CmpTypeGt, CmpTypeGte, CmpTypeLt, CmpTypeLte, CmpTypeIn) {
		return errors.Parameter.AddMsg("动态条件类型 类型不支持:" + string(t))
	}
	if len(values) == 0 {
		return errors.Parameter.AddMsg("动态条件类型 需要填写参数")
	}
	if utils.SliceIn(t, CmpTypeIn, CmpTypeBtw) && len(values) != 2 {
		return errors.Parameter.AddMsgf("动态条件类型:%v 需要填写2个参数:%v", string(t), values)
	}
	return nil
}

// 判断条件是否成立
func (t Conditions) IsHit(ctx context.Context, repo TermRepo) bool {
	if len(t.Terms) == 0 {
		return true
	}
	var finalIsHit = false
	for _, v := range t.Terms {
		isHit := v.IsHit(ctx, repo)
		if isHit && t.Type == TermConditionTypeOr {
			return true
		}
		//如果没有命中又是or条件,或者命中了但是是and条件,则需要继续判断
		finalIsHit = isHit //如果是or,每个都返回false那就是false
	}
	return finalIsHit
}

func (t *Term) IsHit(ctx context.Context, repo TermRepo) bool {
	var isHit bool
	switch t.ColumnType {
	case TermColumnTypeProperty:
		isHit = t.Property.IsHit(ctx, t.ColumnType, repo)
	}

	return isHit
}

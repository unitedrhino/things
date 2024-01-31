// Package scene 触发条件
package scene

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"time"
)

type Terms []*Term

type TermConditionType string

const (
	TermConditionTypeOr  TermConditionType = "or"
	TermConditionTypeAnd TermConditionType = "and"
)

type Term struct {
	ColumnType        TermColumnType    `json:"columnType"`        //字段类型 property:属性 event:事件 sysTime:系统时间
	ColumnSchema      *ColumnSchema     `json:"columnSchema"`      //物模型类型
	ColumnTime        *TimeRange        `json:"columnTime"`        //时间类型 只支持后面几种特殊字符:*  - ,
	NextCondition     TermConditionType `json:"netCondition"`      //和下个条件的关联类型  or  and
	ChildrenCondition TermConditionType `json:"childrenCondition"` //和嵌套条件的关联类型  or  and
	Terms             Terms             `json:"terms"`             //嵌套条件
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
	if !utils.SliceIn(t.NextCondition, TermConditionTypeOr, TermConditionTypeAnd) {
		return errors.Parameter.AddMsg("触发条件中的下个条件的关联类型不支持的类型:" + string(t.NextCondition))
	}
	if !utils.SliceIn(t.ChildrenCondition, TermConditionTypeOr, TermConditionTypeAnd, "") {
		return errors.Parameter.AddMsg("触发条件中的嵌套条件的关联类型不支持的类型:" + string(t.ChildrenCondition))
	}
	for i := range t.Terms {
		err := t.Terms[i].Validate()
		if err != nil {
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
func (t Terms) IsHit(ctx context.Context, repo TermRepo) bool {
	if len(t) == 0 {
		return true
	}
	var nextCondition = TermConditionTypeOr
	var finalIsHit = false
	for _, v := range t {
		isHit := v.IsHit(ctx, repo)
		if !isHit && nextCondition == TermConditionTypeAnd {
			return false
		}
		if isHit && nextCondition == TermConditionTypeOr {
			return true
		}
		//如果没有命中又是or条件,或者命中了但是是and条件,则需要继续判断
		finalIsHit = isHit //如果是or,每个都返回false那就是false
		nextCondition = v.NextCondition
	}
	return finalIsHit
}

func (t *Term) IsHit(ctx context.Context, repo TermRepo) bool {
	var isHit bool
	switch t.ColumnType {
	case TermColumnTypeProperty, TermColumnTypeEvent:
		isHit = t.ColumnSchema.IsHit(ctx, t.ColumnType, repo)
	case TermColumnTypeSysTime:
		isHit = t.ColumnTime.IsHit(time.Now())
	}
	if !isHit && t.ChildrenCondition == TermConditionTypeAnd { //如果没满足,如果是and条件直接返回false即可
		return false
	}
	return t.Terms.IsHit(ctx, repo)
}

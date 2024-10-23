package scene

import (
	"gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"golang.org/x/exp/constraints"
	"strings"
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

type Cmps []Cmp

// 只有结构体类型需要用这个参数
type Cmp struct {
	Column   string   `json:"column"`   //结构体类型,数组及事件需要填写 dataID为aa 结构体类型需要填写:aa.bb 数组类型需要填写 aa.1  数组结构体需要填写 aa.1.abc  事件需要填写 过滤的参数 aa.bb
	TermType CmpType  `json:"termType"` //动态条件类型  eq: 相等  not:不相等  btw:在xx之间  gt: 大于  gte:大于等于 lt:小于  lte:小于等于   in:在xx值之间
	Values   []string `json:"values"`   //比较条件列表
}

type Compare struct {
	DataID   string   `json:"dataID"`             //选择为属性或事件时需要填该字段 属性的id及事件的id 如果是数组类型也支持aaa.123
	DataName string   `json:"dataName"`           //对应的物模型定义,只读
	TermType CmpType  `json:"termType,omitempty"` //动态条件类型  eq: 相等  not:不相等  btw:在xx之间  gt: 大于  gte:大于等于 lt:小于  lte:小于等于   in:在xx值之间
	Values   []string `json:"values,omitempty"`   //比较条件列表
	Terms    Cmps     `json:"terms,omitempty"`    //如果需要多个条件,则可以填写到这里来
	Param    string   `json:"-"`                  //触发的参数
}

func (c Compare) EventValidate(s *schema.Event) error {
	if c.TermType != "" {
		if err := c.TermType.Validate(c.Values); err != nil {
			return err
		}
	} else {
		if err := c.Terms.EventValidate(s); err != nil {
			return err
		}
	}
	return nil
}

func (c Compare) PropertyValidate(s *schema.Property) error {
	if c.TermType != "" {
		if err := c.TermType.Validate(c.Values); err != nil {
			return err
		}
	} else {
		if err := c.Terms.PropertyValidate(s); err != nil {
			return err
		}
	}
	return nil
}

func (c Compare) PropertyIsHit(s *schema.Property, dataID string, param any) (hit bool) {
	defer func() {
		if hit {
			c.Param = cast.ToString(param)
		}
	}()
	if c.TermType != "" {
		if c.DataID != dataID {
			return false
		}
		hit = c.TermType.IsHit(s.Define.Type, param, c.Values)
		return hit
	}
	for _, cmp := range c.Terms {
		hit = func() bool {
			dataIDs := strings.Split(cmp.Column, ".")
			if len(dataIDs) == 1 {
				return cmp.TermType.IsHit(s.Define.Type, param, cmp.Values)
			}
			define := s.Define
			if define.Type == schema.DataTypeArray { //数组类型需要填写比较的 点 ,如 aa.12
				if define.ArrayInfo.Type != schema.DataTypeStruct {
					return cmp.TermType.IsHit(s.Define.Type, param, cmp.Values)
				}
				define = *define.ArrayInfo
				dataIDs = append(dataIDs[:1], dataIDs[2:]...)
			}
			if define.Type != schema.DataTypeStruct {
				return false
			}
			p, ok := param.(map[string]any)
			if !ok {
				return false
			}
			val := p[dataIDs[1]]
			if val == nil {
				return false
			}
			return cmp.TermType.IsHit(define.Spec[dataIDs[1]].DataType.Type, val, cmp.Values)
		}()
		if !hit {
			return false
		}
	}
	return true
}

func (c Compare) EventIsHit(s *schema.Event, dataID string, param any) (hit bool) {
	defer func() {
		if hit {
			c.Param = cast.ToString(param)
		}
	}()
	if c.TermType != "" {
		if c.DataID != dataID {
			return false
		}
		p, ok := param.(map[string]any)
		if !ok {
			return false
		}
		val := p[dataID]
		if val == nil {
			return false
		}
		return c.TermType.IsHit(s.Param[dataID].Define.Type, val, c.Values)
	}
	for _, cmp := range c.Terms {
		hit = func() bool {
			dataIDs := strings.Split(cmp.Column, ".")
			p, ok := param.(map[string]any)
			if !ok {
				return false
			}
			val := p[dataIDs[1]]
			if val == nil {
				return false
			}
			return cmp.TermType.IsHit(s.Param[dataIDs[1]].Define.Type, val, cmp.Values)
		}()
		if !hit {
			return false
		}
	}
	return true
}

func (c Cmps) EventValidate(s *schema.Event) error {
	if len(c) == 0 {
		return errors.Parameter.AddMsg("至少填写一个比较条件")
	}
	for _, cmp := range c {
		err := cmp.TermType.Validate(cmp.Values)
		if err != nil {
			return err
		}
		dataIDs := strings.Split(cmp.Column, ".")
		if len(dataIDs) != 2 || dataIDs[0] != s.Identifier || s.Param[dataIDs[1]] == nil {
			return errors.Parameter.AddMsgf("比较条件的column不正确:%v", cmp.Column)
		}
	}
	return nil
}

func (c Cmps) PropertyValidate(s *schema.Property) error {
	if len(c) == 0 {
		return errors.Parameter.AddMsg("至少填写一个比较条件")
	}
	for _, cmp := range c {
		err := cmp.TermType.Validate(cmp.Values)
		if err != nil {
			return err
		}
		dataIDs := strings.Split(cmp.Column, ".")
		if len(cmp.Column) == 0 || dataIDs[0] != s.Identifier {
			return errors.Parameter.AddMsgf("比较条件的column不正确:%v", cmp.Column)
		}

		define := s.Define
		if define.Type == schema.DataTypeArray { //数组类型需要填写比较的 点 ,如 aa.12
			pos := cast.ToInt64(dataIDs[1])
			if pos > cast.ToInt64(define.Max) {
				return errors.Parameter.AddMsgf("超过数组的长度:%v 填写的值为:%v", define.Max, pos)
			}
			if define.ArrayInfo.Type == schema.DataTypeStruct {
				if len(dataIDs) != 3 { //aaa.123.bbb
					return errors.Parameter.AddMsgf("结构体数组需要填写比较的具体字段:%v 如:aaa.123.bbb", cmp.Column)
				}
				//把pos去掉后面校验结构体类型
				define = *define.ArrayInfo
				dataIDs = append(dataIDs[:1], dataIDs[2:]...)
			}
		}
		if define.Type == schema.DataTypeStruct {
			if len(dataIDs) != 2 {
				return errors.Parameter.AddMsgf("结构体需要填写比较的具体字段:%v 如:aaa.bbb", cmp.Column)
			}
			if define.Spec[dataIDs[1]] == nil {
				return errors.Parameter.AddMsgf("比较条件的column不正确:%v", cmp.Column)
			}
		}

	}
	return nil
}

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

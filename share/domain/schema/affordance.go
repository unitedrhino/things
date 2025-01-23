package schema

import "gitee.com/unitedrhino/share/utils"

// 物模型功能类型 1:property属性 2:event事件 3:action行为
type AffordanceType int64

const (
	//物模型功能类型：1-property 属性
	AffordanceTypeProperty AffordanceType = 1
	//物模型功能类型：2-event 事件
	AffordanceTypeEvent AffordanceType = 2
	//物模型功能类型：3-action 行为
	AffordanceTypeAction AffordanceType = 3
)

func (m AffordanceType) String() string {
	switch m {
	case AffordanceTypeProperty:
		return "property"
	case AffordanceTypeEvent:
		return "event"
	case AffordanceTypeAction:
		return "action"

	}
	return ""
}

type (
	PropertyAffordance struct {
		IsUseShadow bool         `json:"isUseShadow"` //是否使用设备影子
		IsNoRecord  bool         `json:"isNoRecord"`  //不存储历史记录
		Define      Define       `json:"define"`      //数据定义
		Mode        PropertyMode `json:"mode"`        //读写类型: 1:r(只读) 2:rw(可读可写)
	}
	EventAffordance struct {
		Type   EventType `json:"type"`   //事件类型: 信息:info  告警alert  故障:fault
		Params Params    `json:"params"` //事件参数
	}
	ActionAffordance struct {
		Dir    ActionDir `json:"dir"`    //调用方向
		Input  Params    `json:"input"`  //调用参数
		Output Params    `json:"output"` //返回参数
	}
)

func DoToPropertyAffordance(in *Property) *PropertyAffordance {
	return utils.Copy[PropertyAffordance](in)
}

func DoToAffordanceStr(in any) string {
	switch in.(type) {
	case *Property:
		a := utils.Copy[PropertyAffordance](in)
		return utils.MarshalNoErr(a)
	case *Event:
		a := utils.Copy[EventAffordance](in)
		return utils.MarshalNoErr(a)
	case *Action:
		a := utils.Copy[ActionAffordance](in)
		return utils.MarshalNoErr(a)
	default:
		panic("不支持的类型")
	}

}

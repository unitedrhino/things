package schema

import (
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"strings"
	"time"
)

type (
	// Model 物模型协议-数据模板定义
	Model struct {
		Version    string      `json:"version,omitempty"`    //版本号
		Properties Properties  `json:"properties,omitempty"` //属性
		Events     Events      `json:"events,omitempty"`     //事件
		Actions    Actions     `json:"actions,omitempty"`    //行为
		Profile    Profile     `json:"profile"`              //配置信息
		Property   PropertyMap `json:"-"`                    //内部使用,使用map加速匹配,key为id
		Event      EventMap    `json:"-"`                    //内部使用,使用map加速匹配,key为id
		Action     ActionMap   `json:"-"`                    //内部使用,使用map加速匹配,key为id
	}
	CommonParam struct {
		Identifier        string     `json:"identifier"`        //标识符 (统一)
		Tag               Tag        `json:"tag"`               //物模型标签 1:自定义 2:可选 3:必选  必选不可删除
		Name              string     `json:"name"`              //功能名称
		Desc              string     `json:"desc"`              //描述
		Required          bool       `json:"required"`          //是否必须
		ExtendConfig      string     `json:"extendConfig"`      //拓展参数,json格式
		IsCanSceneLinkage int64      `json:"isCanSceneLinkage"` //是否支持场景联动控制 (统一)
		IsShareAuthPerm   int64      `json:"isShareAuthPerm"`   // 分享是否需要校验权限 (统一)
		FuncGroup         int64      `json:"funcGroup"`         // 功能分类: 1:普通功能 2:系统功能
		ControlMode       int64      `json:"controlMode"`       //控制模式: 1: 可以群控,可以单控  2:只能单控
		UserPerm          int64      `json:"userPerm"`          //用户权限操作: 1:r(只读) 3:rw(可读可写)
		RecordMode        RecordMode `json:"recordMode"`        // 1(默认) 记录历史记录 2 只记录差异值 3 不记录历史记录
		IsPassword        def.Bool   `json:"isPassword"`        //是否是密码类型,密码类型需要加掩码
		Order             int64      `json:"order"`             //排序
	}
	/*配置信息*/
	Profile struct {
		ProductID string `json:"productID"` //产品ID
	}

	/*内部使用*/
	PropertyMap map[string]*Property
	EventMap    map[string]*Event
	ActionMap   map[string]*Action

	/*结构体说明*/
	Spec struct {
		Identifier string `json:"identifier"` //参数标识符
		Name       string `json:"name"`       //参数名称
		DataType   Define `json:"dataType"`   //参数定义
	}
	Specs []Spec

	/*参数*/
	Param struct {
		Identifier string `json:"identifier"`       //参数标识符
		Name       string `json:"name"`             //参数名称
		Define     Define `json:"define,omitempty"` //参数定义
	}
	Params []Param

	/*事件*/
	Event struct {
		CommonParam
		Type   EventType         `json:"type"`   //事件类型: 1:信息:info  2:告警alert  3:故障:fault
		Params Params            `json:"params"` //事件参数
		Param  map[string]*Param `json:"-"`      //内部使用,使用map加速匹配,key为id
	}
	Events []Event

	/*行为*/
	Action struct {
		CommonParam
		Dir    ActionDir         `json:"dir"`    //调用方向
		Input  Params            `json:"input"`  //调用参数
		Output Params            `json:"output"` //返回参数
		In     map[string]*Param `json:"-"`      //内部使用,使用map加速匹配,key为id
		Out    map[string]*Param `json:"-"`      //内部使用,使用map加速匹配,key为id
	}
	Actions []Action

	/*属性*/
	Property struct {
		CommonParam
		Mode        PropertyMode `json:"mode"`        //读写类型:rw(可读可写) r(只读)
		Define      Define       `json:"define"`      //数据定义
		IsUseShadow bool         `json:"isUseShadow"` //是否使用设备影子
		IsNoRecord  bool         `json:"isNoRecord"`  //不存储历史记录
	}
	Properties []Property

	/*数据类型定义*/
	Define struct {
		Type      DataType          `json:"type"`                //参数类型:bool int string struct float timestamp array enum
		Mapping   map[string]string `json:"mapping,omitempty"`   //枚举及bool类型:bool enum
		Min       string            `json:"min,omitempty"`       //数值最小值:int  float
		Max       string            `json:"max,omitempty"`       //数值最大值:int string float array
		Start     string            `json:"start,omitempty"`     //初始值:int float
		Step      string            `json:"step,omitempty"`      //步长:int float
		Unit      string            `json:"unit,omitempty"`      //单位:int float
		Specs     Specs             `json:"specs,omitempty"`     //结构体:struct
		ArrayInfo *Define           `json:"arrayInfo,omitempty"` //数组:array
		Spec      map[string]*Spec  `json:"-"`                   //内部使用,使用map加速匹配,key为id
	}
)

func (m *Model) String() string {
	tls, _ := json.Marshal(m)
	return string(tls)
}

func (m *Model) ToSimple() *ModelSimple {
	var ret ModelSimple
	for _, p := range m.Properties {
		ret.Properties = append(ret.Properties, PropertySimple{
			Identifier: p.Identifier,
			Name:       p.Name,
			Mode:       p.Mode,
			Define:     p.Define,
		})
	}
	for _, e := range m.Events {
		ret.Events = append(ret.Events, EventSimple{
			Identifier: e.Identifier,
			Name:       e.Name,
			Type:       e.Type,
			Params:     e.Params.ToSimple(),
		})
	}
	for _, a := range m.Actions {
		ret.Actions = append(ret.Actions, ActionSimple{
			Identifier: a.Identifier,
			Name:       a.Name,
			Output:     a.Output.ToSimple(),
			Input:      a.Input.ToSimple(),
		})
	}
	return &ret
}

func (p Param) ToSimple() ParamSimple {
	return ParamSimple{
		Identifier: p.Identifier,
		Name:       p.Name,
		Define:     p.Define.ToSimple(),
	}
}
func (p Params) ToSimple() ParamSimples {
	var ret ParamSimples
	for _, v := range p {
		ret = append(ret, v.ToSimple())
	}
	return ret
}

func (d Define) ToSimple() Define {
	ret := d
	ret.Min = ""
	ret.Max = ""
	ret.Start = ""
	ret.Step = ""
	ret.Unit = ""
	ret.Specs = nil
	ret.ArrayInfo = nil
	if d.Type == DataTypeBool {
		ret.Mapping = nil
	}
	return ret
}

func (d *Define) String() string {
	if d == nil {
		return "{}"
	}
	def, _ := json.Marshal(d)
	return string(def)
}

func (p *PropertyMap) GetMapWithIDs(datas ...string) map[string]*Property {
	var ids = map[string]*Property{}
	for _, data := range datas {
		id, _, ok := GetArray(data)
		if ok {
			v := (*p)[id]
			if v != nil {
				ids[data] = v
			}
			continue
		}
		v := (*p)[data]
		if v != nil {
			if v.Define.Type == DataTypeArray {
				for i := 0; i < cast.ToInt(v.Define.Max); i++ {
					ids[GenArray(data, i)] = v
				}
				continue
			}
			ids[data] = v
		}
	}
	return ids
}

func (p *PropertyMap) GetMap() map[string]*Property {
	var ids = map[string]*Property{}
	for _, v := range *p {
		switch v.Define.Type {
		case DataTypeArray:
			for i := 0; i < cast.ToInt(v.Define.Max); i++ {
				ids[GenArray(v.Identifier, i)] = v
			}
		default:
			ids[v.Identifier] = v
		}
	}
	return ids
}

func (p *PropertyMap) GetIDs() []string {
	var ids []string
	for _, v := range *p {
		switch v.Define.Type {
		case DataTypeArray:
			for i := 0; i < cast.ToInt(v.Define.Max); i++ {
				ids = append(ids, GenArray(v.Identifier, i))
			}
		default:
			ids = append(ids, v.Identifier)
		}
	}
	return ids
}

func (d *Define) GetValueDesc(value any) string {
	switch d.Type {
	case DataTypeBool:
		if utils.ToBool(value) {
			return d.Mapping["1"]
		}
		return d.Mapping["0"]
	case DataTypeTimestamp:
		return time.Unix(0, cast.ToInt64(value)).Format("2006-01-02 15:04:05.000")
	case DataTypeStruct:
		var descs []string
		vv := utils.ToStringMap(value)
		for k, v := range vv {
			s := d.Spec[k]
			if s != nil {
				descs = append(descs, fmt.Sprintf("%s:%s", s.Name, s.DataType.GetValueDesc(v)))
			}
		}
		return strings.Join(descs, ",")
	case DataTypeEnum:
		return d.Mapping[utils.ToString(value)]
	case DataTypeInt, DataTypeFloat:
		return fmt.Sprintf("%v%s", value, d.Unit)
	default:
		return cast.ToString(value)
	}
}

func (d *Define) GetDefaultValue() (retAny any, err error) {
	switch d.Type {
	case DataTypeBool:
		if d.Start != "" {
			return cast.ToInt64(d.Mapping[d.Start]), nil
		}
		return 0, nil
	case DataTypeInt:
		return cast.ToInt64(d.Start), nil
	case DataTypeString:
		return d.Start, nil
	case DataTypeStruct:
		var ret = map[string]any{}
		for _, v := range d.Specs {
			ret[v.Identifier], err = v.DataType.GetDefaultValue()
		}
		return ret, err
	case DataTypeFloat:
		return cast.ToFloat64(d.Start), nil
	case DataTypeTimestamp:
		return int64(0), nil
	case DataTypeArray:
		return d.ArrayInfo.GetDefaultValue()
	case DataTypeEnum:
		if d.Start != "" {
			return cast.ToInt64(d.Mapping[d.Start]), nil
		}
		var keys []int64
		for k := range d.Mapping {
			keys = append(keys, cast.ToInt64(k))
		}
		return utils.Min(keys), nil
	}
	return nil, errors.Parameter.AddMsgf("数据类型:%v 不支持", d.Type)
}

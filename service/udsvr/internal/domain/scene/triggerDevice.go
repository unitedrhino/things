package scene

import (
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/application"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"reflect"
)

type DeviceSelector string

const (
	DeviceSelectorAll   DeviceSelector = "all"
	DeviceSelectorFixed DeviceSelector = "fixed"
)

type TriggerDevices []*TriggerDevice

type TriggerDevice struct {
	ProductID       string                  `json:"productID"`       //产品id
	Selector        DeviceSelector          `json:"selector"`        //设备选择方式  all: 全部 fixed:指定的设备
	SelectorValues  []string                `json:"selectorValues"`  //选择的列表  选择的列表, fixed类型是设备名列表
	Operator        DeviceOperationOperator `json:"operator"`        //触发类型  connected:上线 disConnected:下线 reportProperty:属性上报 reportEvent: 事件上报
	OperationSchema *OperationSchema        `json:"operationSchema"` //物模型类型的具体操作 reportProperty:属性上报 reportEvent: 事件上报
}

type DeviceOperationOperator string

const (
	DeviceOperationOperatorConnected      DeviceOperationOperator = "connected"
	DeviceOperationOperatorDisConnected   DeviceOperationOperator = "disConnected"
	DeviceOperationOperatorReportProperty DeviceOperationOperator = "reportProperty"
)

type OperationSchema struct {
	DataID    []string   `json:"dataID"`    //选择为属性或事件时需要填该字段 属性的id及事件的id aa.bb.cc
	TermType  CmpType    `json:"termType"`  //动态条件类型  eq: 相等  not:不相等  btw:在xx之间  gt: 大于  gte:大于等于 lt:小于  lte:小于等于   in:在xx值之间
	Values    []string   `json:"values"`    //比较条件列表
	StateKeep *StateKeep `json:"stateKeep"` //状态保持 todo
}

func (t *TriggerDevice) Validate() error {
	if t == nil {
		return nil
	}
	if t.ProductID == "" {
		return errors.Parameter.AddMsg("设备触发类型产品未选择,产品id为:" + t.ProductID)
	}
	if !utils.SliceIn(t.Selector, DeviceSelectorAll, DeviceSelectorFixed) {
		return errors.Parameter.AddMsg("设备触发类型设备选择方式不支持:" + string(t.Selector))
	}
	if !utils.SliceIn(t.Operator, DeviceOperationOperatorConnected, DeviceOperationOperatorDisConnected, DeviceOperationOperatorReportProperty) {
		return errors.Parameter.AddMsg("设备触发的触发类型不支持:" + string(t.Operator))
	}
	return nil
}
func (t TriggerDevices) Validate() error {
	if len(t) == 0 {
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

func (o *OperationSchema) Validate() error {
	if o == nil {
		return nil
	}
	if len(o.DataID) == 0 {
		return errors.Parameter.AddMsg("触发设备类型中的标识符需要填写")
	}
	if err := o.StateKeep.Validate(); err != nil {
		return err
	}
	if err := o.TermType.Validate(o.Values); err != nil {
		return err
	}
	return nil
}

// IsTrigger 判断触发器是否命中
func (t TriggerDevices) IsTriggerWithConn(device devices.Core, operator DeviceOperationOperator) bool {
	for _, d := range t {
		//需要排除不是该设备的触发类型
		if d.Operator != operator {
			continue
		}
		if d.Selector == DeviceSelectorAll {
			return true
		}
		for _, d := range d.SelectorValues {
			if d == device.DeviceName {
				return true
			}
		}
	}
	return false
}

// IsTrigger 判断触发器是否命中属性上报类型
func (t TriggerDevices) IsTriggerWithProperty(reportInfo *application.PropertyReport) bool {
	for _, d := range t {
		//需要排除不是该设备的触发类型
		if d.Operator != DeviceOperationOperatorReportProperty {
			continue
		}
		//判断设备是否命中
		if d.Selector != DeviceSelectorAll {
			var hit bool
			for _, d := range d.SelectorValues {
				if d == reportInfo.Device.DeviceName {
					hit = true
					break
				}
			}
			if !hit {
				return false
			}
		}
		if d.OperationSchema.IsHit(reportInfo.Identifier, reportInfo.Param) {
			return true
		}
	}
	return false
}

func (o *OperationSchema) IsHit(dataID string, param any) bool {
	if o.DataID[0] != dataID {
		return false
	}
	var val = param
	dataType := schema.DataType(reflect.TypeOf(param).String())
	if dataType == schema.DataTypeStruct {
		if len(o.DataID) < 2 { //必须指定到结构体的成员
			return false
		}
		st := param.(application.StructValue)
		v, ok := st[o.DataID[1]]
		if !ok { //如果没有获取到该结构体的属性
			return false
		}
		val = v
	}
	return o.TermType.IsHit(dataType, val, o.Values)
}

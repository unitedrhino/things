package scene

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
)

type TriggerDeviceSelector string

const (
	TriggerDeviceSelectorAll   TriggerDeviceSelector = "all"
	TriggerDeviceSelectorFixed TriggerDeviceSelector = "fixed"
)

type TriggerDevices []*TriggerDevice

type TriggerDevice struct {
	ProductID       string                  `json:"productID"`       //产品id
	Selector        TriggerDeviceSelector   `json:"selector"`        //设备选择方式  all: 全部 fixed:指定的设备
	SelectorValues  []string                `json:"selectorValues"`  //选择的列表  选择的列表, fixed类型是设备名列表
	Operator        DeviceOperationOperator `json:"operator"`        //触发类型  online:上线 offline:下线 reportProperty:属性上报 reportEvent: 事件上报
	OperationSchema *OperationSchema        `json:"operationSchema"` //物模型类型的具体操作 reportProperty:属性上报 reportEvent: 事件上报
}

type DeviceOperationOperator string

const (
	DeviceOperationOperatorConnected      DeviceOperationOperator = "connected"
	DeviceOperationOperatorDisConnected   DeviceOperationOperator = "disConnected"
	DeviceOperationOperatorReportProperty DeviceOperationOperator = "reportProperty"
	DeviceOperationOperatorReportEvent    DeviceOperationOperator = "reportEvent"
)

type OperationSchema struct {
	DataID    []string   `json:"dataID"`    //选择为属性或事件时需要填该字段 属性的id及事件的id aa.bb.cc
	TermType  TermType   `json:"termType"`  //动态条件类型  eq: 相等  not:不相等  btw:在xx之间  gt: 大于  gte:大于等于 lt:小于  lte:小于等于   in:在xx值之间
	Values    []string   `json:"values"`    //比较条件列表
	StateKeep *StateKeep `json:"stateKeep"` //状态保持
}

func (t *TriggerDevice) Validate() error {
	if t == nil {
		return nil
	}
	if t.ProductID == "" {
		return errors.Parameter.AddMsg("设备触发类型产品未选择,产品id为:" + t.ProductID)
	}
	if !utils.SliceIn(t.Selector, TriggerDeviceSelectorAll, TriggerDeviceSelectorFixed) {
		return errors.Parameter.AddMsg("设备触发类型设备选择方式不支持:" + string(t.Selector))
	}
	if !utils.SliceIn(t.Operator, DeviceOperationOperatorConnected, DeviceOperationOperatorDisConnected, DeviceOperationOperatorReportProperty, DeviceOperationOperatorReportEvent) {
		return errors.Parameter.AddMsg("设备触发的触发类型不支持:" + string(t.Operator))
	}
	return nil
}
func (t TriggerDevices) Validate() error {
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

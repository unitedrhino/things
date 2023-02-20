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

type TriggerDevice struct {
	ProductID      string                `json:"productID"`      //产品id
	Selector       TriggerDeviceSelector `json:"selector"`       //设备选择方式  all: 全部 fixed:指定的设备
	SelectorValues []string              `json:"selectorValues"` //选择的列表  选择的列表, fixed类型是设备名列表
	Operation      DeviceOperation       `json:"operation"`
}

type DeviceOperationOperator string

const (
	DeviceOperationOperatorOnline         DeviceOperationOperator = "online"
	DeviceOperationOperatorOffline        DeviceOperationOperator = "offline"
	DeviceOperationOperatorReportProperty DeviceOperationOperator = "reportProperty"
	DeviceOperationOperatorReportEvent    DeviceOperationOperator = "reportEvent"
)

type DeviceOperation struct {
	Operator DeviceOperationOperator `json:"operator"` //触发类型  online:上线 offline:下线 reportProperty:属性上报 reportEvent: 事件上报
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
	return t.Operation.Validate()
}

func (d *DeviceOperation) Validate() error {
	if d == nil {
		return nil
	}
	if !utils.SliceIn(d.Operator, DeviceOperationOperatorOnline, DeviceOperationOperatorOffline, DeviceOperationOperatorReportProperty, DeviceOperationOperatorReportEvent) {
		return errors.Parameter.AddMsg("设备触发的触发类型不支持:" + string(d.Operator))
	}
	return nil
}

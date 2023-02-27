package scene

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
)

type TermColumnType string

const (
	TermColumnTypeProperty TermColumnType = "property"
	TermColumnTypeEvent    TermColumnType = "event"
	//TermColumnTypeReportTime TermColumnType = "reportTime"
	TermColumnTypeSysTime TermColumnType = "sysTime"
)

// ColumnSchema 物模型类型 属性,事件
type ColumnSchema struct {
	ProductID  string   `json:"productID"` //产品id
	DeviceName string   `json:"deviceName"`
	DataID     []string `json:"dataID"`   //属性的id及事件的id aa.bb.cc
	TermType   TermType `json:"termType"` //动态条件类型  eq: 相等  not:不相等  btw:在xx之间  gt: 大于  gte:大于等于 lt:小于  lte:小于等于   in:在xx值之间
	Values     []string `json:"values"`   //条件值 参数根据动态条件类型会有多个参数
}

func (t TermColumnType) Validate() error {
	if !utils.SliceIn(t, TermColumnTypeProperty, TermColumnTypeEvent, TermColumnTypeSysTime) {
		return errors.Parameter.AddMsg("条件类型不支持:" + string(t))
	}
	return nil
}

func (c *ColumnSchema) Validate() error {
	if c == nil {
		return nil
	}
	if err := c.TermType.Validate(c.Values); err != nil {
		return err
	}
	if c.ProductID == "" {
		return errors.Parameter.AddMsg("触发设备类型中的产品id需要填写")
	}
	if c.DeviceName == "" {
		return errors.Parameter.AddMsg("触发设备类型中的设备名需要填写")
	}
	if len(c.DataID) == 0 {
		return errors.Parameter.AddMsg("触发设备类型中的标识符需要填写")
	}

	return nil
}

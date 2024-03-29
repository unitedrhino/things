package scene

import (
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/application"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"reflect"
)

type SelectType = string

const (
	SelectorDeviceAll SelectType = "allDevice"   //产品下的所有设备
	SelectDeviceFixed SelectType = "fixedDevice" //产品下的指定设备
	SelectGroup       SelectType = "group"       //某个设备组的设备(支持不同产品的设备,但是需要有共同的公共物模型)
)

type TriggerDevices []*TriggerDevice

type TriggerDevice struct {
	AreaID      int64             `json:"areaID,string"`         //涉及到的区域ID
	ProductID   string            `json:"productID,omitempty"`   //产品id
	SelectType  SelectType        `json:"selectType"`            //设备选择方式  all: 全部 fixed:指定的设备
	GroupID     int64             `json:"groupID,omitempty"`     //分组id
	DeviceName  string            `json:"deviceName,omitempty"`  //选择的列表  选择的列表, fixedDevice类型是设备名列表
	DeviceAlias string            `json:"deviceAlias,omitempty"` //设备别名,只读
	Type        TriggerDeviceType `json:"type,omitempty"`        //触发类型  connected:上线 disConnected:下线 reportProperty:属性上报 reportEvent: 事件上报
	//Schema      *TriggerDeviceSchema `json:"schema,omitempty"`      //物模型类型的具体操作 reportProperty:属性上报 reportEvent: 事件上报
	DataID   string   `json:"dataID"`   //选择为属性或事件时需要填该字段 属性的id及事件的id aa.bb.cc
	DataName string   `json:"dataName"` //对应的物模型定义,只读
	TermType CmpType  `json:"termType"` //动态条件类型  eq: 相等  not:不相等  btw:在xx之间  gt: 大于  gte:大于等于 lt:小于  lte:小于等于   in:在xx值之间
	Values   []string `json:"values"`   //比较条件列表
}

type TriggerDeviceType string

const (
	TriggerDeviceTypeConnected      TriggerDeviceType = "connected"
	TriggerDeviceTypeDisConnected   TriggerDeviceType = "disConnected"
	TriggerDeviceTypeReportProperty TriggerDeviceType = "reportProperty"
)

func (t *TriggerDevice) Validate(repo ValidateRepo) error {
	if t == nil {
		return nil
	}
	if t.ProductID == "" {
		return errors.Parameter.AddMsg("设备触发类型产品未选择,产品id为:" + t.ProductID)
	}
	if !utils.SliceIn(t.SelectType, SelectorDeviceAll, SelectDeviceFixed) {
		return errors.Parameter.AddMsg("设备触发类型设备选择方式不支持:" + string(t.SelectType))
	}
	if !utils.SliceIn(t.Type, TriggerDeviceTypeConnected, TriggerDeviceTypeDisConnected, TriggerDeviceTypeReportProperty) {
		return errors.Parameter.AddMsgf("设备触发的触发类型不支持:%s", string(t.Type))
	}
	if t.Type == TriggerDeviceTypeReportProperty {
		if len(t.DataID) == 0 {
			return errors.Parameter.AddMsg("触发设备类型中的标识符需要填写")
		}
		v, err := repo.ProductSchemaCache.GetData(repo.Ctx, t.ProductID)
		if err != nil {
			return err
		}
		p := v.Property[t.DataID]
		if p == nil {
			return errors.Parameter.AddMsg("dataID不存在")
		}
		t.DataName = p.Name
		if err := t.TermType.Validate(t.Values); err != nil {
			return err
		}
	}
	t.DeviceAlias = GetDeviceAlias(repo.Ctx, repo.DeviceCache, t.ProductID, t.DeviceName)
	return nil
}
func (t TriggerDevices) Validate(repo ValidateRepo) error {
	if len(t) == 0 {
		return nil
	}
	for _, v := range t {
		err := v.Validate(repo)
		if err != nil {
			return err
		}
	}
	return nil
}

// IsTrigger 判断触发器是否命中
func (t TriggerDevices) IsTriggerWithConn(device devices.Core, operator TriggerDeviceType) bool {
	for _, d := range t {
		//需要排除不是该设备的触发类型
		if d.Type != operator {
			continue
		}
		if d.SelectType == SelectorDeviceAll {
			return true
		}
		if d.DeviceName == device.DeviceName {
			return true
		}

	}
	return false
}

// IsTrigger 判断触发器是否命中属性上报类型
func (t TriggerDevices) IsTriggerWithProperty(reportInfo *application.PropertyReport) bool {
	for _, d := range t {
		//需要排除不是该设备的触发类型
		if d.Type != TriggerDeviceTypeReportProperty {
			continue
		}
		//判断设备是否命中
		if d.SelectType != SelectorDeviceAll {
			var hit bool
			if d.DeviceName == reportInfo.Device.DeviceName {
				hit = true
				break
			}

			if !hit {
				return false
			}
		}
		if d.IsHit(reportInfo.Identifier, reportInfo.Param) {
			return true
		}
	}
	return false
}

func (t *TriggerDevice) IsHit(dataID string, param any) bool {
	if t.DataID != dataID {
		return false
	}
	var val = param
	dataType := schema.DataType(reflect.TypeOf(param).String())
	if dataType == schema.DataTypeStruct {
		if len(t.DataID) < 2 { //必须指定到结构体的成员
			return false
		}
		st := param.(application.StructValue)
		v, ok := st[t.DataID]
		if !ok { //如果没有获取到该结构体的属性
			return false
		}
		val = v
	}
	return t.TermType.IsHit(dataType, val, t.Values)
}

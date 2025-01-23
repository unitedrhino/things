package scene

import (
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/dmExport"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/application"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"strings"
)

type SelectType = string

const (
	SelectorDeviceAll SelectType = "all"   //产品下的所有设备
	SelectDeviceFixed SelectType = "fixed" //产品下的指定设备
	SelectArea        SelectType = "area"  //某个区域下的设备
	SelectGroup       SelectType = "group" //某个设备组的设备(支持不同产品的设备,但是需要有共同的公共物模型)
)

type TriggerDevices []*TriggerDevice

type TriggerDevice struct {
	StateKeep        *StateKeep        `json:"stateKeep,omitempty"`   //状态保持
	ProductName      string            `json:"productName,omitempty"` //产品名称,只读
	ProductID        string            `json:"productID,omitempty"`   //产品id
	SelectType       SelectType        `json:"selectType"`            //设备选择方式  all: 全部 fixed:指定的设备
	GroupID          int64             `json:"groupID,omitempty"`     //分组id
	DeviceName       string            `json:"deviceName,omitempty"`  //选择的列表  选择的列表, fixedDevice类型是设备名列表
	DeviceAlias      string            `json:"deviceAlias,omitempty"` //设备别名,只读
	Type             TriggerDeviceType `json:"type,omitempty"`        //触发类型  connected:上线 disConnected:下线 reportProperty:属性上报 reportEvent: 事件上报
	SchemaAffordance string            `json:"schemaAffordance"`      //只读,返回物模型定义
	Compare
	Body string `json:"body,omitempty"` //自定义字段
}

type TriggerDeviceType = string

const (
	TriggerDeviceTypeConnected      TriggerDeviceType = "connected"
	TriggerDeviceTypeDisConnected   TriggerDeviceType = "disConnected"
	TriggerDeviceTypePropertyReport TriggerDeviceType = "propertyReport"
	TriggerDeviceTypeEventReport    TriggerDeviceType = "eventReport"
)

func (t *TriggerDevice) Validate(repo CheckRepo, tt *Trigger) error {
	if t == nil {
		return nil
	}
	if repo.Info.DeviceMode == DeviceModeSingle {
		t.ProductID = repo.Info.ProductID
		t.DeviceName = repo.Info.DeviceName
		t.SelectType = SelectDeviceFixed
	}
	err := t.StateKeep.Validate()
	if err != nil {
		return err
	}
	if t.ProductID == "" {
		return errors.Parameter.AddMsg("设备触发类型产品未选择,产品id为:" + t.ProductID)
	}
	if !utils.SliceIn(t.SelectType, SelectorDeviceAll, SelectDeviceFixed) {
		return errors.Parameter.AddMsg("设备触发类型设备选择方式不支持:" + string(t.SelectType))
	}
	if !utils.SliceIn(t.Type, TriggerDeviceTypeConnected, TriggerDeviceTypeDisConnected, TriggerDeviceTypePropertyReport, TriggerDeviceTypeEventReport) {
		return errors.Parameter.AddMsgf("设备触发的触发类型不支持:%s", string(t.Type))
	}
	uc := ctxs.GetUserCtx(repo.Ctx)
	switch t.SelectType {
	case SelectorDeviceAll:
		if !uc.IsAdmin {
			pa, ok := uc.ProjectAuth[uc.ProjectID]
			if !ok {
				return errors.Permissions.WithMsg("无项目权限")
			}
			if tt.AreaID <= def.NotClassified && pa.AuthType == def.AuthRead {
				return errors.Permissions.WithMsg("只有项目管理员可以制定全局的规则")
			}
			if pa.AuthType == def.AuthRead && pa.Area[tt.AreaID] == 0 {
				return errors.Permissions.WithMsg("无区域权限")
			}
		}
	default:
		_, err = dmExport.SchemaAccess(repo.Ctx, repo.DeviceCache, repo.UserShareCache, def.AuthRead, devices.Core{
			ProductID:  t.ProductID,
			DeviceName: t.DeviceName,
		}, nil)
		if err != nil {
			return err
		}
	}

	switch t.Type {
	case TriggerDeviceTypeEventReport:
		if len(t.DataID) == 0 {
			return errors.Parameter.AddMsg("触发设备类型中的标识符需要填写")
		}
		v, err := repo.SchemaCache.GetData(repo.Ctx, devices.Core{ProductID: t.ProductID, DeviceName: t.DeviceName})
		if err != nil {
			return err
		}
		dataIDs := strings.Split(t.DataID, ".")
		p := v.Event[dataIDs[0]]
		if p == nil {
			return errors.Parameter.AddMsg("dataID不存在")
		}
		if t.DataName == "" {
			t.DataName = p.Name
		}
		t.SchemaAffordance = schema.DoToAffordanceStr(p)
		err = t.Compare.EventValidate(p)
		if err != nil {
			return err
		}
		if t.TermType != "" {
			if err := t.TermType.Validate(t.Values); err != nil {
				return err
			}
		} else {
			if err := t.Terms.EventValidate(p); err != nil {
				return err
			}
		}
	case TriggerDeviceTypePropertyReport:
		if len(t.DataID) == 0 {
			return errors.Parameter.AddMsg("触发设备类型中的标识符需要填写")
		}
		v, err := repo.SchemaCache.GetData(repo.Ctx, devices.Core{ProductID: t.ProductID, DeviceName: t.DeviceName})
		if err != nil {
			return err
		}
		dataIDs := strings.Split(t.DataID, ".")
		p := v.Property[dataIDs[0]]
		if p == nil {
			return errors.Parameter.AddMsg("dataID不存在")
		}
		if t.DataName == "" {
			t.DataName = p.Name
		}
		err = t.Compare.PropertyValidate(p)
		if err != nil {
			return err
		}
		t.SchemaAffordance = schema.DoToAffordanceStr(p)

	}

	if repo.Info.DeviceMode != DeviceModeSingle {
		t.DeviceAlias = GetDeviceAlias(repo.Ctx, repo.DeviceCache, t.ProductID, t.DeviceName)
	}
	pi, err := repo.ProductCache.GetData(repo.Ctx, t.ProductID)
	if err != nil {
		return err
	}
	t.ProductName = pi.ProductName
	return nil
}

//func (t TriggerDevices) Validate(repo CheckRepo) error {
//	if len(t) == 0 {
//		return nil
//	}
//	for _, v := range t {
//		err := v.Validate(repo)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

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
func (t TriggerDevices) IsTriggerWithProperty(model *schema.Model, reportInfo *application.PropertyReport) bool {
	for _, d := range t {
		//需要排除不是该设备的触发类型
		if d.Type != TriggerDeviceTypePropertyReport {
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
		if d.IsHit(model, reportInfo.Identifier, reportInfo.Param) {
			return true
		}
	}
	return false
}

func (t *TriggerDevice) IsHit(model *schema.Model, dataID string, param any) bool {
	if t.DataID != dataID {
		return false
	}
	dataIDs := strings.Split(t.DataID, ".")
	switch t.Type {
	case TriggerDeviceTypePropertyReport:
		property := model.Property[dataIDs[0]]
		if property == nil {
			return false
		}
		return t.PropertyIsHit(property, dataID, param)
	case TriggerDeviceTypeEventReport:
		e := model.Event[dataIDs[0]]
		if e == nil {
			return false
		}
		return t.EventIsHit(e, dataID, param)
	}
	return false
}

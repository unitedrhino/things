package scene

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	deviceinteract "github.com/i-Things/things/service/dmsvr/client/deviceinteract"
	devicemanage "github.com/i-Things/things/service/dmsvr/client/devicemanage"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

type ActionDeviceType = string

const (
	ActionDeviceTypePropertyControl ActionDeviceType = "propertyControl"
	ActionDeviceTypeAction          ActionDeviceType = "action"
)

type ActionDevice struct {
	ProjectID        int64            `json:"-"`                     //项目id
	AreaID           int64            `json:"areaID,string"`         //涉及到的区域ID
	AreaName         string           `json:"areaName"`              //区域的名字
	ProductID        string           `json:"productID"`             //产品id
	SelectType       SelectType       `json:"selectType"`            //设备选择方式
	DeviceName       string           `json:"deviceName"`            //选择的设备列表 指定设备的时候才需要填写(如果设备换到其他区域里,这里删除该设备)
	DeviceAlias      string           `json:"deviceAlias,omitempty"` //设备别名,只读
	GroupID          int64            `json:"groupID"`               //分组id
	Type             ActionDeviceType `json:"type"`                  // 云端向设备发起属性控制: propertyControl  应用调用设备行为:action  todo:通知设备上报
	DataID           string           `json:"dataID"`                // 属性的id及事件的id,不填则取values里面的
	DataName         string           `json:"dataName"`              //对应的物模型定义,只读
	SchemaAffordance string           `json:"schemaAffordance"`      //只读,返回物模型定义
	Value            string           `json:"value"`                 //传的值
	Values           DeviceValues     `json:"values"`                //如果需要控制多个参数
	Body             string           `json:"body,omitempty"`        //自定义字段
}

type DeviceValues = []*DeviceValue
type DeviceValue struct {
	DataID           string `json:"dataID"`           // 属性的id及事件的id
	DataName         string `json:"dataName"`         //对应的物模型定义,只读
	SchemaAffordance string `json:"schemaAffordance"` //只读,返回物模型定义
	Value            string `json:"value"`            //传的值
}

func (a *ActionDevice) Validate(repo ValidateRepo) error {
	if repo.Info.DeviceMode == DeviceModeSingle {
		a.ProductID = repo.Info.ProductID
		a.DeviceName = repo.Info.DeviceName
		a.SelectType = SelectDeviceFixed
	}
	if a.ProductID == "" {
		return errors.Parameter.AddMsgf("产品id不能为空:%v", a.ProductID)
	}
	if !utils.SliceIn(a.SelectType, SelectorDeviceAll, SelectArea, SelectDeviceFixed, SelectGroup) {
		return errors.Parameter.AddMsg("执行的设备选择方式不支持:" + string(a.SelectType))
	}
	if !utils.SliceIn(a.Type, ActionDeviceTypePropertyControl, ActionDeviceTypeAction) {
		return errors.Parameter.AddMsg("云端向设备发起属性控制的方式不支持:" + string(a.Type))
	}
	if a.DataID == "" && len(a.Values) == 0 { //todo 这里需要添加校验,是否存在
		return errors.Parameter.AddMsg("dataID不能为空")
	}
	if repo.Info.DeviceMode != DeviceModeSingle {
		a.DeviceAlias = GetDeviceAlias(repo.Ctx, repo.DeviceCache, a.ProductID, a.DeviceName)
	}
	v, err := repo.ProductSchemaCache.GetData(repo.Ctx, a.ProductID)
	if err != nil {
		return err
	}
	if a.DataID != "" {
		p := v.Property[a.DataID]
		if p == nil {
			return errors.Parameter.AddMsg("dataID不存在")
		}
		if a.DataName == "" {
			a.DataName = p.Name
		}
		a.SchemaAffordance = schema.DoToAffordanceStr(p)
		if a.Value == "" {
			return errors.Parameter.AddMsg("传的值不能为空:%v")
		}
	} else if len(a.Values) != 0 {
		for _, val := range a.Values {
			p := v.Property[val.DataID]
			if p == nil {
				return errors.Parameter.AddMsg("dataID不存在")
			}
			if val.DataName == "" {
				val.DataName = p.Name
			}
			val.SchemaAffordance = schema.DoToAffordanceStr(p)
			if val.Value == "" {
				return errors.Parameter.AddMsg("传的值不能为空:%v")
			}
		}
	}

	return nil
}

var limitChan chan struct{}

func init() {
	limitChan = make(chan struct{}, 500) //设备执行限制并发数为500
}
func (a *ActionDevice) Execute(ctx context.Context, repo ActionRepo) error {
	var (
		executeFunc func(device devices.Core) error
		deviceList  []devices.Core
	)

	toData := func() string {
		if a.DataID != "" {
			var data = map[string]any{
				a.DataID: a.Value,
			}
			ret, _ := json.Marshal(data)
			return string(ret)
		}
		var data = map[string]any{}
		for _, val := range a.Values {
			data[val.DataID] = val.Value
		}
		ret, _ := json.Marshal(data)
		return string(ret)
	}
	switch a.Type {
	case ActionDeviceTypePropertyControl:
		executeFunc = func(device devices.Core) error {
			_, err := repo.DeviceInteract.PropertyControlSend(ctx, &deviceinteract.PropertyControlSendReq{
				IsAsync:    true,
				ProductID:  device.ProductID,
				DeviceName: device.DeviceName,
				Data:       toData(),
			})
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.DeviceInfoIndex SendProperty:%#v err:%v", utils.FuncName(), a, err)
				return err
			}
			return nil
		}
	case ActionDeviceTypeAction:
		executeFunc = func(device devices.Core) error {
			_, err := repo.DeviceInteract.ActionSend(ctx, &deviceinteract.ActionSendReq{
				IsAsync:     true,
				ProductID:   device.ProductID,
				DeviceName:  device.DeviceName,
				ActionID:    a.DataID,
				InputParams: a.Value})
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.DeviceInfoIndex SendAction:%#v err:%v", utils.FuncName(), a, err)
				return err
			}
			return nil
		}
	}
	switch a.SelectType {
	case SelectDeviceFixed:
		deviceList = append(deviceList, devices.Core{
			ProductID:  a.ProductID,
			DeviceName: a.DeviceName,
		})
	case SelectorDeviceAll:
		var areaIDs []int64
		if repo.Info.AreaID != 0 {
			areaIDs = []int64{repo.Info.AreaID}
		}
		ret, err := repo.DeviceM.DeviceInfoIndex(ctx, &devicemanage.DeviceInfoIndexReq{
			AreaIDs:   areaIDs,
			ProductID: a.ProductID,
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.DeviceInfoIndex ActionDevice:%#v err:%v", utils.FuncName(), a, err)
			return err
		}
		for _, v := range ret.List {
			deviceList = append(deviceList, devices.Core{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
			})
		}
	case SelectGroup:
		ret, err := repo.DeviceG.GroupDeviceIndex(ctx, &dm.GroupDeviceIndexReq{GroupID: a.GroupID})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.GroupDeviceIndex ActionDevice:%#v err:%v", utils.FuncName(), a, err)
			return err
		}
		for _, v := range ret.List {
			deviceList = append(deviceList, devices.Core{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
			})
		}
	}
	wait := sync.WaitGroup{}
	for _, device := range deviceList {
		wait.Add(1)
		go func(device devices.Core) {
			defer wait.Done()
			{ //限制并发数,避免打崩
				limitChan <- struct{}{}
				defer func() {
					<-limitChan
				}()
			}
			err := executeFunc(device)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.DeviceInfoIndex device:%v execute:%#v err:%v", utils.FuncName(), device, a, err)
				//return err
			}
		}(device)
	}
	wait.Wait()
	return nil
}

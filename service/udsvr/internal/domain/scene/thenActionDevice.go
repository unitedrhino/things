package scene

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	deviceinteract "github.com/i-Things/things/service/dmsvr/client/deviceinteract"
	devicemanage "github.com/i-Things/things/service/dmsvr/client/devicemanage"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

type ActionDeviceType string

const (
	ActionDeviceTypePropertyControl ActionDeviceType = "propertyControl"
	ActionDeviceTypeAction          ActionDeviceType = "action"
)

type ActionDevice struct {
	ProductID   string           `json:"productID"`   //产品id
	SelectType  SelectType       `json:"selector"`    //设备选择方式   fixed:指定的设备
	DeviceNames []string         `json:"deviceNames"` //选择的设备列表 指定设备的时候才需要填写
	GroupID     int64            `json:"groupID"`     //分组id
	Type        ActionDeviceType `json:"type"`        // 云端向设备发起属性控制: propertyControl  应用调用设备行为:action  todo:通知设备上报
	DataID      string           `json:"dataID"`      // 属性的id及事件的id
	Value       string           `json:"value"`       //传的值
}

func (a *ActionDevice) Validate() error {
	if a.ProductID == "" {
		return errors.Parameter.AddMsgf("产品id不能为空:%v", a.ProductID)
	}
	if !utils.SliceIn(a.SelectType, SelectorDeviceAll, SelectDeviceFixed, SelectGroup) {
		return errors.Parameter.AddMsg("执行的设备选择方式不支持:" + string(a.SelectType))
	}
	if !utils.SliceIn(a.Type, ActionDeviceTypePropertyControl, ActionDeviceTypeAction) {
		return errors.Parameter.AddMsg("云端向设备发起属性控制的方式不支持:" + string(a.Type))
	}
	if a.DataID == "" { //todo 这里需要添加校验,是否存在
		return errors.Parameter.AddMsg("dataID不能为空")
	}
	if a.Value == "" {
		return errors.Parameter.AddMsg("传的值不能为空:%v")
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
		toCores     func(productID string, deviceNames []string) []devices.Core
	)
	toData := func(dataID string, Value any) string {
		var data = map[string]any{
			dataID: Value,
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
				Data:       toData(a.DataID, a.Value), //todo 这里需要根据dataID来生成
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
		deviceList = toCores(a.ProductID, a.DeviceNames)
	case SelectorDeviceAll:
		ret, err := repo.DeviceM.DeviceInfoIndex(ctx, &devicemanage.DeviceInfoIndexReq{
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

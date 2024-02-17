package scene

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	deviceinteract "github.com/i-Things/things/service/dmsvr/client/deviceinteract"
	devicemanage "github.com/i-Things/things/service/dmsvr/client/devicemanage"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

type ActionDeviceType string

const (
	ActionDeviceTypePropertyControl ActionDeviceType = "propertyControl"
	ActionDeviceTypeAction          ActionDeviceType = "action"
)

type ActionDevice struct {
	ProductID      string           `json:"productID"`      //产品id
	Selector       DeviceSelector   `json:"selector"`       //设备选择方式   fixed:指定的设备
	SelectorValues []string         `json:"selectorValues"` //选择的列表  选择的列表, fixed类型是设备名列表
	Type           ActionDeviceType `json:"type"`           // 云端向设备发起属性控制: propertyControl  应用调用设备行为:action  todo:通知设备上报
	DataID         string           `json:"dataID"`         // 属性的id及事件的id
	Value          string           `json:"value"`          //传的值
}

func (a *ActionDevice) Validate() error {
	if a.ProductID == "" {
		return errors.Parameter.AddMsgf("产品id不能为空:%v", a.ProductID)
	}
	if !utils.SliceIn(a.Selector, DeviceSelectorAll, DeviceSelectorFixed) {
		return errors.Parameter.AddMsg("执行的设备选择方式不支持:" + string(a.Selector))
	}
	if !utils.SliceIn(a.Type, ActionDeviceTypePropertyControl, ActionDeviceTypeAction) {
		return errors.Parameter.AddMsg("云端向设备发起属性控制的方式不支持:" + string(a.Type))
	}
	if a.Value == "" {
		return errors.Parameter.AddMsgf("传的值不能为空:%v", a.Value)
	}
	return nil
}

func (a *ActionDevice) Execute(ctx context.Context, repo ActionRepo) error {
	var (
		executeFunc func(productID, deviceName string) error
		deviceList  []string
	)

	switch a.Type {
	case ActionDeviceTypePropertyControl:
		executeFunc = func(productID, deviceName string) error {
			_, err := repo.DeviceInteract.SendPropertyControl(ctx, &deviceinteract.SendPropertyControlReq{
				IsAsync:    true,
				ProductID:  productID,
				DeviceName: deviceName,
				Data:       a.Value, //todo 这里需要根据dataID来生成
			})
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.DeviceInfoIndex SendProperty:%#v err:%v", utils.FuncName(), a, err)
				return err
			}
			return nil
		}
	case ActionDeviceTypeAction:
		executeFunc = func(productID, deviceName string) error {
			_, err := repo.DeviceInteract.SendAction(ctx, &deviceinteract.SendActionReq{
				IsAsync:     true,
				ProductID:   productID,
				DeviceName:  deviceName,
				ActionID:    a.DataID,
				InputParams: a.Value})
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.DeviceInfoIndex SendAction:%#v err:%v", utils.FuncName(), a, err)
				return err
			}
			return nil
		}
	}
	if a.Selector == DeviceSelectorFixed {
		deviceList = a.SelectorValues
	} else {
		ret, err := repo.DeviceM.DeviceInfoIndex(ctx, &devicemanage.DeviceInfoIndexReq{
			ProductID: a.ProductID,
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.DeviceInfoIndex ActionDevice:%#v err:%v", utils.FuncName(), a, err)
			return err
		}
		for _, v := range ret.List {
			deviceList = append(deviceList, v.DeviceName)
		}
	}
	wait := sync.WaitGroup{}
	for _, device := range deviceList {
		wait.Add(1)
		go func(device string) {
			defer wait.Done()
			err := executeFunc(a.ProductID, device)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.DeviceInfoIndex device:%v execute:%#v err:%v", utils.FuncName(), device, a, err)
				//return err
			}
		}(device)
	}
	wait.Wait()
	return nil
}

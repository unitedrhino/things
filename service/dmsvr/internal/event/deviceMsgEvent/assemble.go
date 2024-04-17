package deviceMsgEvent

import (
	"gitee.com/i-Things/share/domain/deviceMsg/msgGateway"
	"gitee.com/i-Things/share/domain/deviceMsg/msgThing"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToDmDevicesCore(devices []*msgGateway.Device) (ret []*dm.DeviceCore) {
	for _, v := range devices {
		ret = append(ret, &dm.DeviceCore{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	return
}

func ToDmDevicesBind(devices []*msgGateway.Device) (ret []*dm.DeviceGatewayBindDevice, err error) {
	for _, v := range devices {
		if v.Signature == "" || v.SignMethod == "" {
			return nil, errors.Parameter.AddMsgf("产品ID为:%v,设备名为:%v 的设备没有填写签名", v.ProductID, v.DeviceName)
		}
		ret = append(ret, &dm.DeviceGatewayBindDevice{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
			Sign: &dm.DeviceGatewaySign{
				Signature:  v.Signature,
				Random:     v.Random,
				Timestamp:  v.Timestamp,
				SignMethod: v.SignMethod,
			},
		})
	}
	return
}

func ToDmDevicesInfoReq(diDeviceBasicInfoDo *msgThing.DeviceBasicInfo) (dmDeviceInfoReq *dm.DeviceInfo) {
	var position *dm.Point
	if p := diDeviceBasicInfoDo.Position; p != nil {
		gcp := utils.PositionToMars(*p)
		position = &dm.Point{Longitude: gcp.Longitude, Latitude: gcp.Latitude}
	}

	return &dm.DeviceInfo{
		ProductID:  diDeviceBasicInfoDo.ProductID,
		DeviceName: diDeviceBasicInfoDo.DeviceName,
		Imei:       diDeviceBasicInfoDo.Imei,
		Mac:        diDeviceBasicInfoDo.Mac,
		Version:    utils.ToRpcNullString(diDeviceBasicInfoDo.Version),
		HardInfo:   diDeviceBasicInfoDo.HardInfo,
		SoftInfo:   diDeviceBasicInfoDo.SoftInfo,
		Position:   position,
		Tags:       diDeviceBasicInfoDo.Tags,
		Rssi:       utils.ToRpcNullInt64(diDeviceBasicInfoDo.Rssi),
		Iccid:      utils.ToRpcNullString(diDeviceBasicInfoDo.Iccid),
	}
}

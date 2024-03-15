package deviceMsgEvent

import (
	"gitee.com/i-Things/share/domain/deviceMsg/msgGateway"
	"gitee.com/i-Things/share/domain/deviceMsg/msgThing"
	"gitee.com/i-Things/share/domain/schema"
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

func ToParamValues(tp map[string]msgThing.Param) (map[string]any, error) {
	ret := make(map[string]any, len(tp))
	var err error
	for k, v := range tp {
		ret[k], err = ToParamValue(v)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func ToParamValue(p msgThing.Param) (any, error) {
	var ret any
	var err error
	switch p.Value.Type {
	case schema.DataTypeStruct:
		v, ok := p.Value.Value.(map[string]msgThing.Param)
		if ok == false {
			return ret, errors.Parameter.AddMsgf("struct Param is not find")
		}
		val := make(map[string]any, len(v)+1)
		for _, tp := range v {
			val[tp.Identifier], err = ToParamValue(tp)
			if err != nil {
				return ret, err
			}
		}
		ret = val
		return ret, nil
	case schema.DataTypeArray:
		array, ok := p.Value.Value.([]any)
		if ok == false {
			return ret, errors.Parameter.AddMsgf("array Param is not find")
		}
		val := make([]any, 0, len(array)+1)
		for _, value := range array {
			switch value.(type) {
			case map[string]msgThing.Param:
				valMap := make(map[string]any, len(array)+1)
				for _, tp := range value.(map[string]msgThing.Param) {
					valMap[tp.Identifier], err = ToParamValue(tp)
					return ret, err
				}
				val = append(val, valMap)
			default:
				val = append(val, value)
			}
		}
		ret = val
		return ret, nil
	default:
		ret = p.Value.Value
		return ret, nil
	}
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

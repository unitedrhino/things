package deviceMsgEvent

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/shared/domain/application"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgGateway"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgThing"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
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

func ToParamValues(tp map[string]msgThing.Param) map[string]application.ParamValue {
	ret := make(map[string]application.ParamValue, len(tp))
	for k, v := range tp {
		ret[k] = ToParamValue(v)
	}
	return ret
}

func ToParamValue(p msgThing.Param) application.ParamValue {
	var ret application.ParamValue
	ret.Type = p.Value.Type
	switch p.Value.Type {
	case schema.DataTypeStruct:
		v, ok := p.Value.Value.(map[string]msgThing.Param)
		if ok == false {
			panic("struct Param is not find")
		}
		val := make(map[string]application.ParamValue, len(v)+1)
		for _, tp := range v {
			val[tp.Identifier] = ToParamValue(tp)
		}
		ret.Value = val
		return ret
	case schema.DataTypeArray:
		array, ok := p.Value.Value.([]any)
		if ok == false {
			panic("array Param is not find")
		}
		val := make([]any, 0, len(array)+1)
		for _, value := range array {
			switch value.(type) {
			case map[string]msgThing.Param:
				valMap := make(map[string]application.ParamValue, len(array)+1)
				for _, tp := range value.(map[string]msgThing.Param) {
					valMap[tp.Identifier] = ToParamValue(tp)
				}
				val = append(val, valMap)
			default:
				val = append(val, value)
			}
		}
		ret.Value = val
		return ret
	default:
		ret.Value = p.Value.Value
		return ret
	}
}

func ToDmDevicesInfoReq(diDeviceBasicInfoDo *msgThing.DeviceBasicInfo) (dmDeviceInfoReq *dm.DeviceInfo) {
	var position *dm.Point
	if p := diDeviceBasicInfoDo.Position; p != nil {
		gcp := utils.PositionToBaidu(p.CoordinateSystem, p.Longitude, p.Latitude)
		position = &dm.Point{Longitude: gcp.Lon, Latitude: gcp.Lat}
	}

	//TODO 拓展 DeviceInfo 支持新的 基础信息字段
	return &dm.DeviceInfo{
		ProductID:  diDeviceBasicInfoDo.ProductID,
		DeviceName: diDeviceBasicInfoDo.DeviceName,
		Imei:       &wrappers.StringValue{Value: diDeviceBasicInfoDo.Imei},
		Mac:        &wrappers.StringValue{Value: diDeviceBasicInfoDo.Mac},
		Version:    &wrappers.StringValue{Value: diDeviceBasicInfoDo.Version},
		HardInfo:   &wrappers.StringValue{Value: diDeviceBasicInfoDo.HardInfo},
		SoftInfo:   &wrappers.StringValue{Value: diDeviceBasicInfoDo.SoftInfo},
		Position:   position,
		Tags:       diDeviceBasicInfoDo.Tags,
	}
}

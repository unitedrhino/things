package devicemanagelogic

import (
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/events"
	"gitee.com/i-Things/share/utils"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg/msgGateway"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToDeviceInfo(in *relationDB.DmDeviceInfo) *dm.DeviceInfo {
	if in.IsOnline == def.Unknown {
		in.IsOnline = def.False
	}
	if in.LogLevel == def.Unknown {
		in.LogLevel = def.LogClose
	}

	return &dm.DeviceInfo{
		ProductID:      in.ProductID,
		DeviceName:     in.DeviceName,
		ProjectID:      int64(in.ProjectID),
		AreaID:         int64(in.AreaID),
		DeviceAlias:    &wrappers.StringValue{Value: in.DeviceAlias},
		MobileOperator: in.MobileOperator,
		Phone:          utils.ToRpcNullString(in.Phone),
		Iccid:          utils.ToRpcNullString(in.Iccid),
		Secret:         in.Secret,
		Cert:           in.Cert,
		Imei:           in.Imei,
		Mac:            in.Mac,
		Version:        &wrappers.StringValue{Value: in.Version},
		HardInfo:       in.HardInfo,
		SoftInfo:       in.SoftInfo,
		Position:       logic.ToDmPoint(&in.Position),
		Address:        &wrappers.StringValue{Value: in.Address},
		Tags:           in.Tags,
		SchemaAlias:    in.SchemaAlias,
		IsOnline:       in.IsOnline,
		FirstLogin:     utils.GetNullTime(in.FirstLogin),
		LastLogin:      utils.GetNullTime(in.LastLogin),
		LogLevel:       in.LogLevel,
		CreatedTime:    in.CreatedTime.Unix(),
	}
}

func BindToDeviceCoreDos(in []*dm.DeviceGatewayBindDevice) (ret []*devices.Core) {
	for _, v := range in {
		ret = append(ret, &devices.Core{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	return
}

func ToDeviceCoreDos(in []*dm.DeviceCore) (ret []*devices.Core) {
	for _, v := range in {
		ret = append(ret, &devices.Core{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	return
}

func ToDeviceCoreEvents(in []*dm.DeviceCore) (ret []*events.DeviceCore) {
	for _, v := range in {
		ret = append(ret, &events.DeviceCore{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	return
}

func ToDeviceCoreDo(core *dm.DeviceCore) *devices.Core {
	if core == nil {
		return nil
	}
	return &devices.Core{
		ProductID:  core.ProductID,
		DeviceName: core.DeviceName,
	}
}
func BindToDeviceCoreEvents(in []*dm.DeviceGatewayBindDevice) (ret []*events.DeviceCore) {
	for _, v := range in {
		ret = append(ret, &events.DeviceCore{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	return
}

func ToGatewayDevice(gateway *devices.Core, subDevice []*devices.Core) (ret []*relationDB.DmGatewayDevice) {
	for _, v := range subDevice {
		ret = append(ret, &relationDB.DmGatewayDevice{
			GatewayProductID:  gateway.ProductID,
			GatewayDeviceName: gateway.DeviceName,
			ProductID:         v.ProductID,
			DeviceName:        v.DeviceName,
		})
	}
	return
}

func ToGatewayPayload(status int32, in []*devices.Core) *msgGateway.GatewayPayload {
	var ret = msgGateway.GatewayPayload{Status: status}
	for _, v := range in {
		ret.Devices = append(ret.Devices, &msgGateway.Device{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	return &ret
}

func ToTimeRange(in *dm.TimeRange) *def.TimeRange {
	if in == nil {
		return nil
	}
	return &def.TimeRange{
		Start: in.Start,
		End:   in.End,
	}
}

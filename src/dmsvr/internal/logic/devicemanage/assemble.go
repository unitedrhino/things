package devicemanagelogic

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToDeviceInfo(di *relationDB.DmDeviceInfo) *dm.DeviceInfo {
	if di.IsOnline == def.Unknown {
		di.IsOnline = def.False
	}
	if di.LogLevel == def.Unknown {
		di.LogLevel = def.LogClose
	}

	return &dm.DeviceInfo{
		ProductID:      di.ProductID,
		DeviceName:     di.DeviceName,
		ProjectID:      utils.ToRpcNullInt64(di.ProjectID),
		AreaID:         utils.ToRpcNullInt64(di.AreaID),
		DeviceAlias:    &wrappers.StringValue{Value: di.DeviceAlias},
		MobileOperator: di.MobileOperator,
		Phone:          utils.ToRpcNullString(di.Phone),
		Iccid:          utils.ToRpcNullString(di.Iccid),
		Secret:         di.Secret,
		Cert:           di.Cert,
		Imei:           di.Imei,
		Mac:            di.Mac,
		Version:        &wrappers.StringValue{Value: di.Version},
		HardInfo:       di.HardInfo,
		SoftInfo:       di.SoftInfo,
		Position:       logic.ToDmPoint(&di.Position),
		Address:        &wrappers.StringValue{Value: di.Address},
		Tags:           di.Tags,
		IsOnline:       di.IsOnline,
		FirstLogin:     utils.GetNullTime(di.FirstLogin),
		LastLogin:      utils.GetNullTime(di.LastLogin),
		LogLevel:       di.LogLevel,
		CreatedTime:    di.CreatedTime.Unix(),
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

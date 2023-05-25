package devicemanagelogic

import (
	"encoding/json"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	mysql "github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToDeviceInfo(di *mysql.DmDeviceInfo) *dm.DeviceInfo {
	var (
		tags map[string]string
	)

	_ = json.Unmarshal([]byte(di.Tags), &tags)

	if di.IsOnline == def.Unknown {
		di.IsOnline = def.False
	}
	if di.LogLevel == def.Unknown {
		di.LogLevel = def.LogClose
	}

	var Longitude float64
	var Latitude float64
	Longitude, Latitude = utils.GetPositionValue(di.Position)

	return &dm.DeviceInfo{
		ProductID:   di.ProductID,
		DeviceName:  di.DeviceName,
		DeviceAlias: &wrappers.StringValue{Value: di.DeviceAlias},
		Secret:      di.Secret,
		Cert:        di.Cert,
		Imei:        di.Imei,
		Mac:         di.Mac,
		Version:     &wrappers.StringValue{Value: di.Version},
		HardInfo:    di.HardInfo,
		SoftInfo:    di.SoftInfo,
		Position:    &dm.Point{Longitude: Longitude, Latitude: Latitude},
		Address:     &wrappers.StringValue{Value: di.Address},
		Tags:        tags,
		IsOnline:    di.IsOnline,
		FirstLogin:  utils.GetNullTime(di.FirstLogin),
		LastLogin:   utils.GetNullTime(di.LastLogin),
		LogLevel:    di.LogLevel,
		CreatedTime: di.CreatedTime.Unix(),
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

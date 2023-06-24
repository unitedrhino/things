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

func FillDeviceInfo(in *dm.DeviceInfo, di *mysql.DmDeviceInfo) {
	if in.Tags != nil {
		tags, err := json.Marshal(in.Tags)
		if err == nil {
			di.Tags = string(tags)
		}
	} else {
		di.Tags = "{}"
	}

	if in.LogLevel != def.Unknown {
		di.LogLevel = def.LogClose
	}
	if in.Address != nil {
		di.Address = in.Address.Value
	}

	if in.DeviceAlias != nil {
		di.DeviceAlias = in.DeviceAlias.Value
	}
	if in.UserID != 0 {
		di.UserID = in.UserID
	}
	if in.MobileOperator != 0 {
		di.MobileOperator = in.MobileOperator
	}
	if in.Iccid != nil {
		di.Iccid = utils.AnyToNullString(in.Iccid)
	}
	if in.Phone != nil {
		di.Phone = utils.AnyToNullString(in.Phone)
	}
}

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
		ProductID:      di.ProductID,
		DeviceName:     di.DeviceName,
		ProjectID:      utils.ToRpcNullInt64(di.ProjectID),
		AreaID:         utils.ToRpcNullInt64(di.AreaID),
		DeviceAlias:    &wrappers.StringValue{Value: di.DeviceAlias},
		UserID:         di.UserID,
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
		Position:       &dm.Point{Longitude: Longitude, Latitude: Latitude},
		Address:        &wrappers.StringValue{Value: di.Address},
		Tags:           tags,
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

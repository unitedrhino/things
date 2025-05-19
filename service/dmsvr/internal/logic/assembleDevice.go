package logic

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/events"
	"gitee.com/unitedrhino/share/oss/common"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgGateway"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/zeromicro/go-zero/core/logx"
)

func ToDeviceInfo(ctx context.Context, svcCtx *svc.ServiceContext, in *relationDB.DmDeviceInfo) *dm.DeviceInfo {
	if in == nil {
		return nil
	}
	if in.IsOnline == def.Unknown {
		in.IsOnline = def.False
	}
	if in.LogLevel == def.Unknown {
		in.LogLevel = def.LogClose
	}
	var (
		productName string
		deviceType  int64 = def.DeviceTypeDevice
		netType     int64 = def.Unknown
		ProductImg  string
		CategoryID  int64
	)
	pi, err := svcCtx.ProductCache.GetData(ctx, in.ProductID)
	if err == nil {
		deviceType = pi.DeviceType
		productName = pi.ProductName
		netType = pi.NetType
		ProductImg = pi.ProductImg
		CategoryID = pi.CategoryID
	}
	//return utils.Copy[dm.DeviceInfo](in)
	if in.DeviceImg != "" {
		in.DeviceImg, err = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, in.DeviceImg, 60*60, common.OptionKv{})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
		}
	}
	if in.File != "" {
		in.File, err = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, in.File, 60*60, common.OptionKv{})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
		}
	}
	return &dm.DeviceInfo{
		Id:                 in.ID,
		TenantCode:         string(in.TenantCode),
		ProductID:          in.ProductID,
		DeviceName:         in.DeviceName,
		ProjectID:          int64(in.ProjectID),
		AreaID:             int64(in.AreaID),
		AreaIDPath:         string(in.AreaIDPath),
		DeviceAlias:        &wrappers.StringValue{Value: in.DeviceAlias},
		MobileOperator:     in.MobileOperator,
		Phone:              utils.ToRpcNullString(in.Phone),
		Iccid:              utils.ToRpcNullString(in.Iccid),
		Secret:             in.Secret,
		Cert:               in.Cert,
		Imei:               in.Imei,
		Mac:                in.Mac,
		LastIp:             in.LastIp,
		Desc:               utils.ToRpcNullString(in.Desc),
		Version:            &wrappers.StringValue{Value: in.Version},
		HardInfo:           in.HardInfo,
		SoftInfo:           in.SoftInfo,
		Position:           ToDmPoint(&in.Position),
		Address:            &wrappers.StringValue{Value: in.Address},
		Rssi:               &wrappers.Int64Value{Value: in.Rssi},
		Tags:               in.Tags,
		SchemaAlias:        in.SchemaAlias,
		IsOnline:           in.IsOnline,
		IsEnable:           in.IsEnable,
		FirstLogin:         utils.GetNullTime(in.FirstLogin),
		FirstBind:          utils.GetNullTime(in.FirstBind),
		LastBind:           utils.GetNullTime(in.LastBind),
		LastLogin:          utils.GetNullTime(in.LastLogin),
		ExpTime:            utils.TimeToNullInt(in.ExpTime),
		File:               in.File,
		DeviceImg:          in.DeviceImg,
		LogLevel:           in.LogLevel,
		CreatedTime:        in.CreatedTime.Unix(),
		ProtocolConf:       in.ProtocolConf,
		SubProtocolConf:    in.SubProtocolConf,
		Status:             in.Status,
		ProductName:        productName,
		DeviceType:         deviceType,
		RatedPower:         in.RatedPower,
		NetType:            netType,
		Distributor:        utils.Copy[dm.IDPathWithUpdate](&in.Distributor),
		NeedConfirmVersion: in.NeedConfirmVersion,
		NeedConfirmJobID:   in.NeedConfirmJobID,
		ProductImg:         ProductImg,
		UserID:             in.UserID,
		Sort:               in.Sort,
		CategoryID:         CategoryID,
		BelongGroup:        utils.CopyMap2[dm.IDsInfo](in.BelongGroup),
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

func ToGatewayPayload(status def.GatewayStatus, in []*devices.Core) *msgGateway.GatewayPayload {
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

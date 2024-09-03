package deviceinteractlogic

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
)

func CheckIsOnline(ctx context.Context, svcCtx *svc.ServiceContext, core devices.Core) (protocolCode string, err error) {
	info, err := svcCtx.ProductCache.GetData(ctx, core.ProductID)
	if err != nil {
		return "", err
	}
	dev, err := svcCtx.DeviceCache.GetData(ctx, core)
	if err != nil {
		return info.ProtocolCode, err
	}
	if dev.IsOnline == def.False {
		return info.ProtocolCode, errors.NotOnline
	}
	//if dev.IsEnable == def.False { //未启用不能控制
	//	return info.ProtocolCode, errors.NotEnable
	//}

	if info.DeviceType != def.DeviceTypeSubset {
		return info.ProtocolCode, nil
	}
	//子设备需要查询网关的在线状态
	g, err := relationDB.NewGatewayDeviceRepo(ctx).FindOneByFilter(ctx, relationDB.GatewayDeviceFilter{SubDevice: &devices.Core{
		ProductID:  dev.ProductID,
		DeviceName: dev.DeviceName,
	}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return info.ProtocolCode, errors.NotFind.AddMsg("子设备未绑定网关")
		}
		return info.ProtocolCode, err
	}
	di, err := svcCtx.DeviceCache.GetData(ctx, devices.Core{ProductID: g.ProductID, DeviceName: g.DeviceName})
	if err != nil {
		return info.ProtocolCode, err
	}
	if di.IsOnline == def.True {
		return info.ProtocolCode, nil
	}
	return info.ProtocolCode, errors.NotOnline.WithMsg("网关未在线")
}

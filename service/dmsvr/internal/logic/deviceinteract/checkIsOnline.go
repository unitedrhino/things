package deviceinteractlogic

import (
	"context"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/devices"
	"gitee.com/i-Things/core/shared/errors"
	devicemanage "github.com/i-Things/things/service/dmsvr/internal/server/devicemanage"
	productmanage "github.com/i-Things/things/service/dmsvr/internal/server/productmanage"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func CheckIsOnline(ctx context.Context, svcCtx *svc.ServiceContext, core devices.Core) error {
	dev, err := devicemanage.NewDeviceManageServer(svcCtx).DeviceInfoRead(ctx, &dm.DeviceInfoReadReq{
		ProductID:  core.ProductID,
		DeviceName: core.DeviceName,
	})
	if err != nil {
		return err
	}
	if dev.IsOnline == def.False {
		return errors.NotOnline
	}
	info, err := productmanage.NewProductManageServer(svcCtx).ProductInfoRead(ctx, &dm.ProductInfoReadReq{ProductID: core.ProductID})
	if err != nil {
		return err
	}
	if info.DeviceType != def.DeviceTypeSubset {
		return nil
	}
	//子设备需要查询网关的在线状态
	gateways, err := devicemanage.NewDeviceManageServer(svcCtx).DeviceGatewayIndex(ctx, &dm.DeviceGatewayIndexReq{SubDevice: &dm.DeviceCore{
		ProductID:  dev.ProductID,
		DeviceName: dev.DeviceName,
	}})
	if err != nil {
		return err
	}
	if len(gateways.List) == 0 {
		return errors.NotFind.AddMsg("子设备未绑定网关")
	}
	for _, g := range gateways.List {
		if g.IsOnline == def.True {
			return nil
		}
	}
	return errors.NotOnline.AddMsg("网关未在线")
}

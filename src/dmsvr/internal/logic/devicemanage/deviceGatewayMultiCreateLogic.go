package devicemanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceGatewayMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceGatewayMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceGatewayMultiCreateLogic {
	return &DeviceGatewayMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建分组设备
func (l *DeviceGatewayMultiCreateLogic) DeviceGatewayMultiCreate(in *dm.DeviceGatewayMultiCreateReq) (*dm.Response, error) {
	{ //检查是否是网关类型
		pi, err := l.svcCtx.ProductInfo.FindOne(l.ctx, in.GatewayProductID)
		if err != nil {
			if err == mysql.ErrNotFound {
				return nil, errors.Parameter.AddDetail("not find GatewayProductID id:" + cast.ToString(in.GatewayProductID))
			}
			return nil, errors.Database.AddDetail(err)
		}
		if pi.DeviceType != def.DeviceTypeGateway {
			return nil, errors.Parameter.AddMsg("子设备只能由网关设备进行绑定")
		}
	}
	{ //检查设备是否都是子设备类型
		var (
			deviceProductList []string
			deviceProductMap  = map[string]struct{}{}
		)
		for _, v := range in.List {
			deviceProductMap[v.ProductID] = struct{}{}
		}
		deviceProductList = utils.SetToSlice(deviceProductMap)
		products, err := l.svcCtx.ProductInfo.FindByFilter(l.ctx, mysql.ProductFilter{
			ProductIDs: deviceProductList,
		}, nil)
		if err != nil {
			return nil, errors.Database.AddDetail(err)
		}
		for _, v := range products {
			if v.DeviceType != def.DeviceTypeSubset {
				return nil, errors.Parameter.AddMsg("网关只能绑定子设备类型")
			}
		}
	}

	err := l.svcCtx.Gateway.CreateList(l.ctx, &devices.Core{
		ProductID:  in.GatewayProductID,
		DeviceName: in.GatewayDeviceName,
	}, ToDeviceCoreDos(in.List))
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.DataUpdate.DeviceGatewayUpdate(l.ctx, &events.GatewayUpdateInfo{
		GatewayProductID:  in.GatewayProductID,
		GatewayDeviceName: in.GatewayDeviceName,
		Status:            def.GatewayBind,
		Devices:           ToDeviceCoreEvents(in.List),
	})
	if err != nil {
		l.Errorf("%s.DeviceGatewayUpdate err=%+v", utils.FuncName(), err)
	}
	return &dm.Response{}, nil
}

package devicemanagelogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceGatewayIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	GdDB *relationDB.GatewayDeviceRepo
}

func NewDeviceGatewayIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceGatewayIndexLogic {
	return &DeviceGatewayIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		GdDB:   relationDB.NewGatewayDeviceRepo(ctx),
	}
}

// 获取分组设备信息列表
func (l *DeviceGatewayIndexLogic) DeviceGatewayIndex(in *dm.DeviceGatewayIndexReq) (*dm.DeviceGatewayIndexResp, error) {
	f := relationDB.GatewayDeviceFilter{
		Gateway:   logic.ToDeviceCoreDo(in.Gateway),
		SubDevice: logic.ToDeviceCoreDo(in.SubDevice),
	}

	size, err := l.GdDB.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}

	di, err := l.GdDB.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	fmt.Println(di, size)
	var list []*dm.DeviceCore
	if f.Gateway != nil {
		for _, v := range di {
			list = append(list, &dm.DeviceCore{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
			})
		}
	} else {
		for _, v := range di {
			list = append(list, &dm.DeviceCore{
				ProductID:  v.GatewayProductID,
				DeviceName: v.GatewayDeviceName,
			})
		}

	}
	return &dm.DeviceGatewayIndexResp{
		List:  list,
		Total: size,
	}, nil
	//info := make([]*dm.DeviceInfo, 0, len(di))
	//for _, v := range di {
	//	info = append(info, logic.ToDeviceInfo(v))
	//}
	//
	//return &dm.DeviceGatewayIndexResp{List: info, Total: size}, nil
}

package devicemanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"time"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceOnlineMultiFixLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceOnlineMultiFixLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceOnlineMultiFixLogic {
	return &DeviceOnlineMultiFixLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceOnlineMultiFixLogic) DeviceOnlineMultiFix(in *dm.DeviceOnlineMultiFixReq) (*dm.Empty, error) {
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		log := logx.WithContext(ctx)
		for _, device := range in.Devices {
			ld := device.Device
			di, err := l.svcCtx.DeviceCache.GetData(ctx, devices.Core{
				ProductID:  ld.ProductID,
				DeviceName: ld.DeviceName,
			})
			if err != nil {
				log.Error(err)
				continue
			}
			//暂时只做离线修复
			var updates = map[string]any{"is_online": def.True, "last_login": time.UnixMilli(device.ConnectAt), "status": def.DeviceStatusOnline}
			if di.FirstLogin == 0 {
				updates["first_login"] = time.UnixMilli(device.ConnectAt)
			}
			err = relationDB.NewDeviceInfoRepo(ctx).UpdateWithField(ctx,
				relationDB.DeviceFilter{Cores: []*devices.Core{{ProductID: ld.ProductID, DeviceName: ld.DeviceName}}}, updates)
			if err != nil {
				log.Error(err)
			}
			err = l.svcCtx.DeviceCache.SetData(ctx, devices.Core{
				ProductID:  ld.ProductID,
				DeviceName: ld.DeviceName,
			}, nil)
			if err != nil {
				log.Error(err)
			}
		}
	})

	return &dm.Empty{}, nil
}

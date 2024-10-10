package devicemanagelogic

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	DiDB *relationDB.DeviceInfoRepo
}

func NewDeviceInfoCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoCountLogic {
	return &DeviceInfoCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
	}
}

// 设备计数
func (l *DeviceInfoCountLogic) DeviceInfoCount(in *dm.DeviceInfoCountReq) (*dm.DeviceInfoCount, error) {
	f := relationDB.DeviceFilter{
		LastLoginTime: logic.ToTimeRange(in.TimeRange),
		AreaIDs:       in.AreaIDs,
	}
	if len(in.GroupIDs) != 0 {
		gds, err := relationDB.NewGroupDeviceRepo(l.ctx).FindByFilter(l.ctx, relationDB.GroupDeviceFilter{
			GroupIDs: in.GroupIDs,
		}, nil)
		if err != nil || len(gds) == 0 {
			return &dm.DeviceInfoCount{}, err
		}
		for _, v := range gds {
			f.Cores = append(f.Cores, &devices.Core{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
			})
		}
	}
	diCount, err := l.DiDB.CountGroupByField(
		l.ctx, f, "is_online")
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.NotFind
		}
		return nil, err
	}

	onlineCount := diCount[fmt.Sprintf("%d", def.DeviceStatusOnline-1)]
	offlineCount := diCount[fmt.Sprintf("%d", def.DeviceStatusOffline-1)]
	InactiveCount := diCount[fmt.Sprintf("%d", def.DeviceStatusInactive-1)]
	var allCount int64
	for _, v := range diCount {
		allCount += v
	}

	return &dm.DeviceInfoCount{
		Total:    allCount,
		Online:   onlineCount,
		Offline:  offlineCount,
		Inactive: InactiveCount,
		Unknown:  allCount - onlineCount - offlineCount - InactiveCount,
	}, nil
}

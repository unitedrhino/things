package devicemanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/devices"
	"gitee.com/i-Things/core/shared/errors"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

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
		LastLoginTime: ToTimeRange(in.TimeRange),
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

	onlineCount := diCount[fmt.Sprintf("%d", def.DeviceStatusOnline)]
	offlineCount := diCount[fmt.Sprintf("%d", def.DeviceStatusOffline)]
	InactiveCount := diCount[fmt.Sprintf("%d", def.DeviceStatusInactive)]
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

package devicemanagelogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceInfoCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoCountLogic {
	return &DeviceInfoCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 设备计数
func (l *DeviceInfoCountLogic) DeviceInfoCount(in *dm.DeviceInfoCountReq) (*dm.DeviceInfoCountResp, error) {
	diCount, err := l.svcCtx.DeviceInfo.CountGroupByField(
		l.ctx,
		mysql.DeviceFilter{
			LastLoginTime: struct {
				Start int64
				End   int64
			}{Start: in.StartTime, End: in.EndTime},
		},
		"isOnline")
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.NotFind
		}
		return nil, err
	}

	onlineCount := diCount[fmt.Sprintf("%d", def.DeviceStatusOnline)]
	offlineCount := diCount[fmt.Sprintf("%d", def.DeviceStatusOffline)]
	var allCount int64
	for _, v := range diCount {
		allCount += v
	}

	return &dm.DeviceInfoCountResp{
		Online:  onlineCount,
		Offline: offlineCount,
		Unknown: allCount - onlineCount - offlineCount,
	}, nil
}

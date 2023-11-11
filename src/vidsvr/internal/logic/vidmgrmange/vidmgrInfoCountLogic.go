package vidmgrmangelogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrInfoCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrInfoRepo
}

func NewVidmgrInfoCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrInfoCountLogic {
	return &VidmgrInfoCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrtInfoRepo(ctx),
	}
}

// 获取服务统计  在线，离线，未激活
func (l *VidmgrInfoCountLogic) VidmgrInfoCount(in *vid.VidmgrInfoCountReq) (*vid.VidmgrInfoCountResp, error) {
	// todo: add your logic here and delete this line
	VidmgrCount, err := l.PiDB.CountVidmgrByField(
		l.ctx,
		relationDB.VidmgrFilter{
			LastLoginTime: struct {
				Start int64
				End   int64
			}{Start: in.StartTime, End: in.EndTime},
		},
		"status")
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.NotFind
		}
		return nil, err
	}

	onlineCount := VidmgrCount[fmt.Sprintf("%d", def.DeviceStatusOnline)]
	offlineCount := VidmgrCount[fmt.Sprintf("%d", def.DeviceStatusOffline)]
	InactiveCount := VidmgrCount[fmt.Sprintf("%d", def.DeviceStatusInactive)]
	var allCount int64
	for _, v := range VidmgrCount {
		allCount += v
	}
	fmt.Printf("onlineCount == %d\n", onlineCount)
	fmt.Printf("onlineCount == %d\n", offlineCount)
	fmt.Printf("onlineCount == %d\n", InactiveCount)

	return &vid.VidmgrInfoCountResp{
		Online:   onlineCount,
		Offline:  offlineCount,
		Inactive: InactiveCount,
	}, nil

}

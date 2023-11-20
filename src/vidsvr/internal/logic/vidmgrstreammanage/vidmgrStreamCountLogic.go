package vidmgrstreammanagelogic

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

type VidmgrStreamCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrStreamRepo
}

func NewVidmgrStreamCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrStreamCountLogic {
	return &VidmgrStreamCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrStreamRepo(ctx),
	}
}

// 统计流 在线，离线
func (l *VidmgrStreamCountLogic) VidmgrStreamCount(in *vid.VidmgrStreamCountReq) (*vid.VidmgrStreamCountResp, error) {
	// todo: add your logic here and delete this line
	//只需要查看当前的在线标识 就可以了
	StreamCount, err := l.PiDB.CountStreamByField(
		l.ctx,
		relationDB.VidmgrStreamFilter{
			LastLoginTime: struct {
				Start int64
				End   int64
			}{Start: in.StartTime, End: in.EndTime},
		}, "is_online")
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.NotFind
		}
		return nil, err
	}
	onlineCount := StreamCount[fmt.Sprintf("%d", def.DeviceStatusOnline)]
	offlineCount := StreamCount[fmt.Sprintf("%d", def.DeviceStatusOffline)]
	var allCount int64
	for _, v := range StreamCount {
		allCount += v
	}
	fmt.Printf("onlineCount == %d\n", onlineCount)
	fmt.Printf("onlineCount == %d\n", offlineCount)

	return &vid.VidmgrStreamCountResp{
		Online:  onlineCount,
		Offline: offlineCount,
	}, nil

}

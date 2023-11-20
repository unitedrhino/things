package vidmgrstreammanagelogic

import (
	"context"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrStreamReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrStreamRepo
}

func NewVidmgrStreamReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrStreamReadLogic {
	return &VidmgrStreamReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrStreamRepo(ctx),
	}
}

// 获取流信息详情
func (l *VidmgrStreamReadLogic) VidmgrStreamRead(in *vid.VidmgrStreamReadReq) (*vid.VidmgrStream, error) {
	// todo: add your logic here and delete this line
	pi, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.VidmgrStreamFilter{
		StreamIDs: []int64{in.StreamID},
	})
	if err != nil {
		return nil, err
	}
	return ToRpcConvVidmgrStream(pi), nil
}

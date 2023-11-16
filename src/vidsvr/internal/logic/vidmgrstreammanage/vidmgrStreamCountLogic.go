package vidmgrstreammanagelogic

import (
	"context"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrStreamCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrStreamCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrStreamCountLogic {
	return &VidmgrStreamCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 统计流 在线，离线，未激活
func (l *VidmgrStreamCountLogic) VidmgrStreamCount(in *vid.VidmgrStreamCountReq) (*vid.VidmgrStreamCountResp, error) {
	// todo: add your logic here and delete this line

	return &vid.VidmgrStreamCountResp{}, nil
}

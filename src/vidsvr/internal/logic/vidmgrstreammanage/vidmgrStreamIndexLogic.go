package vidmgrstreammanagelogic

import (
	"context"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrStreamIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrStreamIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrStreamIndexLogic {
	return &VidmgrStreamIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取流列表
func (l *VidmgrStreamIndexLogic) VidmgrStreamIndex(in *vid.VidmgrStreamIndexReq) (*vid.VidmgrStreamIndexResp, error) {
	// todo: add your logic here and delete this line

	return &vid.VidmgrStreamIndexResp{}, nil
}

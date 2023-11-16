package vidmgrstreammanagelogic

import (
	"context"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrStreamUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrStreamUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrStreamUpdateLogic {
	return &VidmgrStreamUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 流更新
func (l *VidmgrStreamUpdateLogic) VidmgrStreamUpdate(in *vid.VidmgrStream) (*vid.Response, error) {
	// todo: add your logic here and delete this line

	return &vid.Response{}, nil
}

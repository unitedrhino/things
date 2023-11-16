package vidmgrstreammanagelogic

import (
	"context"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrStreamCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrStreamCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrStreamCreateLogic {
	return &VidmgrStreamCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 流添加
func (l *VidmgrStreamCreateLogic) VidmgrStreamCreate(in *vid.VidmgrStream) (*vid.Response, error) {
	// todo: add your logic here and delete this line

	return &vid.Response{}, nil
}

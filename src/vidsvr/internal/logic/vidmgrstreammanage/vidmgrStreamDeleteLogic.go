package vidmgrstreammanagelogic

import (
	"context"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrStreamDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrStreamDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrStreamDeleteLogic {
	return &VidmgrStreamDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除流
func (l *VidmgrStreamDeleteLogic) VidmgrStreamDelete(in *vid.VidmgrStreamDeleteReq) (*vid.Response, error) {
	// todo: add your logic here and delete this line

	return &vid.Response{}, nil
}

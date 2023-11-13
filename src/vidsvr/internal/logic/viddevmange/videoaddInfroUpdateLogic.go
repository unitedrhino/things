package viddevmangelogic

import (
	"context"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoaddInfroUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoaddInfroUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoaddInfroUpdateLogic {
	return &VideoaddInfroUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 流更新
func (l *VideoaddInfroUpdateLogic) VideoaddInfroUpdate(in *vid.ViddevInfo) (*vid.Response, error) {
	// todo: add your logic here and delete this line

	return &vid.Response{}, nil
}

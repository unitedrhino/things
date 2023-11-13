package viddevmangelogic

import (
	"context"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoaddInfroCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoaddInfroCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoaddInfroCreateLogic {
	return &VideoaddInfroCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 流添加
func (l *VideoaddInfroCreateLogic) VideoaddInfroCreate(in *vid.ViddevInfo) (*vid.Response, error) {
	// todo: add your logic here and delete this line

	return &vid.Response{}, nil
}

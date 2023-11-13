package viddevmangelogic

import (
	"context"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoaddInfroCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoaddInfroCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoaddInfroCountLogic {
	return &VideoaddInfroCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 统计流 在线，离线，未激活
func (l *VideoaddInfroCountLogic) VideoaddInfroCount(in *vid.ViddevInfo) (*vid.Response, error) {
	// todo: add your logic here and delete this line

	return &vid.Response{}, nil
}

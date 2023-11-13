package viddevmangelogic

import (
	"context"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoaddInfroReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoaddInfroReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoaddInfroReadLogic {
	return &VideoaddInfroReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取流信息详情
func (l *VideoaddInfroReadLogic) VideoaddInfroRead(in *vid.ViddevInfo) (*vid.Response, error) {
	// todo: add your logic here and delete this line

	return &vid.Response{}, nil
}

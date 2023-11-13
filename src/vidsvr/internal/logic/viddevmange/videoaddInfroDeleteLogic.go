package viddevmangelogic

import (
	"context"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoaddInfroDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoaddInfroDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoaddInfroDeleteLogic {
	return &VideoaddInfroDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除流
func (l *VideoaddInfroDeleteLogic) VideoaddInfroDelete(in *vid.ViddevInfo) (*vid.Response, error) {
	// todo: add your logic here and delete this line

	return &vid.Response{}, nil
}

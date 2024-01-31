package project

import (
	"context"

	"github.com/i-Things/things/service/viewsvr/internal/svc"
	"github.com/i-Things/things/service/viewsvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.ProjectInfoReadReq) (resp *types.ProjectInfo, err error) {
	// todo: add your logic here and delete this line

	return
}

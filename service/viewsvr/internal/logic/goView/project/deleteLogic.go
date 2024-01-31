package project

import (
	"context"
	"github.com/i-Things/things/service/viewsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/viewsvr/internal/svc"
	"github.com/i-Things/things/service/viewsvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.ProjectInfoDeleteReq) error {
	err := relationDB.NewProjectInfoRepo(l.ctx).Delete(l.ctx, req.ID)
	if err != nil {
		return err
	}
	err = relationDB.NewProjectDetailRepo(l.ctx).Delete(l.ctx, req.ID)
	return err
}

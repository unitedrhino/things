package detail

import (
	"context"
	"github.com/i-Things/things/service/viewsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/viewsvr/internal/svc"
	"github.com/i-Things/things/service/viewsvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.ProjectDetail) error {
	pd, err := relationDB.NewProjectDetailRepo(l.ctx).FindOne(l.ctx, req.ID)
	if err != nil {
		return err
	}
	pd.Content = req.Content
	err = relationDB.NewProjectDetailRepo(l.ctx).Update(l.ctx, pd)
	return err
}

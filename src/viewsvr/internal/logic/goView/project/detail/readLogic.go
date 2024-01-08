package detail

import (
	"context"
	"github.com/i-Things/things/src/viewsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/viewsvr/internal/svc"
	"github.com/i-Things/things/src/viewsvr/internal/types"

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

func (l *ReadLogic) Read(req *types.ProjectInfoReadReq) (resp *types.ProjectDetail, err error) {
	pd, err := relationDB.NewProjectDetailRepo(l.ctx).FindOne(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &types.ProjectDetail{
		ID:      pd.ProjectID,
		Content: pd.Content,
	}, err
}

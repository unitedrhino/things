package project

import (
	"context"
	"github.com/i-Things/things/service/viewsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/viewsvr/internal/svc"
	"github.com/i-Things/things/service/viewsvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.ProjectInfo) error {
	pi := relationDB.ViewProjectInfo{
		IndexImage: req.IndexImage,
		Name:       req.Name,
		Desc:       req.Desc,
		//	CreatedUserID: req.CreatedUserID,
		Status: req.Status,
	}
	err := relationDB.NewProjectInfoRepo(l.ctx).Insert(l.ctx, &pi)
	if err != nil {
		return err
	}
	err = relationDB.NewProjectDetailRepo(l.ctx).Insert(l.ctx, &relationDB.ViewProjectDetail{
		ProjectID: pi.ID,
	})
	return err
}

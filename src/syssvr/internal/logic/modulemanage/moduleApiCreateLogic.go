package modulemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleApiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleApiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleApiCreateLogic {
	return &ModuleApiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleApiCreateLogic) ModuleApiCreate(in *sys.ApiInfo) (*sys.WithID, error) {
	if err := CheckModule(l.ctx, in.ModuleCode); err != nil {
		return nil, err
	}
	po := logic.ToApiInfoPo(in)
	po.ID = 0
	err := relationDB.NewApiInfoRepo(l.ctx).Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	return &sys.WithID{Id: po.ID}, nil

}

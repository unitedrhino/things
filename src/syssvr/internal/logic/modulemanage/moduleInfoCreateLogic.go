package modulemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleInfoCreateLogic {
	return &ModuleInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleInfoCreateLogic) ModuleInfoCreate(in *sys.ModuleInfo) (*sys.WithID, error) {
	if in.Type == 0 {
		in.Type = 1
	}
	if in.Order == 0 {
		in.Order = 1
	}
	if in.HideInMenu == 0 {
		in.HideInMenu = 1
	}
	po := logic.ToModuleInfoPo(in)
	relationDB.NewModuleInfoRepo(l.ctx).Insert(l.ctx, po)
	return &sys.WithID{Id: po.ID}, nil
}

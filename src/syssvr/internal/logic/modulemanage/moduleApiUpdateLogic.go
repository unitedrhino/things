package modulemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleApiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleApiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleApiUpdateLogic {
	return &ModuleApiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleApiUpdateLogic) ModuleApiUpdate(in *sys.ApiInfo) (*sys.Response, error) {
	err := relationDB.NewApiInfoRepo(l.ctx).Update(l.ctx, logic.ToApiInfoPo(in))
	return &sys.Response{}, err
}

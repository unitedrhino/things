package modulemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleInfoUpdateLogic {
	return &ModuleInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleInfoUpdateLogic) ModuleInfoUpdate(in *sys.ModuleInfo) (*sys.Response, error) {
	err := relationDB.NewModuleInfoRepo(l.ctx).Update(l.ctx, logic.ToModuleInfoPo(in))
	return &sys.Response{}, err
}

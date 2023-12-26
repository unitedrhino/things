package modulemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleMenuUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleMenuUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleMenuUpdateLogic {
	return &ModuleMenuUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleMenuUpdateLogic) ModuleMenuUpdate(in *sys.MenuInfo) (*sys.Response, error) {
	err := relationDB.NewMenuInfoRepo(l.ctx).Update(l.ctx, logic.ToMenuInfoPo(in))
	return &sys.Response{}, err
}

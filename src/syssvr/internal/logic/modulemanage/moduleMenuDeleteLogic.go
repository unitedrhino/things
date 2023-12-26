package modulemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleMenuDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleMenuDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleMenuDeleteLogic {
	return &ModuleMenuDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleMenuDeleteLogic) ModuleMenuDelete(in *sys.WithID) (*sys.Response, error) {
	err := relationDB.NewMenuInfoRepo(l.ctx).Delete(l.ctx, in.Id)
	return &sys.Response{}, err
}

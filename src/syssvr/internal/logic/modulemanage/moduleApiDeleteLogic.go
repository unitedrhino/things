package modulemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleApiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleApiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleApiDeleteLogic {
	return &ModuleApiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleApiDeleteLogic) ModuleApiDelete(in *sys.WithID) (*sys.Response, error) {
	err := relationDB.NewApiInfoRepo(l.ctx).Delete(l.ctx, in.Id)
	return &sys.Response{}, err
}

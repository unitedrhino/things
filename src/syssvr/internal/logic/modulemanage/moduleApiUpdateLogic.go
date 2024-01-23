package modulemanagelogic

import (
	"context"
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
	old, err := relationDB.NewApiInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	old.Route = in.Route
	old.Method = in.Method
	old.Name = in.Name
	old.BusinessType = in.BusinessType
	old.Group = in.Group
	old.IsNeedAuth = in.IsNeedAuth
	old.Desc = in.Desc
	err = relationDB.NewApiInfoRepo(l.ctx).Update(l.ctx, old)
	return &sys.Response{}, err
}

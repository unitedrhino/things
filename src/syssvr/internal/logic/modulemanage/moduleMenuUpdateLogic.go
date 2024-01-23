package modulemanagelogic

import (
	"context"
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
	old, err := relationDB.NewMenuInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	old.Type = in.Type
	old.Order = in.Order
	old.Name = in.Name
	old.Path = in.Path
	old.Component = in.Component
	old.Icon = in.Icon
	old.Redirect = in.Redirect
	old.Body = in.Body.Value
	old.HideInMenu = in.HideInMenu
	err = relationDB.NewMenuInfoRepo(l.ctx).Update(l.ctx, old)
	return &sys.Response{}, err
}

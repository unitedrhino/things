package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppApiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppApiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppApiUpdateLogic {
	return &TenantAppApiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppApiUpdateLogic) TenantAppApiUpdate(in *sys.TenantApiInfo) (*sys.Response, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	old, err := relationDB.NewTenantAppApiRepo(l.ctx).FindOne(l.ctx, in.Info.Id)
	if err != nil {
		return nil, err
	}
	old.Route = in.Info.Route
	old.Method = in.Info.Method
	old.Name = in.Info.Name
	old.BusinessType = in.Info.BusinessType
	old.Group = in.Info.Group
	old.Desc = in.Info.Desc
	err = relationDB.NewTenantAppApiRepo(l.ctx).Update(l.ctx, old)
	return &sys.Response{}, err
}

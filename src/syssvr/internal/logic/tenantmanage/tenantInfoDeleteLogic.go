package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantInfoDeleteLogic {
	return &TenantInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除区域
func (l *TenantInfoDeleteLogic) TenantInfoDelete(in *sys.WithIDCode) (*sys.Response, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	err := relationDB.NewTenantInfoRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.TenantInfoFilter{Codes: []string{in.Code}, ID: in.Id})
	return &sys.Response{}, err
}

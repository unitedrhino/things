package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppDeleteLogic {
	return &TenantAppDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppDeleteLogic) TenantAppDelete(in *sys.TenantAppWithIDOrCode) (*sys.Response, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	f := relationDB.TenantAppFilter{
		TenantCode: in.Code,
		Codes:      []string{in.AppCode},
	}
	if in.AppCode != "" {
		f.Codes = []string{in.AppCode}
	}
	if in.Id != 0 {
		f.IDs = []int64{in.Id}
	}
	err := relationDB.NewTenantAppRepo(l.ctx).DeleteByFilter(l.ctx, f)
	return &sys.Response{}, err
}

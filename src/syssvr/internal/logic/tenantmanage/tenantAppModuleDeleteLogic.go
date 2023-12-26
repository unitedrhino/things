package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppModuleDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppModuleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppModuleDeleteLogic {
	return &TenantAppModuleDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppModuleDeleteLogic) TenantAppModuleDelete(in *sys.TenantModuleWithIDOrCode) (*sys.Response, error) {
	f := relationDB.TenantAppModuleFilter{ID: in.Id, TenantCode: in.Code}
	if in.AppCode != "" {
		f.AppCodes = []string{in.AppCode}
	}
	if in.ModuleCode != "" {
		f.ModuleCodes = []string{in.ModuleCode}
	}
	err := relationDB.NewTenantAppModuleRepo(l.ctx).DeleteByFilter(l.ctx, f)

	return &sys.Response{}, err
}

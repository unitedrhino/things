package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppModuleIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppModuleIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppModuleIndexLogic {
	return &TenantAppModuleIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppModuleIndexLogic) TenantAppModuleIndex(in *sys.TenantModuleIndexReq) (*sys.TenantModuleIndexResp, error) {
	if err := ctxs.IsRoot(l.ctx); err == nil {
		ctxs.GetUserCtx(l.ctx).AllTenant = true
		defer func() {
			ctxs.GetUserCtx(l.ctx).AllTenant = false
		}()
	}
	ret, err := relationDB.NewTenantAppModuleRepo(l.ctx).FindByFilter(l.ctx, relationDB.TenantAppModuleFilter{TenantCode: in.Code, AppCodes: []string{in.AppCode}}, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	var moduleCodes []string
	for _, v := range ret {
		moduleCodes = append(moduleCodes, v.ModuleCode)
	}
	return &sys.TenantModuleIndexResp{ModuleCodes: moduleCodes}, nil
}

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

type TenantAppApiIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppApiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppApiIndexLogic {
	return &TenantAppApiIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppApiIndexLogic) TenantAppApiIndex(in *sys.TenantAppApiIndexReq) (*sys.TenantAppApiIndexResp, error) {
	if err := ctxs.IsRoot(l.ctx); err == nil {
		ctxs.GetUserCtx(l.ctx).AllTenant = true
		defer func() {
			ctxs.GetUserCtx(l.ctx).AllTenant = false
		}()
	}

	f := relationDB.TenantAppApiFilter{
		TenantCode: in.Code,
		AppCode:    in.AppCode,
		ModuleCode: in.ModuleCode,
	}
	resp, err := relationDB.NewTenantAppApiRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := relationDB.NewTenantAppApiRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	info := make([]*sys.TenantApiInfo, 0, len(resp))
	for _, v := range resp {
		info = append(info, logic.ToTenantAppApiInfoPb(v))
	}

	return &sys.TenantAppApiIndexResp{
		Total: total,
		List:  info,
	}, nil
}

package module

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/tenant/app"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiCreateLogic {
	return &MultiCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiCreateLogic) MultiCreate(req *types.TenantAppCreateReq) error {
	_, err := l.svcCtx.TenantRpc.TenantAppModuleMultiCreate(l.ctx, &sys.TenantAppCreateReq{
		Code:    req.Code,
		AppCode: req.AppCode,
		Modules: app.ToTenantAppModulesPb(req.Modules),
	})
	return err
}

package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppModuleMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppModuleMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppModuleMultiCreateLogic {
	return &TenantAppModuleMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppModuleMultiCreateLogic) TenantAppModuleMultiCreate(in *sys.TenantAppCreateReq) (*sys.Response, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	uc := ctxs.GetUserCtx(l.ctx)
	uc.AllTenant = true
	defer func() { uc.AllTenant = false }()
	conn := stores.GetTenantConn(l.ctx)
	err := conn.Transaction(func(tx *gorm.DB) error {
		for _, module := range in.Modules {
			err := ModuleCreate(l.ctx, tx, in.Code, in.AppCode, module.Code, module.MenuIDs, module.ApiIDs)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return &sys.Response{}, err
}

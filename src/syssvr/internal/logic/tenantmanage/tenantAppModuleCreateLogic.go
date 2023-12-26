package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppModuleCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppModuleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppModuleCreateLogic {
	return &TenantAppModuleCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppModuleCreateLogic) TenantAppModuleCreate(in *sys.TenantModuleCreateReq) (*sys.Response, error) {
	conn := stores.GetTenantConn(l.ctx)
	err := conn.Transaction(func(tx *gorm.DB) error {
		err := ModuleCreate(l.ctx, tx, in.Code, in.AppCode, in.ModuleCode, in.MenuIDs, in.ApiIDs)
		return err
	})
	return &sys.Response{}, err
}

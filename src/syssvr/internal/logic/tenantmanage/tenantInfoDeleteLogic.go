package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"gorm.io/gorm"

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
	f := relationDB.TenantInfoFilter{ID: in.Id}
	if in.Code != "" {
		f.Codes = []string{in.Code}
	}
	conn := stores.GetTenantConn(l.ctx)
	err := conn.Transaction(func(tx *gorm.DB) error {
		tir := relationDB.NewTenantInfoRepo(tx)
		ti, err := tir.FindOneByFilter(l.ctx, f)
		if err != nil {
			return err
		}
		err = relationDB.NewAppInfoRepo(tx).DeleteByFilter(l.ctx, relationDB.AppInfoFilter{Codes: []string{ti.Code}})
		if err != nil {
			return err
		}
		err = relationDB.NewTenantAppModuleRepo(tx).DeleteByFilter(l.ctx, relationDB.TenantAppModuleFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewTenantAppApiRepo(tx).DeleteByFilter(l.ctx, relationDB.TenantAppApiFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewTenantAppMenuRepo(tx).DeleteByFilter(l.ctx, relationDB.TenantAppMenuFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewUserInfoRepo(tx).DeleteByFilter(l.ctx, relationDB.UserInfoFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewRoleModuleRepo(tx).DeleteByFilter(l.ctx, relationDB.RoleModuleFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewRoleApiRepo(tx).DeleteByFilter(l.ctx, relationDB.RoleApiFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewRoleApiRepo(tx).DeleteByFilter(l.ctx, relationDB.RoleApiFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewRoleMenuRepo(tx).DeleteByFilter(l.ctx, relationDB.RoleMenuFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewTenantInfoRepo(l.ctx).DeleteByFilter(l.ctx, f)
		if err != nil {
			return err
		}
		return nil
	})

	return &sys.Response{}, err
}

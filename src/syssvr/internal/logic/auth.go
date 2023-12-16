package logic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
)

// 租户号为空则为该用户的租户
func IsSupperAdmin(ctx context.Context, tenantCode string) error {
	uc := ctxs.GetUserCtx(ctx)
	if tenantCode == "" {
		tenantCode = uc.TenantCode
	}
	if uc.TenantCode != tenantCode {
		return errors.Permissions.AddMsgf("只有%s的超管才有权限", tenantCode)
	}
	ti, err := relationDB.NewTenantInfoRepo(ctx).FindOneByFilter(ctx, relationDB.TenantInfoFilter{Codes: []string{uc.TenantCode}})
	if err != nil {
		return err
	}
	if ti.AdminUserID != uc.UserID {
		return errors.Permissions.AddMsgf("只有%s的超管才有权限", tenantCode)
	}
	return nil
}

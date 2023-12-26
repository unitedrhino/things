package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
)

func CheckModule(ctx context.Context, tenantCode, appCode, moduleCode string) error {
	c, err := relationDB.NewTenantAppModuleRepo(ctx).CountByFilter(ctx, relationDB.TenantAppModuleFilter{TenantCode: tenantCode, AppCodes: []string{appCode}, ModuleCodes: []string{moduleCode}})
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.Parameter.AddMsgf("moduleCode not find:%v", moduleCode)
	}
	return nil
}

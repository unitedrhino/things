package app

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToTenantAppModulePb(in *types.TenantAppModule) *sys.TenantAppModule {
	return &sys.TenantAppModule{
		Code:    in.Code,
		MenuIDs: in.MenuIDs,
	}
}
func ToTenantAppModulesPb(in []*types.TenantAppModule) (ret []*sys.TenantAppModule) {
	for _, v := range in {
		ret = append(ret, ToTenantAppModulePb(v))
	}
	return
}

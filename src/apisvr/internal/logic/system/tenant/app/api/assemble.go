package api

import (
	"github.com/i-Things/things/src/apisvr/internal/logic/system/module/api"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToTenantAppApiTypes(in *sys.TenantApiInfo) *types.TenantApiInfo {
	if in == nil || in.Info == nil {
		return nil
	}
	return &types.TenantApiInfo{
		TemplateID: in.TemplateID,
		Code:       in.Code,
		AppCode:    in.AppCode,
		ApiInfo:    *api.ToApiInfoTypes(in.Info),
	}
}
func ToTenantAppApisTypes(in []*sys.TenantApiInfo) (ret []*types.TenantApiInfo) {
	for _, v := range in {
		ret = append(ret, ToTenantAppApiTypes(v))
	}
	return
}

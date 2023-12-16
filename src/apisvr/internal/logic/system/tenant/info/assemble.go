package info

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToTenantInfoRpc(in *types.TenantInfo) *sys.TenantInfo {
	if in == nil {
		return nil
	}
	return &sys.TenantInfo{
		Id:          in.ID,
		Name:        in.Name,
		Code:        in.Code,
		AdminUserID: in.AdminUserID,
		Desc:        utils.ToRpcNullString(in.Desc),
		BaseUrl:     in.BaseUrl,
		LogoUrl:     in.LogoUrl,
	}
}

func ToTenantInfoTypes(in *sys.TenantInfo) *types.TenantInfo {
	if in == nil {
		return nil
	}
	return &types.TenantInfo{
		ID:          in.Id,
		Name:        in.Name,
		Code:        in.Code,
		AdminUserID: in.AdminUserID,
		Desc:        utils.ToNullString(in.Desc),
		BaseUrl:     in.BaseUrl,
		LogoUrl:     in.LogoUrl,
	}
}

func ToTenantInfosTypes(in []*sys.TenantInfo) []*types.TenantInfo {
	var ret []*types.TenantInfo
	for _, v := range in {
		ret = append(ret, ToTenantInfoTypes(v))
	}
	return ret
}

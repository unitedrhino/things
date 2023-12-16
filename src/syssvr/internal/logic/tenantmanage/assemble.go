package tenantmanagelogic

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToTenantInfoRpc(in *relationDB.SysTenantInfo) *sys.TenantInfo {
	if in == nil {
		return nil
	}
	return &sys.TenantInfo{
		Id:          in.ID,
		Code:        in.Code,
		Name:        in.Name,
		AdminUserID: in.AdminUserID,
		BaseUrl:     in.BaseUrl,
		LogoUrl:     in.LogoUrl,
		Desc:        utils.ToRpcNullString(in.Desc),
	}
}
func ToTenantInfosRpc(in []*relationDB.SysTenantInfo) (ret []*sys.TenantInfo) {
	for _, v := range in {
		ret = append(ret, ToTenantInfoRpc(v))
	}
	return
}

func ToTenantInfoPo(in *sys.TenantInfo) *relationDB.SysTenantInfo {
	if in == nil {
		return nil
	}
	return &relationDB.SysTenantInfo{
		ID:          in.Id,
		Code:        in.Code,
		Name:        in.Name,
		AdminUserID: in.AdminUserID,
		BaseUrl:     in.BaseUrl,
		LogoUrl:     in.LogoUrl,
		Desc:        utils.ToEmptyString(in.Desc),
	}
}

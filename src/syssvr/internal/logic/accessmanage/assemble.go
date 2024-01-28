package accessmanagelogic

import (
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToApiInfoPo(in *sys.ApiInfo) *relationDB.SysApiInfo {
	if in == nil {
		return nil
	}
	return &relationDB.SysApiInfo{
		ID:           in.Id,
		AccessCode:   in.AccessCode,
		Method:       in.Method,
		Route:        in.Route,
		Name:         in.Name,
		BusinessType: in.BusinessType,
		Desc:         in.Desc,
		IsAuthTenant: in.IsAuthTenant,
	}
}

func ToApiInfoPb(in *relationDB.SysApiInfo) *sys.ApiInfo {
	if in == nil {
		return nil
	}
	return &sys.ApiInfo{
		Id:           in.ID,
		AccessCode:   in.AccessCode,
		Method:       in.Method,
		Route:        in.Route,
		Name:         in.Name,
		BusinessType: in.BusinessType,
		Desc:         in.Desc,
		IsAuthTenant: in.IsAuthTenant,
	}
}
func ToAccessPo(in *sys.AccessInfo) *relationDB.SysAccessInfo {
	if in == nil {
		return nil
	}
	return &relationDB.SysAccessInfo{
		ID:         in.Id,
		Name:       in.Name,
		Code:       in.Code,
		Group:      in.Group,
		IsNeedAuth: in.IsNeedAuth,
		Desc:       in.Desc,
	}
}

func ToAccessPb(in *relationDB.SysAccessInfo) *sys.AccessInfo {
	if in == nil {
		return nil
	}
	var apis []*sys.ApiInfo
	if len(in.Apis) != 0 {
		for _, v := range in.Apis {
			apis = append(apis, ToApiInfoPb(v))
		}
	}
	return &sys.AccessInfo{
		Id:         in.ID,
		Name:       in.Name,
		Code:       in.Code,
		Group:      in.Group,
		IsNeedAuth: in.IsNeedAuth,
		Desc:       in.Desc,
		Apis:       apis,
	}
}

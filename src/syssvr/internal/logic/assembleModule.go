package logic

import (
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToModuleInfoPo(in *sys.ModuleInfo) *relationDB.SysModuleInfo {
	if in == nil {
		return nil
	}
	return &relationDB.SysModuleInfo{
		ID:         in.Id,
		Code:       in.Code,
		Type:       in.Type,
		SubType:    in.SubType,
		Order:      in.Order,
		Name:       in.Name,
		Path:       in.Path,
		Url:        in.Url,
		Icon:       in.Icon,
		Body:       in.Body.GetValue(),
		HideInMenu: in.HideInMenu,
		Desc:       in.Desc.GetValue(),
	}
}

func ToModuleInfoPb(in *relationDB.SysModuleInfo) *sys.ModuleInfo {
	if in == nil {
		return nil
	}
	return &sys.ModuleInfo{
		Id:         in.ID,
		Code:       in.Code,
		Type:       in.Type,
		SubType:    in.SubType,
		Order:      in.Order,
		Name:       in.Name,
		Path:       in.Path,
		Url:        in.Url,
		Icon:       in.Icon,
		Body:       utils.ToRpcNullString(in.Body),
		HideInMenu: in.HideInMenu,
		Desc:       utils.ToRpcNullString(in.Desc),
	}
}

func ToModuleInfosPb(in []*relationDB.SysModuleInfo) (ret []*sys.ModuleInfo) {
	for _, v := range in {
		ret = append(ret, ToModuleInfoPb(v))
	}
	return
}

func ToApiInfoPo(in *sys.ApiInfo) *relationDB.SysApiInfo {
	if in == nil {
		return nil
	}
	return &relationDB.SysApiInfo{
		ID:           in.Id,
		AccessCode:   in.AccessCode,
		Route:        in.Route,
		Method:       in.Method,
		Name:         in.Name,
		BusinessType: in.BusinessType,
		IsAuthTenant: in.IsAuthTenant,
		Desc:         in.Desc,
	}
}

//func ToTenantApiInfoPo(in *sys.TenantApiInfo) *relationDB.SysTenantAppApi {
//	if in == nil || in.Info == nil {
//		return nil
//	}
//	return &relationDB.SysTenantAppApi{
//		TempLateID:   in.TemplateID,
//		TenantCode:   stores.TenantCode(in.Code),
//		AppCode:      in.AppCode,
//		SysModuleApi: *ToApiInfoPo(in.Info),
//	}
//}

func ToApiInfoPb(in *relationDB.SysApiInfo) *sys.ApiInfo {
	if in == nil {
		return nil
	}
	return &sys.ApiInfo{
		Id:           in.ID,
		AccessCode:   in.AccessCode,
		Route:        in.Route,
		Method:       in.Method,
		Name:         in.Name,
		BusinessType: in.BusinessType,
		IsAuthTenant: in.IsAuthTenant,
		Desc:         in.Desc,
	}
}

func ToMenuInfoPo(in *sys.MenuInfo) *relationDB.SysModuleMenu {
	if in == nil {
		return nil
	}
	return &relationDB.SysModuleMenu{
		ID:         in.Id,
		ModuleCode: in.ModuleCode,
		ParentID:   in.ParentID,
		Type:       in.Type,
		Order:      in.Order,
		Name:       in.Name,
		Path:       in.Path,
		Component:  in.Component,
		Icon:       in.Icon,
		Redirect:   in.Redirect,
		Body:       in.Body.GetValue(),
		HideInMenu: in.HideInMenu,
	}
}

func ToTenantAppMenuPo(in *sys.TenantAppMenu) *relationDB.SysTenantAppMenu {
	if in == nil || in.Info == nil {
		return nil
	}
	return &relationDB.SysTenantAppMenu{
		TempLateID:    in.TemplateID,
		TenantCode:    stores.TenantCode(in.Code),
		AppCode:       in.AppCode,
		SysModuleMenu: *ToMenuInfoPo(in.Info),
	}
}

func ToMenuInfoPb(in *relationDB.SysModuleMenu) *sys.MenuInfo {
	if in == nil {
		return nil
	}
	return &sys.MenuInfo{
		Id:         in.ID,
		ModuleCode: in.ModuleCode,
		ParentID:   in.ParentID,
		Type:       in.Type,
		Order:      in.Order,
		Name:       in.Name,
		Path:       in.Path,
		Component:  in.Component,
		Icon:       in.Icon,
		Redirect:   in.Redirect,
		Body:       utils.ToRpcNullString(in.Body),
		HideInMenu: in.HideInMenu,
	}
}

func ToTenantAppMenuInfoPb(in *relationDB.SysTenantAppMenu) *sys.TenantAppMenu {
	if in == nil {
		return nil
	}
	return &sys.TenantAppMenu{
		TemplateID: in.TempLateID,
		Code:       string(in.TenantCode),
		AppCode:    in.AppCode,
		Info:       ToMenuInfoPb(&in.SysModuleMenu),
	}
}

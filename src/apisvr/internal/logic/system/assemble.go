package system

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ProjectInfoToApi(pb *sys.ProjectInfo) *types.ProjectInfo {
	return &types.ProjectInfo{
		CreatedTime: pb.CreatedTime,
		ProjectID:   pb.ProjectID,
		ProjectName: pb.ProjectName,
		AdminUserID: pb.AdminUserID,
		Desc:        utils.ToNullString(pb.Desc),
	}
}
func ProjectInfosToApi(pb []*sys.ProjectInfo) (ret []*types.ProjectInfo) {
	for _, v := range pb {
		ret = append(ret, ProjectInfoToApi(v))
	}
	return
}

func ToMenuInfoApi(i *sys.MenuInfo) *types.MenuInfo {
	return &types.MenuInfo{
		ModuleCode: i.ModuleCode,
		ID:         i.Id,
		Name:       i.Name,
		ParentID:   i.ParentID,
		Type:       i.Type,
		Path:       i.Path,
		Component:  i.Component,
		Icon:       i.Icon,
		Redirect:   i.Redirect,
		CreateTime: i.CreateTime,
		Order:      i.Order,
		HideInMenu: i.HideInMenu,
		Body:       utils.ToNullString(i.Body),
		Children:   ToMenuInfosApi(i.Children),
	}
}
func ToMenuInfosApi(i []*sys.MenuInfo) (ret []*types.MenuInfo) {
	if i == nil {
		return nil
	}
	for _, v := range i {
		ret = append(ret, ToMenuInfoApi(v))
	}
	return
}

func ToTenantAppMenuApi(i *sys.TenantAppMenu) *types.TenantAppMenu {
	if i == nil {
		return nil
	}
	return &types.TenantAppMenu{
		TemplateID: i.TemplateID,
		Code:       i.Code,
		AppCode:    i.AppCode,
		MenuInfo:   *ToMenuInfoApi(i.Info),
		Children:   ToTenantAppMenusApi(i.Children),
	}
}
func ToTenantAppMenusApi(i []*sys.TenantAppMenu) (ret []*types.TenantAppMenu) {
	for _, v := range i {
		ret = append(ret, ToTenantAppMenuApi(v))
	}
	return
}

func ToSysWithIDCode(in *types.WithIDOrCode) *sys.WithIDCode {
	return &sys.WithIDCode{
		Id:   in.ID,
		Code: in.Code,
	}
}

func ToTenantInfoRpc(in *types.TenantInfo) *sys.TenantInfo {
	if in == nil {
		return nil
	}
	return &sys.TenantInfo{
		Id:          in.ID,
		Name:        in.Name,
		Code:        in.Code,
		AdminUserID: in.AdminUserID,
		AdminRoleID: in.AdminRoleID,
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
		AdminRoleID: in.AdminRoleID,
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

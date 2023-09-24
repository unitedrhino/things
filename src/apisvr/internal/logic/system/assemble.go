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
		CompanyName: utils.ToNullString(pb.CompanyName),
		UserID:      pb.UserID,
		Region:      utils.ToNullString(pb.Region),
		Address:     utils.ToNullString(pb.Address),
		Desc:        utils.ToNullString(pb.Desc),
	}
}

func ToMenuInfoApi(i *sys.MenuData) *types.MenuData {
	return &types.MenuData{
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
	}
}

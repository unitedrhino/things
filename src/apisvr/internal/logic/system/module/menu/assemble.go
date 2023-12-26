package menu

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToMenuInfoRpc(in *types.MenuInfo) *sys.MenuInfo {
	if in == nil {
		return nil
	}
	return &sys.MenuInfo{
		Id:         in.ID,
		Name:       in.Name,
		ParentID:   in.ParentID,
		Type:       in.Type,
		Path:       in.Path,
		Component:  in.Component,
		Icon:       in.Icon,
		Redirect:   in.Redirect,
		Order:      in.Order,
		HideInMenu: in.HideInMenu,
		Body:       utils.ToRpcNullString(in.Body),
		ModuleCode: in.ModuleCode,
	}
}
func ToMenuInfosRpc(in []*types.MenuInfo) (ret []*sys.MenuInfo) {
	for _, v := range in {
		ret = append(ret, ToMenuInfoRpc(v))
	}
	return
}

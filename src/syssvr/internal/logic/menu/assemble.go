package menulogic

import (
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func MenuInfoToPb(ui *mysql.MenuInfo) *sys.MenuData {
	return &sys.MenuData{
		Id:         ui.Id,
		Name:       ui.Name,
		ParentID:   ui.ParentID.Int64,
		Type:       ui.Type.Int64,
		Path:       ui.Path,
		Component:  ui.Component,
		Icon:       ui.Icon,
		Redirect:   ui.Redirect,
		CreateTime: ui.CreatedTime.Unix(),
		Order:      ui.Order.Int64,
		HideInMenu: ui.HideInMenu,
	}
}

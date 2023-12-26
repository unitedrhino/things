package logic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

func findMissingParentIds(menuInfos []*relationDB.SysModuleMenu) map[int64]bool {
	missingParentIds := make(map[int64]bool)
	ids := make(map[int64]bool)
	for _, menu := range menuInfos {
		ids[menu.ID] = true
	}
	for _, menu := range menuInfos {
		if !ids[menu.ParentID] && menu.ParentID != def.RootNode {
			missingParentIds[menu.ParentID] = true
		}
	}
	return missingParentIds
}

func CheckMissingParentIdMenuIndex(ctx context.Context, menuInfos []*relationDB.SysModuleMenu) []*relationDB.SysModuleMenu {
	var MenuInfos []*relationDB.SysModuleMenu
	missingParentIds := findMissingParentIds(menuInfos)
	if len(missingParentIds) > 0 {
		menuIDs := lo.Keys(missingParentIds)
		menuInfo, err := relationDB.NewMenuInfoRepo(ctx).FindByFilter(ctx, relationDB.MenuInfoFilter{MenuIDs: menuIDs}, nil)
		if err != nil {
			logx.WithContext(ctx).Errorf("MenuIndex find menu_info err,menuIds:%d,err:%v", menuIDs, err)
			return MenuInfos
		}
		MenuInfos = append(MenuInfos, menuInfo...)
	}

	return MenuInfos
}

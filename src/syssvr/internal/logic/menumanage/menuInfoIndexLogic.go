package menumanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/samber/lo"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	MiDB *relationDB.MenuInfoRepo
	RiDB *relationDB.RoleInfoRepo
}

func NewMenuInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuInfoIndexLogic {
	return &MenuInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		MiDB:   relationDB.NewMenuInfoRepo(ctx),
		RiDB:   relationDB.NewRoleInfoRepo(ctx),
	}
}

func findMissingParentIds(menuInfos []*relationDB.SysMenuInfo) map[int64]bool {
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

func (l *MenuInfoIndexLogic) checkMissingParentIdMenuIndex(menuInfos []*relationDB.SysMenuInfo) []*relationDB.SysMenuInfo {
	var MenuInfos []*relationDB.SysMenuInfo
	missingParentIds := findMissingParentIds(menuInfos)
	if len(missingParentIds) > 0 {
		menuIDs := lo.Keys(missingParentIds)
		menuInfo, err := l.MiDB.FindByFilter(l.ctx, relationDB.MenuInfoFilter{MenuIds: menuIDs}, nil)
		if err != nil {
			l.Errorf("MenuIndex find menu_info err,menuIds:%d,err:%v", menuIDs, err)
			return MenuInfos
		}
		MenuInfos = append(MenuInfos, menuInfo...)
	}

	return MenuInfos
}

func (l *MenuInfoIndexLogic) MenuInfoIndex(in *sys.MenuInfoIndexReq) (*sys.MenuInfoIndexResp, error) {
	info := make([]*sys.MenuInfo, 0)

	//获取完整菜单列表
	mes, err := l.MiDB.FindByFilter(l.ctx, relationDB.MenuInfoFilter{
		AppCode: in.AppCode,
	}, nil)
	if err != nil {
		return nil, err
	}
	if !in.IsRetTree {
		for _, me := range mes {
			info = append(info, logic.MenuInfoToPb(me))
		}
		return &sys.MenuInfoIndexResp{
			List: info,
		}, nil
	}
	var idMap = make(map[int64]*sys.MenuInfo, len(mes))
	var pidMap = make(map[int64][]*sys.MenuInfo, len(mes))

	for _, me := range mes {
		idMap[me.ID] = logic.MenuInfoToPb(me)
		if me.ParentID == def.RootNode { //root节点
			info = append(info, idMap[me.ID])
		}
	}
	for _, me := range idMap {
		pidMap[me.ParentID] = append(pidMap[me.ParentID], me)
	}
	for _, v := range info {
		FillChildren(v, pidMap)
	}
	return &sys.MenuInfoIndexResp{
		List: info,
	}, nil
}
func FillChildren(node *sys.MenuInfo, pidMap map[int64][]*sys.MenuInfo) {
	cs := pidMap[node.Id]
	if len(cs) != 0 { //如果子节点不为空
		node.Children = cs
		for _, c := range cs {
			FillChildren(c, pidMap)
		}
	}
}

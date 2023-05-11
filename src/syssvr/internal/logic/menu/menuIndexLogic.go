package menulogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuIndexLogic {
	return &MenuIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func findMissingParentIds(menuInfos []*mysql.SysMenuInfo) map[int64]bool {
	missingParentIds := make(map[int64]bool)
	ids := make(map[int64]bool)
	for _, menu := range menuInfos {
		ids[menu.Id] = true
	}
	for _, menu := range menuInfos {
		if !ids[menu.ParentID] && menu.ParentID != 1 {
			missingParentIds[menu.ParentID] = true
		}
	}
	return missingParentIds
}

func (l *MenuIndexLogic) checkMissingParentIdMenuIndex(menuInfos []*mysql.SysMenuInfo) []*mysql.SysMenuInfo {
	var MenuInfos []*mysql.SysMenuInfo
	missingParentIds := findMissingParentIds(menuInfos)
	if len(missingParentIds) > 0 {
		for k, _ := range missingParentIds {
			menuInfo, err := l.svcCtx.MenuInfoModle.FindOne(l.ctx, k)
			if err != nil {
				l.Errorf("MenuIndex find menu_info err,menuIds:%d,err:%v", k, err)
				continue
			}
			MenuInfos = append(MenuInfos, menuInfo)
		}
	}

	return MenuInfos
}

func (l *MenuIndexLogic) MenuIndex(in *sys.MenuIndexReq) (*sys.MenuIndexResp, error) {
	info := make([]*sys.MenuData, 0)
	if in.Role != 0 {
		//获取角色关联的菜单列表
		menuIds, err := l.svcCtx.RoleModel.IndexRoleIDMenuID(in.Role)
		if err != nil {
			return nil, errors.Database.AddDetail(err)
		}
		menuInfos, err := l.svcCtx.MenuModel.Index(l.ctx, &mysql.MenuIndexFilter{MenuIds: menuIds})
		if err != nil {
			l.Errorf("MenuIndex find menu_info err,menuIds:%v,err:%v", menuIds, err)
			return nil, errors.Database.AddDetail(err)
		}
		for _, v := range menuInfos {
			info = append(info, MenuInfoToPb(v))
		}
		//查看缺失的父菜单Id
		missingMenuInfos := l.checkMissingParentIdMenuIndex(menuInfos)
		if len(missingMenuInfos) > 0 {
			for _, v := range missingMenuInfos {
				info = append(info, MenuInfoToPb(v))
			}
		}
	} else {
		//获取完整菜单列表
		mes, err := l.svcCtx.MenuModel.Index(l.ctx, &mysql.MenuIndexFilter{
			Name: in.Name,
			Path: in.Path,
		})
		if err != nil {
			return nil, errors.Database.AddDetail(err)
		}
		for _, me := range mes {
			info = append(info, MenuInfoToPb(me))
		}
	}

	return &sys.MenuIndexResp{
		List: info,
	}, nil
}

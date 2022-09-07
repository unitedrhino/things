package menulogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"

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

func (l *MenuIndexLogic) MenuIndex(in *sys.MenuIndexReq) (*sys.MenuIndexResp, error) {
	info := make([]*sys.MenuData, 0)
	if in.Role != 0 {
		//获取角色关联的菜单列表
		menuIds, err := l.svcCtx.RoleModel.IndexRoleIDMenuID(in.Role)
		if err != nil {
			return nil, errors.Database.AddDetail(err)
		}
		for _, v := range menuIds {
			menuInfo, err := l.svcCtx.MenuInfoModle.FindOne(l.ctx, v)
			if err != nil {
				return nil, errors.Database.AddDetail(err)
			}
			info = append(info, &sys.MenuData{
				Id:         menuInfo.Id,
				Name:       menuInfo.Name,
				ParentID:   menuInfo.ParentID,
				Type:       menuInfo.Type,
				Path:       menuInfo.Path,
				Component:  menuInfo.Component,
				Icon:       menuInfo.Icon,
				Redirect:   menuInfo.Redirect,
				CreateTime: menuInfo.CreatedTime.Unix(),
				Order:      menuInfo.Order,
			})
		}
	} else {
		//获取完整菜单列表
		mes, err := l.svcCtx.MenuModel.Index(in)
		if err != nil {
			return nil, errors.Database.AddDetail(err)
		}
		for _, me := range mes {
			info = append(info, &sys.MenuData{
				Id:         me.Id,
				Name:       me.Name,
				ParentID:   me.ParentID,
				Type:       me.Type,
				Path:       me.Path,
				Component:  me.Component,
				Icon:       me.Icon,
				Redirect:   me.Redirect,
				CreateTime: me.CreatedTime.Unix(),
				Order:      me.Order,
			})
		}
	}

	return &sys.MenuIndexResp{
		List: info,
	}, nil

	return &sys.MenuIndexResp{}, nil
}

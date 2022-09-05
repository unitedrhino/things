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
	info := make([]*sys.MenuIndexData, 0)
	var total int64
	if in.RoleFlag == 1 {
		//获取角色关联的菜单列表
		menuids, err := l.svcCtx.RoleModel.IndexRoleIDMenuID(in.Role)
		if err != nil {
			return nil, errors.Database.AddDetail(err)
		}
		total = int64(len(menuids))
		for _, v := range menuids {
			menuinfo, err := l.svcCtx.MenuInfoModle.FindOne(l.ctx, v)
			if err != nil {
				return nil, errors.Database.AddDetail(err)
			}
			info = append(info, &sys.MenuIndexData{
				Id:         menuinfo.Id,
				Name:       menuinfo.Name,
				ParentID:   menuinfo.ParentID,
				Type:       menuinfo.Type,
				Path:       menuinfo.Path,
				Component:  menuinfo.Component,
				Icon:       menuinfo.Icon,
				Redirect:   menuinfo.Redirect,
				CreateTime: menuinfo.CreatedTime.Unix(),
				Order:      menuinfo.Order,
			})
		}
	} else {
		//获取完整菜单列表
		mes, count, err := l.svcCtx.MenuModel.Index(in)
		if err != nil {
			return nil, errors.Database.AddDetail(err)
		}
		total = count
		for _, me := range mes {
			info = append(info, &sys.MenuIndexData{
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
		List:  info,
		Total: total,
	}, nil

	return &sys.MenuIndexResp{}, nil
}

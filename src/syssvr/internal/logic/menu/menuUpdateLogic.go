package menulogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuUpdateLogic {
	return &MenuUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MenuUpdateLogic) MenuUpdate(in *sys.MenuUpdateReq) (*sys.Response, error) {
	mi, err := l.svcCtx.MenuInfoModel.FindOne(l.ctx, in.Id)
	if err != nil {
		l.Logger.Error("UserInfoModel.FindOne err , sql:%s", l.svcCtx)
		return nil, err
	}
	if in.Type != 1 && in.Type != 2 && in.Type != 3 {
		in.Type = mi.Type
	}
	if in.Order == 0 {
		in.Order = mi.Order
	}
	if in.HideInMenu == 0 {
		in.HideInMenu = mi.HideInMenu
	}

	if in.ParentID == 0 {
		in.ParentID = mi.ParentID
	}

	err = l.svcCtx.MenuInfoModel.Update(l.ctx, &mysql.SysMenuInfo{
		Id:            in.Id,
		ParentID:      in.ParentID,
		Type:          in.Type,
		Order:         in.Order,
		Name:          in.Name,
		Path:          in.Path,
		Component:     in.Component,
		Icon:          in.Icon,
		Redirect:      in.Redirect,
		BackgroundUrl: "",
		HideInMenu:    in.HideInMenu,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &sys.Response{}, nil
}

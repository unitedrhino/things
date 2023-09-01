package menulogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	MiDB *relationDB.MenuInfoRepo
}

func NewMenuUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuUpdateLogic {
	return &MenuUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		MiDB:   relationDB.NewMenuInfoRepo(ctx),
	}
}

func (l *MenuUpdateLogic) MenuUpdate(in *sys.MenuUpdateReq) (*sys.Response, error) {
	mi, err := l.MiDB.FindOne(l.ctx, in.Id)
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

	err = l.MiDB.Update(l.ctx, &relationDB.SysMenuInfo{
		ID:            in.Id,
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
		return nil, err
	}
	return &sys.Response{}, nil
}

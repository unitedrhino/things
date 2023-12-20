package menumanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	MiDB *relationDB.MenuInfoRepo
}

func NewMenuInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuInfoCreateLogic {
	return &MenuInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		MiDB:   relationDB.NewMenuInfoRepo(ctx),
	}
}

func (l *MenuInfoCreateLogic) MenuInfoCreate(in *sys.MenuInfo) (*sys.WithID, error) {

	if in.Type == 0 {
		in.Type = 1
	}
	if in.ParentID == 0 {
		in.ParentID = 1
	}
	if in.Order == 0 {
		in.Order = 1
	}
	if in.HideInMenu == 0 {
		in.HideInMenu = 1
	}
	po := relationDB.SysMenuInfo{
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
	}
	err := l.MiDB.Insert(l.ctx, &po)
	if err != nil {
		return nil, err
	}
	return &sys.WithID{Id: po.ID}, nil
}

package menumanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	MiDB *relationDB.MenuInfoRepo
}

func NewMenuInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuInfoUpdateLogic {
	return &MenuInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		MiDB:   relationDB.NewMenuInfoRepo(ctx),
	}
}

func (l *MenuInfoUpdateLogic) MenuInfoUpdate(in *sys.MenuInfo) (*sys.Response, error) {
	mi, err := l.MiDB.FindOne(l.ctx, in.Id)
	if err != nil {
		l.Logger.Error("UserInfoModel.FindOne err , sql:%s", l.svcCtx)
		return nil, err
	}
	if in.Type != 0 {
		mi.Type = in.Type
	}
	if in.Order != 0 {
		mi.Order = in.Order
	}
	if in.HideInMenu != 0 {
		mi.HideInMenu = in.HideInMenu
	}

	if in.ParentID != 0 {
		mi.ParentID = in.ParentID
	}
	if in.Component != "" {
		mi.Component = in.Component
	}
	if in.Name != "" {
		mi.Name = in.Name
	}
	if in.Icon != "" {
		mi.Icon = in.Icon
	}
	if in.Path != "" {
		mi.Path = in.Path
	}

	err = l.MiDB.Update(l.ctx, mi)
	if err != nil {
		return nil, err
	}
	return &sys.Response{}, nil
}

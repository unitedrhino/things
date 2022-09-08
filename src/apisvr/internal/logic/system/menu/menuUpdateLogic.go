package menu

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMenuUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuUpdateLogic {
	return &MenuUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuUpdateLogic) MenuUpdate(req *types.MenuUpdateReq) error {
	resp, err := l.svcCtx.MenuRpc.MenuUpdate(l.ctx, &sys.MenuUpdateReq{
		Id:         req.ID,
		Name:       req.Name,
		ParentID:   req.ParentID,
		Type:       req.Type,
		Path:       req.Path,
		Component:  req.Component,
		Icon:       req.Icon,
		Redirect:   req.Redirect,
		Order:      req.Order,
		HideInMenu: req.HideInMenu,
	})
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.rpc.MenuUpdate req=%v err=%+v", utils.FuncName(), req, err)
		return err
	}
	if resp == nil {
		l.Errorf("%s.rpc.MenuUpdate return nil req=%+v", utils.FuncName(), req)
		return errors.System.AddDetail("MenuUpdate rpc return nil")
	}
	return nil
}

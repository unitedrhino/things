package menu

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMenuCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuCreateLogic {
	return &MenuCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuCreateLogic) MenuCreate(req *types.MenuCreateReq) error {
	resp, err := l.svcCtx.MenuRpc.MenuCreate(l.ctx, &sys.MenuCreateReq{
		Name:       req.Name,
		ParentID:   req.ParentID,
		Type:       req.Type,
		Path:       req.Path,
		Component:  req.Component,
		Icon:       req.Icon,
		Redirect:   req.Redirect,
		Order:      req.Order,
		HideInMenu: req.HideInMenu,
		Role:       ctxs.GetUserCtx(l.ctx).Role,
	})
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.rpc.MenuCreate req=%v err=%+v", utils.FuncName(), req, err)
		return err
	}
	if resp == nil {
		l.Errorf("%s rpc.MenuCreate return nil req=%+v", utils.FuncName(), req)
		return errors.System.AddDetail("MenuCreate rpc return nil")
	}
	return nil
}

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

type MulitUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMulitUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MulitUpdateLogic {
	return &MulitUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MulitUpdateLogic) MulitUpdate(req *types.RoleMenuMultiUpdateReq) error {
	resp, err := l.svcCtx.RoleRpc.RoleMenuMultiUpdate(l.ctx, &sys.RoleMenuMultiUpdateReq{
		Id:         req.ID,
		AppCode:    req.AppCode,
		ModuleCode: req.ModuleCode,
		MenuIDs:    req.MenuIDs,
	})
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.rpc.RoleMenuUpdate req=%v err=%+v", utils.FuncName(), req, err)
		return err
	}
	if resp == nil {
		l.Errorf("%s.rpc.RoleMenuUpdate return nil req=%+v", utils.FuncName(), req)
		return errors.System.AddDetail("RoleMenuUpdate rpc return nil")
	}
	return nil
}

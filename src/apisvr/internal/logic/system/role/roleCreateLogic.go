package role

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleCreateLogic {
	return &RoleCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleCreateLogic) RoleCreate(req *types.RoleCreateReq) error {
	resp, err := l.svcCtx.RoleRpc.RoleCreate(l.ctx, &sys.RoleCreateReq{
		Name:   req.Name,
		Remark: req.Remark,
		Status: req.Status,
	})
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("[%s]|rpc.RoleCreate|req=%v|err=%+v", utils.FuncName(), req, err)
		return err
	}
	if resp == nil {
		l.Errorf("%s|rpc.RoleCreate|return nil|req=%+v", utils.FuncName(), req)
		return errors.System.AddDetail("RoleCreate rpc return nil")
	}
	return nil
}

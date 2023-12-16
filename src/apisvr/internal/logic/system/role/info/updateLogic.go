package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.RoleInfo) error {
	resp, err := l.svcCtx.RoleRpc.RoleInfoUpdate(l.ctx, &sys.RoleInfo{
		Id:     req.ID,
		Name:   req.Name,
		Desc:   req.Desc,
		Status: req.Status,
	})
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.rpc.RoleUpdate req=%v err=%v", utils.FuncName(), req, err)
		return err
	}
	if resp == nil {
		l.Errorf("%s.rpc.RoleUpdate return nil req=%v", utils.FuncName(), req)
		return errors.System.AddDetail("RoleUpdate rpc return nil")
	}
	return nil
}

package api

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"

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

func (l *UpdateLogic) Update(req *types.ApiInfo) error {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return err
	}
	_, err := l.svcCtx.ModuleRpc.ModuleApiUpdate(l.ctx, ToApiInfoRpc(req))
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.rpc.ApiUpdate req=%v err=%+v", utils.FuncName(), req, err)
		return err
	}
	return nil
}

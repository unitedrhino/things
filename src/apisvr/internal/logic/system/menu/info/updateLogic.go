package info

import (
	"context"
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

func (l *UpdateLogic) Update(req *types.MenuInfo) error {
	_, err := l.svcCtx.MenuRpc.MenuInfoUpdate(l.ctx, ToMenuInfoRpc(req))
	if err != nil {
		l.Errorf("%s.rpc.MenuUpdate req=%v err=%+v", utils.FuncName(), req, err)
		return err
	}
	return nil
}

package info

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

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

func (l *UpdateLogic) Update(req *types.SceneInfoUpdateReq) error {
	_, err := l.svcCtx.Rule.SceneInfoUpdate(l.ctx, ToScenePb(&req.SceneInfo))
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.SceneInfoUpdate req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}

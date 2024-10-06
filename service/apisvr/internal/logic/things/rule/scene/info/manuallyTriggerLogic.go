package info

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/udsvr/pb/ud"

	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ManuallyTriggerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManuallyTriggerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManuallyTriggerLogic {
	return &ManuallyTriggerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ManuallyTriggerLogic) ManuallyTrigger(req *types.WithID) error {
	_, err := l.svcCtx.Rule.SceneManuallyTrigger(l.ctx, &ud.WithID{Id: req.ID})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.SceneManuallyTrigger req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}

package scene

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/rulesvr/pb/rule"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiUpdateLogic {
	return &MultiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiUpdateLogic) MultiUpdate(req *types.AlarmSceneMultiUpdateReq) error {
	_, err := l.svcCtx.Alarm.AlarmSceneMultiUpdate(l.ctx, &rule.AlarmSceneMultiUpdateReq{
		AlarmID:  req.AlarmID,
		SceneIDs: req.SceneIDs,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.AlarmSceneMultiUpdate req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}

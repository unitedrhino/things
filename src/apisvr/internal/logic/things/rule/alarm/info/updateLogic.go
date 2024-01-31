package info

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

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

func (l *UpdateLogic) Update(req *types.AlarmInfoUpdateReq) error {
	_, err := l.svcCtx.Alarm.AlarmInfoUpdate(l.ctx, &rule.AlarmInfo{
		Id:     req.ID,
		Name:   req.Name,
		Status: req.Status,
		Desc:   req.Desc,
		Level:  req.Level,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.AlarmInfoUpdate req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}

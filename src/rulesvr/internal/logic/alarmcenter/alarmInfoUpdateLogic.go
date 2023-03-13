package alarmcenterlogic

import (
	"context"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoUpdateLogic {
	return &AlarmInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmInfoUpdateLogic) AlarmInfoUpdate(in *rule.AlarmInfo) (*rule.Response, error) {
	// todo: add your logic here and delete this line

	return &rule.Response{}, nil
}

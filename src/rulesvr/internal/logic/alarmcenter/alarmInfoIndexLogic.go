package alarmcenterlogic

import (
	"context"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoIndexLogic {
	return &AlarmInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmInfoIndexLogic) AlarmInfoIndex(in *rule.AlarmInfoIndexReq) (*rule.Response, error) {
	// todo: add your logic here and delete this line

	return &rule.Response{}, nil
}

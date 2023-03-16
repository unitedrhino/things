package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmLogIndexLogic {
	return &AlarmLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 告警日志
func (l *AlarmLogIndexLogic) AlarmLogIndex(in *rule.AlarmLogIndexReq) (*rule.AlarmLogIndexResp, error) {

	return &rule.AlarmLogIndexResp{}, nil
}

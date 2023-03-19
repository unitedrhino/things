package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/src/rulesvr/internal/domain/alarm"
	"github.com/i-Things/things/src/rulesvr/internal/logic"
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
	var (
		info []*rule.AlarmLogInfo
		size int64
		err  error
	)
	filter := alarm.LogFilter{
		Time: ToTimeRange(in.TimeRange)}
	size, err = l.svcCtx.AlarmLogRepo.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	di, err := l.svcCtx.AlarmLogRepo.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	info = make([]*rule.AlarmLogInfo, 0, len(di))
	for _, v := range di {
		info = append(info, ToAlarmLog(v))
	}
	return &rule.AlarmLogIndexResp{
		List:  info,
		Total: size,
	}, nil
}

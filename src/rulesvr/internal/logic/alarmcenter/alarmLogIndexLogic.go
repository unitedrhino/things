package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/src/rulesvr/internal/logic"
	"github.com/i-Things/things/src/rulesvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AlDB *relationDB.AlarmLogRepo
}

func NewAlarmLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmLogIndexLogic {
	return &AlarmLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AlDB:   relationDB.NewAlarmLogRepo(ctx),
	}
}

// 告警日志
func (l *AlarmLogIndexLogic) AlarmLogIndex(in *rule.AlarmLogIndexReq) (*rule.AlarmLogIndexResp, error) {
	var (
		info []*rule.AlarmLog
		size int64
		err  error
	)
	filter := relationDB.AlarmLogFilter{
		AlarmRecordID: in.AlarmRecordID,
		Time:          ToTimeRange(in.TimeRange)}
	size, err = l.AlDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	di, err := l.AlDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	info = make([]*rule.AlarmLog, 0, len(di))
	for _, v := range di {
		info = append(info, ToAlarmLog(v))
	}
	return &rule.AlarmLogIndexResp{
		List:  info,
		Total: size,
	}, nil
}

package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/service/rulesvr/internal/logic"
	"github.com/i-Things/things/service/rulesvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/rulesvr/internal/svc"
	"github.com/i-Things/things/service/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmRecordIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	ArDB *relationDB.AlarmRecordRepo
}

func NewAlarmRecordIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmRecordIndexLogic {
	return &AlarmRecordIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ArDB:   relationDB.NewAlarmRecordRepo(ctx),
	}
}

// 告警记录
func (l *AlarmRecordIndexLogic) AlarmRecordIndex(in *rule.AlarmRecordIndexReq) (*rule.AlarmRecordIndexResp, error) {
	var (
		info []*rule.AlarmRecord
		size int64
		err  error
	)
	filter := relationDB.AlarmRecordFilter{
		AlarmID: in.AlarmID,
		Time:    ToTimeRange(in.TimeRange)}
	size, err = l.ArDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	di, err := l.ArDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	info = make([]*rule.AlarmRecord, 0, len(di))
	for _, v := range di {
		info = append(info, ToAlarmRecord(v))
	}
	return &rule.AlarmRecordIndexResp{
		List:  info,
		Total: size,
	}, nil
}

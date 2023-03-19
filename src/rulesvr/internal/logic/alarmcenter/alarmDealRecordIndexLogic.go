package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/src/rulesvr/internal/domain/alarm"
	"github.com/i-Things/things/src/rulesvr/internal/logic"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmDealRecordIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmDealRecordIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmDealRecordIndexLogic {
	return &AlarmDealRecordIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmDealRecordIndexLogic) AlarmDealRecordIndex(in *rule.AlarmDealRecordIndexReq) (*rule.AlarmDealRecordIndexResp, error) {
	var (
		info []*rule.AlarmDeal
		size int64
		err  error
	)
	filter := alarm.DealRecordFilter{
		Time: ToTimeRange(in.TimeRange)}
	size, err = l.svcCtx.AlarmDealRecordRepo.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	di, err := l.svcCtx.AlarmDealRecordRepo.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	info = make([]*rule.AlarmDeal, 0, len(di))
	for _, v := range di {
		info = append(info, ToAlarmDealRecord(v))
	}
	return &rule.AlarmDealRecordIndexResp{
		List:  info,
		Total: size,
	}, nil
}

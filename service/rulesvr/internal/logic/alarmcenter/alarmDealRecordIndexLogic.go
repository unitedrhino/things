package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/service/rulesvr/internal/logic"
	"github.com/i-Things/things/service/rulesvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/rulesvr/internal/svc"
	"github.com/i-Things/things/service/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmDealRecordIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AdrDB *relationDB.AlarmDealRecordRepo
}

func NewAlarmDealRecordIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmDealRecordIndexLogic {
	return &AlarmDealRecordIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AdrDB:  relationDB.NewAlarmDealRecordRepo(ctx),
	}
}

func (l *AlarmDealRecordIndexLogic) AlarmDealRecordIndex(in *rule.AlarmDealRecordIndexReq) (*rule.AlarmDealRecordIndexResp, error) {
	var (
		info []*rule.AlarmDeal
		size int64
		err  error
	)
	filter := relationDB.AlarmDealRecordFilter{
		AlarmRecordID: in.AlarmRecordID,
		Time:          ToTimeRange(in.TimeRange)}
	size, err = l.AdrDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	di, err := l.AdrDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
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

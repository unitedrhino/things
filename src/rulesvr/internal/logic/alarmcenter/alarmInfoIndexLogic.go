package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/src/rulesvr/internal/domain/alarm"
	"github.com/i-Things/things/src/rulesvr/internal/logic"
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

func (l *AlarmInfoIndexLogic) AlarmInfoIndex(in *rule.AlarmInfoIndexReq) (*rule.AlarmInfoIndexResp, error) {
	var (
		info []*rule.AlarmInfo
		size int64
		err  error
	)
	filter := alarm.InfoFilter{
		Name: in.Name, SceneID: in.SceneID, AlarmIDs: in.AlarmIDs}
	size, err = l.svcCtx.AlarmInfoRepo.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	di, err := l.svcCtx.AlarmInfoRepo.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	info = make([]*rule.AlarmInfo, 0, len(di))
	for _, v := range di {
		info = append(info, ToAlarmInfo(v))
	}
	return &rule.AlarmInfoIndexResp{
		List:  info,
		Total: size,
	}, nil
}

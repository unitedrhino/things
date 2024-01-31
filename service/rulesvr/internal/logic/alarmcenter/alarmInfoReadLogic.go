package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/service/rulesvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/rulesvr/internal/svc"
	"github.com/i-Things/things/service/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AlarmInfoRepo
}

func NewAlarmInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoReadLogic {
	return &AlarmInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAlarmInfoRepo(ctx),
	}
}

func (l *AlarmInfoReadLogic) AlarmInfoRead(in *rule.WithID) (*rule.AlarmInfo, error) {
	di, err := l.AiDB.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return ToAlarmInfo(di), nil
}

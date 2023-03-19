package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/repo/mysql"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoCreateLogic {
	return &AlarmInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmInfoCreateLogic) AlarmInfoCreate(in *rule.AlarmInfo) (*rule.Response, error) {
	_, err := l.svcCtx.AlarmInfoRepo.FindOneByName(l.ctx, in.Name)
	if !(err == mysql.ErrNotFound) {
		return nil, errors.Parameter.AddMsg("告警名称重复").AddDetail(err)
	}
	in.LastAlarm = 0
	in.DealState = 1
	_, err = l.svcCtx.AlarmInfoRepo.Insert(l.ctx, ToAlarmInfoPo(in))
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &rule.Response{}, nil
}

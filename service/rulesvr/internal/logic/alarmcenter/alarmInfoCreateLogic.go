package alarmcenterlogic

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"github.com/i-Things/things/service/rulesvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/rulesvr/internal/svc"
	"github.com/i-Things/things/service/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AlarmInfoRepo
}

func NewAlarmInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoCreateLogic {
	return &AlarmInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAlarmInfoRepo(ctx),
	}
}

func (l *AlarmInfoCreateLogic) AlarmInfoCreate(in *rule.AlarmInfo) (*rule.WithID, error) {
	_, err := l.AiDB.FindOneByFilter(l.ctx, relationDB.AlarmInfoFilter{Name: in.Name})
	if !(errors.Cmp(err, errors.NotFind)) {
		return nil, errors.Parameter.AddMsg("告警名称重复").AddDetail(err)
	}
	db := ToAlarmInfoPo(in)
	err = l.AiDB.Insert(l.ctx, db)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &rule.WithID{Id: db.ID}, nil
}

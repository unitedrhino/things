package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/repo/mysql"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoReadLogic {
	return &AlarmInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmInfoReadLogic) AlarmInfoRead(in *rule.AlarmInfoReadReq) (*rule.AlarmInfo, error) {
	di, err := l.svcCtx.AlarmInfoRepo.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.NotFind
		}
		return nil, err
	}
	return ToAlarmInfo(di), nil
}

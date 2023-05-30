package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/repo/mysql"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoUpdateLogic {
	return &AlarmInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmInfoUpdateLogic) Update(old *mysql.RuleAlarmInfo, in *rule.AlarmInfo) *mysql.RuleAlarmInfo {
	old.Name = in.Name
	old.State = in.State
	old.Level = in.Level
	old.Desc = in.Desc
	return old
}

func (l *AlarmInfoUpdateLogic) AlarmInfoUpdate(in *rule.AlarmInfo) (*rule.Empty, error) {
	old, err := l.svcCtx.AlarmInfoRepo.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.AlarmInfoRepo.Update(l.ctx, l.Update(old, in))
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &rule.Empty{}, nil
}

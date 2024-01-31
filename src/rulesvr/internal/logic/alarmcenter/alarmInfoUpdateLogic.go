package alarmcenterlogic

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AlDB *relationDB.AlarmInfoRepo
}

func NewAlarmInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoUpdateLogic {
	return &AlarmInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AlDB:   relationDB.NewAlarmInfoRepo(ctx),
	}
}

func (l *AlarmInfoUpdateLogic) Update(old *relationDB.RuleAlarmInfo, in *rule.AlarmInfo) *relationDB.RuleAlarmInfo {
	old.Name = in.Name
	old.Status = in.Status
	old.Level = in.Level
	old.Desc = in.Desc
	return old
}

func (l *AlarmInfoUpdateLogic) AlarmInfoUpdate(in *rule.AlarmInfo) (*rule.Empty, error) {
	old, err := l.AlDB.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	err = l.AlDB.Update(l.ctx, l.Update(old, in))
	if err != nil {
		return nil, err
	}
	return &rule.Empty{}, nil
}

package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/domain/alarm"
	"github.com/i-Things/things/src/rulesvr/internal/repo/mysql"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmDealRecordCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmDealRecordCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmDealRecordCreateLogic {
	return &AlarmDealRecordCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 告警处理记录
func (l *AlarmDealRecordCreateLogic) AlarmDealRecordCreate(in *rule.AlarmDealRecordCreateReq) (*rule.WithID, error) {
	ai, err := l.svcCtx.AlarmRecordRepo.FindOne(l.ctx, in.AlarmRecordID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddMsg("该告警目前不处于告警中,不能处理告警")
		}
		return nil, errors.Database.AddDetail(err)
	}
	if ai.DealState != alarm.DealStateAlarming {
		return nil, errors.Parameter.AddMsg("该告警目前不处于告警中,不能处理告警")
	}
	_, err = l.svcCtx.AlarmDealRecordRepo.Insert(l.ctx, &mysql.RuleAlarmDealRecord{
		AlarmRecordID: in.AlarmRecordID,
		Result:        in.Result,
		Type:          in.Type,
		AlarmTime:     ai.LastAlarm,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	ai.DealState = alarm.DealStateAlarmed
	err = l.svcCtx.AlarmRecordRepo.Update(l.ctx, ai)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &rule.WithID{}, nil
}

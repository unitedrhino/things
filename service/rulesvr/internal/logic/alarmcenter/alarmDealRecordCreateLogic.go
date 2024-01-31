package alarmcenterlogic

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/rulesvr/internal/domain/alarm"
	"github.com/i-Things/things/service/rulesvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/rulesvr/internal/svc"
	"github.com/i-Things/things/service/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmDealRecordCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	ArDB  *relationDB.AlarmRecordRepo
	AdrDB *relationDB.AlarmDealRecordRepo
}

func NewAlarmDealRecordCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmDealRecordCreateLogic {
	return &AlarmDealRecordCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ArDB:   relationDB.NewAlarmRecordRepo(ctx),
		AdrDB:  relationDB.NewAlarmDealRecordRepo(ctx),
	}
}

// 告警处理记录
func (l *AlarmDealRecordCreateLogic) AlarmDealRecordCreate(in *rule.AlarmDealRecordCreateReq) (*rule.WithID, error) {
	ai, err := l.ArDB.FindOne(l.ctx, in.AlarmRecordID)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsg("该告警目前不处于告警中,不能处理告警")
		}
		return nil, errors.Database.AddDetail(err)
	}
	if ai.DealState != alarm.DealStateAlarming {
		return nil, errors.Parameter.AddMsg("该告警目前不处于告警中,不能处理告警")
	}
	err = l.AdrDB.Insert(l.ctx, &relationDB.RuleAlarmDealRecord{
		AlarmRecordID: in.AlarmRecordID,
		Result:        in.Result,
		Type:          in.Type,
		AlarmTime:     ai.LastAlarm,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	ai.DealState = alarm.DealStateAlarmed
	err = l.ArDB.Update(l.ctx, ai)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &rule.WithID{}, nil
}

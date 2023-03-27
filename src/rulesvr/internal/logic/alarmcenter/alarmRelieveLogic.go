package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/rulesvr/internal/domain/alarm"
	"github.com/i-Things/things/src/rulesvr/internal/repo/mysql"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmRelieveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmRelieveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmRelieveLogic {
	return &AlarmRelieveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmRelieveLogic) AlarmRelieve(in *rule.AlarmRelieveReq) (*rule.Response, error) {
	//调这个接口默认都是场景联动调用的
	alarms, err := l.svcCtx.AlarmInfoRepo.FindByFilter(l.ctx, alarm.InfoFilter{SceneID: in.SceneID}, nil)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	for _, a := range alarms {
		err := l.HandleOne(in, a)
		if err != nil {
			l.Errorf("%s.AlarmTrigger Alarm:%#v err:%v", utils.FuncName(), a, err)
		}
	}

	return &rule.Response{}, nil
}

func (l *AlarmRelieveLogic) HandleOne(in *rule.AlarmRelieveReq, alarmInfo *mysql.RuleAlarmInfo) error {
	var recordID int64
	ars, err := l.svcCtx.AlarmRecordRepo.FindByFilter(l.ctx, alarm.RecordFilter{
		AlarmID: alarmInfo.Id,
	}, nil)
	if err != nil {
		return errors.Database.AddDetail(err)
	}
	if len(ars) == 0 { //第一次触发
		return nil
	}
	for _, ar := range ars {
		if ar.DealState != alarm.DealStateAlarming {
			continue
		}
		ar.DealState = alarm.DealStateAlarmed
		err := l.svcCtx.AlarmRecordRepo.Update(l.ctx, ar)
		if err != nil {
			return errors.Database.AddDetail(err)
		}
		_, err = l.svcCtx.AlarmDealRecordRepo.Insert(l.ctx, &mysql.RuleAlarmDealRecord{
			AlarmRecordID: recordID,
			Result:        "场景触发解除告警",
			Type:          alarm.DealTypeSystem,
			AlarmTime:     ar.LastAlarm,
		})
		if err != nil {
			return errors.Database.AddDetail(err)
		}
	}

	return nil
}

package repoComplex

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/rulesvr/internal/domain/alarm"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/src/rulesvr/internal/repo/mysql"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type SceneAlarm struct {
	svcCtx *svc.ServiceContext
}

func NewSceneAlarm(svcCtx *svc.ServiceContext) *SceneAlarm {
	return &SceneAlarm{
		svcCtx: svcCtx,
	}
}

func (l *SceneAlarm) AlarmTrigger(ctx context.Context, in scene.AlarmTrigger) error {
	//调这个接口默认都是场景联动调用的
	alarms, err := l.svcCtx.AlarmInfoRepo.FindByFilter(ctx, alarm.InfoFilter{SceneID: in.SceneID}, nil)
	if err != nil {
		return errors.Database.AddDetail(err)
	}
	for _, a := range alarms {
		err := l.HandleTrigger(ctx, in, a)
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.AlarmTrigger Alarm:%#v err:%v", utils.FuncName(), a, err)
		}
	}
	return nil
}

func (l *SceneAlarm) AlarmRelieve(ctx context.Context, in scene.AlarmRelieve) error {
	//调这个接口默认都是场景联动调用的
	alarms, err := l.svcCtx.AlarmInfoRepo.FindByFilter(ctx, alarm.InfoFilter{SceneID: in.SceneID}, nil)
	if err != nil {
		return errors.Database.AddDetail(err)
	}
	for _, a := range alarms {
		err := l.HandleRelieve(ctx, in, a)
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.AlarmTrigger Alarm:%#v err:%v", utils.FuncName(), a, err)
		}
	}
	return nil
}
func (l *SceneAlarm) HandleRelieve(ctx context.Context, in scene.AlarmRelieve, alarmInfo *mysql.RuleAlarmInfo) error {
	var recordID int64
	ars, err := l.svcCtx.AlarmRecordRepo.FindByFilter(ctx, alarm.RecordFilter{
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
		err := l.svcCtx.AlarmRecordRepo.Update(ctx, ar)
		if err != nil {
			return errors.Database.AddDetail(err)
		}
		_, err = l.svcCtx.AlarmDealRecordRepo.Insert(ctx, &mysql.RuleAlarmDealRecord{
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

func (l *SceneAlarm) HandleTrigger(ctx context.Context, in scene.AlarmTrigger, alarmInfo *mysql.RuleAlarmInfo) error {
	var recordID int64
	ars, err := l.svcCtx.AlarmRecordRepo.FindByFilter(ctx, alarm.RecordFilter{
		AlarmID:     alarmInfo.Id,
		TriggerType: in.TriggerType,
		ProductID:   in.ProductID,
		DeviceName:  in.DeviceName,
	}, nil)
	if err != nil {
		return errors.Database.AddDetail(err)
	}
	if len(ars) == 0 { //第一次触发
		ret, err := l.svcCtx.AlarmRecordRepo.Insert(ctx, &mysql.RuleAlarmRecord{
			AlarmID:     alarmInfo.Id,
			TriggerType: in.TriggerType,
			ProductID:   in.ProductID,
			DeviceName:  in.DeviceName,
			Level:       alarmInfo.Level,
			SceneName:   in.SceneName,
			SceneID:     in.SceneID,
			DealState:   alarm.DealStateAlarming,
			LastAlarm:   time.Now(),
		})
		if err != nil {
			return errors.Database.AddDetail(err)
		}
		recordID, _ = ret.RowsAffected()
	} else {
		ar := ars[0]
		ar.LastAlarm = time.Now()
		ar.DealState = alarm.DealStateAlarming
		l.svcCtx.AlarmRecordRepo.Update(ctx, ar)
	}
	_, err = l.svcCtx.AlarmLogRepo.Insert(ctx, &mysql.RuleAlarmLog{
		AlarmRecordID: recordID,
		Serial:        in.Serial,
		SceneName:     in.SceneName,
		SceneID:       in.SceneID,
		Desc:          in.Desc,
	})
	if err != nil {
		return errors.Database.AddDetail(err)
	}
	return nil
}

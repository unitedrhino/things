package repoComplex

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/rulesvr/internal/domain/alarm"
	"github.com/i-Things/things/service/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/service/rulesvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/rulesvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type SceneAlarm struct {
	svcCtx *svc.ServiceContext
	AiDB   *relationDB.AlarmInfoRepo
	ArDB   *relationDB.AlarmRecordRepo
	AdrDB  *relationDB.AlarmDealRecordRepo
	AlDB   *relationDB.AlarmLogRepo
}

func NewSceneAlarm(svcCtx *svc.ServiceContext, ctx context.Context) *SceneAlarm {
	return &SceneAlarm{
		svcCtx: svcCtx,
		AiDB:   relationDB.NewAlarmInfoRepo(ctx),
		ArDB:   relationDB.NewAlarmRecordRepo(ctx),
		AdrDB:  relationDB.NewAlarmDealRecordRepo(ctx),
		AlDB:   relationDB.NewAlarmLogRepo(ctx),
	}
}

func (l *SceneAlarm) AlarmTrigger(ctx context.Context, in scene.TriggerSerial) error {
	//调这个接口默认都是场景联动调用的
	alarms, err := l.AiDB.FindByFilter(ctx, relationDB.AlarmInfoFilter{SceneID: in.SceneID}, nil)
	if err != nil {
		return errors.Database.AddDetail(err)
	}
	for _, a := range alarms {
		err := l.HandleTrigger(ctx, in, a)
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.TriggerSerial Alarm:%#v err:%v", utils.FuncName(), a, err)
		}
	}
	return nil
}

func (l *SceneAlarm) AlarmRelieve(ctx context.Context, in scene.AlarmRelieve) error {
	//调这个接口默认都是场景联动调用的
	alarms, err := l.AiDB.FindByFilter(ctx, relationDB.AlarmInfoFilter{SceneID: in.SceneID}, nil)
	if err != nil {
		return errors.Database.AddDetail(err)
	}
	for _, a := range alarms {
		err := l.HandleRelieve(ctx, in, a)
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.TriggerSerial Alarm:%#v err:%v", utils.FuncName(), a, err)
		}
	}
	return nil
}
func (l *SceneAlarm) HandleRelieve(ctx context.Context, in scene.AlarmRelieve, alarmInfo *relationDB.RuleAlarmInfo) error {
	var recordID int64
	ars, err := l.ArDB.FindByFilter(ctx, relationDB.AlarmRecordFilter{
		AlarmID: alarmInfo.ID,
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
		err := l.ArDB.Update(ctx, ar)
		if err != nil {
			return errors.Database.AddDetail(err)
		}
		err = l.AdrDB.Insert(ctx, &relationDB.RuleAlarmDealRecord{
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

func (l *SceneAlarm) HandleTrigger(ctx context.Context, in scene.TriggerSerial, alarmInfo *relationDB.RuleAlarmInfo) error {
	var recordID int64
	ars, err := l.ArDB.FindByFilter(ctx, relationDB.AlarmRecordFilter{
		AlarmID:     alarmInfo.ID,
		TriggerType: in.TriggerType,
		ProductID:   in.Device.ProductID,
		DeviceName:  in.Device.DeviceName,
	}, nil)
	if err != nil {
		return errors.Database.AddDetail(err)
	}
	if len(ars) == 0 { //第一次触发
		db := relationDB.RuleAlarmRecord{
			AlarmID:     alarmInfo.ID,
			TriggerType: in.TriggerType,
			ProductID:   in.Device.ProductID,
			DeviceName:  in.Device.DeviceName,
			Level:       alarmInfo.Level,
			SceneName:   in.SceneName,
			SceneID:     in.SceneID,
			DealState:   alarm.DealStateAlarming,
			LastAlarm:   time.Now(),
		}
		err := l.ArDB.Insert(ctx, &db)
		if err != nil {
			return errors.Database.AddDetail(err)
		}
		recordID = db.ID
	} else {
		ar := ars[0]
		ar.LastAlarm = time.Now()
		ar.DealState = alarm.DealStateAlarming
		l.ArDB.Update(ctx, ar)
	}
	err = l.AlDB.Insert(ctx, &relationDB.RuleAlarmLog{
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

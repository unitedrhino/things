package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/rulesvr/internal/domain/alarm"
	"github.com/i-Things/things/src/rulesvr/internal/repo/relationDB"
	"time"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmTriggerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AlarmInfoRepo
	ArDB *relationDB.AlarmRecordRepo
	AlDB *relationDB.AlarmLogRepo
}

func NewAlarmTriggerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmTriggerLogic {
	return &AlarmTriggerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAlarmInfoRepo(ctx),
		ArDB:   relationDB.NewAlarmRecordRepo(ctx),
		AlDB:   relationDB.NewAlarmLogRepo(ctx),
	}
}

// 告警触发
func (l *AlarmTriggerLogic) AlarmTrigger(in *rule.AlarmTriggerReq) (*rule.WithID, error) {
	//调这个接口默认都是场景联动调用的
	alarms, err := l.AiDB.FindByFilter(l.ctx, relationDB.AlarmInfoFilter{SceneID: in.SceneID}, nil)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	for _, a := range alarms {
		err := l.HandleOne(in, a)
		if err != nil {
			l.Errorf("%s.AlarmTrigger Alarm:%#v err:%v", utils.FuncName(), a, err)
		}
	}
	return &rule.WithID{}, nil
}
func (l *AlarmTriggerLogic) HandleOne(in *rule.AlarmTriggerReq, alarmInfo *relationDB.RuleAlarmInfo) error {
	var recordID int64
	ars, err := l.ArDB.FindByFilter(l.ctx, relationDB.AlarmRecordFilter{
		AlarmID:     alarmInfo.ID,
		TriggerType: in.TriggerType,
		ProductID:   in.ProductID,
		DeviceName:  in.DeviceName,
	}, nil)
	if err != nil {
		return errors.Database.AddDetail(err)
	}
	if len(ars) == 0 { //第一次触发
		db := relationDB.RuleAlarmRecord{
			AlarmID:     alarmInfo.ID,
			TriggerType: in.TriggerType,
			ProductID:   in.ProductID,
			DeviceName:  in.DeviceName,
			Level:       alarmInfo.Level,
			SceneName:   in.SceneName,
			SceneID:     in.SceneID,
			DealState:   alarm.DealStateAlarming,
			LastAlarm:   time.Now(),
		}
		err := l.ArDB.Insert(l.ctx, &db)
		if err != nil {
			return errors.Database.AddDetail(err)
		}
		recordID = db.ID
	} else {
		ar := ars[0]
		ar.LastAlarm = time.Now()
		ar.DealState = alarm.DealStateAlarming
		l.ArDB.Update(l.ctx, ar)
	}
	err = l.AlDB.Insert(l.ctx, &relationDB.RuleAlarmLog{
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

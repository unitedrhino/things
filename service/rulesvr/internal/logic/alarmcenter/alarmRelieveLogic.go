package alarmcenterlogic

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/rulesvr/internal/domain/alarm"
	"github.com/i-Things/things/service/rulesvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/rulesvr/internal/svc"
	"github.com/i-Things/things/service/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmRelieveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB  *relationDB.AlarmInfoRepo
	ArDB  *relationDB.AlarmRecordRepo
	AdrDB *relationDB.AlarmDealRecordRepo
}

func NewAlarmRelieveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmRelieveLogic {
	return &AlarmRelieveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAlarmInfoRepo(ctx),
		ArDB:   relationDB.NewAlarmRecordRepo(ctx),
		AdrDB:  relationDB.NewAlarmDealRecordRepo(ctx),
	}
}

func (l *AlarmRelieveLogic) AlarmRelieve(in *rule.AlarmRelieveReq) (*rule.WithID, error) {
	//调这个接口默认都是场景联动调用的
	alarms, err := l.AiDB.FindByFilter(l.ctx, relationDB.AlarmInfoFilter{SceneID: in.SceneID}, nil)
	if err != nil {
		return nil, err
	}
	for _, a := range alarms {
		err := l.HandleOne(in, a)
		if err != nil {
			l.Errorf("%s.AlarmTrigger Alarm:%#v err:%v", utils.FuncName(), a, err)
		}
	}

	return &rule.WithID{}, nil
}

func (l *AlarmRelieveLogic) HandleOne(in *rule.AlarmRelieveReq, alarmInfo *relationDB.RuleAlarmInfo) error {
	var recordID int64
	ars, err := l.ArDB.FindByFilter(l.ctx, relationDB.AlarmRecordFilter{
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
		err := l.ArDB.Update(l.ctx, ar)
		if err != nil {
			return errors.Database.AddDetail(err)
		}
		err = l.AdrDB.Insert(l.ctx, &relationDB.RuleAlarmDealRecord{
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

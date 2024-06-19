package rulelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"time"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmRecordCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmRecordCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmRecordCreateLogic {
	return &AlarmRecordCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmRecordCreateLogic) AlarmRecordCreate(in *ud.AlarmRecordCreateReq) (*ud.Empty, error) {
	pos, err := relationDB.NewAlarmSceneRepo(l.ctx).FindByFilter(l.ctx, relationDB.AlarmSceneFilter{SceneID: in.SceneID, WithAlarmInfo: true}, nil)
	if err != nil {
		return nil, err
	}
	if len(pos) == 0 {
		return nil, err
	}
	for _, alarm := range pos {
		if alarm.AlarmInfo.Status == def.False {
			continue
		}
		switch in.Mode {
		case scene.ActionAlarmModeRelieve:
			err = relationDB.NewAlarmRecordRepo(l.ctx).UpdateWithField(l.ctx,
				relationDB.AlarmRecordFilter{AlarmID: alarm.AlarmID, DealStatus: scene.AlarmDealStatusWaring},
				map[string]any{
					"deal_status": scene.AlarmDealStatusIgnore,
					"desc":        fmt.Sprintf("场景:%v(%v)解除告警", in.SceneName, in.SceneID),
				})
			if err != nil {
				return nil, err
			}
		case scene.ActionAlarmModeTrigger:
			po, err := relationDB.NewAlarmRecordRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.AlarmRecordFilter{
				AlarmID:      alarm.AlarmID,
				DealStatuses: []scene.AlarmDealStatus{scene.AlarmDealStatusWaring}, //还处在报警中,新增次数即可
			})
			err = errors.NotFind //先不开重复
			if err != nil {
				if errors.Cmp(err, errors.NotFind) { //直接创建并且进行通知
					po := relationDB.UdAlarmRecord{
						TenantCode:  alarm.TenantCode,
						ProjectID:   alarm.ProjectID,
						AlarmID:     alarm.AlarmID,
						AlarmName:   alarm.AlarmInfo.Name,
						TriggerType: in.TriggerType,
						ProductID:   in.ProductID,
						DeviceName:  in.DeviceName,
						Level:       alarm.AlarmInfo.Level,
						SceneName:   alarm.SceneInfo.Name,
						SceneID:     alarm.SceneID,
						DealStatus:  scene.AlarmDealStatusWaring,
						Desc:        "自动化触发告警",
						AlarmCount:  1,
						LastAlarm:   time.Now(),
					}
					err = relationDB.NewAlarmRecordRepo(l.ctx).Insert(l.ctx, &po)
					if err != nil {
						return nil, err
					}
					for _, notify := range alarm.AlarmInfo.Notifies {
						_, err := l.svcCtx.NotifyM.NotifyConfigSend(l.ctx, &sys.NotifyConfigSendReq{
							UserIDs:    alarm.AlarmInfo.UserIDs,
							Accounts:   alarm.AlarmInfo.Accounts,
							NotifyCode: def.NotifyCodeDeviceAlarm,
							TemplateID: notify.TemplateID,
							Type:       notify.Type,
							Params: map[string]string{
								"productID":  in.ProductID,
								"deviceName": in.DeviceName,
								"sceneName":  in.SceneName,
							},
						})
						if err != nil {
							l.Error(err)
							continue
						}
					}
					return &ud.Empty{}, err
				}
				return nil, err
			}
			po.LastAlarm = time.Now()
			po.AlarmCount++
			err = relationDB.NewAlarmRecordRepo(l.ctx).Update(l.ctx, po)
			return &ud.Empty{}, err
		}
	}
	return &ud.Empty{}, nil
}

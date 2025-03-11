package rulelogic

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/service/udsvr/internal/domain/alarm"
	"gitee.com/unitedrhino/things/service/udsvr/internal/domain/scene"
	"gitee.com/unitedrhino/things/service/udsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/topics"
	"time"

	"gitee.com/unitedrhino/things/service/udsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"

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
	if in.AlarmID == 0 && in.SceneID == 0 {
		return nil, errors.Parameter.AddMsg("请填写告警ID或触发的场景ID")
	}
	alarmHandle := func(alarmInfo *relationDB.UdAlarmInfo, sceneID int64) error {
		if alarmInfo.Status == def.False {
			return nil
		}
		err := l.svcCtx.FastEvent.Publish(l.ctx, fmt.Sprintf(topics.UdRuleAlarmNotify, in.Mode), utils.Copy[alarm.Notify](in))
		if err != nil {
			l.Error(err)
		}
		switch in.Mode {
		case scene.ActionAlarmModeRelieve:
			err = relationDB.NewAlarmRecordRepo(l.ctx).UpdateWithField(l.ctx,
				relationDB.AlarmRecordFilter{AlarmID: alarmInfo.ID, ProductID: in.ProductID, DeviceName: in.DeviceName, DealStatus: scene.AlarmDealStatusWaring},
				map[string]any{
					"deal_status": scene.AlarmDealStatusIgnore,
					"desc":        fmt.Sprintf("场景:%v(%v)解除告警", in.SceneName, in.SceneID),
				})
			if err != nil {
				return err
			}
		case scene.ActionAlarmModeTrigger:
			err := func() error {
				po, err := relationDB.NewAlarmRecordRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.AlarmRecordFilter{
					AlarmID: alarmInfo.ID, ProductID: in.ProductID, DeviceName: in.DeviceName,
					DealStatuses: []scene.AlarmDealStatus{scene.AlarmDealStatusWaring}, //还处在报警中,新增次数即可
				})
				var di *dm.DeviceInfo
				if in.DeviceName != "" {
					di, err = l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName})
					if err != nil {
						l.Error(err)
					}

				}
				defer func() {
					if po != nil {
						var params = map[string]string{
							"productID":   in.ProductID,
							"deviceName":  in.DeviceName,
							"sceneName":   in.SceneName,
							"deviceAlias": in.DeviceAlias,
						}
						if di != nil {
							p := utils.ToStringMapString(di)
							params = p
							params["sceneName"] = in.SceneName
						}
						for _, notify := range alarmInfo.Notifies {
							_, err := l.svcCtx.NotifyM.NotifyConfigSend(l.ctx, &sys.NotifyConfigSendReq{
								UserIDs:    alarmInfo.UserIDs,
								Accounts:   alarmInfo.Accounts,
								NotifyCode: def.NotifyCodeDeviceAlarm,
								TemplateID: notify.TemplateID,
								Type:       notify.Type,
								Params:     params,
							})
							if err != nil {
								l.Error(err)
								return
							}
						}
					}
				}()
				if err != nil {
					if !errors.Cmp(err, errors.NotFind) { //直接创建并且进行通知
						l.Errorf("NewAlarmRecordFind a:%v err:%v", alarmInfo, err)
						return err
					}
					po = &relationDB.UdAlarmRecord{
						TenantCode:     alarmInfo.TenantCode,
						ProjectID:      alarmInfo.ProjectID,
						AlarmID:        alarmInfo.ID,
						AlarmName:      alarmInfo.Name,
						TriggerType:    in.TriggerType,
						TriggerSubType: in.TriggerSubType,
						TriggerDetail:  in.TriggerDetail,
						ProductID:      in.ProductID,
						DeviceName:     in.DeviceName,
						DeviceAlias:    in.DeviceAlias,
						Level:          alarmInfo.Level,
						SceneName:      alarmInfo.Name,
						SceneID:        sceneID,
						DealStatus:     scene.AlarmDealStatusWaring,
						Desc:           fmt.Sprintf("自动化触发告警:%v", in.SceneName),
						AlarmCount:     1,
						LastAlarm:      time.Now(),
					}
					if di != nil {
						po.AreaID = dataType.AreaID(di.AreaID)
						po.AreaIDPath = dataType.AreaIDPath(di.AreaIDPath)
					}
					err = relationDB.NewAlarmRecordRepo(l.ctx).Insert(l.ctx, po)
					if err != nil {
						return err
					}
					return nil
				}
				po.LastAlarm = time.Now()
				po.AlarmCount++
				err = relationDB.NewAlarmRecordRepo(l.ctx).Update(l.ctx, po)
				if err != nil {
					l.Errorf("NewAlarmRecordUpdate a:%v err:%v", alarmInfo, err)
					return err
				}
				return nil
			}()
			if err != nil {
				return err
			}
		}
		return nil
	}

	if in.SceneID != 0 {
		pos, err := relationDB.NewAlarmSceneRepo(l.ctx).FindByFilter(l.ctx, relationDB.AlarmSceneFilter{
			SceneID: in.SceneID, WithAlarmInfo: true, WithSceneInfo: true}, nil)
		if err != nil {
			return nil, err
		}
		if len(pos) == 0 {
			return nil, err
		}
		for _, a := range pos {
			err := alarmHandle(a.AlarmInfo, a.SceneID)
			if err != nil {
				return nil, err
			}
		}
	} else if in.AlarmID != 0 {
		ai, err := relationDB.NewAlarmInfoRepo(l.ctx).FindOne(l.ctx, in.AlarmID)
		if err != nil {
			return nil, err
		}
		alarmHandle(ai, 0)
	}

	if in.DeviceName != "" && in.ProductID != "" {
		di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName})
		if err != nil {
			l.Error(err)
			return &ud.Empty{}, nil
		}
		total, err := relationDB.NewAlarmRecordRepo(l.ctx).CountByFilter(l.ctx, relationDB.AlarmRecordFilter{
			ProductID:    in.ProductID,
			DeviceName:   in.DeviceName,
			DealStatuses: []int64{scene.AlarmDealStatusWaring, scene.AlarmDealStatusInHand},
		})
		if err != nil {
			l.Error(err)
			return &ud.Empty{}, nil
		}
		if total > 0 && utils.SliceIn(di.Status, def.DeviceStatusOnline, def.DeviceStatusOffline) { //告警中的状态
			_, err := l.svcCtx.DeviceM.DeviceInfoUpdate(l.ctx, &dm.DeviceInfo{ProductID: in.ProductID, DeviceName: in.DeviceName, Status: def.DeviceStatusWarming})
			if err != nil {
				l.Error(err)
				return &ud.Empty{}, nil
			}
		}
		if total == 0 && di.Status == def.DeviceStatusWarming {
			_, err := l.svcCtx.DeviceM.DeviceInfoUpdate(l.ctx, &dm.DeviceInfo{ProductID: in.ProductID, DeviceName: in.DeviceName, Status: di.IsOnline + 1})
			if err != nil {
				l.Error(err)
				return &ud.Empty{}, nil
			}
		}
	}
	return &ud.Empty{}, nil
}

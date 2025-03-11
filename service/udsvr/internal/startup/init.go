package startup

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/tools"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/udsvr/internal/domain/alarm"
	"gitee.com/unitedrhino/things/service/udsvr/internal/domain/scene"
	"gitee.com/unitedrhino/things/service/udsvr/internal/event/sceneChangeEvent"
	"gitee.com/unitedrhino/things/service/udsvr/internal/event/timerEvent"
	"gitee.com/unitedrhino/things/service/udsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/udsvr/internal/svc"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/application"
	"gitee.com/unitedrhino/things/share/topics"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

func Init(svcCtx *svc.ServiceContext) {
	tools.InitStore(svcCtx.Config.CacheRedis)
	TimerInit(svcCtx)
	InitEventBus(svcCtx)
}

func InitEventBus(svcCtx *svc.ServiceContext) {
	err := svcCtx.FastEvent.QueueSubscribe(eventBus.CoreProjectInfoDelete, func(ctx context.Context, t time.Time, body []byte) error {
		pi := cast.ToInt64(string(body))
		logx.WithContext(ctx).Infof("CoreProjectInfoDelete value:%v ")
		return sceneChangeEvent.NewHandle(ctx, svcCtx).SceneProjectDelete(pi)
	})
	logx.Must(err)

	err = svcCtx.FastEvent.QueueSubscribe(topics.DmDeviceInfoDelete, func(ctx context.Context, t time.Time, body []byte) error {
		var di devices.Core
		err := json.Unmarshal(body, &di)
		logx.WithContext(ctx).Infof("DmDeviceInfoDelete value:%v err:%v", string(body), err)
		return sceneChangeEvent.NewHandle(ctx, svcCtx).SceneDeviceDelete(di)
	})
	logx.Must(err)
	err = svcCtx.FastEvent.QueueSubscribe(topics.DmDeviceInfoUnbind, func(ctx context.Context, t time.Time, body []byte) error {
		var di devices.Core
		err := json.Unmarshal(body, &di)
		logx.WithContext(ctx).Infof("DmDeviceInfoUnbind value:%v err:%v", string(body), err)
		return sceneChangeEvent.NewHandle(ctx, svcCtx).SceneDeviceDelete(di)
	})
	logx.Must(err)
	err = svcCtx.FastEvent.QueueSubscribe(topics.DmProductInfoUpdate, func(ctx context.Context, t time.Time, body []byte) error {
		var di devices.Core
		err := json.Unmarshal(body, &di)
		logx.WithContext(ctx).Infof("SceneDeviceUpdate value:%v err:%v", string(body), err)
		return sceneChangeEvent.NewHandle(ctx, svcCtx).SceneDeviceUpdate(di)
	})
	logx.Must(err)
	//err := svcCtx.FastEvent.QueueSubscribe(eventBus.UdRuleTimer, func(ctx context.Context, t time.Time, body []byte) error {
	//	if t.Add(2 * time.Second).Before(time.Now()) { //2秒过期时间
	//		return nil
	//	}
	//	th := timerEvent.NewSceneHandle(ctx, svcCtx)
	//	return th.DeviceTimer()
	//})
	//logx.Must(err)

	{
		err = svcCtx.FastEvent.QueueSubscribe(topics.UdRuleTimer, func(ctx context.Context, t time.Time, body []byte) error {
			if t.Before(time.Now().Add(-time.Second * 2)) { //2秒之前的跳过
				return nil
			}
			th := timerEvent.NewSceneHandle(ctxs.WithRoot(ctx), svcCtx)
			var wait sync.WaitGroup
			wait.Add(2)
			utils.Go(ctx, func() {
				err := th.SceneTiming()
				if err != nil {
					logx.WithContext(ctx).Error(err)
				}
			})
			utils.Go(ctx, func() {
				err := th.DeviceTriggerCheck()
				if err != nil {
					logx.WithContext(ctx).Error(err)
				}
			})
			return nil
		})
		logx.Must(err)
	}
	{
		err = svcCtx.FastEvent.QueueSubscribe(topics.UdRuleTimerTenMinutes, func(ctx context.Context, t time.Time, body []byte) error {
			if t.Before(time.Now().Add(-time.Second * 2)) { //2秒之前的跳过
				return nil
			}
			th := timerEvent.NewSceneHandle(ctxs.WithRoot(ctx), svcCtx)
			utils.Go(ctx, func() {
				err := th.SceneTimingTenMinutes()
				if err != nil {
					logx.WithContext(ctx).Error(err)
				}
			})
			return nil
		})
		logx.Must(err)
	}
	err = svcCtx.FastEvent.QueueSubscribe(topics.ApplicationDeviceReportThingPropertyAllDevice, func(ctx context.Context, t time.Time, body []byte) error {
		if t.Before(time.Now().Add(-time.Second * 2)) { //2秒之前的跳过
			return nil
		}
		th := timerEvent.NewSceneHandle(ctxs.WithRoot(ctx), svcCtx)
		var stu application.PropertyReport
		err := utils.Unmarshal(body, &stu)
		if err != nil {
			logx.WithContext(ctx).Errorf("Subscribe.QueueSubscribe.Unmarshal body:%v err:%v", string(body), err)
			return err
		}
		return th.SceneThingPropertyReport(stu)
	})
	logx.Must(err)
	err = svcCtx.FastEvent.QueueSubscribe(topics.ApplicationDeviceReportThingEventAllDevice, func(ctx context.Context, t time.Time, body []byte) error {
		if t.Before(time.Now().Add(-time.Second * 2)) { //2秒之前的跳过
			return nil
		}
		th := timerEvent.NewSceneHandle(ctxs.WithRoot(ctx), svcCtx)
		var stu application.EventReport
		err := utils.Unmarshal(body, &stu)
		if err != nil {
			logx.WithContext(ctx).Errorf("Subscribe.QueueSubscribe.Unmarshal body:%v err:%v", string(body), err)
			return err
		}
		return th.SceneThingEventReport(stu)
	})
	logx.Must(err)
	err = svcCtx.FastEvent.QueueSubscribe(topics.ApplicationDeviceStatusAllDevice, func(ctx context.Context, t time.Time, body []byte) error {
		if t.Before(time.Now().Add(-time.Second * 2)) { //2秒之前的跳过
			return nil
		}
		th := timerEvent.NewSceneHandle(ctxs.WithRoot(ctx), svcCtx)
		var stu application.ConnectMsg
		err := utils.Unmarshal(body, &stu)
		if err != nil {
			logx.WithContext(ctx).Errorf("Subscribe.QueueSubscribe.Unmarshal body:%v err:%v", string(body), err)
			return err
		}
		return th.SceneDeviceOnline(stu)
	})
	logx.Must(err)

	{
		err = svcCtx.FastEvent.QueueSubscribe(eventBus.CoreOpsWorkOrderFinish, func(ctx context.Context, t time.Time, body []byte) error {
			pi := cast.ToInt64(string(body))
			logx.WithContext(ctx).Infof("CoreOpsWorkOrderFinish value:%v err:%v", string(body), err)
			if pi == 0 {
				return nil
			}
			ctx = ctxs.WithRoot(ctx)
			ar, err := relationDB.NewAlarmRecordRepo(ctx).FindOneByFilter(ctx, relationDB.AlarmRecordFilter{WorkOrderID: pi})
			if err != nil {
				logx.WithContext(ctx).Error(err)
				return nil
			}
			ar.DealStatus = scene.AlarmDealStatusProcessed
			err = relationDB.NewAlarmRecordRepo(ctx).Update(ctx, ar)
			if err != nil {
				logx.WithContext(ctx).Error(err)
				return nil
			}
			n := utils.Copy[alarm.Notify](ar)
			n.Mode = scene.ActionAlarmModeRelieve
			err = svcCtx.FastEvent.Publish(ctx, fmt.Sprintf(topics.UdRuleAlarmNotify, scene.ActionAlarmModeRelieve), n)
			if err != nil {
				logx.WithContext(ctx).Error(err)
			}
			//if ar.DeviceName != "" && ar.ProductID != "" {
			//	di, err := svcCtx.DeviceCache.GetData(ctx, devices.Core{ProductID: ar.ProductID, DeviceName: ar.DeviceName})
			//	if err != nil {
			//		logx.WithContext(ctx).Error(err)
			//		return nil
			//	}
			//	total, err := relationDB.NewAlarmRecordRepo(ctx).CountByFilter(ctx, relationDB.AlarmRecordFilter{
			//		ProductID:    ar.ProductID,
			//		DeviceName:   ar.DeviceName,
			//		DealStatuses: []int64{scene.AlarmDealStatusWaring, scene.AlarmDealStatusInHand},
			//	})
			//	if err != nil {
			//		logx.WithContext(ctx).Error(err)
			//		return nil
			//	}
			//	if total == 0 && di.IsAbnormal == def.DeviceStatusWarming {
			//		_, err := svcCtx.DeviceM.DeviceInfoUpdate(ctx, &dm.DeviceInfo{ProductID: ar.ProductID, DeviceName: ar.DeviceName, IsAbnormal: di.IsOnline + 1})
			//		if err != nil {
			//			logx.WithContext(ctx).Error(err)
			//			return nil
			//		}
			//	}
			//}
			return nil
		})
		logx.Must(err)
	}
	err = svcCtx.FastEvent.Start()
	logx.Must(err)
}

func TimerInit(svcCtx *svc.ServiceContext) {
	ctx := context.Background()
	{
		_, err := svcCtx.TimedM.TaskInfoCreate(ctx, &timedmanage.TaskInfo{
			GroupCode: def.TimedUnitedRhinoQueueGroupCode,                             //组编码
			Type:      1,                                                              //任务类型 1 定时任务 2 延时任务
			Name:      "联犀规则引擎定时任务",                                                   // 任务名称
			Code:      "iThingsRuleTimer",                                             //任务编码
			Params:    fmt.Sprintf(`{"topic":"%s","payload":""}`, topics.UdRuleTimer), // 任务参数,延时任务如果没有传任务参数会拿数据库的参数来执行
			CronExpr:  "@every 1s",                                                    // cron执行表达式
			Status:    def.StatusWaitRun,                                              // 状态
			Priority:  3,                                                              //优先级: 10:critical 最高优先级  3: default 普通优先级 1:low 低优先级
		})
		if err != nil && !errors.Cmp(errors.Fmt(err), errors.Duplicate) {
			logx.Must(err)
		}
	}
	{
		_, err := svcCtx.TimedM.TaskInfoCreate(ctx, &timedmanage.TaskInfo{
			GroupCode: def.TimedUnitedRhinoQueueGroupCode,                                       //组编码
			Type:      1,                                                                        //任务类型 1 定时任务 2 延时任务
			Name:      "联犀规则引擎定时任务10分钟",                                                         // 任务名称
			Code:      "UdRuleTimerTenMinutes",                                                  //任务编码
			Params:    fmt.Sprintf(`{"topic":"%s","payload":""}`, topics.UdRuleTimerTenMinutes), // 任务参数,延时任务如果没有传任务参数会拿数据库的参数来执行
			CronExpr:  "@every 10m",                                                             // cron执行表达式
			Status:    def.StatusWaitRun,                                                        // 状态
			Priority:  3,                                                                        //优先级: 10:critical 最高优先级  3: default 普通优先级 1:low 低优先级
		})
		if err != nil && !errors.Cmp(errors.Fmt(err), errors.Duplicate) {
			logx.Must(err)
		}
	}

}

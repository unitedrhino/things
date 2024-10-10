package startup

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/events"
	"gitee.com/unitedrhino/things/service/dgsvr/internal/event/deviceSub"
	"gitee.com/unitedrhino/things/service/dgsvr/internal/event/innerSub"
	"gitee.com/unitedrhino/things/service/dgsvr/internal/event/onlineCheck"
	"gitee.com/unitedrhino/things/service/dgsvr/internal/repo/event/publish/pubDev"
	"gitee.com/unitedrhino/things/service/dgsvr/internal/repo/event/publish/pubInner"
	"gitee.com/unitedrhino/things/service/dgsvr/internal/repo/event/subscribe/subDev"
	"gitee.com/unitedrhino/things/service/dgsvr/internal/repo/event/subscribe/subInner"
	"gitee.com/unitedrhino/things/service/dgsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func Init(svcCtx *svc.ServiceContext) {
	//some init for serviceContext
}

// mqtt and nats client
func PostInit(svcCtx *svc.ServiceContext) {
	dl, err := pubDev.NewPubDev(svcCtx.Config.DevLink)
	logx.Must(err)

	il, err := pubInner.NewPubInner(svcCtx.Config.Event, def.ProtocolCodeIThings, svcCtx.NodeID)
	logx.Must(err)

	svcCtx.PubDev = dl
	svcCtx.PubInner = il
	mc, err := clients.NewMqttClient(svcCtx.Config.DevLink.Mqtt)
	logx.Must(err)
	svcCtx.MqttClient = mc
	sd, err := subDev.NewSubDev(svcCtx.Config.DevLink)
	logx.Must(err)
	err = sd.SubDevMsg(func(ctx context.Context) subDev.DevSubHandle {
		return deviceSub.NewDeviceSubServer(svcCtx, ctx)
	})
	logx.Must(err)
	si, err := subInner.NewSubInner(svcCtx.Config.Event, svcCtx.NodeID)
	logx.Must(err)
	err = si.SubToDevMsg(func(ctx context.Context) subInner.InnerSubHandle {
		return innerSub.NewInnerSubServer(svcCtx, ctx)
	})
	logx.Must(err)
	InitEventBus(svcCtx)
	TimerInit(svcCtx)
}
func InitEventBus(svcCtx *svc.ServiceContext) {
	err := svcCtx.FastEvent.Subscribe(eventBus.DmProductCustomUpdate, func(ctx context.Context, t time.Time, body []byte) error {
		info := events.DeviceUpdateInfo{}
		err := json.Unmarshal(body, &info)
		if err != nil {
			return err
		}
		return svcCtx.Script.ClearCache(ctx, info.ProductID)
	})
	logx.Must(err)
	err = svcCtx.FastEvent.QueueSubscribe(eventBus.DgOnlineTimer, func(ctx context.Context, t time.Time, body []byte) error {
		return onlineCheck.NewOnlineCheckEvent(svcCtx, ctx).Check()
	})
	logx.Must(err)
	err = svcCtx.FastEvent.Start()
	logx.Must(err)
}

func TimerInit(svcCtx *svc.ServiceContext) {
	ctx := context.Background()
	_, err := svcCtx.TimedM.TaskInfoCreate(ctx, &timedmanage.TaskInfo{
		GroupCode: def.TimedIThingsQueueGroupCode,                                     //组编码
		Type:      1,                                                                  //任务类型 1 定时任务 2 延时任务
		Name:      "iThings协议网关定时处理",                                                  // 任务名称
		Code:      "iThingsDgOnlineTimer",                                             //任务编码
		Params:    fmt.Sprintf(`{"topic":"%s","payload":""}`, eventBus.DgOnlineTimer), // 任务参数,延时任务如果没有传任务参数会拿数据库的参数来执行
		CronExpr:  "@every 5m",                                                        // cron执行表达式
		Status:    def.StatusWaitRun,                                                  // 状态
		Priority:  3,                                                                  //优先级: 10:critical 最高优先级  3: default 普通优先级 1:low 低优先级
	})
	if err != nil && !errors.Cmp(errors.Fmt(err), errors.Duplicate) {
		logx.Must(err)
	}
}

package serverEvent

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/core/service/syssvr/sysExport"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/application"
	"gitee.com/i-Things/share/domain/deviceAuth"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/domain/deviceMsg/msgThing"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/events/topics"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/dmExport"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceLog"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceStatus"
	"github.com/i-Things/things/service/dmsvr/internal/domain/serverDo"
	deviceinteractlogic "github.com/i-Things/things/service/dmsvr/internal/logic/deviceinteract"
	"github.com/i-Things/things/service/dmsvr/internal/repo/cache"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type ServerHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewServerHandle(ctx context.Context, svcCtx *svc.ServiceContext) *ServerHandle {
	return &ServerHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

// 定时处理设备在线状态改变,需要过滤设备在线抖动的问题
func (l *ServerHandle) OnlineStatusHandle() error {
	ok, err := l.svcCtx.DeviceStatus.Lock(l.ctx)
	if err != nil {
		return err
	}
	if !ok { //没抢到锁的话需要
		return nil
	}
	defer func() {
		l.svcCtx.DeviceStatus.UnLock(l.ctx)
	}()
	devs, err := l.svcCtx.DeviceStatus.GetDevices(l.ctx)
	if err != nil {
		return err
	}
	if len(devs) == 0 {
		return nil
	}
	var now = time.Now()
	var clientIDSet = map[string][]*deviceStatus.ConnectMsg{}
	for _, v := range devs {
		if _, ok := clientIDSet[v.ClientID]; ok {
			clientIDSet[v.ClientID] = append(clientIDSet[v.ClientID], v)
		} else {
			clientIDSet[v.ClientID] = []*deviceStatus.ConnectMsg{v}
		}
	}
	var removeList []*deviceStatus.ConnectMsg
	var insertList []*deviceStatus.ConnectMsg
	var t = now.Add(-time.Second * 5)
	var older = now.Add(-time.Second * 10) //10秒以外的直接入库,可能是服务挂了引起
	for _, v := range clientIDSet {
		if len(v) == 1 && v[0].Timestamp.Before(t) { //如果5秒过去了,还没有反复的登入登出,则入库
			removeList = append(removeList, v...)
			insertList = append(insertList, v...)
			continue
		}
		if len(v) > 1 {
			var (
				connected    []*deviceStatus.ConnectMsg
				disconnected []*deviceStatus.ConnectMsg
			)
			for _, one := range v {
				if one.Timestamp.Before(older) { //历史的直接入库即可
					removeList = append(removeList, one)
					insertList = append(insertList, one)
					continue
				}
				if one.Action == devices.ActionConnected {
					connected = append(connected, one)
				} else {
					disconnected = append(disconnected, one)
				}
			}
			//如果存在同时上线和下线的情况,则需要过滤了
			var hasShake bool
			for len(connected) > 0 && len(disconnected) > 0 {
				hasShake = true
				//删除最后一个
				removeList = append(removeList, connected[len(connected)-1], disconnected[len(disconnected)-1])
				//更新缓存
				connected = connected[:len(connected)-1]
				disconnected = disconnected[:len(disconnected)-1]
			}
			if hasShake {
				l.Errorf("设备上下线出现抖动症状:设备信息: connected:%v disconnected:%v", utils.Fmt(connected), utils.Fmt(disconnected))
			}
			var conns = connected
			conns = append(conns, disconnected...)
			if len(conns) > 0 {
				for _, one := range conns {
					if one.Timestamp.Before(t) { //如果5秒过去了,且已经过滤过重复的登录登出状态,则直接入库即可,没有到时间的后续再继续检查
						removeList = append(removeList, v...)
						insertList = append(insertList, v...)
					}
				}
			}
		}
	}

	if len(insertList) == 0 && len(removeList) == 0 {
		return nil
	}
	l.Infof("insertList:%v removeList:%v", utils.Fmt(insertList), utils.Fmt(removeList))
	//入库异步处理
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		var ( //这里是最后更新数据库状态的设备列表
			OffLineDevices []*devices.Core
		)
		var log = logx.WithContext(ctx)
		for _, msg := range insertList {
			status := int64(def.ConnectedStatus)
			if msg.Action == devices.ActionDisconnected {
				status = def.DisConnectedStatus
			}
			ld, err := deviceAuth.GetClientIDInfo(msg.ClientID)
			if err != nil {
				log.Error(err)
				continue
			}
			err = l.svcCtx.StatusRepo.Insert(ctx, &deviceLog.Status{
				ProductID:  ld.ProductID,
				Status:     status,
				Timestamp:  msg.Timestamp, // 操作时间
				DeviceName: ld.DeviceName,
			})
			if err != nil {
				log.Errorf("%s.HubLogRepo.insert productID:%v deviceName:%v err:%v",
					utils.FuncName(), ld.ProductID, ld.DeviceName, err)
			}
			appMsg := application.ConnectMsg{
				Device: devices.Core{
					ProductID:  ld.ProductID,
					DeviceName: ld.DeviceName,
				},
				Status:    status,
				Timestamp: msg.Timestamp.UnixMilli(),
			}
			utils.Go(ctx, func() {
				err = l.svcCtx.WebHook.Publish(l.svcCtx.WithDeviceTenant(l.ctx, appMsg.Device), sysExport.CodeDmDeviceConn, appMsg)
				if err != nil {
					l.Error(err)
				}
				err = l.svcCtx.UserSubscribe.Publish(l.ctx, def.UserSubscribeDeviceConn, appMsg, map[string]any{
					"productID":  ld.ProductID,
					"deviceName": ld.DeviceName,
				}, map[string]any{
					"productID": ld.ProductID,
				}, map[string]any{})
				if err != nil {
					l.Error(err)
				}
			})

			if status == def.ConnectedStatus {
				di, err := relationDB.NewDeviceInfoRepo(ctx).FindOneByFilter(ctx, relationDB.DeviceFilter{Cores: []*devices.Core{{ProductID: ld.ProductID, DeviceName: ld.DeviceName}}})
				if err != nil {
					log.Error(err)
					continue
				}
				if di.IsOnline == def.True {
					log.Infof("already online:%#v", msg)
					continue
				}
				var updates = map[string]any{"is_online": def.True, "last_login": msg.Timestamp, "status": def.DeviceStatusOnline}
				if di.FirstLogin.Valid == false {
					updates["first_login"] = msg.Timestamp
				}
				err = relationDB.NewDeviceInfoRepo(ctx).UpdateWithField(ctx,
					relationDB.DeviceFilter{Cores: []*devices.Core{{ProductID: ld.ProductID, DeviceName: ld.DeviceName}}}, updates)
				if err != nil {
					log.Error(err)
				}
				err = l.svcCtx.DeviceCache.SetData(ctx, dmExport.GenDeviceInfoKey(ld.ProductID, ld.DeviceName), nil)
				if err != nil {
					log.Error(err)
				}

				err = l.svcCtx.PubApp.DeviceStatusConnected(ctx, appMsg)
				if err != nil {
					l.Errorf("%s.pubApp productID:%v deviceName:%v err:%v",
						utils.FuncName(), ld.ProductID, ld.DeviceName, err)
				}
			} else {
				OffLineDevices = append(OffLineDevices, &devices.Core{
					ProductID:  ld.ProductID,
					DeviceName: ld.DeviceName,
				})
				err = l.svcCtx.PubApp.DeviceStatusDisConnected(ctx, appMsg)
				if err != nil {
					l.Errorf("%s.pubApp productID:%v deviceName:%v err:%v",
						utils.FuncName(), ld.ProductID, ld.DeviceName, err)
				}
			}

		}
		diDB := relationDB.NewDeviceInfoRepo(ctx)
		if len(OffLineDevices) > 0 {
			err = diDB.UpdateOfflineStatus(ctx, relationDB.DeviceFilter{Cores: OffLineDevices})
			if err != nil {
				log.Error(err)
			}
			for _, v := range OffLineDevices { //清除缓存
				err := l.svcCtx.DeviceCache.SetData(ctx, dmExport.GenDeviceInfoKey(v.ProductID, v.DeviceName), nil)
				if err != nil {
					log.Error(err)
				}
			}
		}
	})
	if len(removeList) > 0 {
		err = l.svcCtx.DeviceStatus.DelDevices(l.ctx, removeList...)
	}
	return err
}

func (l *ServerHandle) ActionCheck(in *deviceMsg.PublishMsg) error {
	l.Infof("ActionCheck req:%v", in)
	var req msgThing.Req
	var option serverDo.SendOption
	json.Unmarshal(in.Payload, &req)
	if in.Explain != "" {
		json.Unmarshal([]byte(in.Explain), &option)
	}
	resp, err := cache.GetDeviceMsg[msgThing.Resp](l.ctx, l.svcCtx.Cache, deviceMsg.RespMsg, devices.Thing, msgThing.TypeAction,
		devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}, req.MsgToken)
	if err != nil {
		l.Errorf("GetDeviceMsg err:%v", err)
		return err
	}
	if resp != nil { //设备已经回复,不需要管
		return nil
	}
	core := devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}
	sendTime := time.UnixMilli(in.Timestamp)
	now := time.Now()
	sendMsg := func(err error) {
		devErr := errors.Fmt(err)
		ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
			l.Infof("DeviceThingActionReport.Action device:%v,req:%v", core, req)
			//应用事件通知-设备物模型事件上报通知 ↓↓↓
			err := l.svcCtx.PubApp.DeviceThingActionReport(ctx, application.ActionReport{
				Device: core, Timestamp: now.UnixMilli(), ReqType: deviceMsg.ReqMsg, MsgToken: req.MsgToken,
				ActionID: req.ActionID, Dir: schema.ActionDirUp, Code: devErr.GetCode(), Status: devErr.GetMsg(),
			})
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.DeviceThingActionReport.Action  req:%v,err:%v", utils.FuncName(), utils.Fmt(req), err)
			}
		})
	}
	if now.After(sendTime.Add(time.Duration(option.TimeoutToFail) * time.Second)) {
		//过期了 发送失败
		sendMsg(errors.TimeOut)
		return nil
	}
	_, err = deviceinteractlogic.CheckIsOnline(l.ctx, l.svcCtx, core)
	if err != nil {
		sendMsg(err)
		return nil
	}
	err = l.svcCtx.PubDev.PublishToDev(l.ctx, in)
	if err != nil {
		return err
	}
	payload, _ := json.Marshal(in)
	_, err = l.svcCtx.TimedM.TaskCancel(l.ctx, &timedmanage.TaskWithTaskID{TaskID: req.MsgToken})
	if err != nil { //重复创建一个taskID会报错,需要先删除原来的任务
		l.Errorf("TaskSend err:%v", err)
	}
	_, err = l.svcCtx.TimedM.TaskSend(l.ctx, &timedjob.TaskSendReq{
		GroupCode: def.TimedIThingsQueueGroupCode,
		Code:      "disvr-action-check-delay",
		Option: &timedjob.TaskSendOption{
			ProcessIn: option.RetryInterval,
			Deadline:  sendTime.Add(time.Duration(option.TimeoutToFail) * time.Second).Unix(),
			TaskID:    req.MsgToken,
		},
		ParamQueue: &timedjob.TaskParamQueue{
			Topic:   topics.DmActionCheckDelay,
			Payload: string(payload),
		},
	})
	if err != nil {
		l.Errorf("TaskSend err:%v", err)
	}
	return nil
}

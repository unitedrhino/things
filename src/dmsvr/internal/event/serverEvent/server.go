package serverEvent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/application"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgThing"
	"github.com/i-Things/things/src/dmsvr/internal/domain/serverDo"
	deviceinteractlogic "github.com/i-Things/things/src/dmsvr/internal/logic/deviceinteract"
	"github.com/i-Things/things/src/dmsvr/internal/repo/cache"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"
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

func (l *ServerHandle) ActionCheck(in *deviceMsg.PublishMsg) error {
	l.Infof("ActionCheck req:%v", in)
	jsonStr, _ := json.Marshal(in)
	fmt.Println("[---deviceMsg---] ", string(jsonStr))
	//fmt.Println("[***]", "func (l *ServerHandle) ActionCheck(in *deviceMsg.PublishMsg)")
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
		utils.GoNewCtx(l.ctx, func(ctx context.Context) {
			l.Infof("DeviceThingActionReport.Action device:%v,req:%v", core, req)
			//应用事件通知-设备物模型事件上报通知 ↓↓↓
			err := l.svcCtx.PubApp.DeviceThingActionReport(ctx, application.ActionReport{
				Device: core, Timestamp: now.UnixMilli(), ReqType: deviceMsg.ReqMsg, MsgToken: req.MsgToken,
				ActionID: req.ActionID, Dir: schema.ActionDirUp, Code: devErr.GetCode(), Status: devErr.Msg,
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
	err = deviceinteractlogic.CheckIsOnline(l.ctx, l.svcCtx, core)
	if err != nil {
		sendMsg(err)
		return nil
	}
	err = l.svcCtx.PubDev.PublishToDev(l.ctx, in)
	if err != nil {
		return err
	}
	payload, _ := json.Marshal(in)
	_, err = l.svcCtx.TimedM.TaskSend(l.ctx, &timedjob.TaskSendReq{
		GroupCode: def.TimedIThingsQueueGroupCode,
		Code:      "disvr-action-check-delay",
		Option: &timedjob.TaskSendOption{
			ProcessIn: option.RetryInterval,
			Deadline:  sendTime.Add(time.Duration(option.TimeoutToFail) * time.Second).Unix(),
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

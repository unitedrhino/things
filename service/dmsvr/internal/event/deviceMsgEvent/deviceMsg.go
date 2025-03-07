package deviceMsgEvent

//设备的发布,连接及断连处理
import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceStatus"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/application"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type DeviceMsgHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewDeviceMsgHandle(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceMsgHandle {
	return &DeviceMsgHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (l *DeviceMsgHandle) Gateway(msg *deviceMsg.PublishMsg) error {
	l.Infof("%s req=%v", utils.FuncName(), msg)
	msg = l.svcCtx.ScriptTrans.UpBeforeTrans(l.ctx, msg)
	l.FixDisconnect(msg)
	resp, err := NewGatewayLogic(l.ctx, l.svcCtx).Handle(msg)
	l.deviceResp(resp)
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		err := l.svcCtx.ScriptTrans.UpAfterTrans(l.ctx, msg, resp)
		if err != nil {
			l.Error(msg, resp, err)
		}
	})
	l.Infof("%s req:%v resp:%v err:%v", utils.FuncName(), msg, resp, err)
	return err
}

func (l *DeviceMsgHandle) Thing(msg *deviceMsg.PublishMsg) error {
	l.Infof("%s req=%v", utils.FuncName(), msg)
	l.FixDisconnect(msg)
	startTime := time.Now()
	msg = l.svcCtx.ScriptTrans.UpBeforeTrans(l.ctx, msg)
	resp, err := NewThingLogic(l.ctx, l.svcCtx).Handle(msg)
	l.deviceResp(resp)
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		err := l.svcCtx.ScriptTrans.UpAfterTrans(l.ctx, msg, resp)
		if err != nil {
			l.Error(msg, resp, err)
		}
	})
	l.WithDuration(time.Now().Sub(startTime)).Infof("%s startTime:%v req:%v resp:%v err:%v",
		utils.FuncName(), startTime, msg, resp, err)
	return err
}

func (l *DeviceMsgHandle) Ota(msg *deviceMsg.PublishMsg) error {
	msg = l.svcCtx.ScriptTrans.UpBeforeTrans(l.ctx, msg)
	l.Infof("%s req=%v", utils.FuncName(), msg)
	l.FixDisconnect(msg)
	resp, err := NewOtaLogic(l.ctx, l.svcCtx).Handle(msg)
	l.deviceResp(resp)
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		err := l.svcCtx.ScriptTrans.UpAfterTrans(l.ctx, msg, resp)
		if err != nil {
			l.Error(msg, resp, err)
		}
	})
	l.Infof("%s req:%v resp:%v err:%v", utils.FuncName(), msg, resp, err)
	return err
}

func (l *DeviceMsgHandle) Shadow(msg *deviceMsg.PublishMsg) error {
	msg = l.svcCtx.ScriptTrans.UpBeforeTrans(l.ctx, msg)
	l.Infof("%s req=%v", utils.FuncName(), msg)
	l.FixDisconnect(msg)
	resp, err := NewShadowLogic(l.ctx, l.svcCtx).Handle(msg)
	l.deviceResp(resp)
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		err := l.svcCtx.ScriptTrans.UpAfterTrans(l.ctx, msg, resp)
		if err != nil {
			l.Error(msg, resp, err)
		}
	})
	l.Infof("%s req:%v resp:%v err:%v", utils.FuncName(), msg, resp, err)
	return err
}

func (l *DeviceMsgHandle) Config(msg *deviceMsg.PublishMsg) error {
	msg = l.svcCtx.ScriptTrans.UpBeforeTrans(l.ctx, msg)
	l.Infof("%s req=%v", utils.FuncName(), msg)
	l.FixDisconnect(msg)
	respMsg, err := NewConfigLogic(l.ctx, l.svcCtx).Handle(msg)
	l.deviceResp(respMsg)
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		err := l.svcCtx.ScriptTrans.UpAfterTrans(l.ctx, msg, respMsg)
		if err != nil {
			l.Error(msg, respMsg, err)
		}
	})
	l.Infof("%s req:%v resp:%v err:%v", utils.FuncName(), msg, respMsg, err)
	return err
}

func (l *DeviceMsgHandle) Ext(msg *deviceMsg.PublishMsg) error {
	msg = l.svcCtx.ScriptTrans.UpBeforeTrans(l.ctx, msg)
	l.Infof("%s req=%v", utils.FuncName(), msg)
	l.FixDisconnect(msg)
	resp, err := NewExtLogic(l.ctx, l.svcCtx).Handle(msg)
	l.deviceResp(resp)
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		err := l.svcCtx.ScriptTrans.UpAfterTrans(l.ctx, msg, resp)
		if err != nil {
			l.Error(msg, resp, err)
		}
	})
	l.Infof("%s req:%v resp:%v err:%v", utils.FuncName(), msg, resp, err)
	return err
}

func (l *DeviceMsgHandle) SDKLog(msg *deviceMsg.PublishMsg) error {
	msg = l.svcCtx.ScriptTrans.UpBeforeTrans(l.ctx, msg)
	l.FixDisconnect(msg)
	respMsg, err := NewSDKLogLogic(l.ctx, l.svcCtx).Handle(msg)
	l.deviceResp(respMsg)
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		err := l.svcCtx.ScriptTrans.UpAfterTrans(l.ctx, msg, respMsg)
		if err != nil {
			l.Error(msg, respMsg, err)
		}
	})
	l.Infof("%s req:%v resp:%v err:%v", utils.FuncName(), msg, respMsg, err)
	return err
}

func (l *DeviceMsgHandle) Connected(msg *deviceStatus.ConnectMsg) error {
	l.Infof("%s req=%v", utils.FuncName(), msg)
	return NewConnectedLogic(l.ctx, l.svcCtx).Handle(msg)
}

func (l *DeviceMsgHandle) Disconnected(msg *deviceStatus.ConnectMsg) error {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(msg))
	return NewDisconnectedLogic(l.ctx, l.svcCtx).Handle(msg)
}
func (l *DeviceMsgHandle) deviceResp(respMsg *deviceMsg.PublishMsg) {
	if respMsg == nil {
		return
	}
	startTime := time.Now()
	respMsg = l.svcCtx.ScriptTrans.DownBeforeTrans(l.ctx, respMsg)
	er := l.svcCtx.PubDev.PublishToDev(l.ctx, respMsg)
	if er != nil {
		l.Errorf("DeviceMsgHandle.deviceResp.PublishToDev failure err:%v", er)
		return
	}
	l.WithDuration(time.Now().Sub(startTime)).Infof("DeviceMsgHandle.deviceResp startTime:%v msg:%v", startTime, respMsg)
}
func (l *DeviceMsgHandle) FixDisconnect(msg *deviceMsg.PublishMsg) {
	return
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		di, err := l.svcCtx.DeviceCache.GetData(ctx, devices.Core{
			ProductID:  msg.ProductID,
			DeviceName: msg.DeviceName,
		})
		if err != nil {
			l.Error(err)
			return
		}
		if di.IsOnline != def.True { //如果不在线但是上报了需要调整为在线状态
			var updates = map[string]any{"is_online": def.True, "last_login": time.UnixMilli(msg.Timestamp), "status": def.DeviceStatusOnline}
			if di.FirstLogin == 0 {
				updates["first_login"] = time.UnixMilli(msg.Timestamp)
			}
			err = relationDB.NewDeviceInfoRepo(ctx).UpdateWithField(ctx,
				relationDB.DeviceFilter{Cores: []*devices.Core{{ProductID: msg.ProductID, DeviceName: msg.DeviceName}}}, updates)
			if err != nil {
				l.Error(err)
			}
			err = l.svcCtx.DeviceCache.SetData(ctx, devices.Core{
				ProductID:  msg.ProductID,
				DeviceName: msg.DeviceName,
			}, nil)
			if err != nil {
				l.Error(err)
			}
			appMsg := application.ConnectMsg{
				Device: devices.Core{
					ProductID:  msg.ProductID,
					DeviceName: msg.DeviceName,
				},
				Status:    def.ConnectedStatus,
				Timestamp: msg.Timestamp,
			}
			err = l.svcCtx.PubApp.DeviceStatusConnected(ctx, appMsg)
			if err != nil {
				l.Errorf("%s.pubApp productID:%v deviceName:%v err:%v",
					utils.FuncName(), msg.ProductID, msg.DeviceName, err)
			}
		}
	})
}

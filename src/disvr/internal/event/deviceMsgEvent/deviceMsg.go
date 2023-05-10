package deviceMsgEvent

//设备的发布,连接及断连处理
import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceStatus"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
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
	resp, err := NewGatewayLogic(l.ctx, l.svcCtx).Handle(msg)
	l.deviceResp(resp)
	l.Infof("%s req:%v resp:%v err:%v", utils.FuncName(), msg, resp, err)
	return err
}

func (l *DeviceMsgHandle) Thing(msg *deviceMsg.PublishMsg) error {
	l.Infof("%s req=%v", utils.FuncName(), msg)
	resp, err := NewThingLogic(l.ctx, l.svcCtx).Handle(msg)
	l.deviceResp(resp)
	l.Infof("%s req:%v resp:%v err:%v", utils.FuncName(), msg, resp, err)
	return err
}

func (l *DeviceMsgHandle) Ota(msg *deviceMsg.PublishMsg) error {
	l.Infof("%s req=%v", utils.FuncName(), msg)
	resp, err := NewOtaLogic(l.ctx, l.svcCtx).Handle(msg)
	l.deviceResp(resp)
	l.Infof("%s req:%v resp:%v err:%v", utils.FuncName(), msg, resp, err)
	return err
}

func (l *DeviceMsgHandle) Shadow(msg *deviceMsg.PublishMsg) error {
	l.Infof("%s req=%v", utils.FuncName(), msg)
	resp, err := NewShadowLogic(l.ctx, l.svcCtx).Handle(msg)
	l.deviceResp(resp)
	l.Infof("%s req:%v resp:%v err:%v", utils.FuncName(), msg, resp, err)
	return err
}

func (l *DeviceMsgHandle) Config(msg *deviceMsg.PublishMsg) error {
	l.Infof("%s req=%v", utils.FuncName(), msg)
	respMsg, err := NewConfigLogic(l.ctx, l.svcCtx).Handle(msg)
	l.deviceResp(respMsg)
	l.Infof("%s req:%v resp:%v err:%v", utils.FuncName(), msg, respMsg, err)
	return err
}

func (l *DeviceMsgHandle) SDKLog(msg *deviceMsg.PublishMsg) error {
	respMsg, err := NewSDKLogLogic(l.ctx, l.svcCtx).Handle(msg)
	l.deviceResp(respMsg)
	l.Infof("%s req:%v resp:%v err:%v", utils.FuncName(), msg, respMsg, err)
	return err
}

func (l *DeviceMsgHandle) Connected(msg *deviceStatus.ConnectMsg) error {
	l.Infof("%s req=%v", utils.FuncName(), msg)
	return NewConnectedLogic(l.ctx, l.svcCtx).Handle(msg)
}

func (l *DeviceMsgHandle) Disconnected(msg *deviceStatus.ConnectMsg) error {
	l.Infof("%s req=%v", utils.FuncName(), msg)
	return NewDisconnectedLogic(l.ctx, l.svcCtx).Handle(msg)
}
func (l *DeviceMsgHandle) deviceResp(respMsg *deviceMsg.PublishMsg) {
	if respMsg == nil {
		return
	}
	er := l.svcCtx.PubDev.PublishToDev(l.ctx, respMsg)
	if er != nil {
		l.Errorf("DeviceMsgHandle.deviceResp.PublishToDev failure err:%v", er)
		return
	}
	l.Infof("DeviceMsgHandle.deviceResp.PublishToDev msg:%v", respMsg)
}

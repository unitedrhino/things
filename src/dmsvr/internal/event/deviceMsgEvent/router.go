package deviceMsgEvent

//设备的发布,连接及断连处理
import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
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

func (l *DeviceMsgHandle) Thing(msg *device.PublishMsg) error {
	l.Infof("%s|req=%+v", utils.FuncName(), utils.GetJson(msg))
	return NewThingLogic(l.ctx, l.svcCtx).Handle(msg)
}

func (l *DeviceMsgHandle) Ota(msg *device.PublishMsg) error {
	l.Infof("%s|req=%+v", utils.FuncName(), utils.GetJson(msg))
	return NewOtaLogic(l.ctx, l.svcCtx).Handle(msg)
}

func (l *DeviceMsgHandle) Shadow(msg *device.PublishMsg) error {
	l.Infof("%s|req=%+v", utils.FuncName(), utils.GetJson(msg))
	return NewShadowLogic(l.ctx, l.svcCtx).Handle(msg)
}

func (l *DeviceMsgHandle) Config(msg *device.PublishMsg) error {
	l.Infof("%s|req=%+v", utils.FuncName(), utils.GetJson(msg))
	return NewConfigLogic(l.ctx, l.svcCtx).Handle(msg)
}

func (l *DeviceMsgHandle) SDKLog(msg *device.PublishMsg) error {
	l.Infof("%s|req=%+v", utils.FuncName(), utils.GetJson(msg))
	return NewSDKLogLogic(l.ctx, l.svcCtx).Handle(msg)
}

func (l *DeviceMsgHandle) Connected(msg *device.ConnectMsg) error {
	l.Infof("%s|req=%+v", utils.FuncName(), utils.GetJson(msg))
	return NewConnectedLogic(l.ctx, l.svcCtx).Handle(msg)
}

func (l *DeviceMsgHandle) Disconnected(msg *device.ConnectMsg) error {
	l.Infof("%s|req=%+v", utils.FuncName(), utils.GetJson(msg))
	return NewDisconnectedLogic(l.ctx, l.svcCtx).Handle(msg)
}

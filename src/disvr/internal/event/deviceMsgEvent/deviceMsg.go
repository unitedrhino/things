package deviceMsgEvent

//设备的发布,连接及断连处理
import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
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

func (l *DeviceMsgHandle) Thing(msg *deviceMsg.PublishMsg) error {
	l.Infof("%s|req=%v", utils.FuncName(), msg)
	return NewThingLogic(l.ctx, l.svcCtx).Handle(msg)
}

func (l *DeviceMsgHandle) Ota(msg *deviceMsg.PublishMsg) error {
	l.Infof("%s|req=%v", utils.FuncName(), msg)
	return NewOtaLogic(l.ctx, l.svcCtx).Handle(msg)
}

func (l *DeviceMsgHandle) Shadow(msg *deviceMsg.PublishMsg) error {
	l.Infof("%s|req=%v", utils.FuncName(), msg)
	return NewShadowLogic(l.ctx, l.svcCtx).Handle(msg)
}

func (l *DeviceMsgHandle) Config(msg *deviceMsg.PublishMsg) error {
	l.Infof("%s|req=%v", utils.FuncName(), msg)
	return NewConfigLogic(l.ctx, l.svcCtx).Handle(msg)
}

func (l *DeviceMsgHandle) SDKLog(msg *deviceMsg.PublishMsg) error {
	l.Infof("%s|req=%v", utils.FuncName(), msg)
	return NewSDKLogLogic(l.ctx, l.svcCtx).Handle(msg)
}

func (l *DeviceMsgHandle) Connected(msg *deviceMsg.ConnectMsg) error {
	l.Infof("%s|req=%v", utils.FuncName(), msg)
	return NewConnectedLogic(l.ctx, l.svcCtx).Handle(msg)
}

func (l *DeviceMsgHandle) Disconnected(msg *deviceMsg.ConnectMsg) error {
	l.Infof("%s|req=%v", utils.FuncName(), msg)
	return NewDisconnectedLogic(l.ctx, l.svcCtx).Handle(msg)
}

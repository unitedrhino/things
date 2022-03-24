package logic

import (
	"context"
	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/dmsvr/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendDeviceMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendDeviceMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendDeviceMsgLogic {
	return &SendDeviceMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

const (
	//Connect mqtt connect
	Connect = "connect"
	//Publish mqtt publish
	Publish = "publish"
	//Subscribe mqtt sub
	Subscribe = "subscribe"
	//Unsubscribe mqtt sub
	Unsubscribe = "unsubscribe"
	//Disconnect mqtt disconenct
	Disconnect = "disconnect"
)

// 设备端发送信息
func (l *SendDeviceMsgLogic) SendDeviceMsg(in *dm.SendDeviceMsgReq) (*dm.SendDeviceMsgResp, error) {
	// todo 去除kafka后该等待实现该接口
	//msg := &deviceSend.Elements{
	//	ClientID:  in.ClientID,
	//	Username:  in.Username,
	//	Topic:     in.Topic,
	//	Payload:   in.Payload,
	//	timestamp: in.timestamp,
	//	Action:    in.Action,
	//}
	//switch in.Action {
	//case Connect:
	//	return &dm.SendDeviceMsgResp{}, eventDevSub.NewDeviceMsgHandle(l.ctx,l.svcCtx).Disconnected(msg)
	//case Disconnect:
	//	return &dm.SendDeviceMsgResp{}, deviceSend.NewDisConnectLogic(l.ctx, l.svcCtx).Handle(msg)
	//case Publish:
	//	return &dm.SendDeviceMsgResp{}, deviceSend.NewPublishLogic(l.ctx, l.svcCtx).Handle(msg)
	//default:
	//	l.Errorf("SendDeviceMsgLogic|SendDeviceMsg|action is invalid:%v", in.Action)
	//	return nil, errors.Parameter.WithMsgf("action is invalid:%v", in.Action)
	//}
	return &dm.SendDeviceMsgResp{}, nil
}

package logic

import (
	"context"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/src/dmsvr/internal/exchange/logic"
	"github.com/go-things/things/src/dmsvr/internal/exchange/types"

	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/go-things/things/src/dmsvr/internal/svc"

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
	msg := &types.Elements{
		ClientID:  in.ClientID,
		Username:  in.Username,
		Topic:     in.Topic,
		Payload:   in.Payload,
		Timestamp: in.Timestamp,
		Action:    in.Action,
	}
	switch in.Action {
	case Connect:
		return &dm.SendDeviceMsgResp{}, logic.NewConnectLogic(l.ctx, l.svcCtx).Handle(msg)
	case Disconnect:
		return &dm.SendDeviceMsgResp{}, logic.NewDisConnectLogic(l.ctx, l.svcCtx).Handle(msg)
	case Publish:
		return &dm.SendDeviceMsgResp{}, logic.NewPublishLogic(l.ctx, l.svcCtx).Handle(msg)
	default:
		l.Errorf("SendDeviceMsgLogic|SendDeviceMsg|action is invalid:%v", in.Action)
		return nil, errors.Parameter.WithMsgf("action is invalid:%v", in.Action)
	}
	return &dm.SendDeviceMsgResp{}, nil
}

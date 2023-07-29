package logic

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/i-Things/things/src/stocksvr/internal/svc"
	"github.com/i-Things/things/src/stocksvr/types/pb/stock"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgLogic {
	return &SendMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMsgLogic) SendMsg(in *stock.ReqSendMsg) (out *stock.ResSendMsg, err error) {
	out = &stock.ResSendMsg{}
	dmReq := &di.SendMsgReq{
		Topic:   "$thing/down/event/26Md1xEqWfC/t1",
		Payload: []byte(`{"method":"eventReply","clientToken":"123","timestamp":1690095257922,"code":200,"status":"成功"}`),
	}
	if in.Topic != "" {
		dmReq.Topic = in.Topic
	}
	if in.Payload != "" {
		dmReq.Payload = []byte(in.Payload)
	}
	_, err = l.svcCtx.RpcClient.DeviceInteract.SendMsg(l.ctx, dmReq)
	if err != nil {
		l.Errorf("%s.rpc.SendMsg req=%v err=%+v", utils.FuncName(), dmReq, err)
		return
	}

	return
}

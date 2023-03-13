package interact

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgLogic {
	return &SendMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMsgLogic) SendMsg(req *types.DeviceInteractSendMsgReq) error {
	dmReq := &di.SendMsgReq{
		Topic:   req.Topic,
		Payload: []byte(req.Payload),
	}
	_, err := l.svcCtx.DeviceInteract.SendMsg(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.SendMsg req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}

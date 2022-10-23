package deviceinteractlogic

import (
	"context"
	"github.com/i-Things/things/shared/utils"

	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"

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

// 发送消息给设备
func (l *SendMsgLogic) SendMsg(in *di.SendMsgReq) (*di.SendMsgResp, error) {
	l.Infof("%s topic:%v payload:%v", utils.FuncName(), in.GetTopic(), string(in.GetPayload()))
	er := l.svcCtx.PubDev.PublishToDev(l.ctx, in.Topic, in.Payload)
	if er != nil {
		l.Errorf("%s.PublishToDev failure err:%v", utils.FuncName(), er)
		return nil, er
	}
	return &di.SendMsgResp{}, nil
}

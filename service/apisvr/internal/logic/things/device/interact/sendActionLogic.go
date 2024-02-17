package interact

import (
	"context"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type SendActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendActionLogic {
	return &SendActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 调用设备行为
func (l *SendActionLogic) SendAction(req *types.DeviceInteractSendActionReq) (resp *types.DeviceInteractSendActionResp, err error) {
	return NewActionSendLogic(l.ctx, l.svcCtx).ActionSend(req)
}

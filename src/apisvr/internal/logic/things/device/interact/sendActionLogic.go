package interact

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/disvr/pb/di"

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

//调用设备行为
func (l *SendActionLogic) SendAction(req *types.DeviceInteractSendActionReq) (resp *types.DeviceInteractSendActionResp, err error) {
	dmReq := &di.SendActionReq{
		ProductID:   req.ProductID,
		DeviceName:  req.DeviceName,
		ActionID:    req.ActionID,
		InputParams: req.InputParams,
		IsAsync:     req.IsAsync,
	}
	dmResp, err := l.svcCtx.DeviceInteract.SendAction(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.SendAction req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.DeviceInteractSendActionResp{
		ClientToken:  dmResp.ClientToken,
		OutputParams: dmResp.OutputParams,
		Status:       dmResp.Status,
		Code:         dmResp.Code,
	}, nil
}

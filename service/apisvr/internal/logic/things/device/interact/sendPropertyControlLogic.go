package interact

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendPropertyControlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendPropertyControlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPropertyControlLogic {
	return &SendPropertyControlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendPropertyControlLogic) SendPropertyControl(req *types.DeviceInteractSendPropertyReq) (resp *types.DeviceInteractSendPropertyResp, err error) {
	dmReq := &dm.SendPropertyControlReq{
		ProductID:     req.ProductID,
		DeviceName:    req.DeviceName,
		Data:          req.Data,
		IsAsync:       req.IsAsync,
		ShadowControl: req.ShadowControl,
	}
	dmResp, err := l.svcCtx.DeviceInteract.SendPropertyControl(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.SendProperty req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.DeviceInteractSendPropertyResp{
		Code:     dmResp.Code,
		Msg:      dmResp.Msg,
		MsgToken: dmResp.MsgToken,
	}, nil
}

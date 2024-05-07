package interact

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyControlSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPropertyControlSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyControlSendLogic {
	return &PropertyControlSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PropertyControlSendLogic) PropertyControlSend(req *types.DeviceInteractSendPropertyReq) (resp *types.DeviceInteractSendPropertyResp, err error) {
	dmReq := &dm.PropertyControlSendReq{
		ProductID:     req.ProductID,
		DeviceName:    req.DeviceName,
		Data:          req.Data,
		IsAsync:       req.IsAsync,
		ShadowControl: req.ShadowControl,
	}
	dmResp, err := l.svcCtx.DeviceInteract.PropertyControlSend(l.ctx, dmReq)
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

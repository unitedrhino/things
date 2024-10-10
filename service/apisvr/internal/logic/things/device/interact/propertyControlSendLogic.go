package interact

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

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
	dmResp, err := l.svcCtx.DeviceInteract.PropertyControlSend(l.ctx, utils.Copy[dm.PropertyControlSendReq](req))
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

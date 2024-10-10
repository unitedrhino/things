package interact

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActionSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActionSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionSendLogic {
	return &ActionSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ActionSendLogic) ActionSend(req *types.DeviceInteractSendActionReq) (resp *types.DeviceInteractSendActionResp, err error) {
	dmReq := &dm.ActionSendReq{
		ProductID:   req.ProductID,
		DeviceName:  req.DeviceName,
		ActionID:    req.ActionID,
		InputParams: req.InputParams,
		IsAsync:     req.IsAsync,
		Option:      logic.ToDmSendOption(req.Option),
	}
	dmResp, err := l.svcCtx.DeviceInteract.ActionSend(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.SendAction req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.DeviceInteractSendActionResp{
		MsgToken:     dmResp.MsgToken,
		OutputParams: dmResp.OutputParams,
		Msg:          dmResp.Msg,
		Code:         dmResp.Code,
	}, nil
}

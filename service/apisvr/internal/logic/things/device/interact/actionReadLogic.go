package interact

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActionReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActionReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionReadLogic {
	return &ActionReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ActionReadLogic) ActionRead(req *types.DeviceInteractRespReadReq) (resp *types.DeviceInteractSendActionResp, err error) {
	dmReq := &dm.RespReadReq{
		ProductID:  req.ProductID,
		DeviceName: req.DeviceName,
		MsgToken:   req.MsgToken,
	}
	dmResp, err := l.svcCtx.DeviceInteract.ActionRead(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ActionRead req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.DeviceInteractSendActionResp{
		MsgToken:     dmResp.MsgToken,
		OutputParams: dmResp.OutputParams,
		Msg:          dmResp.Msg,
		Code:         dmResp.Code,
	}, nil
}

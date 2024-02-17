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

type PropertyControlReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPropertyControlReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyControlReadLogic {
	return &PropertyControlReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyControlReadLogic) PropertyControlRead(req *types.DeviceInteractRespReadReq) (resp *types.DeviceInteractSendPropertyResp, err error) {
	dmReq := &dm.RespReadReq{
		ProductID:  req.ProductID,
		DeviceName: req.DeviceName,
		MsgToken:   req.MsgToken,
	}
	dmResp, err := l.svcCtx.DeviceInteract.PropertyControlRead(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.PropertyRead req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.DeviceInteractSendPropertyResp{
		Code:     dmResp.Code,
		Msg:      dmResp.Msg,
		MsgToken: dmResp.MsgToken,
	}, nil
}

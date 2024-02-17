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

type GetPropertyLatestReplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPropertyLatestReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPropertyLatestReplyLogic {
	return &GetPropertyLatestReplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPropertyLatestReplyLogic) GetPropertyLatestReply(req *types.DeviceInteractGetPropertyReplyReq) (resp *types.DeviceInteractGetPropertyReplyResp, err error) {
	dmReq := &dm.GetPropertyLatestReplyReq{
		ProductID:  req.ProductID,
		DeviceName: req.DeviceName,
		DataIDs:    req.DataIDs,
	}
	dmResp, err := l.svcCtx.DeviceInteract.GetPropertyLatestReply(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GetPropertyReply req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.DeviceInteractGetPropertyReplyResp{
		Code:      dmResp.Code,
		Msg:       dmResp.Msg,
		MsgToken:  dmResp.MsgToken,
		Timestamp: dmResp.Timestamp,
		Params:    dmResp.Params,
	}, nil
}

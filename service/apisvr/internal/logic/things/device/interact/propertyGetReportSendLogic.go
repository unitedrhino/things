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

type PropertyGetReportSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPropertyGetReportSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyGetReportSendLogic {
	return &PropertyGetReportSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PropertyGetReportSendLogic) PropertyGetReportSend(req *types.DeviceInteractPropertyGetReportSendReq) (resp *types.DeviceInteractPropertyGetReportSendResp, err error) {
	dmReq := &dm.PropertyGetReportSendReq{
		ProductID:  req.ProductID,
		DeviceName: req.DeviceName,
		DataIDs:    req.DataIDs,
	}
	dmResp, err := l.svcCtx.DeviceInteract.PropertyGetReportSend(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GetPropertyReply req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.DeviceInteractPropertyGetReportSendResp{
		Code:      dmResp.Code,
		Msg:       dmResp.Msg,
		MsgToken:  dmResp.MsgToken,
		Timestamp: dmResp.Timestamp,
		Params:    dmResp.Params,
	}, nil
}

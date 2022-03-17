package dc

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dcsvr/dc"

	"github.com/i-Things/things/src/webapi/internal/svc"
	"github.com/i-Things/things/src/webapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendPropertyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendPropertyLogic(ctx context.Context, svcCtx *svc.ServiceContext) SendPropertyLogic {
	return SendPropertyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendPropertyLogic) SendProperty(req types.SendDcPropertyReq) (*types.SendDcPropertyResp, error) {
	l.Infof("SendProperty|req=%+v", req)
	resp, err := l.svcCtx.DcRpc.SendProperty(l.ctx, &dc.SendPropertyReq{
		MemberID:      req.MemberID,
		MemberType:    req.MemberType,
		ProductID:     req.ProductID,
		DeviceName:    req.DeviceName,
		Data:          req.Data,
		DataTimestamp: req.DataTimestamp,
		Method:        req.Method,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.SendProperty|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.SendDcPropertyResp{
		Data:        resp.Data,
		ClientToken: resp.ClientToken, //调用id
		Status:      resp.Status,      //返回状态
		Code:        resp.Code,        //设备返回状态码
	}, nil
}

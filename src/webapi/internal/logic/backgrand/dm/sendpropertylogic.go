package dm

import (
	"context"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/shared/utils"
	"github.com/go-things/things/src/dmsvr/dm"

	"github.com/go-things/things/src/webapi/internal/svc"
	"github.com/go-things/things/src/webapi/internal/types"

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

func (l *SendPropertyLogic) SendProperty(req types.SendDmPropertyReq) (*types.SendDmPropertyResp, error) {
	l.Infof("SendProperty|req=%+v", req)
	resp, err := l.svcCtx.DmRpc.SendProperty(l.ctx, &dm.SendPropertyReq{
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
	return &types.SendDmPropertyResp{
		Data:        resp.Data,
		ClientToken: resp.ClientToken, //调用id
		Status:      resp.Status,      //返回状态
		Code:        resp.Code,        //设备返回状态码
	}, nil
}

package auth5

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic/things/device"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dgsvr/pb/dg"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccessLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAccessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessLogic {
	return &AccessLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccessLogic) Access(req *types.DeviceAuth5AccessReq) (resp *types.DeviceAuth5AccessResp, err error) {
	access := req.Action
	//如果是
	switch req.Action {
	case "subscribe":
		access = devices.Sub
	case "publish":
		access = devices.Pub
	}
	l.ctx = ctxs.WithRoot(l.ctx)
	_, err = l.svcCtx.DeviceA.AccessAuth(l.ctx, &dg.AccessAuthReq{
		Username: req.Username,
		Topic:    req.Topic,
		ClientID: req.ClientID,
		Access:   access,
		Ip:       req.Ip,
	})
	if err == nil {
		return &types.DeviceAuth5AccessResp{Result: "allow"}, nil
	}
	err = device.ThirdProtoAccessAuth(l.ctx, l.svcCtx, &types.DeviceAuthAccessReq{
		Username: req.Username,
		Topic:    req.Topic,
		ClientID: req.ClientID,
		Ip:       req.Ip,
	}, access)
	if err != nil {
		return &types.DeviceAuth5AccessResp{Result: "deny"}, nil
	}
	return &types.DeviceAuth5AccessResp{Result: "allow"}, nil
}

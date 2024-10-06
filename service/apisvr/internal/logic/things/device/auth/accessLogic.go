package auth

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/apisvr/internal/logic/things/device"
	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"
	"gitee.com/i-Things/things/service/dgsvr/pb/dg"

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

func (l *AccessLogic) Access(req *types.DeviceAuthAccessReq) error {
	l.Infof("%s req=%v", utils.FuncName(), req)
	access := req.Access
	//如果是
	switch req.Access {
	case "1":
		access = devices.Sub
	case "2":
		access = devices.Pub
	}
	l.ctx = ctxs.WithRoot(l.ctx)
	_, err := l.svcCtx.DeviceA.AccessAuth(l.ctx, &dg.AccessAuthReq{
		Username: req.Username,
		Topic:    req.Topic,
		ClientID: req.ClientID,
		Access:   access,
		Ip:       req.Ip,
	})
	if err == nil {
		return nil
	}

	return device.ThirdProtoAccessAuth(l.ctx, l.svcCtx, req, access)
}

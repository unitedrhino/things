package auth

import (
	"context"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic/things/device"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dgsvr/pb/dg"

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
	_, err := l.svcCtx.DeviceA.AccessAuth(l.ctx, &dg.AccessAuthReq{
		Username: req.Username,
		Topic:    req.Topic,
		ClientID: req.ClientID,
		Access:   access,
		Ip:       req.Ip,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.AccessAuth req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}

	return device.ThirdProtoAccessAuth(l.ctx, l.svcCtx, req, access)
}

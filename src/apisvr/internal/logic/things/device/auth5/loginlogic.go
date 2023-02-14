package auth5

import (
	"context"
	"encoding/base64"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.DeviceAuth5LoginReq) (resp *types.DeviceAuth5LoginResp, err error) {
	l.Infof("%s req=%+v", utils.FuncName(), req)
	var (
		cert []byte
	)

	// superuser
	_, err = l.svcCtx.DeviceA.RootCheck(l.ctx, &dm.RootCheckReq{
		Username:    req.Username,
		Password:    req.Password,
		ClientID:    req.ClientID,
		Ip:          req.Ip,
		Certificate: cert,
	})
	if err == nil {
		return &types.DeviceAuth5LoginResp{
			Result:      "allow",
			IsSuperuser: true,
		}, nil
	}
	// device auth
	if req.Certificate != "" {
		cert, err = base64.StdEncoding.DecodeString(req.Certificate)
		if err != nil {
			return nil, errors.Parameter.AddDetail("certificate can base64 decode")
		}

	}
	_, err = l.svcCtx.DeviceA.LoginAuth(l.ctx, &dm.LoginAuthReq{Username: req.Username, //用户名
		Password:    req.Password, //密码
		ClientID:    req.ClientID, //clientID
		Ip:          req.Ip,       //访问的ip地址
		Certificate: cert,         //客户端证书
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageDevice req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.DeviceAuth5LoginResp{
		Result:      "allow",
		IsSuperuser: false,
	}, nil
}

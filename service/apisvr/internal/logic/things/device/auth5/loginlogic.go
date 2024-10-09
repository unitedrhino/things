package auth5

import (
	"context"
	"encoding/base64"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/things/service/apisvr/internal/logic/things/device"
	"gitee.com/i-Things/things/service/dgsvr/pb/dg"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"

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
	var (
		cert []byte
	)
	// device auth
	if req.Certificate != "" {
		cert, err = base64.StdEncoding.DecodeString(req.Certificate)
		if err != nil {
			return nil, errors.Parameter.AddDetail("certificate can base64 decode")
		}
	}
	// superuser
	_, err = l.svcCtx.DeviceM.RootCheck(l.ctx, &dm.RootCheckReq{
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
	l.ctx = ctxs.WithRoot(l.ctx)
	_, er := l.svcCtx.DeviceA.LoginAuth(l.ctx, &dg.LoginAuthReq{Username: req.Username, //用户名
		Password:    req.Password, //密码
		ClientID:    req.ClientID, //clientID
		Ip:          req.Ip,       //访问的ip地址
		Certificate: cert,         //客户端证书
	})
	if er == nil {
		return &types.DeviceAuth5LoginResp{
			Result:      "allow",
			IsSuperuser: false,
		}, nil
	}
	err = device.ThirdProtoLoginAuth(l.ctx, l.svcCtx, &types.DeviceAuthLoginReq{
		Username:    req.Username,
		Password:    req.Password,
		ClientID:    req.ClientID,
		Ip:          req.Ip,
		Certificate: req.Certificate,
	}, cert)
	if err != nil {
		l.Errorf("auth5Login iThings err:%v third err:%v", er, err)
		return &types.DeviceAuth5LoginResp{Result: "deny"}, nil
	}
	return &types.DeviceAuth5LoginResp{
		Result:      "allow",
		IsSuperuser: false,
	}, nil
}

package logic

import (
	"context"
	"encoding/base64"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dmsvr/dm"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type LoginAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginAuthLogic {
	return LoginAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginAuthLogic) LoginAuth(req types.LoginAuthReq) (err error) {
	var cert []byte
	if req.Certificate != "" {
		cert, err = base64.StdEncoding.DecodeString(req.Certificate)
		if err != nil {
			return errors.Parameter.AddDetail("certificate can base64 decode")
		}

	}
	_, err = l.svcCtx.DmRpc.LoginAuth(l.ctx, &dm.LoginAuthReq{Username: req.Username, //用户名
		Password:    req.Password, //密码
		ClientID:    req.ClientID, //clientID
		Ip:          req.Ip,       //访问的ip地址
		Certificate: cert,         //客户端证书
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.ManageDevice|req=%v|err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}

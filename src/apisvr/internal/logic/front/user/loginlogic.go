package user

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/usersvr/user"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginLogic {
	return LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req types.LoginReq) (*types.LoginResp, error) {
	l.Infof("Login|req=%+v", req)
	if req.LoginType == "img" {
		if l.svcCtx.Captcha.Verify(req.CodeID, req.Code) == false {
			return nil, errors.Captcha
		}
	}
	resp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		UserID:    req.UserID,
		PwdType:   req.PwdType,
		Password:  req.Password,
		LoginType: req.LoginType,
		Code:      req.Code,
		CodeID:    req.CodeID,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.Login|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	if resp == nil {
		l.Errorf("%s|rpc.RegisterCore|return nil|req=%+v", utils.FuncName(), req)
		return nil, errors.System.AddDetail("register core rpc return nil")
	}
	return &types.LoginResp{
		*types.UserInfoToApi(resp.Info),
		types.JwtToken{
			AccessToken:  resp.Token.AccessToken,
			AccessExpire: resp.Token.AccessExpire,
			RefreshAfter: resp.Token.RefreshAfter,
		},
	}, nil
}

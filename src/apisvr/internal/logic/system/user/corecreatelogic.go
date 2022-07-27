package user

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/users"
	"github.com/i-Things/things/shared/utils"
	"time"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/usersvr/user"
	"github.com/zeromicro/go-zero/core/logx"
)

type CoreCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCoreCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoreCreateLogic {
	return &CoreCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CoreCreateLogic) CoreCreate(req *types.UserCoreCreateReq) (resp *types.UserCoreCreateResp, err error) {
	l.Infof("RegisterCore|req=%+v", req)
	resp1, err1 := l.svcCtx.UserRpc.CoreCreate(l.ctx, &user.UserCoreCreateReq{
		ReqType:  req.ReqType,
		Identity: req.Identity,
		Code:     req.Code,
		CodeID:   req.CodeID,
		Password: req.Password,
		Role:     req.Role,
	})
	if err1 != nil {
		er := errors.Fmt(err1)
		l.Errorf("[%s]|rpc.RegisterCore|req=%v|err=%#v|rpc_err=%+v", utils.FuncName(), req, er, err)
		return &types.UserCoreCreateResp{}, er
	}
	if resp1 == nil {
		l.Errorf("%s|rpc.RegisterCore|return nil|req=%+v", utils.FuncName(), req)
		return &types.UserCoreCreateResp{}, errors.System.AddDetail("register core rpc return nil")
	}
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Rej.AccessExpire
	jwtToken, err := users.GetJwtToken(l.svcCtx.Config.Rej.AccessSecret, now, accessExpire, resp1.Uid)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	return &types.UserCoreCreateResp{
		Uid: resp1.Uid,
		JwtToken: types.JwtToken{
			AccessToken:  jwtToken,
			AccessExpire: now + accessExpire,
			RefreshAfter: now + accessExpire/2,
		},
	}, nil

	return
}

package user

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/usersvr/user"
	"time"

	"github.com/i-Things/things/src/webapi/internal/svc"
	"github.com/i-Things/things/src/webapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterCoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterCoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) RegisterCoreLogic {
	return RegisterCoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterCoreLogic) RegisterCore(req types.RegisterCoreReq) (*types.RegisterCoreResp, error) {
	l.Infof("RegisterCore|req=%+v", req)
	resp, err := l.svcCtx.UserRpc.RegisterCore(l.ctx, &user.RegisterCoreReq{
		ReqType: req.ReqType,
		Note:    req.Note,
		Code:    req.Code,
		CodeID:  req.CodeID,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("[%s]|rpc.RegisterCore|req=%v|err=%#v|rpc_err=%+v", utils.FuncName(), req, er, err)
		return &types.RegisterCoreResp{}, er
	}
	if resp == nil {
		l.Errorf("%s|rpc.RegisterCore|return nil|req=%+v", utils.FuncName(), req)
		return &types.RegisterCoreResp{}, errors.System.AddDetail("register core rpc return nil")
	}
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Rej.AccessExpire
	jwtToken, err := utils.GetJwtToken(l.svcCtx.Config.Rej.AccessSecret, now, accessExpire, resp.Uid)
	if err != nil {
		return nil, errors.System.AddDetail(err.Error())
	}
	return &types.RegisterCoreResp{
		Uid: resp.Uid,
		JwtToken: types.JwtToken{
			AccessToken:  jwtToken,
			AccessExpire: now + accessExpire,
		},
	}, nil
}

package logic

import (
	"context"
	"time"
	"yl/shared/errors"
	"yl/shared/utils"
	"yl/src/user/rpc/user"

	"yl/src/user/api/internal/svc"
	"yl/src/user/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
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
	resp,err:=l.svcCtx.UserRpc.RegisterCore(l.ctx,&user.RegisterCoreReq{
		ReqType: req.ReqType,
		Note: req.Note,
		Code: req.Code,
		CodeID: req.CodeID,
	})
	if err != nil {
		er :=errors.Fmt(err)
		l.Errorf("[%s]|rpc.RegisterCore|req=%v|err=%#v",utils.FuncName(),req,er)
		return &types.RegisterCoreResp{},errors.System.AddDetail(er.Error())
	}
	if resp == nil {
		l.Errorf("%s|rpc.RegisterCore|return nil|req=%v",utils.FuncName(),req)
		return &types.RegisterCoreResp{},errors.System.AddDetail("register core rpc return nil")
	}
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Rej.AccessExpire
	jwtToken, err := utils.GetJwtToken(l.svcCtx.Config.Rej.AccessSecret, now, accessExpire, resp.Uid)
	if err != nil {
		return nil, errors.System.AddDetail(err.Error())
	}
	return &types.RegisterCoreResp{
		Uid: resp.Uid,
		JwtToken:types.JwtToken{
			AccessToken:  jwtToken,
			AccessExpire: now + accessExpire,
		},
	}, nil
}

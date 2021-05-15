package logic

import (
	"context"
	"yl/shared/errors"
	"yl/shared/utils"
	"yl/src/usersvr/user"

	"yl/src/webapi/internal/svc"
	"yl/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
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
	l.Infof("Login|req=%+v",req)
	if req.LoginType == "img"{
		if l.svcCtx.Captcha.Verify(req.CodeID,req.Code) == false {
			return nil,errors.Captcha
		}
	}
	resp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		UserID    :req.UserID,
		PwdType   :req.PwdType,
		Password  :req.Password,
		LoginType :req.LoginType,
		Code      :req.Code,
		CodeID    :req.CodeID,
		})
	if err != nil {
		er :=errors.Fmt(err)
		l.Errorf("%s|rpc.Login|req=%v|err=%+v",utils.FuncName(),req,er)
		return nil,er
	}
	if resp == nil {
		l.Errorf("%s|rpc.RegisterCore|return nil|req=%+v",utils.FuncName(),req)
		return nil,errors.System.AddDetail("register core rpc return nil")
	}
	return &types.LoginResp{
		types.UserInfo{
			Uid        : resp.Info.Uid,
			UserName   : resp.Info.UserName,
			NickName   : resp.Info.NickName,
			InviterUid : resp.Info.InviterUid,
			InviterId  : resp.Info.InviterId,
			Sex        : resp.Info.Sex,
			City       : resp.Info.City,
			Country    : resp.Info.Country,
			Province   : resp.Info.Province,
			Language   : resp.Info.Language,
			HeadImgUrl : resp.Info.HeadImgUrl,
			CreateTime : resp.Info.CreateTime,
		},
		types.JwtToken{
			AccessToken : resp.Token.AccessToken,
			AccessExpire: resp.Token.AccessExpire,
			RefreshAfter: resp.Token.RefreshAfter,
		},
	},nil
}

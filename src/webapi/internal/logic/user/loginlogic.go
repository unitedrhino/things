package logic

import (
	"context"
	"time"
	"yl/shared/errors"
	"yl/shared/utils"
	"yl/src/user/model"
	"yl/src/user/user"

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

func (l *LoginLogic)getRet(uc *model.UserCore)(*types.LoginResp, error){
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	jwtToken, err := utils.GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, uc.Uid)
	if err != nil {
		return nil, err
	}
	ui,err := l.svcCtx.UserInfoModel.FindOne(uc.Uid)
	return &types.LoginResp{
		Info: types.UserInfo{
			Uid         :ui.Uid,
			UserName    :uc.UserName,
			NickName    :ui.NickName,
			InviterUid  :ui.InviterUid,
			InviterId   :ui.InviterId,
			City        :ui.City,
			Country     :ui.Country,
			Province    :ui.Province,
			Language    :ui.Language,
			HeadImgUrl  :ui.Headimgurl,
			CreateTime :ui.CreatedTime.Time.Unix(),
		},
		Token: types.JwtToken{
			AccessToken:  jwtToken,
			AccessExpire: now + accessExpire,
			RefreshAfter: now + accessExpire/2,
		},
	}, nil
}

func (l *LoginLogic) Login(req types.LoginReq) (*types.LoginResp, error) {
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
		l.Errorf("[%s]|rpc.Login|req=%v|err=%#v",utils.FuncName(),req,er)
		return nil,er
	}
	if resp == nil {
		l.Errorf("%s|rpc.RegisterCore|return nil|req=%v",utils.FuncName(),req)
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

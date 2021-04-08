package logic

import (
	"context"
	"time"
	"yl/shared/define"
	"yl/shared/errors"
	"yl/shared/utils"
	"yl/src/user/model"

	"yl/src/user/rpc/internal/svc"
	"yl/src/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic)getRet(uc *model.UserCore)(*user.LoginResp, error){
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.UserToken.AccessExpire
	jwtToken, err := utils.GetJwtToken(l.svcCtx.Config.UserToken.AccessSecret, now, accessExpire, uc.Uid)
	if err != nil {
		return nil, err
	}
	ui,err := l.svcCtx.UserInfoModel.FindOne(uc.Uid)
	return &user.LoginResp{
		Info: &user.UserInfo{
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
		Token: &user.JwtToken{
			AccessToken:  jwtToken,
			AccessExpire: now + accessExpire,
			RefreshAfter: now + accessExpire/2,
		},
	}, nil
}


func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	var uc *model.UserCore
	var err error
	switch in.LoginType {
	case "sms"://暂时不验证
		uc,err=l.svcCtx.UserCoreModel.FindOneByPhone(in.UserID)
	case "img"://暂时不验证
		lt := utils.GetLoginNameType(in.UserID)
		switch lt {
		case define.Phone:
			uc,err=l.svcCtx.UserCoreModel.FindOneByPhone(in.UserID)
		default :
			uc,err=l.svcCtx.UserCoreModel.FindOneByUserName(in.UserID)
		}
	case "wxopen":
		l.Error("wxin not suppost")
	case "wxin":
		l.Error("wxin not suppost")
	default:
		return nil, errors.Parameter
	}
	switch err {
	case nil:
		return l.getRet(uc)
	case model.ErrNotFound:
		return nil, errors.UsernameUnRegister
	default:
		l.Errorf("%s|FindOneByPhone|req=%#v|err=%#v",utils.FuncName(),in,err)
		return nil, errors.System
	}

	return &user.LoginResp{}, nil
}

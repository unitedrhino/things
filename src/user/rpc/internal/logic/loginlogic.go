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
		l.Error(err)
		return nil, errors.System.AddDetail(err.Error()).ToRpc()
	}
	ui,err := l.svcCtx.UserInfoModel.FindOne(uc.Uid)
	if err != nil {
		l.Error(err)
		return nil, errors.System.AddDetail(err.Error()).ToRpc()
	}
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



func (l *LoginLogic)GetUserCore(in *user.LoginReq)(uc *model.UserCore,err error){
	switch in.LoginType {
	case "sms"://暂时不验证
		uc, err = l.svcCtx.UserCoreModel.FindOneByPhone(in.UserID)
	case "img"://暂时不验证
		lt := utils.GetLoginNameType(in.UserID)
		switch lt {
		case define.Phone:
			uc, err = l.svcCtx.UserCoreModel.FindOneByPhone(in.UserID)
		default :
			uc, err = l.svcCtx.UserCoreModel.FindOneByUserName(in.UserID)
		}
	case "wxminip"://微信小程序登录
		auth := l.svcCtx.WxMiniProgram.GetAuth()
		ret, err2 := auth.Code2Session(in.Code)
		if err2 != nil {
			l.Errorf("Code2Session|req=%#v|ret=%#v|err=%#v",in,ret,err2)
			if ret.ErrCode != 0 {
				return nil, errors.Parameter.AddDetail(ret.ErrMsg)
			}
			return nil, errors.System.AddDetail(err2.Error())
		} else if ret.ErrCode != 0 {
			return nil, errors.Parameter.AddDetail(ret.ErrMsg)
		}
		l.Slowf("login|wxminip|ret=%#v",ret)
		uc, err = l.svcCtx.UserCoreModel.FindOneByWechat(ret.UnionID)
	default:
		l.Error("LoginType=%s|not suppost",in.LoginType)
		return nil, errors.Parameter
	}
	l.Slowf("login|uc=%#v|err=%#v",uc,err)
	return uc, err
}


func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	defer func() {
		if p := recover(); p != nil {
			utils.HandleThrow(p)
		}
	}()
	uc,err := l.GetUserCore(in)
	switch err {
	case nil:
		return l.getRet(uc)
	case model.ErrNotFound:
		return nil, errors.UsernameUnRegister.ToRpc()
	default:
		l.Errorf("GetUserCore|req=%#v|err=%#v",in,err)
		return nil,errors.ToRpc(err)
	}
	return &user.LoginResp{}, nil
}

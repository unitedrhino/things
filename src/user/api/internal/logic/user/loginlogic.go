package logic

import (
	"context"
	"database/sql"
	"time"
	"yl/shared/define"
	"yl/shared/errors"
	"yl/shared/utils"
	"yl/src/user/model"

	"yl/src/user/api/internal/svc"
	"yl/src/user/api/internal/types"

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
		UserInfo: types.UserInfo{
			Uid         :ui.Uid,
			UserName    :uc.UserName.String,
			NickName    :ui.NickName.String,
			InviterUid  :ui.InviterUid,
			InviterId   :ui.InviterId.String,
			City        :ui.City.String,
			Country     :ui.Country.String,
			Province    :ui.Province.String,
			Language    :ui.Language.String,
			Headimgurl  :ui.Headimgurl.String,
			CreatedTime :ui.CreatedTime.Time.Unix(),
		},
		JwtToken: types.JwtToken{
			AccessToken:  jwtToken,
			AccessExpire: now + accessExpire,
			RefreshAfter: now + accessExpire/2,
		},
	}, nil
}

func (l *LoginLogic) Login(req types.LoginReq) (*types.LoginResp, error) {
	var uc *model.UserCore
	var err error
	switch req.LoginType {
	case "sms"://暂时不验证
		uc,err=l.svcCtx.UserCoreModel.FindOneByPhone(sql.NullString{String: req.UserID,Valid: true})
	case "image"://暂时不验证
		lt := utils.GetLoginNameType(req.UserID)
		switch lt {
		case define.Phone:
			uc,err=l.svcCtx.UserCoreModel.FindOneByPhone(sql.NullString{String: req.UserID,Valid: true})
		default :
			uc,err=l.svcCtx.UserCoreModel.FindOneByUserName(sql.NullString{String: req.UserID,Valid: true})
		}
	case "wxopen":
		logx.Error("wxin not suppost")
	case "wxin":
		logx.Error("wxin not suppost")
	default:
		return nil, errors.ErrorParameter
	}
	switch err {
	case nil:
		return l.getRet(uc)
	case model.ErrNotFound:
		return nil, errors.ErrorUsernameUnRegister
	default:
		logx.Errorf("%s|FindOneByPhone|req=%#v|err=%#v",utils.FuncName(),req,err)
		return nil, errors.ErrorSystem
	}
	return &types.LoginResp{}, nil
}

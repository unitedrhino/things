package userlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/users"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
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
func (l *LoginLogic) getPwd(in *sys.LoginReq, uc *mysql.UserCore) error {
	//根据密码类型不同做不同处理
	if in.PwdType == 0 {
		//空密码情况暂不考虑
		return errors.UnRegister
	} else if in.PwdType == 1 {
		//明文密码，则对密码做MD5加密后再与数据库密码比对
		//uid_temp := l.svcCtx.UserID.GetSnowflakeId()
		password1 := utils.MakePwd(in.Password, uc.Uid, false) //对密码进行md5加密
		if password1 != uc.Password {
			return errors.Password
		}
	} else if in.PwdType == 2 {
		//md5加密后的密码则通过二次md5加密再对比库中的密码
		password1 := utils.MakePwd(in.Password, uc.Uid, true) //对密码进行md5加密
		if password1 != uc.Password {
			return errors.Password
		}
	} else {
		return errors.UnRegister
	}
	return nil
}

func (l *LoginLogic) getRet(uc *mysql.UserCore) (*sys.LoginResp, error) {
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.UserToken.AccessExpire
	jwtToken, err := users.GetJwtToken(l.svcCtx.Config.UserToken.AccessSecret, now, accessExpire, uc.Uid)
	if err != nil {
		l.Error(err)
		return nil, errors.System.AddDetail(err)
	}
	ui, err := l.svcCtx.UserInfoModel.FindOne(l.ctx, uc.Uid)
	if err != nil {
		l.Errorf("FindOne|UserInfoModel|ui=%+v|err=%+v", ui, err)
		return nil, errors.Database.AddDetail(err)
	}

	uc, err1 := l.svcCtx.UserCoreModel.FindOne(l.ctx, uc.Uid)
	if err1 != nil {
		l.Errorf("FindOne|UserCoreModel|uc=%+v|err=%+v", uc, err)
		return nil, errors.Database.AddDetail(err)
	}

	resp := &sys.LoginResp{
		Info: &sys.UserInfo{
			Uid:        ui.Uid,
			UserName:   uc.UserName,
			NickName:   ui.NickName,
			InviterUid: ui.InviterUid,
			InviterId:  ui.InviterId,
			City:       ui.City,
			Country:    ui.Country,
			Province:   ui.Province,
			Language:   ui.Language,
			HeadImgUrl: ui.HeadImgUrl,

			//CreateTime: ui.CreatedTime.Time.Unix(),
		},
		Token: &sys.JwtToken{
			AccessToken:  jwtToken,
			AccessExpire: now + accessExpire,
			RefreshAfter: now + accessExpire/2,
		},
		Role: uc.AuthorityId,
	}
	l.Infof("Login|getRet=%+v", resp)
	return resp, nil
}

func (l *LoginLogic) GetUserCore(in *sys.LoginReq) (uc *mysql.UserCore, err error) {
	switch in.LoginType {
	case "wxopen":
		uc, err = l.svcCtx.UserCoreModel.FindOneByPhone(l.ctx, in.UserID)
	case "sms": //暂时不验证
		uc, err = l.svcCtx.UserCoreModel.FindOneByPhone(l.ctx, in.UserID)
	case "img": //暂时不验证
		lt := users.GetLoginNameType(in.UserID)
		switch lt {
		case users.Phone:
			//优先用手机号查表，如果未查到，说明可能是纯账号，则用账号字段查表
			uc, err = l.svcCtx.UserCoreModel.FindOneByPhone(l.ctx, in.UserID)
			if err != nil || uc == nil {
				uc, err = l.svcCtx.UserCoreModel.FindOneByUserName(l.ctx, in.UserID)
			}
		default:
			uc, err = l.svcCtx.UserCoreModel.FindOneByUserName(l.ctx, in.UserID)
		}
		if err := l.getPwd(in, uc); err != nil {
			return nil, err
		}

	case "wxminip": //微信小程序登录
		auth := l.svcCtx.WxMiniProgram.GetAuth()
		ret, err2 := auth.Code2Session(in.Code)
		if err2 != nil {
			l.Errorf("Code2Session|req=%#v|ret=%#v|err=%+v", in, ret, err2)
			if ret.ErrCode != 0 {
				return nil, errors.Parameter.AddDetail(ret.ErrMsg)
			}
			return nil, errors.System.AddDetail(err2.Error())
		} else if ret.ErrCode != 0 {
			return nil, errors.Parameter.AddDetail(ret.ErrMsg)
		}
		l.Infof("login|wxminip|ret=%+v", ret)
		uc, err = l.svcCtx.UserCoreModel.FindOneByWechat(l.ctx, ret.UnionID)
	default:
		l.Error("LoginType=%s|not suppost", in.LoginType)
		return nil, errors.Parameter
	}
	l.Infof("login|uc=%#v|err=%+v", uc, err)
	return uc, err
}

func (l *LoginLogic) Login(in *sys.LoginReq) (*sys.LoginResp, error) {
	l.Infof("Login|req=%+v", in)
	uc, err := l.GetUserCore(in)
	switch err {
	case nil:
		if uc.Status != users.NormalStatus {
			return nil, errors.UnRegister
		}
		return l.getRet(uc)
	case mysql.ErrNotFound:
		return nil, errors.UnRegister
	default:
		l.Errorf("GetUserCore|req=%#v|err=%+v", in, err)
		return nil, errors.Database.AddDetail(err)
	}
}

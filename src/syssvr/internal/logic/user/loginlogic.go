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
func (l *LoginLogic) getPwd(in *sys.LoginReq, uc *mysql.UserInfo) error {
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

func (l *LoginLogic) getRet(uc *mysql.UserInfo) (*sys.LoginResp, error) {
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.UserToken.AccessExpire
	jwtToken, err := users.GetJwtToken(l.svcCtx.Config.UserToken.AccessSecret, now, accessExpire, uc.Uid, uc.Role)
	if err != nil {
		l.Error(err)
		return nil, errors.System.AddDetail(err)
	}
	ui, err := l.svcCtx.UserInfoModel.FindOne(l.ctx, uc.Uid)
	if err != nil {
		l.Errorf("%s.FindOne.UserInfoModel ui=%v err=%v",
			utils.FuncName(), utils.Fmt(ui), utils.Fmt(err))
		return nil, errors.Database.AddDetail(err)
	}

	resp := &sys.LoginResp{
		Info: &sys.UserInfo{
			Uid:         ui.Uid,
			UserName:    ui.UserName,
			NickName:    ui.NickName,
			City:        ui.City,
			Country:     ui.Country,
			Province:    ui.Province,
			Language:    ui.Language,
			HeadImgUrl:  ui.HeadImgUrl,
			Email:       ui.Email,
			Phone:       ui.Phone,
			Wechat:      ui.Wechat,
			LastIP:      ui.LastIP,
			RegIP:       ui.RegIP,
			CreatedTime: ui.CreatedTime.Unix(),
			Role:        ui.Role,
			Sex:         ui.Sex,
		},
		Token: &sys.JwtToken{
			AccessToken:  jwtToken,
			AccessExpire: now + accessExpire,
			RefreshAfter: now + accessExpire/2,
		},
	}
	l.Infof("%s getRet=%+v", utils.FuncName(), resp)
	return resp, nil
}

func (l *LoginLogic) GetUserInfo(in *sys.LoginReq) (uc *mysql.UserInfo, err error) {
	switch in.LoginType {
	case "pwd":
		uc, err = l.svcCtx.UserInfoModel.FindOneByUserName(l.ctx, in.UserID)
		if err := l.getPwd(in, uc); err != nil {
			return nil, err
		}
	default:
		l.Error("%s LoginType=%s not support", utils.FuncName(), in.LoginType)
		return nil, errors.Parameter
	}
	l.Infof("%s uc=%#v err=%+v", utils.FuncName(), uc, err)
	return uc, err
}

func (l *LoginLogic) Login(in *sys.LoginReq) (*sys.LoginResp, error) {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))
	uc, err := l.GetUserInfo(in)
	switch err {
	case nil:
		return l.getRet(uc)
	case mysql.ErrNotFound:
		return nil, errors.UnRegister
	default:
		l.Errorf("%s req=%v err=%+v", utils.FuncName(), utils.Fmt(in), err)
		return nil, errors.Database.AddDetail(err)
	}
}

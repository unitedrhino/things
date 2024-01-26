package usermanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/users"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *LoginLogic) getPwd(in *sys.UserLoginReq, uc *relationDB.SysUserInfo) error {
	//根据密码类型不同做不同处理
	if in.PwdType == 0 {
		//空密码情况暂不考虑
		return errors.UnRegister
	} else if in.PwdType == 1 {
		//明文密码，则对密码做MD5加密后再与数据库密码比对
		//uid_temp := l.svcCtx.UserID.GetSnowflakeId()
		password1 := utils.MakePwd(in.Password, uc.UserID, false) //对密码进行md5加密
		if password1 != uc.Password {
			return errors.Password
		}
	} else if in.PwdType == 2 {
		//md5加密后的密码则通过二次md5加密再对比库中的密码
		password1 := utils.MakePwd(in.Password, uc.UserID, true) //对密码进行md5加密
		if password1 != uc.Password {
			return errors.Password
		}
	} else {
		return errors.Password
	}
	return nil
}

func (l *LoginLogic) getRet(ui *relationDB.SysUserInfo, list []*conf.LoginSafeCtlInfo) (*sys.UserLoginResp, error) {
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.UserToken.AccessExpire
	var rolses []int64
	var isAdmin int64 = def.False
	for _, v := range ui.Roles {
		rolses = append(rolses, v.RoleID)
	}

	if ui.Tenant != nil && (utils.SliceIn(ui.Tenant.AdminRoleID, rolses...) || ui.Tenant.AdminUserID == ui.UserID) {
		isAdmin = def.True
	}
	var account = ui.UserName.String
	if account == "" {
		account = ui.Phone.String
	}
	if account == "" {
		account = ui.Email.String
	}
	if account == "" {
		account = cast.ToString(ui.UserID)
	}
	jwtToken, err := users.GetLoginJwtToken(l.svcCtx.Config.UserToken.AccessSecret, now, accessExpire,
		ui.UserID, account, ctxs.GetUserCtx(l.ctx).TenantCode, rolses, ui.IsAllData, isAdmin)
	if err != nil {
		l.Error(err)
		return nil, errors.System.AddDetail(err)
	}

	//登录成功，清除密码错误次数相关redis key
	l.svcCtx.PwdCheck.ClearWrongpassKeys(l.ctx, list)

	InitCacheUserAuthProject(l.ctx, ui.UserID)
	InitCacheUserAuthArea(l.ctx, ui.UserID)

	resp := &sys.UserLoginResp{
		Info: UserInfoToPb(l.ctx, ui, l.svcCtx),
		Token: &sys.JwtToken{
			AccessToken:  jwtToken,
			AccessExpire: now + accessExpire,
			RefreshAfter: now + accessExpire/2,
		},
	}
	l.Infof("%s getRet=%+v", utils.FuncName(), resp)
	return resp, nil
}

func (l *LoginLogic) GetUserInfo(in *sys.UserLoginReq) (uc *relationDB.SysUserInfo, err error) {
	switch in.LoginType {
	case users.RegPwd:
		if l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeImage, def.CaptchaUseLogin, in.CodeID, in.Code) == "" {
			return nil, errors.Captcha
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Accounts: []string{in.Account}, WithRoles: true, WithTenant: true})
		if err != nil {
			return nil, err
		}
		if err = l.getPwd(in, uc); err != nil {
			return nil, err
		}
	case users.RegDingApp:
		cli, err := l.svcCtx.Cm.GetClients(l.ctx, "")
		if err != nil || cli.DingTalk == nil {
			return nil, errors.System.AddDetail(err)
		}
		ret, er := cli.DingTalk.GetUserInfoByCode(in.Code)
		if er != nil {
			return nil, errors.System.AddDetail(er)
		}
		if ret.Code != 0 {
			return nil, errors.Parameter.AddMsgf(ret.Msg)
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{DingTalkUserID: ret.UserInfo.UserId, WithRoles: true, WithTenant: true})
	case users.RegWxMiniP:
		cli, err := l.svcCtx.Cm.GetClients(l.ctx, "")
		if err != nil || cli.MiniProgram == nil {
			return nil, errors.System.AddDetail(err)
		}
		auth := cli.MiniProgram.GetAuth()
		ret, er := auth.Code2SessionContext(l.ctx, in.Code)
		if er != nil {
			return nil, errors.System.AddDetail(er)
		}
		if ret.ErrCode != 0 {
			return nil, errors.Parameter.AddMsgf(ret.ErrMsg)
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{WechatUnionID: ret.UnionID, WechatOpenID: ret.OpenID, WithRoles: true, WithTenant: true})
	case users.RegEmail:
		email := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, def.CaptchaUseLogin, in.CodeID, in.Code)
		if email == "" || email != in.Account {
			return nil, errors.Captcha
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Emails: []string{in.Account}, WithRoles: true, WithTenant: true})
	case users.RegPhone:
		phone := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypePhone, def.CaptchaUseLogin, in.CodeID, in.Code)
		if phone == "" || phone != in.Account {
			return nil, errors.Captcha
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Phones: []string{in.Account}, WithRoles: true, WithTenant: true})

	default:
		l.Error("%s LoginType=%s not support", utils.FuncName(), in.LoginType)
		return nil, errors.Parameter
	}
	l.Infof("%s uc=%#v err=%+v", utils.FuncName(), uc, err)
	return uc, err
}

func (l *LoginLogic) UserLogin(in *sys.UserLoginReq) (*sys.UserLoginResp, error) {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))
	//检查账号是否冻结
	list := l.svcCtx.Config.WrongPasswordCounter.ParseWrongPassConf(in.Account, in.Ip)
	if len(list) > 0 {
		forbiddenTime, f := l.svcCtx.PwdCheck.CheckAccountOrIpForbidden(l.ctx, list)
		if f {
			return nil, errors.Default.AddMsgf("%s %d 分钟", errors.AccountOrIpForbidden.Error(), forbiddenTime/60)
		}
	}
	uc, err := l.GetUserInfo(in)
	if err == nil {
		return l.getRet(uc, list)
	}
	if errors.Cmp(err, errors.NotFind) {
		return nil, errors.UnRegister
	}
	if errors.Cmp(err, errors.Password) {
		ret, err := l.svcCtx.PwdCheck.CheckPasswordTimes(l.ctx, list)
		if err != nil {
			return nil, err
		}
		if ret {
			return nil, errors.UseCaptcha
		}
		return nil, errors.Password
	}
	return nil, err
}

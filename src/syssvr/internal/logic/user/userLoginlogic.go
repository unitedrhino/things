package userlogic

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/users"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/cache"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
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
		return errors.UnRegister
	}
	return nil
}

func (l *LoginLogic) getRet(ui *relationDB.SysUserInfo, store kv.Store, list []*conf.LoginSafeCtlInfo) (*sys.UserLoginResp, error) {
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.UserToken.AccessExpire

	jwtToken, err := users.GetLoginJwtToken(l.svcCtx.Config.UserToken.AccessSecret, now, accessExpire, ui.UserID, ui.Role, ui.IsAllData)
	if err != nil {
		l.Error(err)
		return nil, errors.System.AddDetail(err)
	}

	//登录成功，清除密码错误次数相关redis key
	cache.ClearWrongpassKeys(l.ctx, store, list)

	resp := &sys.UserLoginResp{
		Info: UserInfoToPb(ui),
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
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Accounts: []string{in.Account}})
		if err != nil {
			return nil, err
		}
		if err = l.getPwd(in, uc); err != nil {
			return nil, err
		}
	/*企业版*/
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
		forbiddenTime, f := cache.CheckAccountOrIpForbidden(l.ctx, l.svcCtx.Store, list)
		if f {
			return nil, errors.Default.AddMsgf("%s %d 分钟", errors.AccountOrIpForbidden.Error(), forbiddenTime/60)
		}
	}
	uc, err := l.GetUserInfo(in)
	if err == nil {
		return l.getRet(uc, l.svcCtx.Store, list)
	}
	if errors.Cmp(err, errors.NotFind) {
		return nil, errors.UnRegister
	}
	if errors.Cmp(err, errors.Password) {
		ret, err := cache.CheckCaptchaTimes(l.ctx, l.svcCtx.Store, list)
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

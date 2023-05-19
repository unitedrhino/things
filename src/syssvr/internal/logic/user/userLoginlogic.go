package userlogic

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/users"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"time"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type LoginSafeCtlInfo struct {
	prefix    string
	key       string
	timeout   int
	times     int
	forbidden int
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *LoginLogic) getPwd(in *sys.UserLoginReq, uc *mysql.SysUserInfo) error {
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

func (l *LoginLogic) getRet(uc *mysql.SysUserInfo, store kv.Store, list []*LoginSafeCtlInfo) (*sys.UserLoginResp, error) {
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

	//登录成功，清除密码错误次数相关redis key
	clearWrongpassKeys(store, list)

	resp := &sys.UserLoginResp{
		Info: &sys.UserInfo{
			Uid:         ui.Uid,
			UserName:    ui.UserName.String,
			NickName:    ui.NickName,
			City:        ui.City,
			Country:     ui.Country,
			Province:    ui.Province,
			Language:    ui.Language,
			HeadImgUrl:  ui.HeadImgUrl,
			Email:       ui.Email.String,
			Phone:       ui.Phone.String,
			Wechat:      ui.Wechat.String,
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

func (l *LoginLogic) GetUserInfo(in *sys.UserLoginReq) (uc *mysql.SysUserInfo, err error) {
	switch in.LoginType {
	case "pwd":
		uc, err = l.svcCtx.UserInfoModel.FindOneByUserName(l.ctx, sql.NullString{String: in.UserID, Valid: true})
		if err != nil {
			return nil, err
		}
		if err = l.getPwd(in, uc); err != nil {
			return nil, err
		}
	default:
		l.Error("%s LoginType=%s not support", utils.FuncName(), in.LoginType)
		return nil, errors.Parameter
	}
	l.Infof("%s uc=%#v err=%+v", utils.FuncName(), uc, err)
	return uc, err
}

func clearWrongpassKeys(store kv.Store, list []*LoginSafeCtlInfo) {
	for _, v := range list {
		if v.prefix != "login:wrongPassword:ip:" {
			store.Del(v.key)
		}
	}
}

func parseWrongpassConf(counter conf.WrongPasswordCounter, userID string, ip string) []*LoginSafeCtlInfo {
	var res []*LoginSafeCtlInfo
	res = append(res, &LoginSafeCtlInfo{
		prefix:  "login:wrongPassword:captcha:",
		key:     "login:wrongPassword:captcha:" + userID,
		timeout: 24 * 3600,
		times:   counter.Captcha,
	})

	for i, v := range counter.Account {
		res = append(res, &LoginSafeCtlInfo{
			prefix:    "login:wrongPassword:account:",
			key:       "login:wrongPassword:account:" + cast.ToString(i+1) + ":" + userID,
			timeout:   v.Statistics * 60,
			times:     v.TriggerTimes,
			forbidden: v.ForbiddenTime * 60,
		})
	}
	for i, v := range counter.Ip {
		res = append(res, &LoginSafeCtlInfo{
			prefix:    "login:wrongPassword:ip:",
			key:       "login:wrongPassword:ip:" + cast.ToString(i+1) + ":" + ip,
			timeout:   v.Statistics * 60,
			times:     v.TriggerTimes,
			forbidden: v.ForbiddenTime * 60,
		})
	}

	return res
}

func checkAccountOrIpForbidden(store kv.Store, list []*LoginSafeCtlInfo) (int32, bool) {
	for _, v := range list {
		if v.prefix != "login:wrongPassword:captcha:" {
			ret, err := store.Get(v.key)
			if err != nil {
				continue
			}
			if cast.ToInt(ret) >= v.times {
				return int32(v.forbidden), true
			}
		}
	}
	return 0, false
}

func checkCaptchaTimes(store kv.Store, list []*LoginSafeCtlInfo) (bool, error) {
	for _, v := range list {
		ret, err := store.Get(v.key)
		if ret == "" {
			err = store.Setex(v.key, "1", v.timeout)
			if err != nil {
				return false, errors.Database.AddMsgf("创建 redis key：%s 失败", v.key)
			}
			continue
		}

		_, err = store.Incr(v.key)
		if err != nil {
			return false, errors.Database.AddMsgf("redis key：%s 自增失败", v.key)
		}
		if v.prefix != "login:wrongPassword:captcha:" {
			if cast.ToInt(ret)+1 >= v.times {
				err = store.Setex(v.key, cast.ToString(cast.ToInt(ret)+1), v.forbidden)
				if err != nil {
					return false, errors.Database.AddMsgf("重置 key：%s 时间失败", v.key)
				}
			}
		} else {
			if cast.ToInt(ret)+1 >= v.times {
				return true, nil
			}
		}
	}
	return false, nil
}

func (l *LoginLogic) UserLogin(in *sys.UserLoginReq) (*sys.UserLoginResp, error) {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))

	//检查账号是否冻结
	list := parseWrongpassConf(l.svcCtx.Config.WrongPasswordCounter, in.UserID, in.Ip)
	if len(list) > 0 {
		forbiddenTime, f := checkAccountOrIpForbidden(l.svcCtx.Store, list)
		if f {
			return nil, errors.Default.AddMsgf("%s %d 分钟", errors.AccountOrIpForbidden.Error(), forbiddenTime/60)
		}
	}

	uc, err := l.GetUserInfo(in)
	switch err {
	case nil:
		return l.getRet(uc, l.svcCtx.Store, list)
	case mysql.ErrNotFound:
		return nil, errors.UnRegister
	case errors.Password:
		ret, err := checkCaptchaTimes(l.svcCtx.Store, list)
		if err != nil {
			return nil, err
		}
		if ret {
			return nil, errors.UseCaptcha
		}
		return nil, errors.Password

	default:
		l.Errorf("%s req=%v err=%+v", utils.FuncName(), utils.Fmt(in), err)
		return nil, errors.Database.AddDetail(err)
	}

	return nil, nil
}

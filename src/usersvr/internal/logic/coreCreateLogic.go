package logic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/usersvr/internal/repo/mysql"
	"github.com/i-Things/things/src/usersvr/internal/svc"
	"github.com/i-Things/things/src/usersvr/pb/user"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
)

type CoreCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCoreCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoreCreateLogic {
	return &CoreCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *CoreCreateLogic) getRet(uc *mysql.UserCore) (*user.UserCoreCreateResp, error) {
	l.Infof("register_core|register one success|user_core=%#v", uc)
	return &user.UserCoreCreateResp{
		Uid: uc.Uid,
	}, nil
}

func (l *CoreCreateLogic) handlePhone(in *user.UserCoreCreateReq) (*user.UserCoreCreateResp, error) {
	if !utils.IsMobile(in.Identity) {
		return nil, errors.Parameter.AddDetail("is not phone number")
	}
	if in.CodeID != "6666" {
		return nil, errors.Captcha
	}
	//ip,err:=utils.GetIP(l.r)
	//fmt.Printf("ip=%s|err=%#v\n",ip)
	uc, err := l.svcCtx.UserCoreModel.FindOneByPhone(l.ctx, in.Identity)
	switch err {
	case nil: //如果已经有该账号,如果是注册了第一步,第二步没有注册,那么直接放行
		if uc.Status == def.NotRegistStatus {
			return l.getRet(uc)
		}
		return nil, errors.DuplicateRegister.AddDetail(in.Identity)
	case mysql.ErrNotFound: //如果没有注册过,那么注册账号并进入下一步
		uc := mysql.UserCore{
			Uid:   l.svcCtx.UserID.GetSnowflakeId(),
			Phone: in.Identity,
			//CreatedTime: sql.NullTime{Valid: true, Time: time.Now()},
		}
		_, err := l.svcCtx.UserModel.RegisterCore(uc, mysql.Keys{Key: "phone", Value: uc.Phone})
		if err != nil { //并发情况下有可能重复所以需要再次判断一次
			if err == mysql.ErrDuplicate {
				return nil, errors.DuplicateMobile.AddDetail(in.Identity)
			}
			l.Errorf("handlePhone|Inserts|err=%#v", err)
			break
		}
		return l.getRet(&uc)
	default:
		break
	}
	l.Errorf("handlePhone|err=%#v", err)
	return nil, errors.System.AddDetail(err)
}

func (l *CoreCreateLogic) handleWxminip(in *user.UserCoreCreateReq) (*user.UserCoreCreateResp, error) {
	auth := l.svcCtx.WxMiniProgram.GetAuth()
	ret, err2 := auth.Code2Session(in.Code)
	if err2 != nil {
		l.Errorf("Code2Session|req=%#v|ret=%#v|err=%#v", in, ret, err2)
		if ret.ErrCode != 0 {
			return nil, errors.Parameter.AddDetail(ret.ErrMsg)
		}
		return nil, errors.System.AddDetail(err2.Error())
	} else if ret.ErrCode != 0 {
		return nil, errors.Parameter.AddDetail(ret.ErrMsg)
	}

	uc, err := l.svcCtx.UserCoreModel.FindOneByWechat(l.ctx, ret.UnionID)
	switch err {
	case nil: //如果已经有该账号,如果是注册了第一步,第二步没有注册,那么直接放行
		if uc.Status == def.NotRegistStatus {
			return l.getRet(uc)
		}
		return nil, errors.DuplicateRegister
	case mysql.ErrNotFound: //如果没有注册过,那么注册账号并进入下一步
		uc := mysql.UserCore{
			Uid:    l.svcCtx.UserID.GetSnowflakeId(),
			Wechat: ret.UnionID,
			//CreatedTime: sql.NullTime{Valid: true, Time: time.Now()},
		}
		_, err := l.svcCtx.UserModel.RegisterCore(uc, mysql.Keys{Key: "wechat", Value: uc.Wechat})
		if err != nil {
			if err == mysql.ErrDuplicate {
				return nil, errors.DuplicateRegister.AddDetail(in.Identity)
			}
			l.Errorf("handlePhone|Inserts|err=%#v", err)
			return nil, errors.Database.AddDetail(err)
		}
		return l.getRet(&uc)
	default:
		l.Errorf("handlePhone|FindOneByWechat|err=%#v", err)
		return nil, errors.Database.AddDetail(err)
	}
}

func (l *CoreCreateLogic) CheckPwd(in *user.UserCoreCreateReq) error {
	if l.svcCtx.Config.UserOpt.NeedPassWord &&
		utils.CheckPasswordLever(in.Password) < l.svcCtx.Config.UserOpt.PassLevel {
		return errors.PasswordLevel
	}
	return nil
}

func (l *CoreCreateLogic) handlePassword(in *user.UserCoreCreateReq) (*user.UserCoreCreateResp, error) {
	//首先校验账号格式使用正则表达式，对用户账号做格式校验：只能是大小写字母，数字和下划线，减号
	ret := false
	if ret, _ = regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_-]{6,19}$", in.Identity); !ret {
		return nil, errors.UsernameFormatErr.AddDetail("账号必须以字母开头，且只能包含大小写字母和数字下划线和减号。 长度为6到20位之间")
	}
	err := l.CheckPwd(in)
	if err != nil {
		return nil, err
	}

	//如果是账密，则in.Note为账号
	uc, err := l.svcCtx.UserCoreModel.FindOneByUserName(l.ctx, in.Identity)
	switch err {
	case nil: //如果已经有该账号,如果是注册了第一步,第二步没有注册,那么直接放行
		if uc.Status == def.NotRegistStatus {
			return l.getRet(uc)
		}
		if in.ReqType == "password" {
			return nil, errors.DuplicateUsername.AddDetail(in.Identity)
		} else if in.ReqType == "phone" {
			return nil, errors.DuplicateMobile.AddDetail(in.Identity)
		} else {
			return nil, errors.DuplicateRegister.AddDetail(in.Identity)
		}
	case mysql.ErrNotFound: //如果没有注册过,那么注册账号并进入下一步
		uid_temp := l.svcCtx.UserID.GetSnowflakeId()
		password1 := utils.MakePwd(in.Password, uid_temp, false) //对密码进行md5加密
		uc := mysql.UserCore{
			Uid:         uid_temp,
			UserName:    in.Identity,
			Password:    password1,
			AuthorityId: in.Role,
		}

		_, err := l.svcCtx.UserModel.RegisterCore(uc, mysql.Keys{Key: "userName", Value: in.Identity})
		if err != nil { //并发情况下有可能重复所以需要再次判断一次
			if err == mysql.ErrDuplicate {
				return nil, errors.DuplicateMobile.AddDetail(in.Identity)
			}
			l.Errorf("handlePhone|Inserts|err=%#v", err)
			break
		}

		return &user.UserCoreCreateResp{Uid: uc.Uid}, nil
	default:
		break
	}

	return &user.UserCoreCreateResp{}, nil
}

func (l *CoreCreateLogic) CoreCreate(in *user.UserCoreCreateReq) (*user.UserCoreCreateResp, error) {
	l.Infof("RegisterCore|req=%+v", in)
	switch in.ReqType {
	case "password":
		return l.handlePassword(in)
	case "wxminip":
		return l.handleWxminip(in)
	case "phone":
		return l.handlePhone(in)
	default:
		l.Errorf("%s|ReqType=%s| not suppot yet", utils.FuncName(), in.ReqType)
		return nil, errors.Parameter.AddDetail("reqType not suppot yet :" + in.ReqType)
	}

	return &user.UserCoreCreateResp{}, nil
}

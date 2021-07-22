package logic

import (
	"context"
	"database/sql"
	"gitee.com/godLei6/things/shared/def"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/usersvr/model"
	"time"

	"gitee.com/godLei6/things/src/usersvr/internal/svc"
	"gitee.com/godLei6/things/src/usersvr/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type RegisterCoreLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterCoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterCoreLogic {
	return &RegisterCoreLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterCoreLogic) getRet(uc *model.UserCore) (*user.RegisterCoreResp, error) {
	l.Infof("register_core|register one success|user_core=%#v", uc)
	return &user.RegisterCoreResp{
		Uid: uc.Uid,
	}, nil
}

func (l *RegisterCoreLogic) handlePhone(in *user.RegisterCoreReq) (*user.RegisterCoreResp, error) {
	if !utils.IsMobile(in.Note) {
		return nil, errors.Parameter.AddDetail("is not phone number")
	}
	if in.CodeID != "6666" {
		return nil, errors.Captcha
	}
	//ip,err:=utils.GetIP(l.r)
	//fmt.Printf("ip=%s|err=%#v\n",ip)
	uc, err := l.svcCtx.UserCoreModel.FindOneByPhone(in.Note)
	switch err {
	case nil: //如果已经有该账号,如果是注册了第一步,第二步没有注册,那么直接放行
		if uc.Status == def.NotRegistStatus {
			return l.getRet(uc)
		}
		return nil, errors.DuplicateMobile.AddDetail(in.Note)
	case model.ErrNotFound: //如果没有注册过,那么注册账号并进入下一步
		uc := model.UserCore{
			Uid:         l.svcCtx.UserID.GetSnowflakeId(),
			Phone:       in.Note,
			CreatedTime: sql.NullTime{Valid: true, Time: time.Now()},
		}
		_, err := l.svcCtx.UserModel.RegisterCore(uc, model.Keys{Key: "phone", Value: uc.Phone})
		if err != nil { //并发情况下有可能重复所以需要再次判断一次
			if err == model.ErrDuplicate {
				return nil, errors.DuplicateMobile.AddDetail(in.Note)
			}
			l.Errorf("handlePhone|Inserts|err=%#v", err)
			break
		}
		return l.getRet(&uc)
	default:
		break
	}
	l.Errorf("handlePhone|err=%#v", err)
	return nil, errors.System.AddDetail(err.Error())
}

func (l *RegisterCoreLogic) handleWxminip(in *user.RegisterCoreReq) (*user.RegisterCoreResp, error) {
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

	uc, err := l.svcCtx.UserCoreModel.FindOneByWechat(ret.UnionID)
	switch err {
	case nil: //如果已经有该账号,如果是注册了第一步,第二步没有注册,那么直接放行
		if uc.Status == def.NotRegistStatus {
			return l.getRet(uc)
		}
		return nil, errors.DuplicateRegister
	case model.ErrNotFound: //如果没有注册过,那么注册账号并进入下一步
		uc := model.UserCore{
			Uid:         l.svcCtx.UserID.GetSnowflakeId(),
			Wechat:      ret.UnionID,
			CreatedTime: sql.NullTime{Valid: true, Time: time.Now()},
		}
		_, err := l.svcCtx.UserModel.RegisterCore(uc, model.Keys{Key: "wechat", Value: uc.Wechat})
		if err != nil {
			if err == model.ErrDuplicate {
				return nil, errors.DuplicateRegister.AddDetail(in.Note)
			}
			l.Errorf("handlePhone|Inserts|err=%#v", err)
			return nil, errors.Database.AddDetail(err.Error())
		}
		return l.getRet(&uc)
	default:
		l.Errorf("handlePhone|FindOneByWechat|err=%#v", err)
		return nil, errors.Database.AddDetail(err.Error())
	}
}

func (l *RegisterCoreLogic) RegisterCore(in *user.RegisterCoreReq) (*user.RegisterCoreResp, error) {
	l.Infof("RegisterCore|req=%+v", in)
	switch in.ReqType {
	case "wxminip":
		return l.handleWxminip(in)
	case "phone":
		return l.handlePhone(in)
	default:
		l.Errorf("%s|ReqType=%s| not suppot yet", utils.FuncName(), in.ReqType)
		return nil, errors.Parameter.AddDetail("reqType not suppot yet :" + in.ReqType)
	}
}

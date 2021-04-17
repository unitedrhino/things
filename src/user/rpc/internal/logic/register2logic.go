package logic

import (
	"context"
	"database/sql"
	"time"
	"yl/shared/define"
	"yl/shared/errors"
	"yl/shared/utils"
	"yl/src/user/model"

	"yl/src/user/rpc/internal/svc"
	"yl/src/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type Register2Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegister2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Register2Logic {
	return &Register2Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}


func (l *Register2Logic)register(in *user.Register2Req, uc *model.UserCore)(*user.Register2Resp, error){
	userInfo := model.UserInfo{
		Uid         :in.Info.Uid,
		UserName	:in.Info.UserName,
		NickName    :in.Info.NickName,
		Sex         :in.Info.Sex,
		City        :in.Info.City,
		Country     :in.Info.Country,
		Province    :in.Info.Province,
		Language    :in.Info.Language,
		CreatedTime: sql.NullTime{Valid: true,Time: time.Now()},
	}
	err := l.svcCtx.UserInfoModel.InsertOrUpdate(userInfo)
	if err != nil {
		return nil,errors.Database.AddDetail(err.Error())
	}
	uc.Status = define.NomalStatus
	uc.UserName = in.Info.UserName
	if uc.Password != ""{
		uc.Password = utils.MakePwd(in.Password,uc.Uid,false)
	}
	err = l.svcCtx.UserCoreModel.Update(*uc)
	if err != nil {
		return nil,errors.Database.AddDetail(err.Error())
	}
	return &user.Register2Resp{},nil
}

func (l *Register2Logic) CheckUserCore(in *user.Register2Req)(*model.UserCore, error){
	uc, err := l.svcCtx.UserCoreModel.FindOne(in.Info.Uid)
	switch err{
	case model.ErrNotFound: //如果没有注册过,那么注册账号并进入下一步
		return nil, errors.RegisterOne
	case nil://如果已经有该账号,如果是注册了第一步,第二步没有注册,那么直接放行
		if uc.Status != define.NotRegistStatus{
			return nil, errors.DuplicateRegister
		}
		return uc, nil
	default:
		l.Errorf("%s|FindOne|err=%#v",utils.FuncName(), err)
		return nil, errors.Database.AddDetail(err.Error())
	}
}




func (l *Register2Logic) CheckUserName(in *user.Register2Req) error{
	if in.Info.UserName == ""{//如果有用户名则需要检查密码,如果不需要填用户名则需要检测用户密码
		if l.svcCtx.Config.UserOpt.NeedUserName {
			return errors.NeedUserName
		}
		return nil
	}

	//检查用户名是否重复
	_,err := l.svcCtx.UserCoreModel.FindOneByUserName(in.Info.UserName)
	switch err {
	case nil:
		return errors.DuplicateUsername
	case model.ErrNotFound:
		break
	default:
		return errors.Database.AddDetail(err.Error())
	}
	return nil
}

func (l *Register2Logic) CheckPwd(in *user.Register2Req) error{
	if l.svcCtx.Config.UserOpt.NeedPassWord &&
		utils.CheckPasswordLever(in.Password) < l.svcCtx.Config.UserOpt.PassLevel{
		return errors.PasswordLevel
	}
	return nil
}

func (l *Register2Logic) CheckInfo(in *user.Register2Req)(err error){
	err = l.CheckUserName(in)
	if err != nil {
		return err
	}
	err = l.CheckPwd(in)
	if err != nil {
		return err
	}
	return nil
}

func (l *Register2Logic) Register2(in *user.Register2Req) (*user.Register2Resp, error) {
	err := l.CheckInfo(in)
	if err != nil {
		return nil,err
	}
	uc,err := l.CheckUserCore(in)
	if err != nil {
		return nil,err
	}
	return l.register(in,uc)
}

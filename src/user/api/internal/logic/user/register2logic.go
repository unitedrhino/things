package logic

import (
	"context"
	"database/sql"
	"yl/shared/define"
	"yl/shared/errors"
	"yl/shared/utils"
	"yl/src/user/model"

	"yl/src/user/api/internal/svc"
	"yl/src/user/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type Register2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegister2Logic(ctx context.Context, svcCtx *svc.ServiceContext) Register2Logic {
	return Register2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Register2Logic)register(req types.Register2Req, uc *model.UserCore)(error){
	_,err :=l.svcCtx.UserInfoModel.Insert(model.UserInfo{
		Uid         :req.Uid,
		NickName    :sql.NullString{String: req.NickName,Valid: true},
		InviterUid  :req.InviterUid,
		InviterId   :sql.NullString{String: req.NickName,Valid: true},
		Sex         :req.Sex,
		City        :sql.NullString{String: req.City,Valid: true},
		Country     :sql.NullString{String: req.Country,Valid: true},
		Province    :sql.NullString{String: req.Province,Valid: true},
		Language    :sql.NullString{String: req.Language,Valid: true},
	})
	if err != nil {
		return errors.ErrorSystem
	}
	uc.Status = define.NomalStatus
	uc.UserName = sql.NullString{String: req.UserName,Valid: true}
	uc.Password = sql.NullString{String: utils.MakePwd(req.Password,uc.Uid,false),Valid: true}
	err = l.svcCtx.UserCoreModel.Update(*uc)
	if err != nil {
		return errors.ErrorSystem
	}
	return nil
}


//注册完成后就需要填写用户信息,填写完成后才算注册成功(目前只有手机号注册登录需要走这步)
func (l *Register2Logic) Register2(req types.Register2Req) error {
	uc,err := l.svcCtx.UserCoreModel.FindOne(req.Uid)
	switch err{
	case nil://如果已经有该账号,如果是注册了第一步,第二步没有注册,那么直接放行
		if uc.Status != define.NotRegistStatus{
			return errors.ErrorDuplicateRegister
		}
		return l.register(req,uc)
	case model.ErrNotFound: //如果没有注册过,那么注册账号并进入下一步
		return errors.ErrorRegisterOne
	default:
		break
	}
	return errors.ErrorSystem
}

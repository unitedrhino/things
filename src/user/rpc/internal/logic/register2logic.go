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
	_,err :=l.svcCtx.UserInfoModel.Insert(model.UserInfo{
		Uid         :in.Info.Uid,
		UserName	:in.Info.UserName,
		NickName    :in.Info.NickName,
		InviterUid  :in.Info.InviterUid,
		InviterId   :in.Info.NickName,
		Sex         :in.Info.Sex,
		City        :in.Info.City,
		Country     :in.Info.Country,
		Province    :in.Info.Province,
		Language    :in.Info.Language,
		CreatedTime: sql.NullTime{Valid: true,Time: time.Now()},
	})
	if err != nil {
		return nil,errors.System.AddDetail(err.Error()).ToRpc()
	}
	uc.Status = define.NomalStatus
	uc.UserName = in.UserName
	uc.Password = utils.MakePwd(in.Password,uc.Uid,false)
	err = l.svcCtx.UserCoreModel.Update(*uc)
	if err != nil {
		return nil,errors.System.AddDetail(err.Error()).ToRpc()
	}
	return &user.Register2Resp{},nil
}



func (l *Register2Logic) Register2(in *user.Register2Req) (*user.Register2Resp, error) {
	uc,err := l.svcCtx.UserCoreModel.FindOne(in.Info.Uid)
	switch err{
	case nil://如果已经有该账号,如果是注册了第一步,第二步没有注册,那么直接放行
		if uc.Status != define.NotRegistStatus{
			return nil,errors.DuplicateRegister.ToRpc()
		}
		return l.register(in,uc)
	case model.ErrNotFound: //如果没有注册过,那么注册账号并进入下一步
		return nil,errors.RegisterOne.ToRpc()
	default:
		break
	}
	return nil,errors.System.ToRpc()
}

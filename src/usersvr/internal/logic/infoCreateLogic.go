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
)

type InfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func (l *InfoCreateLogic) register(in *user.UserInfoCreateReq, uc *mysql.UserCore) (*user.Response, error) {
	userInfo := UserInfoToDb(in.Info)
	err := l.svcCtx.UserInfoModel.InsertOrUpdate(l.ctx, *userInfo)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.UserInfoModel.Update(l.ctx, &mysql.UserInfo{
		Uid:        in.Info.Uid,
		UserName:   in.Info.UserName,
		NickName:   in.Info.NickName,
		InviterUid: in.Info.InviterUid,
		InviterId:  in.Info.InviterId,
		Sex:        in.Info.Sex,
		City:       in.Info.City,
		Country:    in.Info.Country,
		Province:   in.Info.Province,
		Language:   in.Info.Language,
		HeadImgUrl: in.Info.HeadImgUrl,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}

	//将user_core表中状态改成1
	err = l.svcCtx.UserCoreModel.Update(l.ctx, &mysql.UserCore{
		Uid:      uc.Uid,
		UserName: uc.UserName,
		Password: uc.Password,
		Email:    uc.Email,
		Phone:    uc.Phone,
		Wechat:   uc.Wechat,
		LastIP:   uc.LastIP,
		RegIP:    uc.RegIP,
		//	CreatedTime: uc.CreatedTime,
		//	UpdatedTime: uc.UpdatedTime,
		//	DeletedTime: uc.DeletedTime,
		Status:      1,
		AuthorityId: uc.AuthorityId,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &user.Response{}, nil
}

func (l *InfoCreateLogic) CheckUserCore(in *user.UserInfoCreateReq) (*mysql.UserCore, error) {
	uc, err := l.svcCtx.UserCoreModel.FindOne(l.ctx, in.Info.Uid)
	switch err {
	case mysql.ErrNotFound: //如果没有注册过,那么注册账号并进入下一步
		return nil, errors.RegisterOne
	case nil: //如果已经有该账号,如果是注册了第一步,第二步没有注册,那么直接放行
		if uc.Status != def.NotRegistStatus {
			return nil, errors.DuplicateRegister
		}
		return uc, nil
	default:
		l.Errorf("%s|FindOne|err=%#v", utils.FuncName(), err)
		return nil, errors.Database.AddDetail(err)
	}
}

func (l *InfoCreateLogic) CheckUserName(in *user.UserInfoCreateReq) error {
	if in.Info.UserName == "" { //如果有用户名则需要检查密码,如果不需要填用户名则需要检测用户密码
		if l.svcCtx.Config.UserOpt.NeedUserName {
			return errors.NeedUserName
		}
		return nil
	}

	//检查用户名是否重复
	_, err := l.svcCtx.UserCoreModel.FindOneByUserName(l.ctx, in.Info.UserName)
	switch err {
	case nil:
		return errors.DuplicateUsername
	case mysql.ErrNotFound:
		break
	default:
		return errors.Database.AddDetail(err)
	}
	return nil
}

func (l *InfoCreateLogic) CheckInfo(in *user.UserInfoCreateReq) (err error) {
	err = l.CheckUserName(in)
	if err != nil {
		return err
	}
	return nil
}

func NewInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoCreateLogic {
	return &InfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InfoCreateLogic) InfoCreate(in *user.UserInfoCreateReq) (*user.Response, error) {
	l.Infof("Register2|req=%+v", in)
	err := l.CheckInfo(in)
	if err != nil {
		return nil, err
	}
	uc, err := l.CheckUserCore(in)
	if err != nil {
		return nil, err
	}
	return l.register(in, uc)

	return &user.Response{}, nil
}

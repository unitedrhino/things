package logic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/usersvr/internal/repo/mysql"
	"github.com/spf13/cast"

	"github.com/i-Things/things/src/usersvr/internal/svc"
	"github.com/i-Things/things/src/usersvr/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModifyUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	ui *mysql.UserInfo
}

func NewModifyUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModifyUserInfoLogic {
	return &ModifyUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModifyUserInfoLogic) HandleUserName(value string) error {

	return nil
}

func (l *ModifyUserInfoLogic) Handle(key, value string) error {
	switch key {
	case "nickName":
		l.ui.NickName = value
		return nil
	case "inviterId":
		l.ui.InviterId = value
		return nil
	case "sex":
		l.ui.Sex = cast.ToInt64(value)
		return nil
	case "city":
		l.ui.City = value
		return nil
	case "country":
		l.ui.Country = value
		return nil
	case "province":
		l.ui.Province = value
		return nil
	case "language":
		l.ui.Language = value
		return nil
	case "headImgUrl":
		l.ui.HeadImgUrl = value
		return nil
	default:
		return errors.Parameter.AddDetail(key + "not support")
	}
}

func (l *ModifyUserInfoLogic) ModifyUserInfo(in *user.ModifyUserInfoReq) (*user.NilResp, error) {
	l.Infof("ModifyUserInfo|req=%+v", in)
	var err error
	l.ui, err = l.svcCtx.UserInfoModel.FindOne(in.Uid)
	if err != nil {
		l.Errorf("ModifyUserInfo|FindOne|uid=%d|err=%+v", in.Uid, err)
		return nil, errors.Database.AddDetail(err.Error())
	}
	for k, v := range in.Info {
		err := l.Handle(k, v)
		if err != nil {
			l.Errorf("ModifyUserInfo|Handle|key=%s|value=%s|ui=%+v|err=%+v", k, v, l.ui, err)
			return nil, err
		}
	}
	err = l.svcCtx.UserInfoModel.Update(*l.ui)
	if err != nil {
		l.Errorf("ModifyUserInfo|Update|ui=%+v|err=%+v", l.ui, err)
		return nil, errors.Database.AddDetail(err.Error())
	}
	l.Infof("ModifyUserInfo|modifyed usersvr info = %+v", l.ui)
	return &user.NilResp{}, nil
}

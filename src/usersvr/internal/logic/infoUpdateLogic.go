package logic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/usersvr/internal/repo/mysql"

	"github.com/i-Things/things/src/usersvr/internal/svc"
	"github.com/i-Things/things/src/usersvr/pb/user"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type InfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	ui *mysql.UserInfo
}

func NewInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoUpdateLogic {
	return &InfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InfoUpdateLogic) InfoUpdate(in *user.UserInfoUpdateReq) (*user.Response, error) {
	l.Infof("ModifyUserInfo|req=%+v", in)
	var err error
	l.ui, err = l.svcCtx.UserInfoModel.FindOne(l.ctx, cast.ToInt64(in.Uid))
	if err != nil {
		l.Errorf("ModifyUserInfo|FindOne|uid=%d|err=%+v", in.Uid, err)
		return nil, errors.Database.AddDetail(err)
	}
	l.ui.UserName = in.UserName
	l.ui.InviterUid = in.InviterUid
	l.ui.Sex = in.Sex
	l.ui.City = in.City
	l.ui.Language = in.Language
	l.ui.HeadImgUrl = in.HeadImgUrl
	l.ui.Province = in.Province
	l.ui.Country = in.Country
	l.ui.InviterId = in.InviterId
	err = l.svcCtx.UserInfoModel.Update(l.ctx, l.ui)
	if err != nil {
		l.Errorf("ModifyUserInfo|Update|ui=%+v|err=%+v", l.ui, err)
		return nil, errors.Database.AddDetail(err)
	}
	l.Infof("ModifyUserInfo|modifyed usersvr info = %+v", l.ui)
	return &user.Response{}, nil
}

package userlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/spf13/cast"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLogic) Update(in *sys.UserUpdateReq) (*sys.Response, error) {
	ui, err := l.svcCtx.UserInfoModel.FindOne(l.ctx, cast.ToInt64(in.Uid))
	if err != nil {
		l.Errorf("ModifyUserInfo|FindOne|uid=%d|err=%+v", in.Uid, err)
		return nil, errors.Database.AddDetail(err)
	}
	ui.UserName = in.UserName
	ui.Email = in.Email
	ui.Phone = in.Phone
	ui.Wechat = in.Wechat
	ui.NickName = in.NickName
	ui.Sex = in.Sex
	ui.City = in.City
	ui.Country = in.Country
	ui.Province = in.Province
	ui.Language = in.Language
	ui.HeadImgUrl = in.GetHeadImgUrl()

	err = l.svcCtx.UserInfoModel.Update(l.ctx, ui)
	if err != nil {
		l.Errorf("ModifyUserInfo|Update|ui=%+v|err=%+v", ui, err)
		return nil, errors.Database.AddDetail(err)
	}
	l.Infof("ModifyUserInfo|modifyed usersvr info = %+v", ui)

	return &sys.Response{}, nil
}

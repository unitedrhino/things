package userlogic

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
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

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLogic) UserUpdate(in *sys.UserUpdateReq) (*sys.Response, error) {
	ui, err := l.svcCtx.UserInfoModel.FindOne(l.ctx, cast.ToInt64(in.Uid))
	if err != nil {
		l.Errorf("%s.FindOne uid=%d err=%v", utils.FuncName(), in.Uid, err)
		return nil, errors.Database.AddDetail(err)
	}
	ui.UserName = sql.NullString{String: in.UserName, Valid: true}
	ui.NickName = in.NickName

	//性別有效才賦值，否則使用旧值
	if ui.Sex == 0 {
		ui.Sex = 1
	}
	if in.Sex == 1 || in.Sex == 2 {
		ui.Sex = in.Sex
	}

	//设置数据超管
	if in.IsAllData == 1 || in.IsAllData == 2 {
		ui.IsAllData = in.IsAllData
	}

	ui.City = in.City
	ui.Country = in.Country
	ui.Province = in.Province
	ui.Language = in.Language
	ui.HeadImgUrl = in.GetHeadImgUrl()
	if in.Role != 0 {
		ui.Role = in.Role
	}

	err = l.svcCtx.UserInfoModel.Update(l.ctx, ui)
	if err != nil {
		l.Errorf("%s.Update ui=%v err=%v", utils.FuncName(), ui, err)
		return nil, errors.Database.AddDetail(err)
	}
	l.Infof("%s.modified usersvr info = %+v", utils.FuncName(), ui)

	return &sys.Response{}, nil
}

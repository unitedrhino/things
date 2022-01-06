package user

import (
	"context"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/shared/utils"
	"github.com/go-things/things/src/usersvr/user"
	"github.com/go-things/things/src/webapi/internal/svc"
	"github.com/go-things/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserInfoLogic {
	return UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (*types.UserInfo, error) {
	userCtx := types.GetUserCtx(l.ctx)
	l.Infof("UserInfo|uid=%d", userCtx.Uid)
	ui, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &user.GetUserInfoReq{Uid: []int64{userCtx.Uid}})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("[%s]|rpc.Login|uid=%v|err=%+v", utils.FuncName(), userCtx.Uid, er)
		return nil, er
	}
	return types.UserInfoToApi(ui.Info[0]), nil
}

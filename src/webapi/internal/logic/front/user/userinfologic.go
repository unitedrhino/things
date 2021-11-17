package user

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/usersvr/user"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

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

func (l *UserInfoLogic) UserInfo(uid int64) (*types.UserInfo, error) {
	l.Infof("UserInfo|uid=%d", uid)
	ui, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &user.GetUserInfoReq{Uid: []int64{uid}})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("[%s]|rpc.Login|uid=%v|err=%+v", utils.FuncName(), uid, er)
		return nil, er
	}

	return types.UserInfoToApi(ui.Info[0]), nil
}

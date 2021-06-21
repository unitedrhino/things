package logic

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/src/usersvr/user"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ModifyUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) ModifyUserInfoLogic {
	return ModifyUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyUserInfoLogic) ModifyUserInfo(req types.ModifyUserInfoReq, uid int64) error {
	l.Infof("ModifyUserInfo|uid=%d|req=%+v", uid, req)
	_, err := l.svcCtx.UserRpc.ModifyUserInfo(l.ctx, &user.ModifyUserInfoReq{Info: req.Info, Uid: uid})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("ModifyUserInfo failure|err=%+v", er)
		return er
	}
	return nil
}

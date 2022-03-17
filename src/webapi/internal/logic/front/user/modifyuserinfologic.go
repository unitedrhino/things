package user

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/usersvr/user"

	"github.com/i-Things/things/src/webapi/internal/svc"
	"github.com/i-Things/things/src/webapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

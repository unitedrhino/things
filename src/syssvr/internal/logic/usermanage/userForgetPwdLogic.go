package usermanagelogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserForgetPwdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserForgetPwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserForgetPwdLogic {
	return &UserForgetPwdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserForgetPwdLogic) UserForgetPwd(in *sys.UserForgetPwdReq) (*sys.Response, error) {
	// todo: add your logic here and delete this line

	return &sys.Response{}, nil
}

package logic

import (
	"context"

	"yl/src/webapi/api/test/internal/svc"
	"yl/src/webapi/api/test/internal/types"

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
	// todo: add your logic here and delete this line

	return &types.UserInfo{}, nil
}

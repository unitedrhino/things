package usermanagelogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProjectMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserProjectMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProjectMultiUpdateLogic {
	return &UserProjectMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserProjectMultiUpdateLogic) UserProjectMultiUpdate(in *sys.UserProjectMultiUpdateReq) (*sys.Response, error) {
	// todo: add your logic here and delete this line

	return &sys.Response{}, nil
}

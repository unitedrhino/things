package usermanagelogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAreaMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAreaMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAreaMultiUpdateLogic {
	return &UserAreaMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAreaMultiUpdateLogic) UserAreaMultiUpdate(in *sys.UserAreaMultiUpdateReq) (*sys.Response, error) {
	// todo: add your logic here and delete this line

	return &sys.Response{}, nil
}

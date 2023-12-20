package usermanagelogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProjectIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserProjectIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProjectIndexLogic {
	return &UserProjectIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserProjectIndexLogic) UserProjectIndex(in *sys.UserProjectIndexReq) (*sys.UserProjectIndexResp, error) {
	// todo: add your logic here and delete this line

	return &sys.UserProjectIndexResp{}, nil
}

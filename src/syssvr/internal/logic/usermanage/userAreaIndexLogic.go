package usermanagelogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAreaIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAreaIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAreaIndexLogic {
	return &UserAreaIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAreaIndexLogic) UserAreaIndex(in *sys.UserAreaIndexReq) (*sys.UserAreaIndexResp, error) {
	// todo: add your logic here and delete this line

	return &sys.UserAreaIndexResp{}, nil
}

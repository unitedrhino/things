package loglogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogIndexLogic {
	return &LoginLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogIndexLogic) LoginLogIndex(in *sys.LoginLogIndexReq) (*sys.LoginLogIndexResp, error) {
	// todo: add your logic here and delete this line

	return &sys.LoginLogIndexResp{}, nil
}

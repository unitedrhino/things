package log

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginIndexLogic {
	return &LoginIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginIndexLogic) LoginIndex(req *types.SysLogLoginIndexReq) (resp *types.SysLogLoginIndexResp, err error) {
	// todo: add your logic here and delete this line

	return
}

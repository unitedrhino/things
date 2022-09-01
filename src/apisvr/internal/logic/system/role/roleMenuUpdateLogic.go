package role

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleMenuUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleMenuUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleMenuUpdateLogic {
	return &RoleMenuUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleMenuUpdateLogic) RoleMenuUpdate(req *types.RoleMenuUpdateReq) error {
	// todo: add your logic here and delete this line

	return nil
}

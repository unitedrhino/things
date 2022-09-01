package role

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleIndexLogic {
	return &RoleIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleIndexLogic) RoleIndex(req *types.RoleIndexReq) (resp *types.RoleIndexResp, err error) {
	// todo: add your logic here and delete this line

	return
}

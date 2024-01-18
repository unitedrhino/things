package role

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/role"
	"github.com/i-Things/things/src/syssvr/pb/sys"

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

func (l *RoleIndexLogic) RoleIndex(req *types.UserRoleIndexReq) (resp *types.UserRoleIndexResp, err error) {
	ret, err := l.svcCtx.UserRpc.UserRoleIndex(l.ctx, &sys.UserRoleIndexReq{UserID: req.UserID})
	if err != nil {
		return nil, err
	}

	return &types.UserRoleIndexResp{
		List:  role.ToRoleInfosTypes(ret.List),
		Total: ret.Total,
	}, nil
}

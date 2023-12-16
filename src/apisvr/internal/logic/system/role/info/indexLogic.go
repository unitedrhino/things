package info

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/role"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.RoleInfoIndexReq) (resp *types.RoleInfoIndexResp, err error) {
	info, err := l.svcCtx.RoleRpc.RoleInfoIndex(l.ctx, &sys.RoleInfoIndexReq{
		Page:   logic.ToSysPageRpc(req.Page),
		Name:   req.Name,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &types.RoleInfoIndexResp{role.ToRoleInfosTypes(info.List), info.Total}, nil
}

package info

import (
	"context"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TreeLogic {
	return &TreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TreeLogic) Tree(req *types.AccessIndexReq) (resp *types.AccessTreeResp, err error) {
	rst, err := l.svcCtx.AccessRpc.AccessInfoIndex(l.ctx, &sys.AccessInfoIndexReq{
		Group:      req.Group,
		Code:       req.Code,
		Name:       req.Name,
		IsNeedAuth: req.IsNeedAuth,
		WithApis:   req.WithApis,
	})
	if err != nil {
		return nil, err
	}
	return &types.AccessTreeResp{
		List: ToAccessGroupInfoTypes(rst.List),
	}, nil
}

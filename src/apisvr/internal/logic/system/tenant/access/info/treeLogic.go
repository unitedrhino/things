package info

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/access/info"
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

func (l *TreeLogic) Tree(req *types.TenantAccessInfoIndexReq) (resp *types.TenantAccessInfoTreeResp, err error) {
	rst, err := l.svcCtx.TenantRpc.TenantAccessIndex(l.ctx, &sys.TenantAccessIndexReq{Code: req.Code})
	if err != nil {
		return nil, err
	}
	ais, err := l.svcCtx.AccessRpc.AccessInfoIndex(l.ctx, &sys.AccessInfoIndexReq{Codes: rst.AccessCodes})
	if err != nil {
		return nil, err
	}
	return &types.TenantAccessInfoTreeResp{
		List:  info.ToAccessGroupInfoTypes(ais.List),
		Total: ais.Total,
	}, nil
}

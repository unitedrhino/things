package app

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/app/info"
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

func (l *IndexLogic) Index(req *types.TenantAppIndexReq) (resp *types.TenantAppIndexResp, err error) {
	ret, err := l.svcCtx.TenantRpc.TenantAppIndex(l.ctx, &sys.TenantAppIndexReq{Code: req.Code})
	return &types.TenantAppIndexResp{
		List: info.ToAppInfosTypes(ret.List),
	}, nil
}

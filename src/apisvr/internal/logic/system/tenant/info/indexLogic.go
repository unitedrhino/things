package info

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/logic/system"
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

func (l *IndexLogic) Index(req *types.TenantInfoIndexReq) (resp *types.TenantInfoIndexResp, err error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.TenantRpc.TenantInfoIndex(l.ctx, &sys.TenantInfoIndexReq{
		Name: req.Name,
		Page: logic.ToSysPageRpc(req.Page),
		Code: req.Code,
	})
	if err != nil {
		return nil, err
	}
	return &types.TenantInfoIndexResp{
		Total: ret.Total,
		List:  system.ToTenantInfosTypes(ret.List),
	}, nil
}

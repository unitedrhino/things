package module

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/logic"
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

func (l *IndexLogic) Index(req *types.TenantModuleIndexReq) (resp *types.TenantModuleIndexResp, err error) {
	ret, err := l.svcCtx.TenantRpc.TenantAppModuleIndex(l.ctx, &sys.TenantModuleIndexReq{
		Page:    logic.ToSysPageRpc(req.Page),
		Code:    req.Code,
		AppCode: req.AppCode,
	})
	if err != nil {
		return nil, err
	}
	return &types.TenantModuleIndexResp{
		ModuleCodes: ret.ModuleCodes,
	}, nil
}

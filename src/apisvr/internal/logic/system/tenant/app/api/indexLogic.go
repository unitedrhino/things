package api

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"

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

func (l *IndexLogic) Index(req *types.TenantAppApiIndexReq) (resp *types.TenantAppApiIndexResp, err error) {
	ret, err := l.svcCtx.TenantRpc.TenantAppApiIndex(l.ctx, &sys.TenantAppApiIndexReq{
		Page:       logic.ToSysPageRpc(req.Page),
		Code:       req.Code,
		AppCode:    req.AppCode,
		ModuleCode: req.ModuleCode,
	})
	if err != nil {
		return nil, err
	}
	var apiInfo []*types.TenantApiInfo
	apiInfo = make([]*types.TenantApiInfo, 0, len(apiInfo))
	for _, i := range ret.List {
		apiInfo = append(apiInfo, ToTenantAppApiTypes(i))
	}
	return &types.TenantAppApiIndexResp{
		Total: ret.Total,
		List:  apiInfo,
	}, nil
}
